# 05-Target

The example shows policies file with different kinds of target.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p target.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=target.yaml
INFO[0000] Parsing policy                                policy=target.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i target.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "first"

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "second"

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "third"

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "fourth"

- effect: NotApplicable

```

PDP logs:
```
...
DEBU[0003] Request context                               context="attributes:
- a.(Address): 192.0.2.1
- x.(String): \"test\"
- c.(Network): 192.0.2.0/24"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- r.(string): \"first\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- x.(String): \"test\"
- c.(Network): 192.0.2.32/28
- a.(Address): 192.0.2.17"
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- r.(string): \"second\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- c.(Network): 192.0.2.32/28
- a.(Address): 192.0.2.33
- x.(String): \"test\""
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- r.(string): \"third\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- c.(Network): 192.0.2.32/28
- a.(Address): 192.0.3.1
- x.(String): \"test\""
DEBU[0003] Response                                      effect=Permit obligations="attributes:
- r.(string): \"fourth\"" reason="<nil>"
DEBU[0003] Request context                               context="attributes:
- x.(String): \"example\"
- c.(Network): 192.0.2.32/28
- a.(Address): 192.0.3.1"
DEBU[0003] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
...
```
