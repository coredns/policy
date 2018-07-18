# 08-Mapper

The example shows policies file with mapper combining algorithm and local content file.

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p mapper.yaml -j content.json
INFO[0000] Starting PDP server                          
INFO[0000] Loading policy                                policy=mapper.yaml
INFO[0000] Parsing policy                                policy=mapper.yaml
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
$ pepcli -i mapper.requests.yaml test
- effect: Deny

- effect: Deny
  obligation:
    - id: "err"
      type: "string"
      value: "Can't calculate policy id"

- effect: Permit
  obligation:
    - id: "p"
      type: "string"
      value: "First PermitNet"

- effect: Deny
  obligation:
    - id: "p"
      type: "string"
      value: "First DenyCom"

- effect: Permit
  obligation:
    - id: "p"
      type: "string"
      value: "Second PermitCom"

- effect: Deny
  obligation:
    - id: "p"
      type: "string"
      value: "Second DenyNet"

- effect: Permit
  obligation:
    - id: "p"
      type: "string"
      value: "External Second"

- effect: Permit
  obligation:
    - id: "p"
      type: "string"
      value: "Internal First"

```

PDP logs:
```
...
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"Unknown\""
DEBU[0033] Response                                      effect=Deny obligations="no attributes" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.com)"
DEBU[0033] Response                                      effect=Deny obligations="attributes:
- err.(string): \"Can't calculate policy id\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"First\"
- d.(Domain): domain(example.net)"
DEBU[0033] Response                                      effect=Permit obligations="attributes:
- p.(string): \"First PermitNet\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"First\"
- d.(Domain): domain(example.com)"
DEBU[0033] Response                                      effect=Deny obligations="attributes:
- p.(string): \"First DenyCom\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"Second\"
- d.(Domain): domain(example.com)"
DEBU[0033] Response                                      effect=Permit obligations="attributes:
- p.(string): \"Second PermitCom\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"Second\"
- d.(Domain): domain(example.net)"
DEBU[0033] Response                                      effect=Deny obligations="attributes:
- p.(string): \"Second DenyNet\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"External\""
DEBU[0033] Response                                      effect=Permit obligations="attributes:
- p.(string): \"External Second\"" reason="<nil>"
DEBU[0033] Request context                               context="content:
- content: no tag

attributes:
- p.(String): \"Internal\""
DEBU[0033] Response                                      effect=Permit obligations="attributes:
- p.(string): \"Internal First\"" reason="<nil>"
...
```
