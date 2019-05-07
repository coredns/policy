package themisplugin

import (
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin(ThemisPluginName, caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	t, err := themisParse(c)

	if err != nil {
		return plugin.Error("themis", err)
	}

	for _, e := range t.engines {
		c.OnStartup(func() error {
			e.trace = dnsserver.GetConfig(c).Handler("trace")
			err := e.connect()
			if err != nil {
				return plugin.Error("themis", err)
			}

			return e.SetupMetrics(c)
		})

		c.OnShutdown(func() error {
			e.closeConn()
			return nil
		})
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		t.next = next
		return t
	})

	return nil
}
