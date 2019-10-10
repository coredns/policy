package opa

import (
	"crypto/tls"
	"net/http"

	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	pkgtls "github.com/coredns/coredns/plugin/pkg/tls"
	"github.com/coredns/policy/plugin/pkg/rqdata"
)

func init() {
	caddy.RegisterPlugin("opa", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	o, err := parse(c)
	if err != nil {
		return err
	}
	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		o.next = next
		return o
	})
	return nil
}

func parse(c *caddy.Controller) (*opa, error) {
	o := newOpa()
	mapping := rqdata.NewMapping("")
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, c.ArgErr()
		}
		name := args[0]
		eng := newEngine(mapping)
		var tlsConfig *tls.Config
		for c.NextBlock() {
			switch c.Val() {
			case "endpoint":
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				eng.endpoint = args[0]
			case "fields":
				args := c.RemainingArgs()
				if len(args) == 0 {
					return nil, c.ArgErr()
				}
				// these fields cannot be validated, because metadata fields are not known at setup time
				eng.fields = args
			case "tls": // cert key cacertfile
				args := c.RemainingArgs()
				if len(args) == 3 {
					var err error
					tlsConfig, err = pkgtls.NewTLSConfigFromArgs(args...)
					if err != nil {
						return nil, err
					}
					tlsConfig.BuildNameToCertificate()
					continue
				}
				return nil, c.ArgErr()
			}
		}
		if eng.endpoint == "" {
			return nil, c.Err("endpoint required")
		}
		eng.client = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
		o.engines[name] = eng
	}
	return o, nil
}
