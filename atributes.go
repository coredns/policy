package policy

const (
	attrNameType       = "type"
	attrNameSourceIP   = "source_ip"
	attrNameDomainName = "domain_name"
	attrNameRedirectTo = "redirect_to"
	attrNameAddress    = "address"
	attrNameRefuse     = "refuse"
	attrNameLog        = "log"
	attrNameDrop       = "drop"

	typeValueQuery    = "query"
	typeValueResponse = "response"

	attrNameDNSQtype     = "dns_qtype"
	attrNamePolicyAction = "policy_action"
)

type confAttrType byte

const (
	confAttrEdns = 1 << iota
	confAttrTransfer
	confAttrDnstap
)

func (a confAttrType) isEnds() bool {
	return a&confAttrEdns != 0
}

func (a confAttrType) isTransfer() bool {
	return a&confAttrTransfer != 0
}

func (a confAttrType) isDnstap() bool {
	return a&confAttrDnstap != 0
}
