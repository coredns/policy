package errors

import (
	"fmt"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("errors", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	handler, err := errorsParse(c)
	if err != nil {
		return plugin.Error("errors", err)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		handler.Next = next
		return handler
	})

	return nil
}

func errorsParse(c *caddy.Controller) (errorHandler, error) {
	handler := errorHandler{}

	i := 0
	for c.Next() {
		if i > 0 {
			return handler, plugin.ErrOnce
		}
		i++

		args := c.RemainingArgs()
		switch len(args) {
		case 0:
		case 1:
			if args[0] != "stdout" {
				return handler, fmt.Errorf("invalid log file: %s", args[0])
			}
		default:
			return handler, c.ArgErr()
		}
	}
	return handler, nil
}
