package policy

import "github.com/infobloxopen/themis/pdp"

const (
	attrNameType         = "type"
	attrNameDomainName   = "domain_name"
	attrNameDNSQtype     = "dns_qtype"
	attrNameSourceIP     = "source_ip"
	attrNameAddress      = "address"
	attrNameLog          = "log"
	attrNameRedirectTo   = "redirect_to"
	attrNameRefuse       = "refuse"
	attrNameDrop         = "drop"
	attrNamePolicyAction = "policy_action"

	typeValueQuery    = "query"
	typeValueResponse = "response"

	ednsAttrsStart = 3
	ipReqAddrPos   = 1
)

func serializeOrPanic(a pdp.AttributeAssignment) string {
	v, err := a.GetValue()
	if err != nil {
		panic(err)
	}

	s, err := v.Serialize()
	if err != nil {
		panic(err)
	}

	return s
}

type custAttr byte

const (
	custAttrEdns = 1 << iota
	custAttrTransfer
	custAttrDnstap
	custAttrMetrics
)

func (a custAttr) isEdns() bool {
	return a&custAttrEdns != 0
}

func (a custAttr) isTransfer() bool {
	return a&custAttrTransfer != 0
}

func (a custAttr) isDnstap() bool {
	return a&custAttrDnstap != 0
}

func (a custAttr) isMetrics() bool {
	return a&custAttrMetrics != 0
}
