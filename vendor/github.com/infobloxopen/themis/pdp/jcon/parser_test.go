package jcon

import (
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/infobloxopen/go-trees/domain"
	"github.com/infobloxopen/themis/pdp"
)

const (
	jsonStream = `{
	"ID": "Test",
	"Items": {
		"first": {
			"type": "set of strings",
			"keys": ["string", "address", "string"],
			"data": {
				"x": {
					"127.0.0.1": {
						"y": {
							"z": false,
							"t": null
						}
					}
				}
			}
		},
		"second": {
			"data": {
				"first": {
					"192.0.2.0/28": {
						"example.com": ["2001:db8::/40", "2001:db8:0100::/40", "2001:db8:0200::/40"],
						"example.net": ["2001:db8:1000::/40", "2001:db8:1100::/40", "2001:db8:1200::/40"]
					},
					"192.0.2.16/28": {
						"example.com": ["2001:db8:2000::/40", "2001:db8:2100::/40", "2001:db8:2200::/40"],
						"example.net": ["2001:db8:3000::/40", "2001:db8:3100::/40", "2001:db8:3200::/40"]
					},
					"192.0.2.32/28": {
						"example.com": ["2001:db8:4000::/40", "2001:db8:4100::/40", "2001:db8:4200::/40"],
						"example.net": ["2001:db8:5000::/40", "2001:db8:5100::/40", "2001:db8:5200::/40"]
					}
				},
				"second": {
					"2001:db8::/36": {
						"example.com": ["2001:db8::/40", "2001:db8:0100::/40", "2001:db8:0200::/40"],
						"example.net": ["2001:db8:1000::/40", "2001:db8:1100::/40", "2001:db8:1200::/40"]
					},
					"2001:db8:1000::/36": {
						"example.com": ["2001:db8:2000::/40", "2001:db8:2100::/40", "2001:db8:2200::/40"],
						"example.net": ["2001:db8:3000::/40", "2001:db8:3100::/40", "2001:db8:3200::/40"]
					},
					"2001:db8:2000::/36": {
						"example.com": ["2001:db8:4000::/40", "2001:db8:4100::/40", "2001:db8:4200::/40"],
						"example.net": ["2001:db8:5000::/40", "2001:db8:5100::/40", "2001:db8:5200::/40"]
					}
				}
			},
			"type": "set of networks",
			"keys": ["string", "network", "domain"]
		},
		"third": {
			"type": {
				"meta": "flags",
				"name": "ft8",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07"]
			},
			"keys": ["string"],
			"data": {
				"first": ["f00", "f02", "f04"],
				"second": ["f01", "f03", "f05"],
				"third": ["f02", "f04", "f06"]
			}
		},
		"fourth": {
			"type": {
				"meta": "flags",
				"name": "ft16",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17"]
			},
			"keys": ["string"],
			"data": {
				"first": ["f00", "f02", "f04"],
				"second": ["f01", "f03", "f05"],
				"third": ["f02", "f04", "f06"]
			}
		},
		"fifth": {
			"type": {
				"meta": "flags",
				"name": "ft32",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37"]
			},
			"keys": ["string"],
			"data": {
				"first": ["f00", "f02", "f04"],
				"second": ["f01", "f03", "f05"],
				"third": ["f02", "f04", "f06"]
			}
		},
		"sixth": {
			"type": {
				"meta": "flags",
				"name": "ft64",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
				          "f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
				          "f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
				          "f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
				          "f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77"]
			},
			"keys": ["string"],
			"data": {
				"first": ["f00", "f02", "f04"],
				"second": ["f01", "f03", "f05"],
				"third": ["f02", "f04", "f06"]
			}
		},
		"seventh": {
			"type": "ft8",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f00", "f02", "f04"],
				"192.0.2.32/28": ["f01", "f03", "f05"],
				"2001:db8::/33": ["f02", "f04", "f06"],
				"2001:db8:8000::1": ["f03", "f05", "f07"]
			}
		},
		"eighth": {
			"type": "ft16",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f00", "f02", "f04"],
				"192.0.2.32/28": ["f01", "f03", "f05"],
				"2001:db8::/33": ["f02", "f04", "f06"],
				"2001:db8:8000::1": ["f03", "f05", "f07"]
			}
		},
		"ninth": {
			"type": "ft32",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f00", "f02", "f04"],
				"192.0.2.32/28": ["f01", "f03", "f05"],
				"2001:db8::/33": ["f02", "f04", "f06"],
				"2001:db8:8000::1": ["f03", "f05", "f07"]
			}
		},
		"tenth": {
			"type": "ft64",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f00", "f02", "f04"],
				"192.0.2.32/28": ["f01", "f03", "f05"],
				"2001:db8::/33": ["f02", "f04", "f06"],
				"2001:db8:8000::1": ["f03", "f05", "f07"]
			}
		},
		"eleventh": {
			"type": "ft8",
			"keys": ["domain"],
			"data": {
				"example.com": ["f00", "f02", "f04"],
				"example.net": ["f01", "f03", "f05"],
				"example.org": ["f02", "f04", "f06"]
			}
		},
		"twelveth": {
			"type": "ft16",
			"keys": ["domain"],
			"data": {
				"example.com": ["f00", "f02", "f04"],
				"example.net": ["f01", "f03", "f05"],
				"example.org": ["f02", "f04", "f06"]
			}
		},
		"thirteenth": {
			"type": "ft32",
			"keys": ["domain"],
			"data": {
				"example.com": ["f00", "f02", "f04"],
				"example.net": ["f01", "f03", "f05"],
				"example.org": ["f02", "f04", "f06"]
			}
		},
		"fourteenth": {
			"type":  "ft64",
			"keys": ["domain"],
			"data": {
				"example.com": ["f00", "f02", "f04"],
				"example.net": ["f01", "f03", "f05"],
				"example.org": ["f02", "f04", "f06"]
			}
		}
	}
}`

	jsonUpdateStream = `[
  {
    "op": "Add",
    "path": ["first", "update"],
    "entity": {
      "type": "set of strings",
      "keys": ["address", "string"],
      "data": {
        "127.0.0.2": {
          "n": {
            "p": false,
            "q": null
          }
        }
      }
    }
  },
  {
    "op": "Delete",
    "path": ["second", "second", "2001:db8:1000::/36", "example.net"]
  },
  {
    "op": "Delete",
    "path": ["second", "second", "2001:db8:1000::/36"]
  },
  {
    "op": "Delete",
    "path": ["third", "third"]
  },
  {
    "op": "Add",
    "path": ["third", "fourth"],
    "entity": {
      "type": "ft8",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["fourth", "third"]
  },
  {
    "op": "Add",
    "path": ["fourth", "fourth"],
    "entity": {
      "type": "ft16",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["fifth", "third"]
  },
  {
    "op": "Add",
    "path": ["fifth", "fourth"],
    "entity": {
      "type": "ft32",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["sixth", "third"]
  },
  {
    "op": "Add",
    "path": ["sixth", "fourth"],
    "entity": {
      "type": "ft64",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["seventh", "192.0.2.32/28"]
  },
  {
    "op": "Add",
    "path": ["seventh", "192.0.2.48/28"],
    "entity": {
      "type": "ft8",
      "data": ["f04", "f06", "f00"]
    }
  },
  {
    "op": "Delete",
    "path": ["seventh", "2001:db8:8000::1"]
  },
  {
    "op": "Add",
    "path": ["seventh", "2001:db8:8000::2"],
    "entity": {
      "type": "ft8",
      "data": ["f05", "f07", "f01"]
    }
  },
  {
    "op": "Delete",
    "path": ["eighth", "192.0.2.32/28"]
  },
  {
    "op": "Add",
    "path": ["eighth", "192.0.2.48/28"],
    "entity": {
      "type": "ft16",
      "data": ["f04", "f06", "f00"]
    }
  },
  {
    "op": "Delete",
    "path": ["eighth", "2001:db8:8000::1"]
  },
  {
    "op": "Add",
    "path": ["eighth", "2001:db8:8000::2"],
    "entity": {
      "type": "ft16",
      "data": ["f05", "f07", "f01"]
    }
  },
  {
    "op": "Delete",
    "path": ["ninth", "192.0.2.32/28"]
  },
  {
    "op": "Add",
    "path": ["ninth", "192.0.2.48/28"],
    "entity": {
      "type": "ft32",
      "data": ["f04", "f06", "f00"]
    }
  },
  {
    "op": "Delete",
    "path": ["ninth", "2001:db8:8000::1"]
  },
  {
    "op": "Add",
    "path": ["ninth", "2001:db8:8000::2"],
    "entity": {
      "type": "ft32",
      "data": ["f05", "f07", "f01"]
    }
  },
  {
    "op": "Delete",
    "path": ["tenth", "192.0.2.32/28"]
  },
  {
    "op": "Add",
    "path": ["tenth", "192.0.2.48/28"],
    "entity": {
      "type": "ft64",
      "data": ["f04", "f06", "f00"]
    }
  },
  {
    "op": "Delete",
    "path": ["tenth", "2001:db8:8000::1"]
  },
  {
    "op": "Add",
    "path": ["tenth", "2001:db8:8000::2"],
    "entity": {
      "type": "ft64",
      "data": ["f05", "f07", "f01"]
    }
  },
  {
    "op": "Delete",
    "path": ["eleventh", "example.org"]
  },
  {
    "op": "Add",
    "path": ["eleventh", "example.gov"],
    "entity": {
      "type": "ft8",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["twelveth", "example.org"]
  },
  {
    "op": "Add",
    "path": ["twelveth", "example.gov"],
    "entity": {
      "type": "ft16",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["thirteenth", "example.org"]
  },
  {
    "op": "Add",
    "path": ["thirteenth", "example.gov"],
    "entity": {
      "type": "ft32",
      "data": ["f03", "f05", "f07"]
    }
  },
  {
    "op": "Delete",
    "path": ["fourteenth", "example.org"]
  },
  {
    "op": "Add",
    "path": ["fourteenth", "example.gov"],
    "entity": {
      "type": "ft64",
      "data": ["f03", "f05", "f07"]
    }
  }
]`

	jsonAllMapsStream = `{
	"ID": "AllMaps",
	"Items": {
		"str-map": {
			"type": "string",
			"keys": ["string"],
			"data": {
				"key-1": "value-1",
				"key-2": "value-2",
				"key-3": "value-3"
			}
		},
		"str8-map": {
			"type": {
				"meta": "flags",
				"name": "ft8",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07"]
			},
			"keys": ["string"],
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			}
		},
		"str16-map": {
			"type": {
				"meta": "flags",
				"name": "ft16",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17"]
			},
			"keys": ["string"],
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			}
		},
		"str32-map": {
			"type": {
				"meta": "flags",
				"name": "ft32",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37"]
			},
			"keys": ["string"],
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			}
		},
		"str64-map": {
			"type": {
				"meta": "flags",
				"name": "ft64",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
				          "f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
				          "f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
				          "f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
				          "f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77"]
			},
			"keys": ["string"],
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			}
		},
		"net-map": {
			"type": "string",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": "value-1",
				"192.0.2.32/28": "value-2",
				"192.0.2.48/28": "value-3"
			}
		},
		"net8-map": {
			"type": "ft8",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			}
		},
		"net16-map": {
			"type": "ft16",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			}
		},
		"net32-map": {
			"type": "ft32",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			}
		},
		"net64-map": {
			"type": "ft64",
			"keys": ["network"],
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			}
		},
		"dom-map": {
			"type": "string",
			"keys": ["domain"],
			"data": {
				"example.com": "value-1",
				"example.net": "value-2",
				"example.org": "value-3"
			}
		},
		"dom8-map": {
			"type":  "ft8",
			"keys": ["domain"],
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			}
		},
		"dom16-map": {
			"type":  "ft16",
			"keys": ["domain"],
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			}
		},
		"dom32-map": {
			"type":  "ft32",
			"keys": ["domain"],
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			}
		},
		"dom64-map": {
			"type":  "ft64",
			"keys": ["domain"],
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			}
		}
	}
}`

	jsonPostprocessAllMapsStream = `{
	"ID": "AllMaps",
	"Items": {
		"str-map": {
			"data": {
				"key-1": "value-1",
				"key-2": "value-2",
				"key-3": "value-3"
			},
			"type": "string",
			"keys": ["string"]
		},
		"str8-map": {
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			},
			"type": {
				"meta": "flags",
				"name": "ft8",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07"]
			},
			"keys": ["string"]
		},
		"str16-map": {
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			},
			"type": {
				"meta": "flags",
				"name": "ft16",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17"]
			},
			"keys": ["string"]
		},
		"str32-map": {
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			},
			"type": {
				"meta": "flags",
				"name": "ft32",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37"]
			},
			"keys": ["string"]
		},
		"str64-map": {
			"data": {
				"key-1": ["f01"],
				"key-2": ["f02", "f04"],
				"key-3": ["f01", "f02"]
			},
			"type": {
				"meta": "flags",
				"name": "ft64",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
				          "f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
				          "f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
				          "f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
				          "f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77"]
			},
			"keys": ["string"]
		},
		"net-map": {
			"data": {
				"192.0.2.16/28": "value-1",
				"192.0.2.32/28": "value-2",
				"192.0.2.48/28": "value-3"
			},
			"type": "string",
			"keys": ["network"]
		},
		"net8-map": {
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			},
			"type": "ft8",
			"keys": ["network"]
		},
		"net16-map": {
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			},
			"type": "ft16",
			"keys": ["network"]
		},
		"net32-map": {
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			},
			"type": "ft32",
			"keys": ["network"]
		},
		"net64-map": {
			"data": {
				"192.0.2.16/28": ["f01"],
				"192.0.2.32/28": ["f02", "f04"],
				"2001:db8::/33": ["f01", "f02"],
				"2001:db8:8000::1": ["f04", "f06"]
			},
			"type": "ft64",
			"keys": ["network"]
		},
		"dom-map": {
			"data": {
				"example.com": "value-1",
				"example.net": "value-2",
				"example.org": "value-3"
			},
			"type": "string",
			"keys": ["domain"]
		},
		"dom8-map": {
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			},
			"type": "ft8",
			"keys": ["domain"]
		},
		"dom16-map": {
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			},
			"type": "ft16",
			"keys": ["domain"]
		},
		"dom32-map": {
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			},
			"type": "ft32",
			"keys": ["domain"]
		},
		"dom64-map": {
			"data": {
				"example.com": ["f01"],
				"example.net": ["f02", "f04"],
				"example.org": ["f01", "f02"]
			},
			"type": "ft64",
			"keys": ["domain"]
		}
	}
}`

	jsonAllValuesStream = `{
	"ID": "AllValues",
	"Items": {
		"boolean": {
			"type": "boolean",
			"keys": ["string"],
			"data": {
				"key": true
			}
		},
		"string": {
			"type": "string",
			"keys": ["string"],
			"data": {
				"key": "value"
			}
		},
        "integer": {
            "type": "integer",
            "keys": ["string"],
            "data": {
                "key": 9.007199254740992e+15
            }
        },
		"address": {
			"type": "address",
			"keys": ["string"],
			"data": {
				"key": "192.0.2.1"
			}
		},
		"network": {
			"type": "network",
			"keys": ["string"],
			"data": {
				"key": "192.0.2.0/24"
			}
		},
		"domain": {
			"type": "domain",
			"keys": ["string"],
			"data": {
				"key": "example.com"
			}
		},
		"[]set of strings": {
			"type": "set of strings",
			"keys": ["string"],
			"data": {
				"key": [
					"1-first",
					"2-second",
					"3-third"
				]
			}
		},
		"{}set of strings": {
			"type": "set of strings",
			"keys": ["string"],
			"data": {
				"key": {
					"1-first": "skip me",
					"2-second": {"skip": "me"},
					"3-third": ["skip", "me"]
				}
			}
		},
		"set of networks": {
			"type": "set of networks",
			"keys": ["string"],
			"data": {
				"key": [
					"192.0.2.16/28",
					"192.0.2.32/28",
					"2001:db8::/32"
				]
			}
		},
		"set of domains": {
			"type": "set of domains",
			"keys": ["string"],
			"data": {
				"key": [
					"example.com",
					"example.net",
					"example.org"
				]
			}
		},
		"list of strings": {
			"type": "list of strings",
			"keys": ["string"],
			"data": {
				"key": [
					"1-first",
					"2-second",
					"3-third"
				]
			}
		},
		"flags8": {
			"type": {
				"meta": "flags",
				"name": "ft8",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07"]
			},
			"keys": ["domain"],
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			}
		},
		"flags16": {
			"type": {
				"meta": "flags",
				"name": "ft16",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17"]
			},
			"keys": ["domain"],
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			}
		},
		"flags32": {
			"type": {
				"meta": "flags",
				"name": "ft32",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37"]
			},
			"keys": ["domain"],
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			}
		},
		"flags64": {
			"type": {
				"meta": "flags",
				"name": "ft64",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
				          "f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
				          "f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
				          "f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
				          "f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77"]
			},
			"keys": ["domain"],
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			}
		}
	}
}`

	jsonPostprocessAllValuesStream = `{
	"ID": "AllValues",
	"Items": {
		"boolean": {
			"data": {
				"key": true
			},
			"type": "boolean",
			"keys": ["string"]
		},
		"string": {
			"data": {
				"key": "value"
			},
			"type": "string",
			"keys": ["string"]
		},
        "integer": {
            "type": "integer",
            "keys": ["string"],
            "data": {
                "key": 9.007199254740992e+15
            }
        },
		"address": {
			"data": {
				"key": "192.0.2.1"
			},
			"type": "address",
			"keys": ["string"]
		},
		"network": {
			"data": {
				"key": "192.0.2.0/24"
			},
			"type": "network",
			"keys": ["string"]
		},
		"domain": {
			"data": {
				"key": "example.com"
			},
			"type": "domain",
			"keys": ["string"]
		},
		"[]set of strings": {
			"data": {
				"key": [
					"1-first",
					"2-second",
					"3-third"
				]
			},
			"type": "set of strings",
			"keys": ["string"]
		},
		"{}set of strings": {
			"data": {
				"key": {
					"1-first": "skip me",
					"2-second": {"skip": "me"},
					"3-third": ["skip", "me"]
				}
			},
			"type": "set of strings",
			"keys": ["string"]
		},
		"set of networks": {
			"data": {
				"key": [
					"192.0.2.16/28",
					"192.0.2.32/28",
					"2001:db8::/32"
				]
			},
			"type": "set of networks",
			"keys": ["string"]
		},
		"set of domains": {
			"data": {
				"key": [
					"example.com",
					"example.net",
					"example.org"
				]
			},
			"type": "set of domains",
			"keys": ["string"]
		},
		"list of strings": {
			"data": {
				"key": [
					"1-first",
					"2-second",
					"3-third"
				]
			},
			"type": "list of strings",
			"keys": ["string"]
		},
		"flags8": {
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			},
			"type": {
				"meta": "flags",
				"name": "ft8",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07"]
			},
			"keys": ["domain"]
		},
		"flags16": {
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			},
			"type": {
				"meta": "flags",
				"name": "ft16",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17"]
			},
			"keys": ["domain"]
		},
		"flags32": {
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			},
			"type": {
				"meta": "flags",
				"name": "ft32",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37"]
			},
			"keys": ["domain"]
		},
		"flags64": {
			"data": {
				"key": ["f00", "f02", "f04", "f06"]
			},
			"type": {
				"meta": "flags",
				"name": "ft64",
				"flags": ["f00", "f01", "f02", "f03", "f04", "f05", "f06", "f07",
				          "f10", "f11", "f12", "f13", "f14", "f15", "f16", "f17",
				          "f20", "f21", "f22", "f23", "f24", "f25", "f26", "f27",
				          "f30", "f31", "f32", "f33", "f34", "f35", "f36", "f37",
				          "f40", "f41", "f42", "f43", "f44", "f45", "f46", "f47",
				          "f50", "f51", "f52", "f53", "f54", "f55", "f56", "f57",
				          "f60", "f61", "f62", "f63", "f64", "f65", "f66", "f67",
				          "f70", "f71", "f72", "f73", "f74", "f75", "f76", "f77"]
			},
			"keys": ["domain"]
		}
	}
}`
)

func TestUnmarshal(t *testing.T) {
	c, err := Unmarshal(strings.NewReader(jsonStream), nil)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		lc, err := c.Get("missing")
		if err == nil {
			t.Errorf("Expected error but got local content item: %#v", lc)
		}

		lc, err = c.Get("first")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			addr, err := pdp.MakeValueFromString(pdp.TypeAddress, "127.0.0.1")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{pdp.MakeStringValue("x"), addr, pdp.MakeStringValue("y")}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"z\",\"t\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected [%s] but got [%s]", e, s)
					}
				}
			}

			addr, err = pdp.MakeValueFromString(pdp.TypeAddress, "127.0.0.2")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{pdp.MakeStringValue("x"), addr, pdp.MakeStringValue("y")}
				r, err := lc.Get(path, nil)
				if err == nil {
					s, err := r.Serialize()
					if err != nil {
						s = err.Error()
					}
					t.Errorf("Expected error but got result %s", s)
				}
			}
		}

		lc, err = c.Get("second")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.4/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{
					pdp.MakeStringValue("first"),
					n,
					pdp.MakeDomainValue(makeTestDN(t, "example.com")),
				}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"2001:db8::/40\",\"2001:db8:100::/40\",\"2001:db8:200::/40\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected [%s] but got [%s]", e, s)
					}
				}
			}

			n, err = pdp.MakeValueFromString(pdp.TypeNetwork, "2001:db8:1000:1::/64")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{
					pdp.MakeStringValue("second"),
					n,
					pdp.MakeDomainValue(makeTestDN(t, "example.net")),
				}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"2001:db8:3000::/40\",\"2001:db8:3100::/40\",\"2001:db8:3200::/40\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected [%s] but got [%s]", e, s)
					}
				}
			}

			n, err = pdp.MakeValueFromString(pdp.TypeNetwork, "2001:db8:3000:1::/64")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{
					pdp.MakeStringValue("second"),
					n,
					pdp.MakeDomainValue(makeTestDN(t, "example.net")),
				}
				r, err := lc.Get(path, nil)
				if err == nil {
					s, err := r.Serialize()
					if err != nil {
						s = err.Error()
					}
					t.Errorf("Expected error but got result %s", s)
				}
			}
		}
	}

	c, err = Unmarshal(strings.NewReader(jsonAllMapsStream), nil)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		lc, err := c.Get("str-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value-2"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("str8-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("str16-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("str32-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("str64-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("net-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "value-2"
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("net8-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"f02\",\"f04\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("net16-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"f02\",\"f04\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("net32-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"f02\",\"f04\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("net64-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "\"f02\",\"f04\""
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("dom-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value-2"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom8-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom16-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom32-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom64-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}
	}

	c, err = Unmarshal(strings.NewReader(jsonPostprocessAllMapsStream), nil)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		lc, err := c.Get("str-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("key-2")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value-2"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("net-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "192.0.2.44/30")
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				path := []pdp.Expression{n}
				r, err := lc.Get(path, nil)
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else {
					e := "value-2"
					s, err := r.Serialize()
					if err != nil {
						t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
					} else if s != e {
						t.Errorf("Expected %q but got %q", e, s)
					}
				}
			}
		}

		lc, err = c.Get("dom-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value-2"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom8-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom16-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom32-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("dom64-map")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeDomainValue(makeTestDN(t, "example.net")),
			}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"f02\",\"f04\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}
	}

	c, err = Unmarshal(strings.NewReader(jsonAllValuesStream), nil)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		path := []pdp.Expression{pdp.MakeStringValue("key")}

		lc, err := c.Get("boolean")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "true"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("string")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("address")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "192.0.2.1"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("network")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "192.0.2.0/24"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("domain")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "example.com"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("[]set of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("{}set of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("set of networks")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"192.0.2.16/28\",\"192.0.2.32/28\",\"2001:db8::/32\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("set of domains")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"example.com\",\"example.net\",\"example.org\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("list of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}
	}

	c, err = Unmarshal(strings.NewReader(jsonPostprocessAllValuesStream), nil)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		path := []pdp.Expression{pdp.MakeStringValue("key")}

		lc, err := c.Get("boolean")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "true"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("string")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "value"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("address")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "192.0.2.1"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("network")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "192.0.2.0/24"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("domain")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "example.com"
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("[]set of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("{}set of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("set of networks")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"192.0.2.16/28\",\"192.0.2.32/28\",\"2001:db8::/32\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("set of domains")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"example.com\",\"example.net\",\"example.org\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}

		lc, err = c.Get("list of strings")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"1-first\",\"2-second\",\"3-third\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected %q but got %q", e, s)
				}
			}
		}
	}
}

func TestUnmarshalUpdate(t *testing.T) {
	s := pdp.NewLocalContentStorage(nil)

	tag := uuid.New()
	c, err := Unmarshal(strings.NewReader(jsonStream), &tag)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		return
	}

	s = s.Add(c)
	tr, err := s.NewTransaction("Test", &tag)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		return
	}

	u, err := UnmarshalUpdate(strings.NewReader(jsonUpdateStream), "Test", tag, uuid.New(), tr.Symbols())
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		return
	}

	err = tr.Apply(u)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		return
	}

	s, err = tr.Commit(s)
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		return
	}

	lc, err := s.Get("Test", "first")
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		addr, err := pdp.MakeValueFromString(pdp.TypeAddress, "127.0.0.2")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{pdp.MakeStringValue("update"), addr, pdp.MakeStringValue("n")}
			r, err := lc.Get(path, nil)
			if err != nil {
				t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
			} else {
				e := "\"p\",\"q\""
				s, err := r.Serialize()
				if err != nil {
					t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
				} else if s != e {
					t.Errorf("Expected [%s] but got [%s]", e, s)
				}
			}
		}
	}

	lc, err = s.Get("Test", "second")
	if err != nil {
		t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
	} else {
		n, err := pdp.MakeValueFromString(pdp.TypeNetwork, "2001:db8:1000:1::/64")
		if err != nil {
			t.Errorf("Expected no error but got (%T):\n\t%s", err, err)
		} else {
			path := []pdp.Expression{
				pdp.MakeStringValue("second"),
				n,
				pdp.MakeDomainValue(makeTestDN(t, "example.com")),
			}
			r, err := lc.Get(path, nil)
			if err == nil {
				s, err := r.Serialize()
				if err != nil {
					s = err.Error()
				}
				t.Errorf("Expected error but got result %s", s)
			}
		}
	}
}

func makeTestDN(t *testing.T, s string) domain.Name {
	d, err := domain.MakeNameFromString(s)
	if err != nil {
		t.Fatalf("can't create domain name from string %q: %s", s, err)
	}

	return d
}
