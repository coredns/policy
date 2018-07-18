package policy

import (
	"log"
	"strconv"
	"time"

	"github.com/coredns/coredns/plugin/dnstap"
	tapmsg "github.com/coredns/coredns/plugin/dnstap/msg"
	"github.com/coredns/coredns/plugin/dnstap/taprw"
	tap "github.com/dnstap/golang-dnstap"
	"github.com/golang/protobuf/proto"
	pb "github.com/infobloxopen/themis/contrib/coredns/policy/dnstap"
	"github.com/miekg/dns"
	"golang.org/x/net/context"

	"github.com/infobloxopen/themis/pdp"
)

var dnstapActionValues [actionsTotal]string

type dnstapSender interface {
	sendCRExtraMsg(w dns.ResponseWriter, msg *dns.Msg, ah *attrHolder)
}

func init() {
	dnstapActionValues[actionInvalid] = strconv.Itoa(int(pb.PolicyAction_INVALID))
	dnstapActionValues[actionRefuse] = strconv.Itoa(int(pb.PolicyAction_REFUSE))
	dnstapActionValues[actionAllow] = strconv.Itoa(int(pb.PolicyAction_PASSTHROUGH))
	dnstapActionValues[actionRedirect] = strconv.Itoa(int(pb.PolicyAction_REDIRECT))
	dnstapActionValues[actionBlock] = strconv.Itoa(int(pb.PolicyAction_NXDOMAIN))
	dnstapActionValues[actionLog] = strconv.Itoa(int(pb.PolicyAction_PASSTHROUGH))
	dnstapActionValues[actionDrop] = strconv.Itoa(int(pb.PolicyAction_DENY))
}

type policyDnstapSender struct {
	ior dnstap.IORoutine
}

func newPolicyDnstapSender(io dnstap.IORoutine) dnstapSender {
	return &policyDnstapSender{ior: io}
}

// sendCRExtraMsg creates Client Response (CR) dnstap Message and writes an array
// of extra attributes to Dnstap.Extra field. Then it asynchronously sends the
// message with IORoutine interface
func (s *policyDnstapSender) sendCRExtraMsg(w dns.ResponseWriter, msg *dns.Msg, ah *attrHolder) {
	if w == nil || msg == nil {
		log.Printf("[ERROR] Failed to create dnstap CR message - no DNS response message found")
		return
	}

	now := time.Now()
	log.Printf("[ERROR] %q", w.RemoteAddr())
	b := tapmsg.New().Time(now).Addr(w.RemoteAddr())
	b.Msg(msg)
	crMsg, err := b.ToClientResponse()
	if err != nil {
		log.Printf("[ERROR] Failed to create dnstap CR message (%v)", err)
		return
	}

	timeNs := uint32(now.Nanosecond())
	crMsg.ResponseTimeNsec = &timeNs
	t := tap.Dnstap_MESSAGE

	var extra []byte
	if ah != nil {
		extra, err = proto.Marshal(&pb.Extra{Attrs: ah.makeDnstapReport()})
		if err != nil {
			log.Printf("[ERROR] Failed to create extra data for dnstap CR message (%v)", err)
		}
	}
	dnstapMsg := tap.Dnstap{Type: &t, Message: crMsg, Extra: extra}
	s.ior.Dnstap(dnstapMsg)
}

func resetCqCr(ctx context.Context) {
	if v := ctx.Value(dnstap.DnstapSendOption); v != nil {
		if so, ok := v.(*taprw.SendOption); ok {
			so.Cq = false
			so.Cr = false
		}
	}
}

func newDnstapAttribute(a pdp.AttributeAssignment) *pb.DnstapAttribute {
	return &pb.DnstapAttribute{
		Id:    a.GetID(),
		Value: serializeOrPanic(a),
	}
}

func newDnstapAttributeFromAction(a byte) *pb.DnstapAttribute {
	return &pb.DnstapAttribute{
		Id:    attrNamePolicyAction,
		Value: dnstapActionValues[a],
	}
}

func newDnstapAttributeFromReqType(ipReq bool) *pb.DnstapAttribute {
	out := &pb.DnstapAttribute{
		Id:    attrNameType,
		Value: typeValueQuery,
	}

	if ipReq {
		out.Value = typeValueResponse
	}

	return out
}

func putAttrsToDnstap(attrs []pdp.AttributeAssignment, out []*pb.DnstapAttribute) int {
	for i, a := range attrs {
		out[i] = newDnstapAttribute(a)
	}

	return len(attrs)
}
