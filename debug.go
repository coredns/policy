package policy

import (
	"strings"

	"github.com/miekg/dns"
)

const (
	queryInfoResolveYes        = "resolve:yes"
	queryInfoResolveNo         = "resolve:no"
	queryInfoResolveSkip       = "resolve:skip"
	queryInfoResolveFailed     = "resolve:failed"
	queryInfoActionPassthrough = "action:passthrough"
)

func (p *policyPlugin) patchDebugMsg(r *dns.Msg) bool {
	if r == nil || len(r.Question) <= 0 {
		return false
	}

	q := r.Question[0]
	if q.Qclass != dns.ClassCHAOS || q.Qtype != dns.TypeTXT || !strings.HasSuffix(q.Name, p.conf.debugSuffix) {
		return false
	}

	q.Name = dns.Fqdn(strings.TrimSuffix(q.Name, p.conf.debugSuffix))
	q.Qtype = dns.TypeA
	q.Qclass = dns.ClassINET

	r.Question[0] = q

	return true
}

func (p *policyPlugin) setDebugQueryPassthroughAnswer(ah *attrHolder, r *dns.Msg) {
	r.Answer = p.makeDebugQueryAnswerRR(ah, queryInfoActionPassthrough)
	r.Ns = nil
}

func (p *policyPlugin) setDebugQueryAnswer(ah *attrHolder, r *dns.Msg, status int) {
	info := getResolveStatusQueryInfo(status)

	action := actionNamePass
	if len(ah.ipRes) <= 0 {
		action = actionNames[ah.action]
	}
	info = appendQueryInfo(info, typeValueQuery, action)

	for _, a := range ah.dnRes {
		info = appendQueryInfo(info, a.GetID(), serializeOrPanic(a))
	}

	if len(ah.ipRes) > 0 {
		info = appendQueryInfo(info, typeValueResponse, actionNames[ah.action])

		for _, a := range ah.ipRes {
			info = appendQueryInfo(info, a.GetID(), serializeOrPanic(a))
		}
	}

	if len(p.conf.debugID) > 0 {
		info = appendQueryInfo(info, "ident", p.conf.debugID)
	}

	r.Answer = p.makeDebugQueryAnswerRR(ah, info)
}

func (p *policyPlugin) makeDebugQueryAnswerRR(ah *attrHolder, info string) []dns.RR {
	return []dns.RR{
		&dns.TXT{
			Hdr: dns.RR_Header{
				Name:   ah.dn + p.conf.debugSuffix,
				Rrtype: dns.TypeTXT,
				Class:  dns.ClassCHAOS,
			},
			Txt: []string{info},
		},
	}
}

func getResolveStatusQueryInfo(status int) string {
	switch status {
	case -1:
		return queryInfoResolveSkip

	case dns.RcodeSuccess:
		return queryInfoResolveYes

	case dns.RcodeServerFailure:
		return queryInfoResolveFailed
	}

	return queryInfoResolveNo
}

func appendQueryInfo(info, key, value string) string {
	return info + "," + key + ":'" + value + "'"
}
