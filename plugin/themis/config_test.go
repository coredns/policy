package themis

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyfile"
)

func TestThemisConfigParse(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		err   error

		endpoints    []string
		options      []*attrSetting
		debugSuffix  *string
		streams      *int
		hotSpot      *bool
		custAttrs    map[string]custAttr
		debugID      *string
		passthrough  []string
		connTimeout  *time.Duration
		autoReqSize  *bool
		maxReqSize   *int
		autoResAttrs *bool
		maxResAttrs  *int
		cacheTTL     *time.Duration
		cacheLimit   *int
	}{
		{
			desc: "InvalidOption",
			input: `.:53 {
						themis NAME {
							error option
						}
					}`,
			err: errors.New("invalid themis plugin option"),
		},
		{
			desc: "NoEndpointArguemnts",
			input: `.:53 {
						themis NAME {
							endpoint
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "SingleEntryEndpoint",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
						}
					}`,
			endpoints: []string{"10.2.4.1:5555"},
		},
		{
			desc: "ThreeEntriesEndpoint",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555 10.2.4.2:5555
						}
					}`,
			endpoints: []string{"10.2.4.1:5555", "10.2.4.2:5555"},
		},
		{
			desc: "InvalidAttrType",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr uid request/uid no-type
						}
					}`,
			err: errors.New("invalid type"),
		},
		{
			desc: "AttrCorrect",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr uid request/uid
						}
					}`,
			options: []*attrSetting{
				{"uid", "request/uid", "string", false},
			},
			custAttrs: map[string]custAttr{
				"uid": custAttrEdns,
			},
		},
		{
			desc: "AttrCorrect2values",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr uid request/uid
							attr ip request/ip address 
						}
					}`,
			options: []*attrSetting{
				{"uid", "request/uid", "string", false},
				{"ip", "request/ip", "address", false},
			},
			custAttrs: map[string]custAttr{
				"uid": custAttrEdns,
				"ip":  custAttrEdns,
			},
		},
		{
			desc: "AttrWithNoLabel",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr my-name
						}
					}`,
			err: errors.New("Invalid attr directive"),
		},
		{
			desc: "NoDebugQuerySuffixArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							debug_query_suffix
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "DebugQuerySuffix",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							debug_query_suffix debug.local.
						}
					}`,
			debugSuffix: newStringPtr("debug.local."),
		},
		{
			desc: "PDPClientStreams",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams 10
						}
					}`,
			streams: newIntPtr(10),
		},
		{
			desc: "InvalidPDPClientStreams",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams Ten
						}
					}`,
			err: errors.New("Could not parse number of streams"),
		},
		{
			desc: "NoPDPClientStreamsArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NegativePDPClientStreams",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams -1
						}
					}`,
			err: errors.New("Expected at least one stream got -1"),
		},
		{
			desc: "PDPClientStreamsWithRoundRobin",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams 10 Round-Robin
						}
					}`,
			streams: newIntPtr(10),
			hotSpot: newBoolPtr(false),
		},
		{
			desc: "PDPClientStreamsWithHotSpot",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams 10 Hot-Spot
						}
					}`,
			streams: newIntPtr(10),
			hotSpot: newBoolPtr(true),
		},
		{
			desc: "InvalidPDPClientStreamsBalancer",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							streams 10 Unknown-Balancer
						}
					}`,
			err: errors.New("Expected round-robin or hot-spot balancing but got Unknown-Balancer"),
		},
		{
			desc: "TransferAttribute",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							transfer themis_id
						}
					}`,
			custAttrs: map[string]custAttr{
				"themis_id": custAttrTransfer,
			},
		},
		{
			desc: "ComplexAttributeConfig",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr uid request/uid
							attr id request/id
							transfer themis_id id
						}
					}`,
			options: []*attrSetting{
				{"uid", "request/uid", "string", false},
				{"id", "request/id", "string", false},
			},
			custAttrs: map[string]custAttr{
				"themis_id": custAttrTransfer,
				"id":        custAttrEdns | custAttrTransfer,
				"uid":       custAttrEdns,
			},
		},
		{
			desc: "NoTransferArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							transfer
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "DebugID",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							metrics
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "ComplexAttributeConfigWithMetrics",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							attr uid request/uid
							attr id request/id
							metrics uid query_id
						}
					}`,
			options: []*attrSetting{
				{"uid", "request/uid", "string", false},
				{"id", "request/id", "string", false},
			},
			custAttrs: map[string]custAttr{
				"id":       custAttrEdns,
				"uid":      custAttrEdns | custAttrMetrics,
				"query_id": custAttrMetrics,
			},
		},
		{
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							debug_id corednsinstance
						}
					}`,
			debugID: newStringPtr("corednsinstance"),
		},
		{
			desc: "NoDebugIDArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							debug_id
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NoConnectionTimeoutArguments",
			input: `.:53 {
						themis NAME {
							connection_timeout
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NoConnectionTimeout",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							connection_timeout no
						}
					}`,
			connTimeout: newDurationPtr(-1),
		},
		{
			desc: "ConnectionTimeout",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							connection_timeout 500ms
						}
					}`,
			connTimeout: newDurationPtr(500 * time.Millisecond),
		},
		{
			desc: "InvalidConnectionTimeout",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							connection_timeout invalid
						}
					}`,
			err: errors.New("Could not parse timeout: time: invalid duration invalid"),
		},
		{
			desc: "Log",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							log
						}
					}`,
		},
		{
			desc: "TrailingLogArgument",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							log stdout
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "MaxRequestSize",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size 128
						}
					}`,
			autoReqSize: newBoolPtr(false),
			maxReqSize:  newIntPtr(128),
		},
		{
			desc: "MaxRequestSize",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size auto
						}
					}`,
			autoReqSize: newBoolPtr(true),
			maxReqSize:  newIntPtr(-1),
		},
		{
			desc: "MaxRequestSize",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size auto 128
						}
					}`,
			autoReqSize: newBoolPtr(true),
			maxReqSize:  newIntPtr(128),
		},
		{
			desc: "NoMaxRequestSizeArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidMaxRequestSize",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size test
						}
					}`,
			err: errors.New("Could not parse PDP request size limit"),
		},
		{
			desc: "OverflowMaxRequestSize",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_request_size 2147483648
						}
					}`,
			err: errors.New("Size limit 2147483648 (> 2147483647) for PDP request is too high"),
		},
		{
			desc: "MaxResponseAttributes",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_response_attributes 128
						}
					}`,
			autoResAttrs: newBoolPtr(false),
			maxResAttrs:  newIntPtr(128),
		},
		{
			desc: "MaxResponseAttributes",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_response_attributes auto
						}
					}`,
			autoResAttrs: newBoolPtr(true),
			maxResAttrs:  newIntPtr(64),
		},
		{
			desc: "NoMaxResponseAttributesArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_response_attributes
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidMaxResponseAttributes",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_response_attributes invalid
						}
					}`,
			err: errors.New("Could not parse PDP response attributes limit"),
		},
		{
			desc: "OverflowMaxResponseAttributes",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							max_response_attributes 2147483648
						}
					}`,
			err: errors.New("Attributes limit 2147483648 (> 2147483647) for PDP response is too high"),
		},
		{
			desc: "NoDecisionCache",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
						}
					}`,
			cacheTTL: newDurationPtr(0),
		},
		{
			desc: "DecisionCache",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache
						}
					}`,
			cacheTTL:   newDurationPtr(10 * time.Minute),
			cacheLimit: newIntPtr(0),
		},
		{
			desc: "DecisionCacheWithTTL",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache 15s
						}
					}`,
			cacheTTL:   newDurationPtr(15 * time.Second),
			cacheLimit: newIntPtr(0),
		},
		{
			desc: "DecisionCacheWithTTLAndLimit",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache 15s 128
						}
					}`,
			cacheTTL:   newDurationPtr(15 * time.Second),
			cacheLimit: newIntPtr(128),
		},
		{
			desc: "TooManyCacheArguments",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache too many of them
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidCacheTTL",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache invalid
						}
					}`,
			err: errors.New("Could not parse decision cache TTL"),
		},
		{
			desc: "WrongCacheTTL",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache -15s
						}
					}`,
			err: errors.New("Can't set decision cache TTL to"),
		},
		{
			desc: "InvalidCacheLimit",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache 15s invalid
						}
					}`,
			err: errors.New("Could not parse decision cache limit"),
		},
		{
			desc: "OverflowCacheLimit",
			input: `.:53 {
						themis NAME {
							endpoint 10.2.4.1:5555
							cache 15s 2147483648
						}
					}`,
			err: errors.New("Cache limit 2147483648 (> 2147483647) is too high"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			blocs, err := caddyfile.Parse("test-file", strings.NewReader(test.input), []string{"themis"})
			d := caddyfile.NewDispenserTokens("themis", blocs[0].Tokens["themis"])
			c := &caddy.Controller{Dispenser: d}
			mw, err := themisParse(c)
			if err != nil {
				if test.err != nil {
					if !strings.Contains(err.Error(), test.err.Error()) {
						t.Errorf("Expected error '%v' but got '%v'\n", test.err, err)
					}
				} else {
					t.Errorf("Expected no error but got '%v'\n", err)
				}
			} else {
				// we consider only one Engine declared
				for _, mwe := range mw.engines {
					if test.err != nil {
						t.Errorf("Expected error '%v' but got 'nil'\n", test.err)
					} else {
						if test.endpoints != nil {
							if len(test.endpoints) != len(mwe.conf.endpoints) {
								t.Errorf("Expected endpoints %v but got %v\n", test.endpoints, mwe.conf.endpoints)
							} else {
								for i := 0; i < len(test.endpoints); i++ {
									if test.endpoints[i] != mwe.conf.endpoints[i] {
										t.Errorf("Expected endpoint '%s' but got '%s'\n",
											test.endpoints[i], mwe.conf.endpoints[i])
									}
								}
							}
						}

						if test.options != nil {
							if len(test.options) != len(mwe.conf.options) {
								t.Errorf("Expected %d Attr options  but got %d",
									len(test.options), len(mwe.conf.options))
							} else {
								for k, testAttr := range test.options {
									mwOpt := mwe.conf.options[k]
									if testAttr.name != mwOpt.name ||
										testAttr.label != mwOpt.label ||
										testAttr.attrType != mwOpt.attrType ||
										testAttr.metrics != mwOpt.metrics {
										t.Errorf("Expected Attr option:\n\t\"%#v\""+
											"\nfor but got:\n\t\"%#v\"",
											*testAttr, *mwOpt)
									}
								}

							}
						}

						if test.debugSuffix != nil && *test.debugSuffix != mwe.conf.debugSuffix {
							t.Errorf("Expected debug suffix %q but got %q", *test.debugSuffix, mwe.conf.debugSuffix)
						}

						if test.streams != nil && *test.streams != mwe.conf.streams {
							t.Errorf("Expected %d streams but got %d", *test.streams, mwe.conf.streams)
						}

						if test.hotSpot != nil && *test.hotSpot != mwe.conf.hotSpot {
							t.Errorf("Expected hotSpot=%v but got %v", *test.hotSpot, mwe.conf.hotSpot)
						}

						if test.custAttrs != nil {
							for k, et := range test.custAttrs {
								at, ok := mwe.conf.custAttrs[k]
								if !ok {
									t.Errorf("Missing conf attribute %q", k)
								} else if et != at {
									t.Errorf("Unexpected type of conf attribute %q; expected=%d, actual=%d", k, et, at)
								}
							}

							for k, at := range mwe.conf.custAttrs {
								if _, ok := test.custAttrs[k]; !ok {
									t.Errorf("Unexpected conf attribute %q=%d", k, at)
								}
							}
						}

						if test.debugID != nil && *test.debugID != mwe.conf.debugID {
							t.Errorf("Expected debug id %q but got %q", *test.debugID, mwe.conf.debugID)
						}

						if test.connTimeout != nil && *test.connTimeout != mwe.conf.connTimeout {
							t.Errorf("Expected connection timeout %s but got %s", *test.connTimeout, mwe.conf.connTimeout)
						}

						if test.autoReqSize != nil && *test.autoReqSize != mwe.conf.autoReqSize {
							t.Errorf("Expected automatic request size %v but got %v",
								*test.autoReqSize, mwe.conf.autoReqSize)
						}

						if test.maxReqSize != nil && *test.maxReqSize != mwe.conf.maxReqSize {
							t.Errorf("Expected request size limit %d but got %d", *test.maxReqSize, mwe.conf.maxReqSize)
						}

						if test.autoResAttrs != nil && *test.autoResAttrs != mwe.conf.autoResAttrs {
							t.Errorf("Expected automatic response attributes %v but got %v",
								*test.autoResAttrs, mwe.conf.autoResAttrs)
						}

						if test.maxResAttrs != nil && *test.maxResAttrs != mwe.conf.maxResAttrs {
							t.Errorf("Expected response attributes limit %d but got %d",
								*test.maxResAttrs, mwe.conf.maxResAttrs)
						}

						if test.cacheTTL != nil && *test.cacheTTL != mwe.conf.cacheTTL {
							t.Errorf("Expected cache TTL %s but got %s", *test.cacheTTL, mwe.conf.cacheTTL)
						}

						if test.cacheLimit != nil && *test.cacheLimit != mwe.conf.cacheLimit {
							t.Errorf("Expected cache limit %d but got %d", *test.cacheLimit, mwe.conf.cacheLimit)
						}
					}
				}
			}
		})
	}
}

func newStringPtr(s string) *string {
	return &s
}

func newIntPtr(n int) *int {
	return &n
}

func newBoolPtr(b bool) *bool {
	return &b
}

func newDurationPtr(d time.Duration) *time.Duration {
	return &d
}
