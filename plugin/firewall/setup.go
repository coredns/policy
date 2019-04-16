package firewall

import (
	"fmt"
	"strings"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/firewall/policy"
	"github.com/coredns/coredns/plugin/firewall/rule"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("firewall", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	fw, err := parse(c)

	if err != nil {
		return plugin.Error("firewall", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		fw.next = next
		return fw
	})

	c.OnStartup(func() error {
		// after all plugin are setup, ensure to have all rules created by enrolling the engines pointed
		// by pending rules
		err := fw.enrollEngines(c)
		if err != nil {
			return err
		}
		for _, loc := range []*rule.List{fw.query, fw.reply} {
			// now that all engines are known, ensure to build the rules for each element of the lists
			err = loc.BuildRules(fw.engines)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return nil
}

func parse(c *caddy.Controller) (*firewall, error) {
	p, err := New()
	if err != nil {
		return nil, fmt.Errorf("cannot create the firewall plugin structure, error : %e", err)
	}

	for c.Next() {
		opts := c.RemainingArgs()
		if len(opts) != 1 {
			return nil, c.Errf("one and only one paramater is expected after firewall : the location of the rulelist. It should be either query or reply")
		}
		location := opts[0]
		var rl *rule.List
		switch location {
		case "query":
			rl = p.query
		case "response":
			rl = p.reply
		default:
			return nil, c.Errf("invalid location of rule list: %s . It should be either query or response", location)
		}
		// check if location already used or not
		for c.NextBlock() {
			r, err := p.parseOptionOrRule(c)
			if err != nil {
				return nil, err
			}
			err = rl.Add(r)
			if err != nil {
				return nil, c.Errf("cannot add a rule to the %s list : %s", location, err)
			}
		}
	}
	return p, nil
}

func (p *firewall) parseOptionOrRule(c *caddy.Controller) (*rule.Element, error) {
	// by default, at least one engine is available : the ExpressionEngine
	e := policy.NewExprEngine()
	switch c.Val() {
	case policy.NameTypes[policy.TypeRefuse]:
		fallthrough
	case policy.NameTypes[policy.TypeAllow]:
		fallthrough
	case policy.NameTypes[policy.TypeBlock]:
		fallthrough
	case policy.NameTypes[policy.TypeDrop]:
		// these 4 direct policy action denotes the actions for the default Engine: ExpressionEngine
		action := c.Val()
		name := ExpressionEngineName
		args := c.RemainingArgs()
		if len(args) < 1 {
			return nil, fmt.Errorf("not enough arguments to build a policy rule, expect allow/refuse/block/drop query/reply <expression>, got %s %s", c.Val(), strings.Join(args, " "))
		}
		params := append([]string{action}, args...)
		r, err := e.BuildRule(params)
		if err != nil {
			return nil, err
		}
		return &rule.Element{"", name, params, r}, nil

	default:
		// we can only suppose it is an engine type(plugin name), name and args
		plugin := c.Val()
		args := c.RemainingArgs()
		if len(args) < 1 {
			return nil, fmt.Errorf("not enough arguments to build a policy rule, expect %s <name-engine>", c.Val())
		}
		name := args[0]
		params := args[1:]
		// as the Engine are not yet knowm, just create a ruleElement with the parameters.The Element will be created later
		return &rule.Element{plugin, name, params, nil}, nil

	}
}

func (p *firewall) enrollEngines(c *caddy.Controller) error {

	var eng = p.query.Engines()
	for n, p := range p.reply.Engines() {
		eng[n] = p
	}
	// remove Expression engine that is built-in
	delete(eng, ExpressionEngineName)

	// retrieve all needed Engines.
	// These are plugins that implements the 'Engineer' interface and are named in one of the firewall lists
	plugins := dnsserver.GetConfig(c).Handlers()
	for _, pl := range plugins {
		if e, ok := pl.(policy.Engineer); ok {
			for n, pn := range eng {
				if pn == pl.Name() {
					re := e.Engine(n)
					if re == nil {
						return c.Errf("missing policy engine for plugin %s and name %s", p.Name(), n)
					}
					p.engines[n] = re
					delete(eng, n)
				}
			}
		}
	}

	// raise error for engines referenced in the list but not present in the Corefile as plugin
	for n := range eng {
		return c.Errf("the policy engine %s is missing", n)
	}
	return nil
}
