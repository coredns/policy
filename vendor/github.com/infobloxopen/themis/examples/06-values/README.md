# 06-Values

The example shows policies file with different immediate values.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p values.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=values.yaml
INFO[0000] Parsing policy                                policy=values.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Opening storage port                          address=":5552"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
```

In other terminal run pepcli:
```
$ pepcli -i values.requests.yaml test
- effect: Permit
  obligation:
    - id: "s"
      type: "string"
      value: "example"

- effect: Permit
  obligation:
    - id: "a"
      type: "address"
      value: "192.0.2.1"

- effect: Permit
  obligation:
    - id: "c"
      type: "network"
      value: "192.0.2.0/28"

- effect: Permit
  obligation:
    - id: "d"
      type: "domain"
      value: "example.net"

    - id: "sd"
      type: "set of domains"
      value: "\"test.com\",\"example.com\""

- effect: Permit
  obligation:
    - id: "ss"
      type: "set of strings"
      value: "\"first\",\"second\""

- effect: Permit
  obligation:
    - id: "sn"
      type: "set of networks"
      value: "\"192.0.2.0/28\",\"192.0.2.16/28\""

- effect: Permit
  obligation:
    - id: "s"
      type: "string"
      value: "first-rule"

- effect: Permit
  obligation:
    - id: "s"
      type: "string"
      value: "second-rule"

```

PDP logs:
```
...
DEBU[0003] Request context                               context="attributes:
- s.(String): \"test\""
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- s.(string): \"example\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- s.(String): \"example\"
- c.(Network): 192.0.2.0/31"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- a.(address): \"192.0.2.1\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- s.(String): \"example\"
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.13"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- c.(network): \"192.0.2.0/28\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- s.(String): \"example\"
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.16
- d.(Domain): domain(example.com)"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- d.(domain): \"example.net\"
- sd.(set of domains): \"\\\"test.com\\\",\\\"example.com\\\"\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- s.(String): \"first\"
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.16
- d.(Domain): domain(example.net)"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- ss.(set of strings): \"\\\"first\\\",\\\"second\\\"\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.16
- d.(Domain): domain(example.net)
- s.(String): \"third\""
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- sn.(set of networks): \"\\\"192.0.2.0/28\\\",\\\"192.0.2.16/28\\\"\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- d.(Domain): domain(example.net)
- s.(String): \"first-rule\"
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.33"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- s.(string): \"first-rule\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- s.(String): \"second-rule\"
- c.(Network): 192.0.2.2/31
- a.(Address): 192.0.2.33
- d.(Domain): domain(example.net)"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- s.(string): \"second-rule\"" reason="<nil>"
...
```
