# 09-Numerical

The example shows policies file with numerical functions.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -pfmt yaml -p numerical.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=numerical.yaml
INFO[0000] Parsing policy                                policy=numerical.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i numerical.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "float"
      value: "2"

- effect: Permit
  obligation:
    - id: "r"
      type: "float"
      value: "10"

- effect: Permit
  obligation:
    - id: "r"
      type: "float"
      value: "3.2"

```

PDP logs:
```
...
DEBU[0014] Request context                               context="attributes:
- actualVal.(Float): 5
- targetVal.(Float): 5"
DEBU[0014] Response                                      effect=Permit obligations="attributes:
- r.(float): \"2\"" reason="<nil>"
DEBU[0014] Request context                               context="attributes:
- actualVal.(Float): 100
- targetVal.(Float): 5"
DEBU[0014] Response                                      effect=Permit obligations="attributes:
- r.(float): \"10\"" reason="<nil>"
DEBU[0014] Request context                               context="attributes:
- actualVal.(Float): 16
- targetVal.(Float): 5"
DEBU[0014] Response                                      effect=Permit obligations="attributes:
- r.(float): \"3.2\"" reason="<nil>"
...
```
