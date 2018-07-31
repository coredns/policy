package policy

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/infobloxopen/go-trees/domain"
	pb "github.com/infobloxopen/themis/contrib/coredns/policy/dnstap"
	"github.com/miekg/dns"

	"github.com/infobloxopen/themis/pdp"
)

const (
	actionInvalid = iota
	actionRefuse
	actionAllow
	actionRedirect
	actionBlock
	actionLog
	actionDrop

	actionsTotal
)

const (
	actionNameInvalid  = "invalid"
	actionNameRefuse   = "refuse"
	actionNameAllow    = "allow"
	actionNameRedirect = "redirect"
	actionNameBlock    = "block"
	actionNameLog      = "log"
	actionNameDrop     = "drop"
	actionNamePass     = "pass"
)

var actionNames [actionsTotal]string

var emptyCtx *pdp.Context

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
	actionNames[actionInvalid] = actionNameInvalid
	actionNames[actionRefuse] = actionNameRefuse
	actionNames[actionAllow] = actionNameAllow
	actionNames[actionRedirect] = actionNameRedirect
	actionNames[actionBlock] = actionNameBlock
	actionNames[actionLog] = actionNameLog
	actionNames[actionDrop] = actionNameDrop

	emptyCtx, _ = pdp.NewContext(nil, 0, nil)
}

func newAttrHolderWithDnReq(w dns.ResponseWriter, r *dns.Msg, optMap map[uint16][]*edns0Opt, ag *AttrGauge) *attrHolder {
	hdrCount := ednsAttrsStart
	qName, qType := getNameAndType(r)
	dn, err := domain.MakeNameFromString(qName)
	if err != nil {
		panic(fmt.Errorf("Can't treat %q as domain name: %s", qName, err))
	}

	srcIP := getRemoteIP(w)
	if srcIP != nil {
		hdrCount++
	}

	ah := &attrHolder{
		dn:    qName,
		dnReq: make([]pdp.AttributeAssignment, hdrCount, 8),
	}

	ah.dnReq[0] = pdp.MakeStringAssignment(attrNameType, typeValueQuery)
	ah.dnReq[1] = pdp.MakeDomainAssignment(attrNameDomainName, dn)
	ah.dnReq[2] = pdp.MakeStringAssignment(attrNameDNSQtype, strconv.FormatUint(uint64(qType), 16))

	if srcIP != nil {
		ah.dnReq[3] = pdp.MakeAddressAssignment(attrNameSourceIP, srcIP)
	}

	extractOptionsFromEDNS0(r, optMap, func(b []byte, opts []*edns0Opt) {
		for _, o := range opts {
			if a, ok := makeAssignmentByType(o, b); ok {
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
	})

	return ah
}

func makeAssignmentByType(o *edns0Opt, b []byte) (pdp.AttributeAssignment, bool) {
	switch o.dataType {
	case typeEDNS0Bytes:
		return pdp.MakeStringAssignment(o.name, string(b)), true

	case typeEDNS0Hex:
		s := o.makeHexString(b)
		return pdp.MakeStringAssignment(o.name, s), s != ""

	case typeEDNS0IP:
		return pdp.MakeAddressAssignment(o.name, net.IP(b)), true
	}

	panic(fmt.Errorf("unknown attribute type %d", o.dataType))
}

func (ah *attrHolder) addDnRes(r *pdp.Response, custAttrs map[string]custAttr) {
	oCount := len(r.Obligations)

	switch r.Effect {
	default:
		log.Printf("[ERROR] PDP Effect: %s, Reason: %s", pdp.EffectNameFromEnum(r.Effect), r.Status)
		ah.action = actionInvalid

	case pdp.EffectPermit:
		ah.action = actionAllow

		i := 0
		for i < oCount {
			o := r.Obligations[i]

			id := o.GetID()
			switch id {
			case attrNameLog:
				ah.action = actionLog

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
		ah.action = actionBlock

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
				ah.action = actionRefuse

			case attrNameRedirectTo:
				ah.addRedirect(o)

			case attrNameDrop:
				ah.action = actionDrop
			}

			i++
		}
	}

	ah.dnRes = r.Obligations[:oCount]
}

func (ah *attrHolder) addRedirect(attr pdp.AttributeAssignment) {
	ah.action = actionRedirect
	dst, err := attr.GetString(emptyCtx)
	if err != nil {
		log.Printf("[ERROR] Action: %s, Destination: %s (%s)", actionNames[ah.action], serializeOrPanic(attr), err)

		ah.action = actionInvalid
		return
	}

	ah.dst = dst
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
		ah.action = actionInvalid

	case pdp.EffectPermit:
		ah.action = actionAllow

		for _, o := range r.Obligations {
			if o.GetID() == attrNameLog {
				ah.action = actionLog
				break
			}
		}

	case pdp.EffectDeny:
		ah.action = actionBlock

		for _, o := range r.Obligations {
			switch o.GetID() {
			case attrNameRefuse:
				ah.action = actionRefuse

			case attrNameRedirectTo:
				ah.addRedirect(o)

			case attrNameDrop:
				ah.action = actionDrop
			}
		}
	}

	ah.ipRes = r.Obligations
}

func (ah *attrHolder) makeDnstapReport() []*pb.DnstapAttribute {
	if ah.action != actionAllow && ah.action != actionInvalid {
		return ah.makeFullDnstapReport()
	}

	edns := ah.dnReq[ednsAttrsStart:]
	dnstap := ah.dnstap

	out := make([]*pb.DnstapAttribute, len(edns)+len(dnstap))
	n := putAttrsToDnstap(edns, out)
	putAttrsToDnstap(dnstap, out[n:])

	return out
}

func (ah *attrHolder) makeFullDnstapReport() []*pb.DnstapAttribute {
	lenIPReq := len(ah.ipReq)
	if lenIPReq > 0 {
		lenIPReq = 1
	}

	out := make([]*pb.DnstapAttribute, len(ah.dnReq)+len(ah.dnRes)+lenIPReq+len(ah.ipRes)+1)

	n := putAttrsToDnstap(ah.dnReq[1:], out)
	n += putAttrsToDnstap(ah.dnRes, out[n:])

	if lenIPReq > 0 {
		out[n] = newDnstapAttribute(ah.ipReq[ipReqAddrPos])
		n++
	}

	n += putAttrsToDnstap(ah.ipRes, out[n:])

	out[n] = newDnstapAttributeFromAction(ah.action)
	n++

	out[n] = newDnstapAttributeFromReqType(lenIPReq > 0)

	return out
}
