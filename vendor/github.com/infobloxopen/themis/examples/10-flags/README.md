# 10-Flags

The example shows how to define custom flags type, use it with selector and update content containing the flags.

## Policy with flags type

Run pdpserver using policy file:
```
$ pdpserver -v 3 -p flags.yaml -j content.json
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=flags.yaml
INFO[0000] Parsing policy                                policy=flags.yaml
INFO[0000] Opening content                               content=content.json
INFO[0000] Parsing content                               content=content.json
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
$ pepcli -i flags.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "red"

    - id: "t"
      type: "list of strings"
      value: "\"red\",\"yellow\",\"indigo\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "orange"

    - id: "t"
      type: "list of strings"
      value: "\"orange\",\"green\",\"blue\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "yellow"

    - id: "t"
      type: "list of strings"
      value: "\"yellow\",\"green\",\"violet\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "green"

    - id: "t"
      type: "list of strings"
      value: "\"green\",\"blue\",\"violet\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "blue"

    - id: "t"
      type: "list of strings"
      value: "\"blue\",\"indigo\",\"violet\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "indigo"

    - id: "t"
      type: "list of strings"
      value: "\"indigo\",\"violet\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "violet"

    - id: "t"
      type: "list of strings"
      value: "\"violet\""

- effect: Deny

```

PDP logs:
```
...
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.red)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"red\"
- t.(list of strings): \"\\\"red\\\",\\\"yellow\\\",\\\"indigo\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.orange)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"orange\"
- t.(list of strings): \"\\\"orange\\\",\\\"green\\\",\\\"blue\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.yellow)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"yellow\"
- t.(list of strings): \"\\\"yellow\\\",\\\"green\\\",\\\"violet\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.green)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"green\"
- t.(list of strings): \"\\\"green\\\",\\\"blue\\\",\\\"violet\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.blue)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"blue\"
- t.(list of strings): \"\\\"blue\\\",\\\"indigo\\\",\\\"violet\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.indigo)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"indigo\"
- t.(list of strings): \"\\\"indigo\\\",\\\"violet\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.violet)"
DEBU[0005] Response                                      effect=Permit obligations="attributes:
- r.(string): \"violet\"
- t.(list of strings): \"\\\"violet\\\"\"" reason="<nil>"
DEBU[0005] Request context                               context="content:
- content: no tag

attributes:
- d.(Domain): domain(example.black)"
DEBU[0005] Response                                      effect=Deny obligations="no attributes" reason="<nil>"
...
```

## Update content with flags values

Run pdpserver using policy file but with no content for now:
```
$ pdpserver -v 3 -p flags.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=flags.yaml
INFO[0000] Parsing policy                                policy=flags.yaml
INFO[0000] Opening control port                          address=":5554"
INFO[0000] Opening storage port                          address=":5552"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
INFO[0000] Opening service port                          address=":5555"
INFO[0000] Serving decision requests
```

Any request to the server now returns an error:
```
$ pepcli -i update-flags.requests.yaml test
- effect: Indeterminate
  reason: "#a2: #a3: Failed to process request: #2e (hidden policy>mapper): Missing content content"

- effect: Indeterminate
  reason: "#a2: #a3: Failed to process request: #2e (hidden policy>mapper): Missing content content"

```

PDP logs:
```
DEBU[0141] Request context                               context="attributes:
- d.(Domain): domain(example.red)"
DEBU[0141] Response                                      effect=Indeterminate obligations="no attributes" reason="#2e (hidden policy>mapper): Missing content content"
DEBU[0141] Request context                               context="attributes:
- d.(Domain): domain(test.red)"
DEBU[0141] Response                                      effect=Indeterminate obligations="no attributes" reason="#2e (hidden policy>mapper): Missing content content"
```

Upload content using papcli:
```
$ papcli -s localhost:5554 -j content.json -vt 823f79f2-0001-4eb2-9ba0-2a8c1b284443
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```

We use "823f79f2-0001-4eb2-9ba0-2a8c1b284443" to tag the content so we can update it later using the string as reference to content version we want to update. The tag can be any valid UUID.

PDP logs:
```
INFO[0313] Got new control request
INFO[0313] Got new data stream
INFO[0313] Got apply command
INFO[0313] New content has been applied                  id=1 tag=823f79f2-0001-4eb2-9ba0-2a8c1b284443
INFO[0313] Got notified about readiness
```

With the content PDP returns some data:
```
$ pepcli -i update-flags.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "red"

    - id: "t"
      type: "list of strings"
      value: "\"red\",\"yellow\",\"indigo\""

- effect: Deny

```

PDP logs:
```
DEBU[0355] Request context                               context="content:
- content: 823f79f2-0001-4eb2-9ba0-2a8c1b284443

attributes:
- d.(Domain): domain(example.red)"
DEBU[0355] Response                                      effect=Permit obligations="attributes:
- r.(string): \"red\"
- t.(list of strings): \"\\\"red\\\",\\\"yellow\\\",\\\"indigo\\\"\"" reason="<nil>"
DEBU[0355] Request context                               context="content:
- content: 823f79f2-0001-4eb2-9ba0-2a8c1b284443

attributes:
- d.(Domain): domain(test.red)"
DEBU[0355] Response                                      effect=Deny obligations="no attributes" reason="<nil>"
```

Let's put some tags for "test.red" and remove some from "example.red":
```
$ papcli -s localhost:5554 -j content-update.json -id content -vf 823f79f2-0001-4eb2-9ba0-2a8c1b284443 -vt 93a17ce2-788d-476f-bd11-a5580a2f35f3
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```

Note that now you don't need to define flags type in update again. Just use it's name defined when content has been uploaded at first time. The update also changes content tag so you can be sure that update goes to right version of content.

PDP logs:
```
INFO[0897] Got new control request
INFO[0897] Got new data stream
DEBU[0897] Content update                                update="content update: 823f79f2-0001-4eb2-9ba0-2a8c1b284443 - 93a17ce2-788d-476f-bd11-a5580a2f35f3
content: \"content\"
commands:
- Delete path (\"domain\"/\"example.red\")
- Add path (\"domain\"/\"example.red\")
- Add path (\"domain\"/\"test.red\")"
INFO[0897] Got apply command
INFO[0897] Content update has been applied               cid=content curr-tag=93a17ce2-788d-476f-bd11-a5580a2f35f3 id=3 prev-tag=823f79f2-0001-4eb2-9ba0-2a8c1b284443
INFO[0897] Got notified about readiness
```

Now PDP respose changed:
```
$ pepcli -i update-flags.requests.yaml test
- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "yellow"

    - id: "t"
      type: "list of strings"
      value: "\"yellow\",\"indigo\""

- effect: Permit
  obligation:
    - id: "r"
      type: "string"
      value: "red"

    - id: "t"
      type: "list of strings"
      value: "\"red\""

```

PDP logs:
```
DEBU[1132] Request context                               context="content:
- content: 93a17ce2-788d-476f-bd11-a5580a2f35f3

attributes:
- d.(Domain): domain(example.red)"
DEBU[1132] Response                                      effect=Permit obligations="attributes:
- r.(string): \"yellow\"
- t.(list of strings): \"\\\"yellow\\\",\\\"indigo\\\"\"" reason="<nil>"
DEBU[1132] Request context                               context="content:
- content: 93a17ce2-788d-476f-bd11-a5580a2f35f3

attributes:
- d.(Domain): domain(test.red)"
DEBU[1132] Response                                      effect=Permit obligations="attributes:
- r.(string): \"red\"
- t.(list of strings): \"\\\"red\\\"\"" reason="<nil>"
```
