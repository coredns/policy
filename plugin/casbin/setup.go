package casbin

import (
	casbin2 "github.com/casbin/casbin/v2"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/policy/plugin/pkg/rqdata"
)

func init() {
	caddy.RegisterPlugin("casbin", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	e, err := parse(c)
	if err != nil {
		return err
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		e.next = next
		return e
	})
	return nil
}

func parse(c *caddy.Controller) (*casbin, error) {
	e := newCasbin()
	mapping := rqdata.NewMapping("")
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, c.ArgErr()
		}
		name := args[0]
		eng := newEngine(mapping)
		for c.NextBlock() {
			switch c.Val() {
			case "model":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				eng.modelPath = args[0]
			case "policy":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				eng.policyPath = args[0]
			}
		}
		if eng.policyPath == "" || eng.modelPath == "" {
			return nil, c.Err("policy path and model path are required")
		}
		if e, err := casbin2.NewEnforcer(eng.modelPath, eng.policyPath); err != nil {
			return nil, err
		} else {
			eng.enforcer = e
		}
		eng.getFields()
		if err := eng.getActionIndex(); err != nil {
			return nil, err
		}
		e.engines[name] = eng
	}
	return e, nil
}
