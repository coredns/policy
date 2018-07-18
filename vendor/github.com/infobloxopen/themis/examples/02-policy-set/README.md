# 02-Policy-Set

The example shows policies file with full featured policy set.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p policy-set.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=policy-set.yaml
INFO[0000] Parsing policy                                policy=policy-set.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i policy-set.requests.yaml test
- effect: Permit
  obligation:
    - id: "a"
      type: "address"
      value: "192.0.2.1"

- effect: Deny
  obligation:
    - id: "a"
      type: "address"
      value: "192.0.2.1"

- effect: NotApplicable

- effect: Indeterminate{P}
  reason: "#99: #9a: Failed to process request: #02 (policy set \"Test Policy Set\">hidden policy set>target>any>all>match>equal>first argument>attr(z.String)): Missing attribute"

```

PDP logs:
```
...
DEBU[0243] Request context                               context="attributes:
- x.(String): \"test\"
- z.(String): \"example\""
DEBU[0243] Response                                      effect=Permit obligations="attributes:
- a.(address): \"192.0.2.1\"" reason="<nil>"
DEBU[0243] Request context                               context="attributes:
- x.(String): \"test\"
- z.(String): \"test\""
DEBU[0243] Response                                      effect=Deny obligations="attributes:
- a.(address): \"192.0.2.1\"" reason="<nil>"
DEBU[0243] Request context                               context="attributes:
- x.(String): \"example\"
- z.(String): \"test\""
DEBU[0243] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
DEBU[0243] Request context                               context="attributes:
- x.(String): \"test\""
DEBU[0243] Response                                      effect="Indeterminate{P}" obligations="no attributes" reason="#02 (policy set \"Test Policy Set\">hidden policy set>target>any>all>match>equal>first argument>attr(z.String)): Missing attribute"
...
```
