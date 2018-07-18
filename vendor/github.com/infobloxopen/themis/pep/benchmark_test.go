package pep

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	_ "github.com/infobloxopen/themis/pdp/selector"
	"github.com/infobloxopen/themis/pdpserver/server"
)

const (
	oneStageBenchmarkPolicySet = `# Policy set for benchmark
attributes:
  k3: domain
  x: string

policies:
  alg:
    id: mapper
    map:
      selector:
        uri: "local:content/second"
        path:
        - attr: k3
        type: list of strings
    default: DefaultRule
    alg: FirstApplicableEffect
  rules:
  - id: DefaultRule
    effect: Deny
    obligations:
    - x:
       val:
         type: string
         content: DefaultRule
  - id: First
    effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: First
  - id: Second
    effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: Second
  - id: Third
    effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: Third
  - id: Fourth
    effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: Fourth
  - id: Fifth
    effect: Permit
    obligations:
    - x:
       val:
         type: string
         content: Fifth
`

	twoStageBenchmarkPolicySet = `# Policy set for benchmark 2-level nesting policy
attributes:
  k2: string
  k3: domain
  x: string

policies:
  alg:
    id: mapper
    map:
      selector:
        uri: "local:content/first"
        path:
        - attr: k2
        type: string
    default: DefaultPolicy

  policies:
  - id: DefaultPolicy
    alg: FirstApplicableEffect
    rules:
    - effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: DefaultPolicy

  - id: P1
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/second"
          path:
          - attr: k3
          type: list of strings
      default: DefaultRule
      alg: FirstApplicableEffect
    rules:
    - id: DefaultRule
      effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: P1.DefaultRule
    - id: First
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P1.First
    - id: Second
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P1.Second

  - id: P2
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/second"
          path:
          - attr: k3
          type: list of strings
      default: DefaultRule
      alg: FirstApplicableEffect
    rules:
    - id: DefaultRule
      effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: P2.DefaultRule
    - id: Second
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P2.Second
    - id: Third
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P2.Third

  - id: P3
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/second"
          path:
          - attr: k3
          type: list of strings
      default: DefaultRule
      alg: FirstApplicableEffect
    rules:
    - id: DefaultRule
      effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: P3.DefaultRule
    - id: Third
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P3.Third
    - id: Fourth
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P3.Fourth

  - id: P4
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/second"
          path:
          - attr: k3
          type: list of strings
      default: DefaultRule
      alg: FirstApplicableEffect
    rules:
    - id: DefaultRule
      effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: P4.DefaultRule
    - id: Fourth
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P4.Fourth
    - id: Fifth
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P4.Fifth

  - id: P5
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/second"
          path:
          - attr: k3
          type: list of strings
      default: DefaultRule
      alg: FirstApplicableEffect
    rules:
    - id: DefaultRule
      effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: P5.DefaultRule
    - id: Fifth
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P5.Fifth
    - id: First
      effect: Permit
      obligations:
      - x:
         val:
           type: string
           content: P5.First
`

	threeStageBenchmarkPolicySet = `# Policy set for benchmark 3-level nesting policy
attributes:
  k1: string
  k2: string
  k3: domain
  x: string

policies:
  alg: FirstApplicableEffect
  policies:
  - target:
    - equal:
      - attr: k1
      - val:
          type: string
          content: "Left"
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/first"
          path:
          - attr: k2
          type: string
      default: DefaultPolicy

    policies:
    - id: DefaultPolicy
      alg: FirstApplicableEffect
      rules:
      - effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: DefaultPolicy

    - id: P1
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P1.DefaultRule
      - id: First
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P1.First
      - id: Second
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P1.Second

    - id: P2
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P2.DefaultRule
      - id: Second
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P2.Second
      - id: Third
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P2.Third

    - id: P3
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P3.DefaultRule
      - id: Third
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P3.Third
      - id: Fourth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P3.Fourth

    - id: P4
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P4.DefaultRule
      - id: Fourth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P4.Fourth
      - id: Fifth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P4.Fifth

    - id: P5
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P5.DefaultRule
      - id: Fifth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P5.Fifth
      - id: First
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P5.First

  - target:
    - equal:
      - attr: k1
      - val:
          type: string
          content: "Right"
    alg:
      id: mapper
      map:
        selector:
          uri: "local:content/first"
          path:
          - attr: k2
          type: string
      default: DefaultPolicy

    policies:
    - id: DefaultPolicy
      alg: FirstApplicableEffect
      rules:
      - effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: DefaultPolicy

    - id: P1
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P1.DefaultRule
      - id: First
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P1.First
      - id: Second
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P1.Second

    - id: P2
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P2.DefaultRule
      - id: Second
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P2.Second
      - id: Third
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P2.Third

    - id: P3
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P3.DefaultRule
      - id: Third
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P3.Third
      - id: Fourth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P3.Fourth

    - id: P4
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P4.DefaultRule
      - id: Fourth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P4.Fourth
      - id: Fifth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P4.Fifth

    - id: P5
      alg:
        id: mapper
        map:
          selector:
            uri: "local:content/second"
            path:
            - attr: k3
            type: list of strings
        default: DefaultRule
        alg: FirstApplicableEffect
      rules:
      - id: DefaultRule
        effect: Deny
        obligations:
        - x:
           val:
             type: string
             content: P5.DefaultRule
      - id: Fifth
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P5.Fifth
      - id: First
        effect: Permit
        obligations:
        - x:
           val:
             type: string
             content: P5.First

  - alg: FirstApplicableEffect
    rules:
    - effect: Deny
      obligations:
      - x:
         val:
           type: string
           content: Root Deny
`

	benchmarkContent = `{
    "id": "content",
    "items": {
        "first": {
            "keys": ["string"],
            "type": "string",
            "data": {
                "First": "P1",
                "Second": "P2",
                "Third": "P3",
                "Fourth": "P4",
                "Fifth": "P5",
                "Sixth": "P6",
                "Seventh": "P7"
            }
        },
        "second": {
            "keys": ["domain"],
            "type": "list of strings",
            "data": {
                "first.example.com": ["First", "Third"],
                "second.example.com": ["Second", "Fourth"],
                "third.example.com": ["Third", "Fifth"],
                "first.test.com": ["Fourth", "Sixth"],
                "second.test.com": ["Fifth", "Seventh"],
                "third.test.com": ["Sixth", "First"],
                "first.example.com": ["Seventh", "Second"],
                "second.example.com": ["Firth", "Fourth"],
                "third.example.com": ["Second", "Fifth"],
                "first.test.com": ["Third", "Sixth"],
                "second.test.com": ["Fourth", "Seventh"],
                "third.test.com": ["Fifth", "First"]
            }
        }
    }
}`
)

type decisionRequest struct {
	Direction string `pdp:"k1"`
	Policy    string `pdp:"k2"`
	Domain    string `pdp:"k3,domain"`
}

type decisionResponse struct {
	Effect int    `pdp:"Effect"`
	Reason error  `pdp:"Reason"`
	X      string `pdp:"x"`
}

func (r decisionResponse) String() string {
	if r.Reason != nil {
		return fmt.Sprintf("Effect: %q, Reason: %q, X: %q",
			pdp.EffectNameFromEnum(r.Effect),
			r.Reason,
			r.X,
		)
	}

	return fmt.Sprintf("Effect: %q, X: %q", pdp.EffectNameFromEnum(r.Effect), r.X)
}

var (
	directionOpts = []string{
		"Left",
		"Right",
	}

	policySetOpts = []string{
		"First",
		"Second",
		"Third",
		"Fourth",
		"Fifth",
		"Sixth",
		"Seventh",
	}

	domainOpts = []string{
		"first.example.com",
		"second.example.com",
		"third.example.com",
		"first.test.com",
		"second.test.com",
		"third.test.com",
		"first.example.com",
		"second.example.com",
		"third.example.com",
		"first.test.com",
		"second.test.com",
		"third.test.com",
	}

	decisionRequests []decisionRequest
	rawRequests      []pb.Msg
)

type testRequest3Keys struct {
	k1 string `pdp:"k1"`
	k2 string `pdp:"k2"`
	k3 string `pdp:"k3,domain"`
}

func init() {
	decisionRequests = make([]decisionRequest, 0x40000)
	for i := range decisionRequests {
		decisionRequests[i] = decisionRequest{
			Direction: directionOpts[rand.Intn(len(directionOpts))],
			Policy:    policySetOpts[rand.Intn(len(policySetOpts))],
			Domain:    domainOpts[rand.Intn(len(domainOpts))],
		}
	}

	rawRequests = make([]pb.Msg, len(decisionRequests))
	for i := range rawRequests {
		b := make([]byte, 128)

		m, err := makeRequest(testRequest3Keys{
			k1: directionOpts[rand.Intn(len(directionOpts))],
			k2: policySetOpts[rand.Intn(len(policySetOpts))],
			k3: domainOpts[rand.Intn(len(domainOpts))],
		}, b)
		if err != nil {
			panic(fmt.Errorf("failed to create %d raw request: %s", i+1, err))
		}

		rawRequests[i] = m
	}

}

func benchmarkPolicySet(name, p string, b *testing.B) {
	pdpServer, _, c := startPDPServer(p, nil, b)
	defer func() {
		c.Close()
		if logs := pdpServer.Stop(); len(logs) > 0 {
			b.Logf("server logs:\n%s", logs)
		}
	}()

	b.Run(name, func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			in := decisionRequests[n%len(decisionRequests)]

			var out decisionResponse
			c.Validate(in, &out)

			if (out.Effect != pdp.EffectDeny &&
				out.Effect != pdp.EffectPermit &&
				out.Effect != pdp.EffectNotApplicable) ||
				out.Reason != nil {
				b.Fatalf("unexpected response: %s", out)
			}
		}
	})
}

func BenchmarkOneStagePolicySet(b *testing.B) {
	benchmarkPolicySet("OneStagePolicySet", oneStageBenchmarkPolicySet, b)
}

func BenchmarkTwoStagePolicySet(b *testing.B) {
	benchmarkPolicySet("TwoStagePolicySet", twoStageBenchmarkPolicySet, b)
}

func BenchmarkThreeStagePolicySet(b *testing.B) {
	benchmarkPolicySet("ThreeStagePolicySet", threeStageBenchmarkPolicySet, b)
}

func BenchmarkUnaryRaw(b *testing.B) {
	pdpServer, _, c := startPDPServer(threeStageBenchmarkPolicySet, nil, b)
	defer func() {
		c.Close()
		if logs := pdpServer.Stop(); len(logs) > 0 {
			b.Logf("server logs:\n%s", logs)
		}
	}()

	name := "UnaryRaw"

	b.Run(name, func(b *testing.B) {
		var (
			out        pdp.Response
			assignment [16]pdp.AttributeAssignment
		)
		for n := 0; n < b.N; n++ {
			in := rawRequests[n%len(rawRequests)]

			out.Obligations = assignment[:]
			c.Validate(in, &out)

			err := assertBenchMsg(&out, "%q request %d", name, n)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkUnaryWithCache(b *testing.B) {
	pdpServer, _, c := startPDPServer(threeStageBenchmarkPolicySet, nil, b,
		WithMaxRequestSize(128),
		WithCacheTTL(15*time.Minute),
	)
	defer func() {
		c.Close()
		if logs := pdpServer.Stop(); len(logs) > 0 {
			b.Logf("server logs:\n%s", logs)
		}
	}()

	name := "UnaryWithCache"

	cc := 10
	var (
		out        pdp.Response
		assignment [16]pdp.AttributeAssignment
	)
	for n := 0; n < cc; n++ {
		in := rawRequests[n%cc]

		out.Obligations = assignment[:]
		c.Validate(in, &out)

		err := assertBenchMsg(&out, "%q request %d", name, n)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Run(name, func(b *testing.B) {
		var (
			out        pdp.Response
			assignment [16]pdp.AttributeAssignment
		)
		for n := 0; n < b.N; n++ {
			in := rawRequests[n%cc]

			out.Obligations = assignment[:]
			c.Validate(in, &out)

			err := assertBenchMsg(&out, "%q request %d", name, n)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func benchmarkStreamingClient(name string, ports []uint16, b *testing.B, opts ...Option) {
	if len(ports) != 0 && len(ports) != 2 {
		b.Fatalf("only 0 for single PDP and 2 for 2 PDP ports supported but got %d", len(ports))
	}

	streams := 96

	opts = append(opts,
		WithStreams(streams),
	)
	pdpSrv, pdpSrvAlt, c := startPDPServer(threeStageBenchmarkPolicySet, ports, b, opts...)
	defer func() {
		c.Close()
		if logs := pdpSrv.Stop(); len(logs) > 0 {
			b.Logf("primary server logs:\n%s", logs)
		}
		if pdpSrvAlt != nil {
			if logs := pdpSrvAlt.Stop(); len(logs) > 0 {
				b.Logf("secondary server logs:\n%s", logs)
			}
		}
	}()

	b.Run(name, func(b *testing.B) {
		assignments := make(chan []pdp.AttributeAssignment, streams)
		for i := 0; i < cap(assignments); i++ {
			assignments <- make([]pdp.AttributeAssignment, 16)
		}

		th := make(chan int, streams)
		for n := 0; n < b.N; n++ {
			th <- 0
			go func(i int) {
				defer func() { <-th }()

				var out pdp.Response

				assignment := <-assignments
				defer func() { assignments <- assignment }()
				out.Obligations = assignment

				c.Validate(rawRequests[i%len(rawRequests)], &out)
				if err := assertBenchMsg(&out, "%q request %d", name, i); err != nil {
					panic(err)
				}
			}(n)
		}
	})
}

func BenchmarkStreamingClient(b *testing.B) {
	benchmarkStreamingClient("StreamingClient", nil, b)
}

func BenchmarkRoundRobinStreamingClient(b *testing.B) {
	benchmarkStreamingClient("RoundRobinStreamingClient",
		[]uint16{
			5555,
			5556,
		},
		b,
		WithRoundRobinBalancer("127.0.0.1:5555", "127.0.0.1:5556"),
	)
}

func BenchmarkHotSpotStreamingClient(b *testing.B) {
	benchmarkStreamingClient("HotSpotStreamingClient",
		[]uint16{
			5555,
			5556,
		},
		b,
		WithHotSpotBalancer("127.0.0.1:5555", "127.0.0.1:5556"),
	)
}

func BenchmarkStreamingClientWithCache(b *testing.B) {
	streams := 96

	pdpServer, _, c := startPDPServer(threeStageBenchmarkPolicySet, nil, b,
		WithStreams(streams),
		WithMaxRequestSize(128),
		WithCacheTTL(15*time.Minute),
	)
	defer func() {
		c.Close()
		if logs := pdpServer.Stop(); len(logs) > 0 {
			b.Logf("server logs:\n%s", logs)
		}
	}()

	name := "StreamingClientWithCache"

	cc := 10 * streams
	var (
		out        pdp.Response
		assignment [16]pdp.AttributeAssignment
	)
	for n := 0; n < cc; n++ {
		in := rawRequests[n%cc]

		out.Obligations = assignment[:]
		c.Validate(in, &out)

		err := assertBenchMsg(&out, "%q request %d", name, n)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.Run(name, func(b *testing.B) {
		assignments := make(chan []pdp.AttributeAssignment, streams)
		for i := 0; i < cap(assignments); i++ {
			assignments <- make([]pdp.AttributeAssignment, 16)
		}

		th := make(chan int, streams)
		for n := 0; n < b.N; n++ {
			th <- 0
			go func(i int) {
				defer func() { <-th }()

				var out pdp.Response

				assignment := <-assignments
				defer func() { assignments <- assignment }()
				out.Obligations = assignment

				c.Validate(rawRequests[i%cc], &out)
				if err := assertBenchMsg(&out, "%q request %d", name, i); err != nil {
					panic(err)
				}
			}(n)
		}
	})
}

func waitForPortOpened(address string) error {
	var (
		c   net.Conn
		err error
	)

	for i := 0; i < 20; i++ {
		after := time.After(500 * time.Millisecond)
		c, err = net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err == nil {
			return c.Close()
		}

		<-after
	}

	return err
}

func waitForPortClosed(address string) error {
	var (
		c   net.Conn
		err error
	)

	for i := 0; i < 20; i++ {
		after := time.After(500 * time.Millisecond)
		c, err = net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err != nil {
			return nil
		}

		c.Close()
		<-after
	}

	return fmt.Errorf("port at %s hasn't been closed yet", address)
}

type loggedServer struct {
	s *server.Server
	b *bytes.Buffer
}

func newServer(opts ...server.Option) *loggedServer {
	s := &loggedServer{
		b: new(bytes.Buffer),
	}

	logger := log.New()
	logger.Out = s.b
	logger.Level = log.ErrorLevel
	opts = append(opts,
		server.WithLogger(logger),
	)

	s.s = server.NewServer(opts...)
	return s
}

func (s *loggedServer) Stop() string {
	s.s.Stop()
	return s.b.String()
}

func startPDPServer(p string, ports []uint16, b *testing.B, opts ...Option) (*loggedServer, *loggedServer, Client) {
	var (
		primary   *loggedServer
		secondary *loggedServer
	)

	service := "127.0.0.1:5555"
	if len(ports) > 0 {
		service = fmt.Sprintf("127.0.0.1:%d", ports[0])
	}

	primary = newServer(
		server.WithServiceAt(service),
	)

	if err := primary.s.ReadPolicies(strings.NewReader(p)); err != nil {
		b.Fatalf("can't read policies: %s", err)
	}

	if err := primary.s.ReadContent(strings.NewReader(benchmarkContent)); err != nil {
		b.Fatalf("can't read content: %s", err)
	}

	if err := waitForPortClosed(service); err != nil {
		b.Fatalf("port still in use: %s", err)
	}
	go func() {
		if err := primary.s.Serve(); err != nil {
			b.Fatalf("primary server failed: %s", err)
		}
	}()

	if err := waitForPortOpened(service); err != nil {
		if logs := primary.Stop(); len(logs) > 0 {
			b.Logf("primary server logs:\n%s", logs)
		}

		b.Fatalf("can't connect to PDP server: %s", err)
	}

	if len(ports) > 1 {
		service := fmt.Sprintf("127.0.0.1:%d", ports[1])
		secondary = newServer(
			server.WithServiceAt(service),
		)

		if err := secondary.s.ReadPolicies(strings.NewReader(p)); err != nil {
			if logs := primary.Stop(); len(logs) > 0 {
				b.Logf("primary server logs:\n%s", logs)
			}
			b.Fatalf("can't read policies: %s", err)
		}

		if err := secondary.s.ReadContent(strings.NewReader(benchmarkContent)); err != nil {
			if logs := primary.Stop(); len(logs) > 0 {
				b.Logf("primary server logs:\n%s", logs)
			}
			b.Fatalf("can't read content: %s", err)
		}

		if err := waitForPortClosed(service); err != nil {
			b.Fatalf("port still in use: %s", err)
		}
		go func() {
			if err := secondary.s.Serve(); err != nil {
				b.Fatalf("secondary server failed: %s", err)
			}
		}()

		if err := waitForPortOpened(service); err != nil {
			if logs := secondary.Stop(); len(logs) > 0 {
				b.Logf("secondary server logs:\n%s", logs)
			}
			if logs := primary.Stop(); len(logs) > 0 {
				b.Logf("primary server logs:\n%s", logs)
			}

			b.Fatalf("can't connect to PDP server: %s", err)
		}
	}

	c := NewClient(opts...)
	if err := c.Connect(service); err != nil {
		if secondary != nil {
			if logs := secondary.Stop(); len(logs) > 0 {
				b.Logf("secondary server logs:\n%s", logs)
			}
		}
		if logs := primary.Stop(); len(logs) > 0 {
			b.Logf("primary server logs:\n%s", logs)
		}

		b.Fatalf("can't connect to PDP server: %s", err)
	}

	return primary, secondary, c
}

func assertBenchMsg(r *pdp.Response, s string, args ...interface{}) error {
	if r.Effect != pdp.EffectDeny && r.Effect != pdp.EffectPermit && r.Effect != pdp.EffectNotApplicable {
		desc := fmt.Sprintf(s, args...)
		return fmt.Errorf("unexpected response effect for %s: %s", desc, pdp.EffectNameFromEnum(r.Effect))
	}

	if r.Status != nil {
		desc := fmt.Sprintf(s, args...)
		return fmt.Errorf("unexpected response status for %s: %s", desc, r.Status)
	}

	return nil
}
