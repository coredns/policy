package yast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/infobloxopen/themis/pdp"
	_ "github.com/infobloxopen/themis/pdp/selector"
)

const (
	invalidYAML = `# Invalid YAML
x:
- one
+ two
- three
`

	invalidRootKeysPolicy = `# Policy with invalid keys
attributes:
  x: string

invalid:
- first
`

	simpleAllPermitPolicy = `# Simple All Permit Policy
policies:
  id: Default
  alg: FirstApplicableEffect
  rules:
  - effect: Permit
`

	policyToUpdate = `# Policy to update
attributes:
  a: string
  b: string
  r: string
policies:
  id: Parent policy set
  alg:
    id: mapper
    map:
      attr: a
    default: Deny policy
  policies:
  - id: Deny policy
    alg: FirstApplicableEffect
    rules:
    - effect: Deny
      obligations:
      - r:
         val:
           type: string
           content: Default Deny Policy
  - id: Parent policy
    alg:
      id: mapper
      map:
        attr: b
      default: Deny rule
    rules:
    - id: Deny rule
      effect: Deny
      obligations:
      - r:
         val:
           type: string
           content: Default Deny rule
    - id: Some rule
      effect: Permit
      obligations:
      - r:
         val:
           type: string
           content: Some rule
  - id: Useless policy
    alg: FirstApplicableEffect
    rules:
    - effect: Deny
      obligations:
      - r:
         val:
           type: string
           content: Useless policy
`

	simpleUpdate = `# Simple several commands update
- op: add
  path:
  - Parent policy set
  entity:
    id: Policy Set
    alg: FirstApplicableEffect
    policies:
    - id: Permit Policy
      alg: FirstApplicableEffect
      rules:
      - id: Permit Rule
        effect: permit
        obligations:
        - r:
           val:
             type: string
             content: First Added Update Item

- op: add
  path:
  - Parent policy set
  entity:
    id: Policy
    alg: FirstApplicableEffect
    rules:
    - id: Permit Rule
      effect: permit
      obligations:
      - r:
         val:
           type: string
           content: Second Added Update Item

- op: add
  path:
  - Parent policy set
  - Parent policy
  entity:
    id: Permit Rule
    effect: permit
    obligations:
    - r:
       val:
         type: string
         content: Third Added Update Item

- op: delete
  path:
  - Parent policy set
  - Useless policy
`

	allFeaturePolicies = `# Policies YAML with all features
types:
  flags8:
    meta: flags
    flags: [f00, f01, f02, f03, f04, f05, f06]
  flags16:
    meta: flags
    flags: [f00, f01, f02, f03, f04, f05, f06, f07,
            f10, f11, f12, f13, f14, f15, f16]
  flags32:
    meta: flags
    flags: [f00, f01, f02, f03, f04, f05, f06, f07,
            f10, f11, f12, f13, f14, f15, f16, f17,
            f20, f21, f22, f23, f24, f25, f26, f27,
            f30, f31, f32, f33, f34, f35, f36]
  flags64:
    meta: flags
    flags: [f00, f01, f02, f03, f04, f05, f06, f07,
            f10, f11, f12, f13, f14, f15, f16, f17,
            f20, f21, f22, f23, f24, f25, f26, f27,
            f30, f31, f32, f33, f34, f35, f36, f37,
            f40, f41, f42, f43, f44, f45, f46, f47,
            f50, f51, f52, f53, f54, f55, f56, f57,
            f60, f61, f62, f63, f64, f65, f66, f67,
            f70, f71, f72, f73, f74, f75, f76]

attributes:
  boolAttr: boolean
  strAttr: string
  intAttr: integer
  floatAttr: float
  minAttr: float
  maxAttr: float
  valAttr: float
  addrAttr: address
  netAttr: network
  domAttr: domain
  ssAttr: set of strings
  snAttr: set of networks
  sdAttr: set of domains
  lsAttr: list of strings

policies:
  # Hidden policy set
  alg: FirstApplicableEffect
  target:
  - equal:
    - attr: strAttr
    - val:
        type: string
        content: string
  - any:
    - contains:
      - val:
          type: network
          content: 192.0.2.0/24
      - attr: addrAttr
    - equal:
      - attr: strAttr
      - val:
         type: string
         content: string
    - all:
      - contains:
        - val:
            type: network
            content: 192.0.2.0/24
        - attr: addrAttr
      - equal:
        - val:
           type: string
           content: string
        - attr: strAttr
  policies:
  - id: Permit
    alg: DenyOverrides
    rules:
    - condition:
       not:
       - and:
         - attr: boolAttr
         - or:
           - contains:
             - attr: netAttr
             - val:
                 type: address
                 content: "192.0.2.1"
           - contains:
             - val:
                 type: network
                 content: "192.0.2.0/24"
             - attr: addrAttr
           - contains:
             - attr: sdAttr
             - val:
                 type: domain
                 content: "example.com"
           - contains:
             - val:
                 type: set of strings
                 content:
                 - first
                 - second
                 - third
             - attr: strAttr
           - contains:
             - val:
                 type: set of networks
                 content:
                 - 192.0.2.16/28
                 - 192.0.2.32/28
                 - 2001:db8::/32
             - attr: addrAttr
           - contains:
             - val:
                 type: set of domains
                 content:
                 - example.com
                 - exmaple.net
                 - example.org
             - attr: domAttr
           - equal:
             - attr: strAttr
             - selector:
                uri: "local:content/content-item"
                type: "string"
                path:
                - attr: netAttr
                - attr: domAttr
      effect: Permit
  - id: Nested Mappers Policy Set
    alg:
      id: Mapper
      map:
        attr: lsAttr
      error: Error
      default: Default
      alg:
        id: Mapper
        map:
          selector:
            uri: "local:content/content-item"
            type: "string"
            path:
            - attr: netAttr
            - attr: netAttr
    policies:
    - id: Default
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
        obligations:
        - strAttr:
           val:
             type: string
             content: "Nested Mappers Policy Set Permit"
    - id: Error
      alg: FirstApplicableEffect
      rules:
      - effect: Deny
        obligations:
        - strAttr:
           val:
             type: string
             content: "Nested Mappers Policy Set Deny"
  - id: Nested Mappers Policy
    alg:
      id: Mapper
      map:
        attr: lsAttr
      error: Error
      default: Default
      alg:
        id: Mapper
        map:
          selector:
            uri: "local:content/content-item"
            type: "string"
            path:
            - attr: netAttr
            - attr: netAttr
    rules:
    - id: Default
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "Nested Mappers Policy Permit"
    - id: Error
      effect: Deny
      obligations:
      - strAttr:
         val:
           type: string
           content: "Nested Mappers Policy Deny"
      - lsAttr:
         val:
           type: list of strings
           content:
             - first
             - second
             - third
      - intAttr:
         val:
           type: integer
           content: 9.007199254740992e+15
    - id: IntEqual
      effect: Deny
      target:
      - equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
      condition:
        equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
    - id: FloatEqual
      effect: Deny
      target:
      - equal:
        - attr: floatAttr
        - val:
            type: float
            content: 0.0
      condition:
        equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
    - id: IntGreater
      effect: Deny
      target:
      - equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
      condition:
        greater:
        - attr: intAttr
        - val:
            type: integer
            content: 0
    - id: FloatGreater
      effect: Deny
      target:
      - equal:
        - attr: floatAttr
        - val:
            type: float
            content: 0.0
      condition:
        greater:
        - attr: floatAttr
        - val:
            type: float
            content: 0
    - id: NumAdd
      effect: Deny
      target:
      - equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
      condition:
        greater:
        - add:
            - attr: intAttr
            - attr: floatAttr
        - val:
            type: integer
            content: 10
    - id: NumSubtract
      effect: Deny
      target:
      - equal:
        - attr: intAttr
        - val:
            type: integer
            content: 0
      condition:
        greater:
        - subtract:
            - attr: floatAttr
            - attr: intAttr
        - val:
            type: float
            content: 10.0
    - id: NumMultiply
      effect: Deny
      target:
      - equal:
        - attr: floatAttr
        - val:
            type: float
            content: 10.0
      condition:
        greater:
        - multiply:
            - attr: floatAttr
            - attr: intAttr
        - val:
            type: float
            content: 10.0
    - id: NumDivide
      effect: Deny
      target:
      - equal:
        - attr: floatAttr
        - val:
            type: float
            content: 10.0
      condition:
        greater:
        - divide:
            - attr: floatAttr
            - attr: intAttr
        - val:
            type: float
            content: 10.0
  - id: Float Range Policy
    alg:
      id: Mapper
      map:
        range:
          - attr: minAttr
          - attr: maxAttr
          - attr: valAttr
      alg: FirstApplicableEffect
    rules:
    - id: Below
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: Below

    - id: Above
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: Above

    - id: Within
      effect: Permit
      obligations:
      - floatAttr:
         divide:
           - attr: valAttr
           - attr: minAttr
  - id: Reodering Mapper Policy Set
    alg:
      id: Mapper
      map:
        attr: lsAttr
      order: Internal
      alg: FirstApplicableEffect
    policies:
    - id: first
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "First Rule"
    - id: second
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "Second Rule"
    - id: third
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "Third Rule"
  - id: Reodering Mapper Policy
    alg:
      id: Mapper
      map:
        attr: lsAttr
      order: Internal
      alg: FirstApplicableEffect
    rules:
    - id: first
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "First Rule"
    - id: second
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "Second Rule"
    - id: third
      effect: Permit
      obligations:
      - strAttr:
         val:
           type: string
           content: "Third Rule"
  - id: Flags
    alg: FirstApplicableEffect
    rules:
    - id: Flags8
      effect: Permit
      obligations:
      - lsAttr:
          list of strings:
          - val:
              type: flags8
              content: [f00, f06]

    - id: Flags16
      effect: Permit
      obligations:
      - lsAttr:
          list of strings:
          - val:
              type: flags16
              content: [f00, f16]

    - id: Flags32
      effect: Permit
      obligations:
      - lsAttr:
          list of strings:
          - val:
              type: flags32
              content: [f00, f36]

    - id: Flags64
      effect: Permit
      obligations:
      - lsAttr:
          list of strings:
          - val:
              type: flags64
              content: [f00, f76]

  - id: Flags8MapperRCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom8-map"
          type: "flags8"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    rules:
    - id: f00
      effect: Permit
    - id: f02
      effect: Deny

  - id: Flags16MapperRCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom16-map"
          type: "flags16"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    rules:
    - id: f00
      effect: Permit
    - id: f02
      effect: Deny

  - id: Flags32MapperRCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom32-map"
          type: "flags32"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    rules:
    - id: f00
      effect: Permit
    - id: f02
      effect: Deny

  - id: Flags64MapperRCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom64-map"
          type: "flags64"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    rules:
    - id: f00
      effect: Permit
    - id: f02
      effect: Deny

  - id: Flags8MapperPCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom8-map"
          type: "flags8"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    policies:
    - id: f00
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
    - id: f02
      alg: FirstApplicableEffect
      rules:
      - effect: Deny

  - id: Flags16MapperPCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom16-map"
          type: "flags16"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    policies:
    - id: f00
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
    - id: f02
      alg: FirstApplicableEffect
      rules:
      - effect: Deny

  - id: Flags32MapperPCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom32-map"
          type: "flags32"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    policies:
    - id: f00
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
    - id: f02
      alg: FirstApplicableEffect
      rules:
      - effect: Deny

  - id: Flags64MapperPCA
    alg:
      id: Mapper
      map:
        selector:
          uri: "local:content/content-item-dom64-map"
          type: "flags64"
          path:
          - attr: domAttr
      order: Internal
      alg: FirstApplicableEffect
    policies:
    - id: f00
      alg: FirstApplicableEffect
      rules:
      - effect: Permit
    - id: f02
      alg: FirstApplicableEffect
      rules:
      - effect: Deny

  - id: TryAndConcat
    alg: FirstApplicableEffect
    rules:
    - id: Try
      effect: Permit
      obligations:
      - lsAttr:
          try:
          - attr: lsAttr
          - val:
              type: list of strings
              content: ["one", "two", "three"]
    - id: Concat
      effect: Permit
      obligations:
      - lsAttr:
          concat:
          - attr: lsAttr
          - val:
              type: flags64
              content: [f00, f76]
`

	missingSubAlgPCAPolicy = `# Policies YAML with missing mapper subalgorithm
types:
  flags:
    meta: flags
    flags: ["Deny", "Permit"]

attributes:
  s: string

policies:
  alg:
    id: mapper
    map:
      selector:
        uri: local:content/map
        type: flags
        path:
        - attr: s
  policies:
  - id: Deny
    alg: FirstApplicableEffect
    rules:
    - effect: Deny
  - id: Permit
    alg: FirstApplicableEffect
    rules:
    - effect: Permit
`

	missingSubAlgRCAPolicy = `# Policies YAML with missing mapper subalgorithm
types:
  flags:
    meta: flags
    flags: ["Deny", "Permit"]

attributes:
  s: string

policies:
  alg:
    id: mapper
    map:
      selector:
        uri: local:content/map
        type: flags
        path:
        - attr: s
  rules:
  - effect: Deny
  - effect: Permit
`
)

func TestUnmarshal(t *testing.T) {
	p := Parser{}
	_, err := p.Unmarshal(strings.NewReader(invalidYAML), nil)
	if err == nil {
		t.Errorf("Expected error for invalid YAML but got nothing")
	}

	_, err = p.Unmarshal(strings.NewReader(invalidRootKeysPolicy), nil)
	if err == nil {
		t.Errorf("Expected error for policy with invalid keys but got nothing")
	} else {
		_, ok := err.(*rootKeysError)
		if !ok {
			t.Errorf("Expected *rootTagsError for policy with invalid keys but got %T (%s)", err, err)
		}
	}

	s, err := p.Unmarshal(strings.NewReader(simpleAllPermitPolicy), nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		p, ok := s.Root().(*pdp.Policy)
		if !ok {
			t.Errorf("Expected policy as root item in Simple All Permit Policy but got %T", p)
		} else {
			PID, ok := p.GetID()
			if !ok {
				t.Errorf("Expected %q as Simple All Permit Policy ID but got hidden policy", "Default")
			} else if PID != "Default" {
				t.Errorf("Expected %q as Simple All Permit Policy ID but got %q", "Default", PID)
			}
		}

		r := s.Root().Calculate(&pdp.Context{})
		if r.Effect != pdp.EffectPermit {
			t.Errorf("Expected permit as a response for Simple All Permit Policy but got %d", r.Effect)
		}
	}

	s, err = p.Unmarshal(strings.NewReader(allFeaturePolicies), nil)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
	} else {
		ctx, err := pdp.NewContext(nil, 5, func(i int) (string, pdp.AttributeValue, error) {
			switch i {
			case 0:
				v, err := pdp.MakeValueFromString(pdp.TypeBoolean, "true")
				if err != nil {
					return "", pdp.UndefinedValue, err
				}

				return "boolAttr", v, nil

			case 1:
				v, err := pdp.MakeValueFromString(pdp.TypeString, "string")
				if err != nil {
					return "", pdp.UndefinedValue, err
				}

				return "strAttr", v, nil

			case 2:
				v, err := pdp.MakeValueFromString(pdp.TypeAddress, "192.0.2.1")
				if err != nil {
					return "", pdp.UndefinedValue, err
				}

				return "addrAttr", v, nil

			case 3:
				v, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.0/24")
				if err != nil {
					return "", pdp.UndefinedValue, err
				}

				return "netAttr", v, nil

			case 4:
				v, err := pdp.MakeValueFromString(pdp.TypeString, "example.com")
				if err != nil {
					return "", pdp.UndefinedValue, err
				}

				return "domAttr", v, nil
			}

			return "", pdp.UndefinedValue, fmt.Errorf("no attribute for index %d", i)
		})
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
		} else {
			r := s.Root().Calculate(ctx)
			effect, o, err := r.Effect, r.Obligations, r.Status
			if effect != pdp.EffectDeny {
				if err != nil {
					t.Errorf("Expected deny as a response for Simple All Permit Policy but got %d (%s)", effect, err)
				} else {
					t.Errorf("Expected deny as a response for Simple All Permit Policy but got %d", effect)
				}
			}

			if len(o) < 1 {
				t.Error("Expected at least one obligation")
			} else {
				_, _, v, err := o[0].Serialize(ctx)
				if err != nil {
					t.Errorf("Expected no error but got %T (%s)", err, err)
				} else {
					e := "Nested Mappers Policy Set Deny"
					if v != e {
						t.Errorf("Expected %q but got %q", e, v)
					}
				}
			}
		}
	}
}

func TestMissingSubAlg(t *testing.T) {
	p := Parser{}
	_, err := p.Unmarshal(strings.NewReader(missingSubAlgPCAPolicy), nil)
	if err == nil {
		t.Errorf("expected *missingPCAError but got no error")
	} else if _, ok := err.(*missingPCAError); !ok {
		t.Errorf("expected *missingPCAError but got %T: %s", err, err)
	}

	p = Parser{}
	_, err = p.Unmarshal(strings.NewReader(missingSubAlgRCAPolicy), nil)
	if err == nil {
		t.Errorf("expected *missingRCAError but got no error")
	} else if _, ok := err.(*missingRCAError); !ok {
		t.Errorf("expected *missingRCAError but got %T: %s", err, err)
	}
}

func TestUnmarshalUpdate(t *testing.T) {
	p := Parser{}
	tag := uuid.New()
	s, err := p.Unmarshal(strings.NewReader(policyToUpdate), &tag)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	attrs := map[string]string{
		"a": "Parent policy",
		"b": "Some rule"}
	assertPolicy(s, attrs, "Some rule", "\"some rule\"", t)

	attrs = map[string]string{"a": "Useless policy"}
	assertPolicy(s, attrs, "Useless policy", "\"useless policy\"", t)

	tr, err := s.NewTransaction(&tag)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	u, err := p.UnmarshalUpdate(strings.NewReader(simpleUpdate), tr.Symbols(), tag, uuid.New())
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	err = tr.Apply(u)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	s, err = tr.Commit()
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	attrs = map[string]string{"a": "Policy Set"}
	assertPolicy(s, attrs, "First Added Update Item", "\"new policy set\"", t)

	attrs = map[string]string{"a": "Policy"}
	assertPolicy(s, attrs, "Second Added Update Item", "\"new policy\"", t)

	attrs = map[string]string{
		"a": "Parent policy",
		"b": "Permit Rule"}
	assertPolicy(s, attrs, "Third Added Update Item", "\"new nested policy set\"", t)

	attrs = map[string]string{"a": "Useless policy"}
	assertPolicy(s, attrs, "Default Deny Policy", "\"deleted useless policy\"", t)
}

func assertPolicy(s *pdp.PolicyStorage, attrs map[string]string, e, desc string, t *testing.T) {
	ctx, err := newStringContext(attrs)
	if err != nil {
		t.Errorf("Expected no error for %s but got %T (%s)", desc, err, err)
		return
	}

	r := s.Root().Calculate(ctx)
	if r.Status != nil {
		t.Errorf("Expected no error for %s but got %T (%s)", desc, r.Status, r.Status)
		return
	}

	o := r.Obligations
	if len(o) < 1 {
		t.Errorf("Expected at least one obligation for %s but got nothing", desc)
		return
	}

	_, _, v, err := o[0].Serialize(ctx)
	if err != nil {
		t.Errorf("Expected no error for %s but got %T (%s)", desc, err, err)
		return
	}

	if v != e {
		t.Errorf("Expected %q for %s but got %q", e, desc, v)
	}
}

func newStringContext(m map[string]string) (*pdp.Context, error) {
	names := make([]string, len(m))
	values := make([]string, len(m))
	i := 0
	for k, v := range m {
		names[i] = k
		values[i] = v
		i++
	}

	return pdp.NewContext(nil, len(m), func(i int) (string, pdp.AttributeValue, error) {
		if i >= len(names) {
			return "", pdp.UndefinedValue, fmt.Errorf("no attribute name for index %d", i)
		}
		n := names[i]

		if i >= len(values) {
			return "", pdp.UndefinedValue, fmt.Errorf("no attribute value for index %d", i)
		}
		v := values[i]

		return n, pdp.MakeStringValue(v), nil
	})
}
