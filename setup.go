package policy

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/dnstap"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("policy", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	policyPlugin, err := policyParse(c)

	if err != nil {
		return plugin.Error("policy", err)
	}

	c.OnStartup(func() error {
		if taph := dnsserver.GetConfig(c).Handler("dnstap"); taph != nil {
			if tapPlugin, ok := taph.(dnstap.Dnstap); ok && tapPlugin.IO != nil {
				policyPlugin.tapIO = newPolicyDnstapSender(tapPlugin.IO)
			}
		}

		policyPlugin.trace = dnsserver.GetConfig(c).Handler("trace")
		err := policyPlugin.connect()
		if err != nil {
			return plugin.Error("policy", err)
		}

		return policyPlugin.SetupMetrics(c)
	})

	c.OnShutdown(func() error {
		policyPlugin.closeConn()
		return nil
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		policyPlugin.next = next
		return policyPlugin
	})

	return nil
}
