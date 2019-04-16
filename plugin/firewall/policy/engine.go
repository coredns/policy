package policy

import (
	"context"

	"github.com/coredns/coredns/request"
)

const (
	// TypeNone policy action is UNKNOWN: no decision, apply the next rule.
	TypeNone = iota
	// TypeRefuse policy action is REFUSE: do not resolve a query and return code REFUSED
	TypeRefuse
	// TypeAllow policy action is ALLOW: continue to resolve query
	TypeAllow
	// TypeBlock policy action is BLOCK: do not resolve a query and return code NXDOMAIN
	TypeBlock
	// TypeDrop policy action is DROP: do not resolve a query and simulate a lost query
	TypeDrop

	// TypeCount total number of actions allowed
	TypeCount
)

// NameTypes keep a mapping of the byte constant to the corresponding name
var NameTypes = map[int]string{
	TypeNone:   "none",
	TypeAllow:  "allow",
	TypeRefuse: "refuse",
	TypeBlock:  "block",
	TypeDrop:   "drop",
}

// Rule defines a policy for continuing DNS query processing.
type Rule interface {
	// Evaluate the rule and return one of the TypeXXX defined above
	//   - TypeNone should be returned if the Rule is not able to decide any action for this query
	//   - otherwise return one of TypeAllow/TypeRefuse/TypeDrop/TypeBlock
	Evaluate(data interface{}) (int, error)
}

// Engine for Firewall plugin
type Engine interface {
	// BuildRules - create a Rule based on args or throw an error, This Rule will be evaluated during processing of DNS Queries
	BuildRule(args []string) (Rule, error) // create a rule based on parameters

	//BuildQueryData generate the data needed to evaluate - for one query - ALL the rules of this Engine
	BuildQueryData(ctx context.Context, state request.Request) (interface{}, error)

	//BuildReplyData generate the data needed to evaluate - for one response - ALL the rules of this Engine
	BuildReplyData(ctx context.Context, state request.Request, queryData interface{}) (interface{}, error)
}

// Engineer allow registration of Policy Engines. One plugin can declare several Engines.
type Engineer interface {
	Engine(name string) Engine
}
