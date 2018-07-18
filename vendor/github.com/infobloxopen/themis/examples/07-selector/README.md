# 07-Selector

The example shows policies file with selector and local content file.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p selector.yaml -j content.json
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=selector.yaml
INFO[0000] Parsing policy                                policy=selector.yaml
INFO[0000] Opening content                               content=content.json
INFO[0000] Parsing content                               content=content.json
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i selector.requests.yaml test
- effect: Permit
  obligation:
    - id: "s"
      type: "string"
      value: "Good"

- effect: Deny
  obligation:
    - id: "s"
      type: "string"
      value: "Bad"

- effect: Permit
  obligation:
    - id: "s"
      type: "string"
      value: "Good"

- effect: Deny
  obligation:
    - id: "s"
      type: "string"
      value: "Bad"

```

PDP logs:
```
...
DEBU[0280] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.com)
- a.(Address): 192.0.2.18"
DEBU[0280] Response                                      effect=Permit obligations="attributes:
- s.(string): \"Good\"" reason="<nil>"
DEBU[0280] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.com)
- a.(Address): 2001:db8:1000::1"
DEBU[0280] Response                                      effect=Deny obligations="attributes:
- s.(string): \"Bad\"" reason="<nil>"
DEBU[0280] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(test.com)
- a.(Address): 192.0.2.50"
DEBU[0280] Response                                      effect=Permit obligations="attributes:
- s.(string): \"Good\"" reason="<nil>"
DEBU[0280] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(test.com)
- a.(Address): 2001:db8:3000::1"
DEBU[0280] Response                                      effect=Deny obligations="attributes:
- s.(string): \"Bad\"" reason="<nil>"
...
```
