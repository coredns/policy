package themisplugin

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/mholt/caddy"
)

var errInvalidOption = errors.New("invalid themis plugin option")

var allowedAttrTypes = map[string]string{
	"string":  "string",
	"domain":  "domain",
	"address": "address"}

type config struct {
	policyFile   string
	contentFiles []string
	endpoints    []string
	options      []*attrSetting
	custAttrs    map[string]custAttr
	debugID      string
	debugSuffix  string
	streams      int
	hotSpot      bool
	connTimeout  time.Duration
	autoReqSize  bool
	maxReqSize   int
	autoResAttrs bool
	maxResAttrs  int
	log          bool
	cacheTTL     time.Duration
	cacheLimit   int
}

func themisParse(c *caddy.Controller) (*ThemisPlugin, error) {
	tp := newThemisPlugin()
	for c.Next() {
		args := c.RemainingArgs()
		if len(args) != 1 {
			return nil, c.Errf("themis plugin should have the format : themis <EngineName> {...} ")
		}
		name := args[0]
		if _, ok := tp.engines[name]; ok {
			return nil, c.Errf("themis plugin with engine name %s is already declared", name)
		}
		p := newThemisEngine()
		for c.NextBlock() {
			if err := p.conf.parseOption(c); err != nil {
				return nil, err
			}
		}
		tp.engines[name] = p
	}
	return tp, nil
}

func (conf *config) parseOption(c *caddy.Controller) error {
	switch c.Val() {
	case "pdp":
		return conf.parsePDP(c)

	case "endpoint":
		return conf.parseEndpoint(c)

	case "attr":
		return conf.parseAttr(c)

	case "debug_query_suffix":
		return conf.parseDebugQuerySuffix(c)

	case "streams":
		return conf.parseStreams(c)

	case "transfer":
		return conf.parseAttributes(c, custAttrTransfer)

	case "metrics":
		return conf.parseAttributes(c, custAttrMetrics)

	case "debug_id":
		return conf.parseDebugID(c)

	case "connection_timeout":
		return conf.parseConnectionTimeout(c)

	case "log":
		return conf.parseLog(c)

	case "max_request_size":
		return conf.parseMaxRequestSize(c)

	case "max_response_attributes":
		return conf.parseMaxResponseAttributes(c)

	case "cache":
		return conf.parseCache(c)
	}

	return errInvalidOption
}

// Usage: pdp policy.[yaml|json] content1 content2...
func (conf *config) parsePDP(c *caddy.Controller) error {
	args := c.RemainingArgs()
	argsLen := len(args)
	if argsLen < 1 {
		return c.ArgErr()
	}

	conf.policyFile = args[0]
	if argsLen > 1 {
		conf.contentFiles = args[1:]
	}
	return nil
}

func (conf *config) parseEndpoint(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) <= 0 {
		return c.ArgErr()
	}

	conf.endpoints = args
	return nil
}

func (conf *config) parseAttr(c *caddy.Controller) error {
	args := c.RemainingArgs()
	// Usage: edns0 <code> <name> [ <dataType> ] [ <size> <start> <end> ].
	// Valid dataTypes are hex (default), bytes, ip.
	// Valid destTypes depend on PDP (default string).
	argsLen := len(args)
	if argsLen != 2 && argsLen != 3 {
		return c.Errf("Invalid attr directive. Expected 2 or 3 arguments but got %d", argsLen)
	}

	name := args[0]
	label := args[1]
	dataType := "string"
	if argsLen > 2 {
		dataType = args[2]
	}

	if _, ok := allowedAttrTypes[strings.ToLower(dataType)]; !ok {
		tp := make([]string, 0)
		for k := range allowedAttrTypes {
			tp = append(tp, k)
		}
		return c.Errf("invalid type %s for an attribute - allowed types are : %s", dataType, strings.Join(tp, ","))
	}
	conf.options = append(conf.options, &attrSetting{name, label, dataType, false})
	conf.custAttrs[name] = conf.custAttrs[name] | custAttrEdns
	return nil
}

func (conf *config) parseAttributes(c *caddy.Controller, a custAttr) error {
	args := c.RemainingArgs()
	if len(args) <= 0 {
		return c.ArgErr()
	}

	for _, item := range args {
		conf.custAttrs[item] = conf.custAttrs[item] | a
	}

	return nil
}

func (conf *config) parseStreams(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) < 1 || len(args) > 2 {
		return c.ArgErr()
	}

	streams, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		return c.Errf("Could not parse number of streams: %s", err)
	}
	if streams < 1 {
		return c.Errf("Expected at least one stream got %d", streams)
	}

	conf.streams = int(streams)

	if len(args) > 1 {
		switch strings.ToLower(args[1]) {
		default:
			return c.Errf("Expected round-robin or hot-spot balancing but got %s", args[1])

		case "round-robin":
			conf.hotSpot = false

		case "hot-spot":
			conf.hotSpot = true
		}
	} else {
		conf.hotSpot = false
	}

	return nil
}

func (conf *config) parseDebugQuerySuffix(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return c.ArgErr()
	}

	conf.debugSuffix = args[0]
	return nil
}

func (conf *config) parseDebugID(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return c.ArgErr()
	}

	conf.debugID = args[0]
	return nil
}

func (conf *config) parseConnectionTimeout(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return c.ArgErr()
	}

	if strings.ToLower(args[0]) == "no" {
		conf.connTimeout = -1
	} else {
		timeout, err := time.ParseDuration(args[0])
		if err != nil {
			return c.Errf("Could not parse timeout: %s", err)
		}

		conf.connTimeout = timeout
	}

	return nil
}

func (conf *config) parseLog(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) != 0 {
		return c.ArgErr()
	}

	conf.log = true
	return nil
}

func (conf *config) parseMaxRequestSize(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) < 1 || len(args) > 2 {
		return c.ArgErr()
	}

	s := ""
	if strings.ToLower(args[0]) == "auto" {
		conf.autoReqSize = true
		if len(args) > 1 {
			s = args[1]
		}
	} else {
		s = args[0]
	}

	if len(s) > 0 {
		size, err := strconv.ParseUint(s, 10, 0)
		if err != nil {
			return c.Errf("Could not parse PDP request size limit: %s", err)
		}

		if size > math.MaxInt32 {
			return c.Errf("Size limit %d (> %d) for PDP request is too high", size, math.MaxInt32)
		}

		conf.maxReqSize = int(size)
	}

	return nil
}

func (conf *config) parseMaxResponseAttributes(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) != 1 {
		return c.ArgErr()
	}

	if strings.ToLower(args[0]) == "auto" {
		conf.autoResAttrs = true
		return nil
	}

	n, err := strconv.ParseUint(args[0], 10, 0)
	if err != nil {
		return c.Errf("Could not parse PDP response attributes limit: %s", err)
	}

	if n > math.MaxInt32 {
		return c.Errf("Attributes limit %d (> %d) for PDP response is too high", n, math.MaxInt32)
	}

	conf.maxResAttrs = int(n)
	return nil
}

func (conf *config) parseCache(c *caddy.Controller) error {
	args := c.RemainingArgs()
	if len(args) > 2 {
		return c.ArgErr()
	}

	if len(args) > 0 {
		ttl, err := time.ParseDuration(args[0])
		if err != nil {
			return c.Errf("Could not parse decision cache TTL: %s", err)
		}

		if ttl <= 0 {
			return c.Errf("Can't set decision cache TTL to %s", ttl)
		}

		conf.cacheTTL = ttl
	} else {
		conf.cacheTTL = 10 * time.Minute
	}

	if len(args) > 1 {
		n, err := strconv.ParseUint(args[1], 10, 0)
		if err != nil {
			return c.Errf("Could not parse decision cache limit: %s", err)
		}

		if n > math.MaxInt32 {
			return c.Errf("Cache limit %d (> %d) is too high", n, math.MaxInt32)
		}

		conf.cacheLimit = int(n)
	}

	return nil
}
