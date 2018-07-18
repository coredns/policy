# 01-Intro

The example shows a couple of simple policies and how to make request decisions on them.

## All Permit Policy

This is a simplest policy which responds with **Permit** to any request. To try the policy run pdpserver:
```
$ pdpserver -v 3 -p all-permit.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=all-permit.yaml
INFO[0000] Parsing policy                                policy=all-permit.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli to make some requests to pdpserver:
```
$ pepcli -i all-permit.requests.yaml test
- effect: Permit

- effect: Permit

- effect: Permit

- effect: Permit

- effect: Permit

- effect: Permit

```

All requests have been permitted as expected. Note that the file with requests contains some attributes and pepcli sends them but pdp ignores all as policy doesn't refer to any.

In terminal with pdpserver you can see how pdpserver processed the requests:
```
...
DEBU[0884] Request context                               context=
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0884] Request context                               context="attributes:
- b.(Boolean): true"
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0884] Request context                               context="attributes:
- s.(String): \"test\""
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0884] Request context                               context="attributes:
- a.(Address): 192.0.2.1"
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0884] Request context                               context="attributes:
- d.(Domain): domain(example.com)"
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0884] Request context                               context="attributes:
- c.(Network): 192.0.2.0/24"
DEBU[0884] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
...
```

## Permit X Test Policy

Another simple example is an example of one attribute policy. Run pdpserver:
```
$ pdpserver -v 3 -p permit-x-test.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=permit-x-test.yaml
INFO[0000] Parsing policy                                policy=permit-x-test.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

And run requests for the policy:
```
$ pepcli -i permit-x-test.requests.yaml test
- effect: Permit

- effect: NotApplicable

```

PDP logs:
```
...
DEBU[0236] Request context                               context="attributes:
- x.(String): \"test\""
DEBU[0236] Response                                      effect=Permit obligations="no attributes" reason="<nil>"
DEBU[0236] Request context                               context="attributes:
- x.(String): \"example\""
DEBU[0236] Response                                      effect=NotApplicable obligations="no attributes" reason="<nil>"
...
```
