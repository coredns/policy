package themisplugin

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"context"
	"strings"

	"github.com/coredns/coredns/plugin/metadata"
	"github.com/coredns/policy/plugin/firewall/policy"
	rq "github.com/coredns/policy/plugin/pkg/rqdata"
	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/themis/pdp"
	"github.com/miekg/dns"
)

var emptyCtx *pdp.Context

type attrSetting struct {
	name     string
	label    string
	attrType string
	metrics  bool
}

type attrHolder struct {
	dn string

	dnReq []pdp.AttributeAssignment
	dnRes []pdp.AttributeAssignment

	transfer []pdp.AttributeAssignment

	ipReq []pdp.AttributeAssignment
	ipRes []pdp.AttributeAssignment

	dnstap []pdp.AttributeAssignment

	action byte
	dst    string
}

func init() {
	emptyCtx, _ = pdp.NewContext(nil, 0, nil)
}

func setAttrRequestValueMetadata(label string) {

}

func newAttrHolderWithContext(ctx context.Context, xtr *rq.Extractor, optMap []*attrSetting, ag *AttrGauge) *attrHolder {

	hdrCount := ednsAttrsStart
	qName, _ := xtr.Value("name")
	qType, _ := xtr.Value("qtype")
	dn, err := domain.MakeNameFromString(qName)
	if err != nil {
		panic(fmt.Errorf("Can't treat %q as domain name: %s", qName, err))
	}

	clientIP, _ := xtr.Value("remote")
	var srcIP = net.IP(nil)
	if clientIP != "" {
		srcIP = net.ParseIP(clientIP)
	}
	if srcIP != nil {
		hdrCount++
	}

	ah := &attrHolder{
		dn:    qName,
		dnReq: make([]pdp.AttributeAssignment, hdrCount, 8),
	}

	ah.dnReq[0] = pdp.MakeStringAssignment(attrNameType, typeValueQuery)
	ah.dnReq[1] = pdp.MakeDomainAssignment(attrNameDomainName, dn)
	ah.dnReq[2] = pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(dns.StringToType[qType]), 16))

	if srcIP != nil {
		ah.dnReq[3] = pdp.MakeAddressAssignment(attrNameSourceIP, srcIP)
	}

	for _, o := range optMap {
		f := metadata.ValueFunc(ctx, o.label)
		if f == nil {
			continue
		}
		value := f()
		if value == "" {
			continue
		}
		if a, ok := makeAssignmentByType(o, value); ok {
			if o.name == attrNameSourceIP && srcIP != nil {
				ah.dnReq[3] = a
			} else {
				ah.dnReq = append(ah.dnReq, a)
				if o.metrics && ag != nil {
					ag.Inc(a)
				}
			}
		}
	}

	return ah
}

func makeAssignmentByType(o *attrSetting, value string) (pdp.AttributeAssignment, bool) {
	switch strings.ToLower(o.attrType) {
	case "address":
		return pdp.MakeAddressAssignment(o.name, net.ParseIP(value)), true
	case "string":
		return pdp.MakeStringAssignment(o.name, value), true
	}
	panic(fmt.Errorf("unknown attribute type %s", o.attrType))
}

func (ah *attrHolder) addDnRes(r *pdp.Response, custAttrs map[string]custAttr) {
	oCount := len(r.Obligations)

	switch r.Effect {
	default:
		log.Printf("[ERROR] PDP Effect: %s, Reason: %s", pdp.EffectNameFromEnum(r.Effect), r.Status)
		ah.action = policy.TypeNone

	case pdp.EffectPermit:
		ah.action = policy.TypeAllow

		i := 0
		for i < oCount {
			o := r.Obligations[i]

			id := o.GetID()
			switch id {
			case attrNameLog:
				//ah.action = policy.TypeLog

			default:
				if t, ok := custAttrs[id]; ok {
					ah.putCustomAttr(o, t)

					if t.isEdns() {
						oCount--
						r.Obligations[i] = r.Obligations[oCount]
						continue
					}
				}
			}

			i++
		}

	case pdp.EffectDeny:
		ah.action = policy.TypeBlock

		i := 0
		for i < oCount {
			o := r.Obligations[i]

			id := o.GetID()
			switch id {
			default:
				if t, ok := custAttrs[id]; ok && t.isEdns() {
					ah.putCustomAttr(o, t)

					oCount--
					r.Obligations[i] = r.Obligations[oCount]
					continue
				}

			case attrNameRefuse:
				ah.action = policy.TypeRefuse

			case attrNameDrop:
				ah.action = policy.TypeDrop
			}

			i++
		}
	}

	ah.dnRes = r.Obligations[:oCount]
}


func (ah *attrHolder) putCustomAttr(attr pdp.AttributeAssignment, f custAttr) {
	if f.isEdns() {
		id := attr.GetID()

		for _, a := range ah.dnReq[ednsAttrsStart:] {
			if id == a.GetID() {
				return
			}
		}

		ah.dnReq = append(ah.dnReq, attr)
	}

	if f.isTransfer() {
		ah.transfer = append(ah.transfer, attr)
	}

	if f.isDnstap() {
		ah.dnstap = append(ah.dnstap, attr)
	}
}

func (ah *attrHolder) prepareResponseFromContext(ctx context.Context, xtr *rq.Extractor) {
	ipResp, _ := xtr.Value("response_ip")
	if ipResp != "" {
		ip := net.ParseIP(ipResp)
		if ip != nil {
			ah.addIPReq(ip)
		}
	}
}

func (ah *attrHolder) addIPReq(ip net.IP) {
	ah.ipReq = append(
		[]pdp.AttributeAssignment{
			pdp.MakeStringAssignment(attrNameType, typeValueResponse),
			pdp.MakeAddressAssignment(attrNameAddress, ip),
		},
		ah.transfer...,
	)
}

func (ah *attrHolder) addIPRes(r *pdp.Response) {
	switch r.Effect {
	default:
		log.Printf("[ERROR] PDP Effect: %s, Reason: %s", pdp.EffectNameFromEnum(r.Effect), r.Status)
		ah.action = policy.TypeNone

	case pdp.EffectPermit:
		ah.action = policy.TypeAllow

		//for _, o := range r.Obligations {
		//	if o.GetID() == attrNameLog {
		//		ah.action = firewall.TypeLog
		//		break
		//	}
		//}

	case pdp.EffectDeny:
		ah.action = policy.TypeBlock

		for _, o := range r.Obligations {
			switch o.GetID() {
			case attrNameRefuse:
				ah.action = policy.TypeRefuse

			//case attrNameRedirectTo:
			//	ah.addRedirect(o)

			case attrNameDrop:
				ah.action = policy.TypeDrop
			}
		}
	}

	ah.ipRes = r.Obligations
}
