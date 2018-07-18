# 04-Rule

The example shows policies file with full featured rule.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p rule.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=rule.yaml
INFO[0000] Parsing policy                                policy=rule.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i rule.requests.yaml test
- effect: Permit
  obligation:
    - id: "a"
      type: "address"
      value: "192.0.2.1"

- effect: NotApplicable

- effect: NotApplicable

- effect: NotApplicable

```

PDP logs:
```
...
DEBU[0034] Request context                               context="attributes:
- x.(String): \"test\"
- c.(Network): 192.0.2.16/28
- b.(Boolean): false"
DEBU[0034] Response                                      effect=Permit obligations="attributes:
- a.(address): \"192.0.2.1\"" reason="<nil>"
DEBU[0034] Request context                               context="attributes:
- x.(String): \"test\"
- c.(Network): 192.0.2.16/28
- b.(Boolean): true"
DEBU[0034] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
DEBU[0034] Request context                               context="attributes:
- x.(String): \"test\"
- c.(Network): 192.0.2.0/24
- b.(Boolean): false"
DEBU[0034] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
DEBU[0034] Request context                               context="attributes:
- x.(String): \"example\""
DEBU[0034] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
...
```
