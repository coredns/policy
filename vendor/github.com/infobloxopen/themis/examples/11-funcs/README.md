# 11-Funcs

The example shows how to use functions such as **try** and **concat**.

## Policy with function

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p funcs.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=funcs.yaml
INFO[0000] Parsing policy                                policy=funcs.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Opening storage port                          address=":5552"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

In other terminal run pepcli:
```
$ pepcli -i funcs.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "test"

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "default"

- effect: Permit
  obligation:
    - id: "ls"
      type: "list of strings"
      value: "\"one\",\"two\",\"three\",\"first\",\"second\",\"third\""

```

PDP logs:
```
...
DEBU[0002] Request context                               context="attributes:
- func.(String): \"try\"
- x.(String): \"test\""
DEBU[0002] Response                                      effect=Permit obligations="attributes:
- r.(string): \"test\"" reason="<nil>"
DEBU[0002] Request context                               context="attributes:
- func.(String): \"try\""
DEBU[0002] Response                                      effect=Permit obligations="attributes:
- r.(string): \"default\"" reason="<nil>"
DEBU[0002] Request context                               context="attributes:
- func.(String): \"concat\""
DEBU[0002] Response                                      effect=Permit obligations="attributes:
- ls.(list of strings): \"\\\"one\\\",\\\"two\\\",\\\"three\\\",\\\"first\\\",\\\"second\\\",\\\"third\\\"\"" reason="<nil>"
...
```
