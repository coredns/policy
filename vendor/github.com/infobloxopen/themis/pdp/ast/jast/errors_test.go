package jast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
)

var testCases = []map[string]string{
	{
		"policy": `
{
  "attributes": {
    "a": "sometype"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownTypeError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "effect": "Permit"
      }
    ],
    "policies": [
      {
        "id": "Default",
        "alg": "FirstApplicableEffect",
        "rules": [
          {
            "effect": "Permit"
          }
        ]
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &policyAmbiguityError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect"
  }
}
`,
		"err": fmt.Sprintf("%T", &policyMissingKeyError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "SomeAlg",
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownRCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingRCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": {
      "id": "mapper",
      "map": {
        "attr": "a"
      },
      "default": "Default"
    },
    "rules": [
      {
        "id": "Error",
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingDefaultRuleRCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": {
      "id": "mapper",
      "map": {
        "attr": "a"
      },
      "error": "Error"
    },
    "rules": [
      {
        "id": "Default",
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingErrorRuleRCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "SomeAlg",
    "policies": [
      {
        "id": "Default",
        "alg": "FirstApplicableEffect",
        "rules": [
          {
            "effect": "Permit"
          }
        ]
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownPCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "policies": [
      {
        "id": "Default",
        "alg": "FirstApplicableEffect",
        "rules": [
          {
            "effect": "Permit"
          }
        ]
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingPCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": {
      "id": "mapper",
      "map": {
        "attr": "a"
      },
      "default": "Default"
    },
    "policies": [
      {
        "id": "Error",
        "alg": "FirstApplicableEffect",
        "rules": [
          {
            "effect": "Deny"
          }
        ]
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingDefaultPolicyPCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": {
      "id": "mapper",
      "map": {
        "attr": "a"
      },
      "error": "Error"
    },
    "policies": [
      {
        "id": "Default",
        "alg": "FirstApplicableEffect",
        "rules": [
          {
            "effect": "Permit"
          }
        ]
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &missingErrorPolicyPCAError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "boolean"
  },
  "policies": {
    "id": "Default",
    "alg": {
      "id": "mapper",
      "map": {
        "attr": "a"
      }
    },
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &mapperArgumentTypeError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "condition": {
          "attr": "a"
        },
        "effect": "Permit"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &conditionTypeError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "effect": "Bye"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownEffectError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "some": [
          {
            "attr": "a"
          },
          {
            "val": {
              "type": "string",
              "content": "string"
            }
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownMatchFunctionError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "boolean"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "contains": [
          {
            "attr": "a"
          },
          {
            "val": {
              "type": "string",
              "content": "string"
            }
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &matchFunctionCastError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "equal": [
          {
            "attr": "a"
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &matchFunctionArgsNumberError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string",
    "b": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "equal": [
          {
            "attr": "a"
          },
          {
            "attr": "b"
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &matchFunctionBothAttrsError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "equal": [
          {
            "val": {
              "type": "string",
              "content": "string"
            }
          },
          {
            "val": {
              "type": "string",
              "content": "string"
            }
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &matchFunctionBothValuesError{}),
	},

	{
		"policy": `
{
  "attributes": {
    "a": "string",
    "b": "string"
  },
  "policies": {
    "id": "Default",
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "equal": [
          {
            "some": "a"
          },
          {
            "some": "b"
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Deny"
      }
    ]
  }
}
`,
		"err": fmt.Sprintf("%T", &unknownFunctionError{}),
	},
}

func TestUnmarshalErrors(t *testing.T) {
	p := Parser{}

	for _, tc := range testCases {
		_, err := p.Unmarshal(strings.NewReader(tc["policy"]), nil)
		if err == nil {
			t.Errorf("Expected %s error but got nothing", tc["err"])
		} else if e := fmt.Sprintf("%T", err); e != tc["err"] {
			t.Errorf("Expected %s error but got %s", tc["err"], e)
		}
	}
}

var testCasesUpdate = []map[string]string{
	{
		"update": `
[
  {
    "op": "some",
    "path": [
      "Parent policy set"
    ],
    "entity": {
      "id": "Policy Set",
      "alg": "FirstApplicableEffect",
      "rules": {
        "effect": "Permit"
      }
    }
  }
]
`,
		"err": fmt.Sprintf("%T", &unknownPolicyUpdateOperationError{}),
	},

	{
		"update": `
[
  {
    "op": "add",
    "path": [
      "Parent policy set"
    ],
    "entity": {
      "id": "Policy Set"
    }
  }
]
`,
		"err": fmt.Sprintf("%T", &entityMissingKeyError{}),
	},

	{
		"update": `
[
  {
    "op": "add",
    "path": [
      "Parent policy set"
    ],
    "entity": {
      "id": "Policy Set",
      "alg": "FirstApplicableEffect",
      "rules": [
        {
          "effect": "Permit"
        }
      ],
      "policies": [
        {
          "id": "Permit Policy",
          "alg": "FirstApplicableEffect",
          "rules": [
            {
              "effect": "Permit"
            }
          ]
        }
      ]
    }
  }
]
`,
		"err": fmt.Sprintf("%T", &entityAmbiguityError{}),
	},
}

func TestUnmarshalUpdateErrors(t *testing.T) {
	p := Parser{}
	tag := uuid.New()
	s, err := p.Unmarshal(strings.NewReader(policyToUpdate), &tag)
	if err != nil {
		t.Errorf("Expected no error but got %T (%s)", err, err)
		return
	}

	for _, tc := range testCasesUpdate {
		tr, err := s.NewTransaction(&tag)
		if err != nil {
			t.Errorf("Expected no error but got %T (%s)", err, err)
			return
		}

		_, err = p.UnmarshalUpdate(strings.NewReader(tc["update"]), tr.Symbols(), tag, uuid.New())
		if err == nil {
			t.Errorf("Expected %s error but got nothing", tc["err"])
			return
		}

		if e := fmt.Sprintf("%T", err); e != tc["err"] {
			t.Errorf("Expected %s error but got %s", tc["err"], e)
		}
	}
}
