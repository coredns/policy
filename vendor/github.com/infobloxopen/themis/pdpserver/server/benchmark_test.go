package server

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/infobloxopen/go-trees/domain"
	log "github.com/sirupsen/logrus"

	"github.com/infobloxopen/themis/pdp"
	pb "github.com/infobloxopen/themis/pdp-service"
	"github.com/infobloxopen/themis/pdp/ast"
	"github.com/infobloxopen/themis/pdp/jcon"
	_ "github.com/infobloxopen/themis/pdp/selector"
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

type requestAttributeValue struct {
	k string
	v pdp.AttributeValue
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

	benchmarkContentStorage          *pdp.LocalContentStorage
	oneStageBenchmarkPolicyStorage   *pdp.PolicyStorage
	twoStageBenchmarkPolicyStorage   *pdp.PolicyStorage
	threeStageBenchmarkPolicyStorage *pdp.PolicyStorage

	benchmarkRequests []*pb.Msg
)

func init() {
	log.SetLevel(log.ErrorLevel)

	c, err := jcon.Unmarshal(strings.NewReader(benchmarkContent), nil)
	if err != nil {
		panic(err)
	}

	benchmarkContentStorage = pdp.NewLocalContentStorage([]*pdp.LocalContent{c})
	parser := ast.NewYAMLParser()

	oneStageBenchmarkPolicyStorage, err = parser.Unmarshal(strings.NewReader(oneStageBenchmarkPolicySet), nil)
	if err != nil {
		panic(err)
	}

	twoStageBenchmarkPolicyStorage, err = parser.Unmarshal(strings.NewReader(twoStageBenchmarkPolicySet), nil)
	if err != nil {
		panic(err)
	}

	threeStageBenchmarkPolicyStorage, err = parser.Unmarshal(strings.NewReader(threeStageBenchmarkPolicySet), nil)
	if err != nil {
		panic(err)
	}

	benchmarkRequests = make([]*pb.Msg, 0x40000)
	for i := range benchmarkRequests {
		b := make([]byte, 128)

		dn, err := domain.MakeNameFromString(domainOpts[rand.Intn(len(domainOpts))])
		if err != nil {
			panic(fmt.Errorf("failed to make domain for request %d: %s", i+1, err))
		}

		n, err := pdp.MarshalRequestAssignmentsToBuffer(b, []pdp.AttributeAssignment{
			pdp.MakeStringAssignment("k1", directionOpts[rand.Intn(len(directionOpts))]),
			pdp.MakeStringAssignment("k2", policySetOpts[rand.Intn(len(policySetOpts))]),
			pdp.MakeDomainAssignment("k3", dn),
		})
		if err != nil {
			panic(fmt.Errorf("failed to marshal request %d: %s", i+1, err))
		}

		benchmarkRequests[i] = &pb.Msg{Body: b[:n]}
	}
}

func benchmarkPolicySet(p *pdp.PolicyStorage, b *testing.B) {
	s := NewServer()
	s.p = p
	s.c = benchmarkContentStorage

	var a [1]pdp.AttributeAssignment

	for n := 0; n < b.N; n++ {
		r, err := s.Validate(nil, benchmarkRequests[n%len(benchmarkRequests)])
		if err != nil {
			b.Fatalf("Expected no error while evaluating policies at %d iteration but got: %s", n+1, err)
		}

		effect, n, err := pdp.UnmarshalResponseToAssignmentsArray(r.Body, a[:])
		if err != nil {
			b.Fatalf("Expected no error while unmarshalling response at %d iteration but got: %s", n+1, err)
		}

		if effect >= pdp.EffectIndeterminate {
			b.Fatalf("Expected specific result of policy evaluation at %d iteration but got %s",
				n+1, pdp.EffectNameFromEnum(effect))
		}
	}
}

func BenchmarkOneStagePolicySet(b *testing.B) {
	benchmarkPolicySet(oneStageBenchmarkPolicyStorage, b)
}

func BenchmarkTwoStagePolicySet(b *testing.B) {
	benchmarkPolicySet(twoStageBenchmarkPolicyStorage, b)
}

func BenchmarkThreeStagePolicySet(b *testing.B) {
	benchmarkPolicySet(threeStageBenchmarkPolicyStorage, b)
}

func benchmarkRawPolicySet(p *pdp.PolicyStorage, b *testing.B) {
	s := NewServer()
	s.p = p
	s.c = benchmarkContentStorage

	var a [1]pdp.AttributeAssignment

	for n := 0; n < b.N; n++ {
		s.RLock()
		p := s.p
		c := s.c
		s.RUnlock()

		r := s.rawValidate(p, c, benchmarkRequests[n%len(benchmarkRequests)].Body)

		effect, n, err := pdp.UnmarshalResponseToAssignmentsArray(r, a[:])
		if err != nil {
			b.Fatalf("Expected no error while unmarshalling response at %d iteration but got: %s", n+1, err)
		}

		if effect >= pdp.EffectIndeterminate {
			b.Fatalf("Expected specific result of policy evaluation at %d iteration but got %s",
				n+1, pdp.EffectNameFromEnum(effect))
		}
	}
}

func BenchmarkRawOneStagePolicySet(b *testing.B) {
	benchmarkRawPolicySet(oneStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawTwoStagePolicySet(b *testing.B) {
	benchmarkRawPolicySet(twoStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawThreeStagePolicySet(b *testing.B) {
	benchmarkRawPolicySet(threeStageBenchmarkPolicyStorage, b)
}

func benchmarkRawPolicySetWithAllocator(p *pdp.PolicyStorage, b *testing.B) {
	s := NewServer()
	s.p = p
	s.c = benchmarkContentStorage

	var a [1]pdp.AttributeAssignment
	buf := []byte{}

	for n := 0; n < b.N; n++ {
		s.RLock()
		p := s.p
		c := s.c
		s.RUnlock()

		r := s.rawValidateWithAllocator(p, c, benchmarkRequests[n%len(benchmarkRequests)].Body, func(n int) ([]byte, error) {
			if len(buf) < n {
				buf = make([]byte, n)
			}

			return buf, nil
		})

		effect, n, err := pdp.UnmarshalResponseToAssignmentsArray(r, a[:])
		if err != nil {
			b.Fatalf("Expected no error while unmarshalling response at %d iteration but got: %s", n+1, err)
		}

		if effect >= pdp.EffectIndeterminate {
			b.Fatalf("Expected specific result of policy evaluation at %d iteration but got %s",
				n+1, pdp.EffectNameFromEnum(effect))
		}
	}
}

func BenchmarkRawOneStagePolicySetWithAllocator(b *testing.B) {
	benchmarkRawPolicySetWithAllocator(oneStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawTwoStagePolicySetWithAllocator(b *testing.B) {
	benchmarkRawPolicySetWithAllocator(twoStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawThreeStagePolicySetWithAllocator(b *testing.B) {
	benchmarkRawPolicySetWithAllocator(threeStageBenchmarkPolicyStorage, b)
}

func benchmarkRawPolicySetToBuffer(p *pdp.PolicyStorage, b *testing.B) {
	s := NewServer()
	s.p = p
	s.c = benchmarkContentStorage

	var a [1]pdp.AttributeAssignment
	var buf [1024]byte

	for n := 0; n < b.N; n++ {
		s.RLock()
		p := s.p
		c := s.c
		s.RUnlock()

		r := s.rawValidateToBuffer(p, c, benchmarkRequests[n%len(benchmarkRequests)].Body, buf[:])

		effect, n, err := pdp.UnmarshalResponseToAssignmentsArray(r, a[:])
		if err != nil {
			b.Fatalf("Expected no error while unmarshalling response at %d iteration but got: %s", n+1, err)
		}

		if effect >= pdp.EffectIndeterminate {
			b.Fatalf("Expected specific result of policy evaluation at %d iteration but got %s",
				n+1, pdp.EffectNameFromEnum(effect))
		}
	}
}

func BenchmarkRawOneStagePolicySetToBuffer(b *testing.B) {
	benchmarkRawPolicySetToBuffer(oneStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawTwoStagePolicySetToBuffer(b *testing.B) {
	benchmarkRawPolicySetToBuffer(twoStageBenchmarkPolicyStorage, b)
}

func BenchmarkRawThreeStagePolicySetToBuffer(b *testing.B) {
	benchmarkRawPolicySetToBuffer(threeStageBenchmarkPolicyStorage, b)
}
