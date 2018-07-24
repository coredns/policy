package policy

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mholt/caddy"
)

func TestPolicyConfigParse(t *testing.T) {
	tests := []struct {
		desc  string
		input string
		err   error

		endpoints    []string
		options      map[uint16][]*edns0Opt
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
			desc: "MissingPolicySection",
			input: `.:53 {
						log stdout
					}`,
			err: errors.New("Policy setup called without keyword 'policy' in Corefile"),
		},
		{
			desc: "InvalidOption",
			input: `.:53 {
						policy {
							error option
						}
					}`,
			err: errors.New("invalid policy plugin option"),
		},
		{
			desc: "NoEndpointArguemnts",
			input: `.:53 {
						policy {
							endpoint
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "SingleEntryEndpoint",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
						}
					}`,
			endpoints: []string{"10.2.4.1:5555"},
		},
		{
			desc: "ThreeEntriesEndpoint",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555 10.2.4.2:5555
						}
					}`,
			endpoints: []string{"10.2.4.1:5555", "10.2.4.2:5555"},
		},
		{
			desc: "InvalidEDNS0Size",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex wrong_size 0 32
						}
					}`,
			err: errors.New("Could not parse EDNS0 data size"),
		},
		{
			desc: "EDNS0Hex",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 32
						}
					}`,
			options: map[uint16][]*edns0Opt{
				0xfff0: {
					&edns0Opt{
						name:     "uid",
						dataType: typeEDNS0Hex,
						size:     32,
						start:    0,
						end:      32},
				},
			},
			custAttrs: map[string]custAttr{
				"uid": custAttrEdns,
			},
		},
		{
			desc: "InvalidEDNS0Code",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 wrong_hex uid hex
						}
					}`,
			err: errors.New("Could not parse EDNS0 code"),
		},
		{
			desc: "InvalidEDNS0StartIndex",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 wrong_offset 32
						}
					}`,
			err: errors.New("Could not parse EDNS0 start index"),
		},
		{
			desc: "InvalidEDNS0EndIndex",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 wrong_size
						}
					}`,
			err: errors.New("Could not parse EDNS0 end index"),
		},
		{
			desc: "EDNS0Hex2",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 16
							edns0 0xfff0 id hex 32 16 32
						}
					}`,
			options: map[uint16][]*edns0Opt{
				0xfff0: {
					&edns0Opt{
						name:     "uid",
						dataType: typeEDNS0Hex,
						size:     32,
						start:    0,
						end:      16},
					&edns0Opt{
						name:     "id",
						dataType: typeEDNS0Hex,
						size:     32,
						start:    16,
						end:      32},
				},
			},
			custAttrs: map[string]custAttr{
				"uid": custAttrEdns,
				"id":  custAttrEdns,
			},
		},
		{
			desc: "InvalidEDNS0StartEndPair",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 16 15
						}
					}`,
			err: errors.New("End index should be > start index"),
		},
		{
			desc: "InvalidEDNS0SizeEndPair",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 33
						}
					}`,
			err: errors.New("End index should be <= size"),
		},
		{
			desc: "NotEnoughEDNS0Arguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff1
						}
					}`,
			err: errors.New("Invalid edns0 directive"),
		},
		{
			desc: "InvalidEDNS0Type",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff1 guid bin
						}
					}`,
			err: errors.New("Could not add EDNS0"),
		},
		{
			desc: "NoDebugQuerySuffixArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							debug_query_suffix
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "DebugQuerySuffix",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							debug_query_suffix debug.local.
						}
					}`,
			debugSuffix: newStringPtr("debug.local."),
		},
		{
			desc: "PDPClientStreams",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							streams 10
						}
					}`,
			streams: newIntPtr(10),
		},
		{
			desc: "InvalidPDPClientStreams",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							streams Ten
						}
					}`,
			err: errors.New("Could not parse number of streams"),
		},
		{
			desc: "NoPDPClientStreamsArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							streams
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NegativePDPClientStreams",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							streams -1
						}
					}`,
			err: errors.New("Expected at least one stream got -1"),
		},
		{
			desc: "PDPClientStreamsWithRoundRobin",
			input: `.:53 {
						policy {
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
						policy {
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
						policy {
							endpoint 10.2.4.1:5555
							streams 10 Unknown-Balancer
						}
					}`,
			err: errors.New("Expected round-robin or hot-spot balancing but got Unknown-Balancer"),
		},
		{
			desc: "TransferAttribute",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							transfer policy_id
						}
					}`,
			custAttrs: map[string]custAttr{
				"policy_id": custAttrTransfer,
			},
		},
		{
			desc: "ComplexAttributeConfig",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 16
							edns0 0xfff1 id
							transfer policy_id id
							dnstap policy_id query_id
						}
					}`,
			options: map[uint16][]*edns0Opt{
				0xfff0: {
					&edns0Opt{
						name:     "uid",
						dataType: typeEDNS0Hex,
						size:     32,
						start:    0,
						end:      16},
				},
				0xfff1: {
					&edns0Opt{
						name:     "id",
						dataType: typeEDNS0Hex,
						size:     0,
						start:    0,
						end:      0},
				},
			},
			custAttrs: map[string]custAttr{
				"policy_id": custAttrTransfer | custAttrDnstap,
				"id":        custAttrEdns | custAttrTransfer,
				"uid":       custAttrEdns,
				"query_id":  custAttrDnstap,
			},
		},
		{
			desc: "NoDNStapArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							dnstap
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NoTransferArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							transfer
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "DebugID",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							metrics
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							edns0 0xfff0 uid hex 32 0 16
							edns0 0xfff1 id
							metrics uid query_id
						}
					}`,
			options: map[uint16][]*edns0Opt{
				0xfff0: {
					&edns0Opt{
						name:     "uid",
						dataType: typeEDNS0Hex,
						size:     32,
						start:    0,
						end:      16,
						metrics:  true,
					},
				},
				0xfff1: {
					&edns0Opt{
						name:     "id",
						dataType: typeEDNS0Hex,
						size:     0,
						start:    0,
						end:      0},
				},
			},
			custAttrs: map[string]custAttr{
				"id":       custAttrEdns,
				"uid":      custAttrEdns | custAttrMetrics,
				"query_id": custAttrMetrics,
			},
		},
		{
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							debug_id corednsinstance
						}
					}`,
			debugID: newStringPtr("corednsinstance"),
		},
		{
			desc: "NoDebugIDArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							debug_id
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "Passthrough",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							passthrough google.com. facebook.org.
						}
					}`,
			passthrough: []string{
				"google.com.",
				"facebook.org.",
			},
		},
		{
			desc: "NoPassthroughArguments",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							passthrough
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NoConnectionTimeoutArguments",
			input: `.:53 {
						policy {
							connection_timeout
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "NoConnectionTimeout",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							connection_timeout no
						}
					}`,
			connTimeout: newDurationPtr(-1),
		},
		{
			desc: "ConnectionTimeout",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							connection_timeout 500ms
						}
					}`,
			connTimeout: newDurationPtr(500 * time.Millisecond),
		},
		{
			desc: "InvalidConnectionTimeout",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							connection_timeout invalid
						}
					}`,
			err: errors.New("Could not parse timeout: time: invalid duration invalid"),
		},
		{
			desc: "Log",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							log
						}
					}`,
		},
		{
			desc: "TrailingLogArgument",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							log stdout
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "MaxRequestSize",
			input: `.:53 {
						policy {
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
						policy {
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
						policy {
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
						policy {
							endpoint 10.2.4.1:5555
							max_request_size
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidMaxRequestSize",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							max_request_size test
						}
					}`,
			err: errors.New("Could not parse PDP request size limit"),
		},
		{
			desc: "OverflowMaxRequestSize",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							max_request_size 2147483648
						}
					}`,
			err: errors.New("Size limit 2147483648 (> 2147483647) for PDP request is too high"),
		},
		{
			desc: "MaxResponseAttributes",
			input: `.:53 {
						policy {
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
						policy {
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
						policy {
							endpoint 10.2.4.1:5555
							max_response_attributes
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidMaxResponseAttributes",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							max_response_attributes invalid
						}
					}`,
			err: errors.New("Could not parse PDP response attributes limit"),
		},
		{
			desc: "OverflowMaxResponseAttributes",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							max_response_attributes 2147483648
						}
					}`,
			err: errors.New("Attributes limit 2147483648 (> 2147483647) for PDP response is too high"),
		},
		{
			desc: "NoDecisionCache",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
						}
					}`,
			cacheTTL: newDurationPtr(0),
		},
		{
			desc: "DecisionCache",
			input: `.:53 {
						policy {
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
						policy {
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
						policy {
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
						policy {
							endpoint 10.2.4.1:5555
							cache too many of them
						}
					}`,
			err: errors.New("Wrong argument count or unexpected line ending"),
		},
		{
			desc: "InvalidCacheTTL",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							cache invalid
						}
					}`,
			err: errors.New("Could not parse decision cache TTL"),
		},
		{
			desc: "WrongCacheTTL",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							cache -15s
						}
					}`,
			err: errors.New("Can't set decision cache TTL to"),
		},
		{
			desc: "InvalidCacheLimit",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							cache 15s invalid
						}
					}`,
			err: errors.New("Could not parse decision cache limit"),
		},
		{
			desc: "OverflowCacheLimit",
			input: `.:53 {
						policy {
							endpoint 10.2.4.1:5555
							cache 15s 2147483648
						}
					}`,
			err: errors.New("Cache limit 2147483648 (> 2147483647) is too high"),
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			c := caddy.NewTestController("dns", test.input)
			mw, err := policyParse(c)
			if err != nil {
				if test.err != nil {
					if !strings.Contains(err.Error(), test.err.Error()) {
						t.Errorf("Expected error '%v' but got '%v'\n", test.err, err)
					}
				} else {
					t.Errorf("Expected no error but got '%v'\n", err)
				}
			} else {
				if test.err != nil {
					t.Errorf("Expected error '%v' but got 'nil'\n", test.err)
				} else {
					if test.endpoints != nil {
						if len(test.endpoints) != len(mw.conf.endpoints) {
							t.Errorf("Expected endpoints %v but got %v\n", test.endpoints, mw.conf.endpoints)
						} else {
							for i := 0; i < len(test.endpoints); i++ {
								if test.endpoints[i] != mw.conf.endpoints[i] {
									t.Errorf("Expected endpoint '%s' but got '%s'\n",
										test.endpoints[i], mw.conf.endpoints[i])
								}
							}
						}
					}

					if test.options != nil {
						for k, testOpts := range test.options {
							if mwOpts, ok := mw.conf.options[k]; ok {
								if len(testOpts) != len(mwOpts) {
									t.Errorf("Expected %d EDNS0 options for 0x%04x but got %d",
										len(testOpts), k, len(mwOpts))
								} else {
									for i, testOpt := range testOpts {
										mwOpt := mwOpts[i]
										if testOpt.name != mwOpt.name ||
											testOpt.dataType != mwOpt.dataType ||
											testOpt.size != mwOpt.size ||
											testOpt.start != mwOpt.start ||
											testOpt.end != mwOpt.end {
											t.Errorf("Expected EDNS0 option:\n\t\"%#v\""+
												"\nfor 0x%04x at %d but got:\n\t\"%#v\"",
												*testOpt, k, i, *mwOpt)
										}
									}
								}
							} else {
								t.Errorf("Expected EDNS0 options 0x%04x but got nothing", k)
							}
						}

						for k := range mw.conf.options {
							if _, ok := test.options[k]; !ok {
								t.Errorf("Got unexpected options 0x%04x", k)
							}
						}
					}

					if test.debugSuffix != nil && *test.debugSuffix != mw.conf.debugSuffix {
						t.Errorf("Expected debug suffix %q but got %q", *test.debugSuffix, mw.conf.debugSuffix)
					}

					if test.streams != nil && *test.streams != mw.conf.streams {
						t.Errorf("Expected %d streams but got %d", *test.streams, mw.conf.streams)
					}

					if test.hotSpot != nil && *test.hotSpot != mw.conf.hotSpot {
						t.Errorf("Expected hotSpot=%v but got %v", *test.hotSpot, mw.conf.hotSpot)
					}

					if test.custAttrs != nil {
						for k, et := range test.custAttrs {
							at, ok := mw.conf.custAttrs[k]
							if !ok {
								t.Errorf("Missing conf attribute %q", k)
							} else if et != at {
								t.Errorf("Unexpected type of conf attribute %q; expected=%d, actual=%d", k, et, at)
							}
						}

						for k, at := range mw.conf.custAttrs {
							if _, ok := test.custAttrs[k]; !ok {
								t.Errorf("Unexpected conf attribute %q=%d", k, at)
							}
						}
					}

					if test.debugID != nil && *test.debugID != mw.conf.debugID {
						t.Errorf("Expected debug id %q but got %q", *test.debugID, mw.conf.debugID)
					}

					if test.passthrough != nil {
						if len(test.passthrough) != len(mw.conf.passthrough) {
							t.Errorf("Expected %d passthrough suffixes but got %d",
								len(test.passthrough), len(mw.conf.passthrough))
						} else {
							for i, s := range test.passthrough {
								if s != mw.conf.passthrough[i] {
									t.Errorf("Expected %q passthrough suffix at %d but got %q",
										s, i, mw.conf.passthrough[i])
								}
							}
						}
					}

					if test.connTimeout != nil && *test.connTimeout != mw.conf.connTimeout {
						t.Errorf("Expected connection timeout %s but got %s", *test.connTimeout, mw.conf.connTimeout)
					}

					if test.autoReqSize != nil && *test.autoReqSize != mw.conf.autoReqSize {
						t.Errorf("Expected automatic request size %v but got %v",
							*test.autoReqSize, mw.conf.autoReqSize)
					}

					if test.maxReqSize != nil && *test.maxReqSize != mw.conf.maxReqSize {
						t.Errorf("Expected request size limit %d but got %d", *test.maxReqSize, mw.conf.maxReqSize)
					}

					if test.autoResAttrs != nil && *test.autoResAttrs != mw.conf.autoResAttrs {
						t.Errorf("Expected automatic response attributes %v but got %v",
							*test.autoResAttrs, mw.conf.autoResAttrs)
					}

					if test.maxResAttrs != nil && *test.maxResAttrs != mw.conf.maxResAttrs {
						t.Errorf("Expected response attributes limit %d but got %d",
							*test.maxResAttrs, mw.conf.maxResAttrs)
					}

					if test.cacheTTL != nil && *test.cacheTTL != mw.conf.cacheTTL {
						t.Errorf("Expected cache TTL %s but got %s", *test.cacheTTL, mw.conf.cacheTTL)
					}

					if test.cacheLimit != nil && *test.cacheLimit != mw.conf.cacheLimit {
						t.Errorf("Expected cache limit %d but got %d", *test.cacheLimit, mw.conf.cacheLimit)
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
