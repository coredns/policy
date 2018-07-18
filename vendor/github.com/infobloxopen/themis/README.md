# Themis
Themis represents a set of tools for managing and enforcing security policies along with framework to create such tools:
- **pdp** - Policy Decision Point (core component of Themis);
- **pdpserver** - standalone application server which runs PDP;
- **proto**, **pdp-service**, **pdp-control** - gRPC protocol definitions and implementations;
- **pep** - golang client package for "service" protocol (Policy Enforcement Point or PEP);
- **pepcli** - CLI application which implements simple PEP and performance measurement tool for PDP server;
- **pdpctr-client** - golang client package for "control" protocol (Policy Administration Point or PAP);
- **papcli** - CLI application which implements simple PAP;
- **egen** - error processing code generator (development tool).

Themis design is inspired by eXtensible Access Control Markup Language (XACML) **[XACML-V3.0]**.

# Policy Decision Point
Policy Decision Point or PDP (according to **[XACML-V3.0]**) is an entity that evaluates applicable policy and renders an authorization decision. Themis provides PDP as a golang package.

## Policy Evaluation
To make a decision PDP evaluates policies it has on request **context**. A request **context** represents set of **attributes** together with local **content** (additional data which can be used by policies during the evaluation). Resulting decision consists of **effect**, **status** (**reason**) and set of **obligations**. Decision effect can be:
- **Deny** - request is denied;
- **Permit** - request is permitted;
- **Not Applicable** - no policy applicable to the request;
- **Indeterminate** - PDP can't evaluate particular effect;
- **IndeterminateD** - PDP can't evaluate effect but if it could it would be **Deny**;
- **IndeterminateP** - PDP can't evaluate effect but if it could it would be **Permit**;
- **IndeterminateDP** - PDP can't evaluate effect but if it could it would be only **Deny** or **Permit**;
In case of any **Indeterminate** effect **status** contains textual representation of an issue.

Some application may require more details for particular decision. For example if application can write log it may need a flag attached to decision which says when to do it. These details can be delivered as **obligations**. The **obligations** are set of attributes like in request context. Each attribute has name, type and value. Attribute name is arbitrary string (request context requires pair of attribute name and type to be unique). Following built-in types are defined:
- **boolean**;
- **string**;
- **integer** - signed 64-bit integer;
- **float** - 64-bit floating point number;
- **address** - IPv4 or IPv6 address;
- **network** - IPv4 or IPv6 network address;
- **domain** - domain name;
- **set of strings** - ordered set of strings;
- **set of domains** - set of domains (unordered);
- **set of networks** - set of IPv4 or IPv6 network addresses (unordered);
- **list of strings**.

**Boolean** value is accepted as "1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false", "False" and serialized to "true" and "false". **Integer** value is a decimal number in range [-9223372036854775808, 9223372036854775807]. **Float** value can be specified using decimal format (e.g. 3.1416) or scientific notation (e.g. 6.022E+23). **Address** accepted in dotted decimal ("192.0.2.1") form or in IPv6 ("2001:db8::68") form and serialized respectively. **Network** is accepted as a CIDR notation IP address and prefix (for example "192.0.2.0/24" or "2001:db8::/32"). **Domain** name is accepted as string of labels separated by dots which satisfies to RFC1035, 2181 and 4343 requirements. **Set of strings**, **set of domains**, **set of networks** and **list of strings** aren't accepted in request context but can appear in response's obligations as comma separated list of values.

User can define her custom type based on **flags** metatype. A value of the type can be any combination of listed flags. PDP allows to define up to 64 flags for a type. Values can't appear in request or returned as obligations.

## Policies
PDP uses YAML based language (YAML Abstract Syntax Tree or YAST) or JSON based language (JSON Abstract Syntax Tree or JAST) to define **policies** and specifically constructed JSON to define local **content**  (JSON Content or JCON). YAST can be converted to JAST (and vise versa) with any YAML to JSON convertor.

### Root
Any **policies** definition consists of policies (required), attributes (optional) and types (optional) sections. Policies section contains root **policy** or **policy set**. **Policy** holds rules under its "rules" field while **policy set** is able to contain both inner policies or policy sets under its "policies" field. For example:

**YAST**
```yaml
# All permit policy (without "attributes" section)
policies:
  alg: FirstApplicableEffect
  rules:
  - effect: Permit
```

**JAST**
```json
{
  "policies": {
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
```

Attributes section contains set of pairs attribute name and type:

**YAST**
```yaml
# Permit if x is "test" otherwise Not Applicable
attributes:
  x: string

policies:
  alg: FirstApplicableEffect
  target:
  - equal:
    - attr: x
    - val:
        type: string
        content: "test"
  rules:
  - effect: Permit
```

**JAST**
```json
{
  "attributes": {
    "x": "string"
  },
  "policies": {
    "alg": "FirstApplicableEffect",
    "target": [
      {
        "equal": [
          {
            "attr": "x"
          },
          {
            "val": {
              "type": "string",
              "content": "test"
            }
          }
        ]
      }
    ],
    "rules": [
      {
        "effect": "Permit"
      }
    ]
  }
}
```

Types section is designed for custom type definitions which now are limited to only **flags** metatype:

**YAST**
```yaml
types:
  colors:
    meta: flags
    flags:
    - red
    - green
    - blue

attributes:
  c: list of strings

policies:
  alg: FirstApplicableEffect
  rules:
  - effect: Permit
    obligations:
    - c:
        list of strings:
        - val:
            type: colors
            content:
            - red
            - blue
```

**JAST**
```json
{
  "types": {
    "colors": {
      "meta": "flags",
      "flags": ["red", "green", "blue"]
    }
  },
  "attributes": {
    "c": "list of strings"
  },
  "policies": {
    "alg": "FirstApplicableEffect",
    "rules": [
      {
        "effect": "Permit",
        "obligations": [
          {
            "c": {
              "list of strings": [
                {
                  "val": {
                    "type": "colors",
                    "content": ["red", "blue"]
                  }
                }
              ]
            }
          }
        ]
      }
    ]
  }
}
```

### Policy Set
Policy Set holds set of **policies** or inner **policy sets** and defines how to combine them. It has following fields:
- **id** - policy id (optional, if not defined policy is hidden);
- **target** - target expression which defines if policy set is applicable to request (optional, if not defined policy set is applicable to any request);
- **policies** - set of inner policies and policy sets;
- **alg** - policy combining algorithm (any of **FirstApplicableEffect**, **DenyOverrides** and **Mapper**);
- **obligations** - set of obligations (optional).

Example of policy set with all its fields (it contains one hidden policy set and one hidden policy):
```yaml
# Policy set with all its fields
...
id: "Test Policy Set"
target: # x == "test"
- equal:
  - attr: x
  - val:
      type: string
      content: "test"
alg: FirstApplicableEffect
policies:
- alg: FirstApplicableEffect
  target:
  - equal: # z == "example"
    - attr: z
    - val:
        type: string
        content: "example"
  policies:
  - alg: FirstApplicableEffect
    rules:
    - effect: Permit
- alg: FirstApplicableEffect
  rules:
  - effect: Deny
obligations:
- a:
   val:
     type: address
     content: "192.0.2.1"
```

### Policy
Policy stores set of rules and defines how to combine them. It has following fields:
- **id** - policy id (optional, if not defined policy is hidden);
- **target** - target expression which defines if policy is applicable to request (optional, if not defined policy is applicable to any request);
- **rules** - set of rules;
- **alg** - rule combining algorithm (the same as for policy set);
- **obligations** - set of obligations (optional).

Here is an example of policy with all fields defined (it contains one hidden rule):
```yaml
# Policy with all its fields
...
id: "Test Policy"
target: # x == "test"
- equal:
  - attr: x
  - val:
      type: string
      content: "test"
alg: FirstApplicableEffect
rules:
- effect: Permit
obligations:
- a:
   val:
     type: address
     content: "192.0.2.1"
```

### Rule
Rule defines decision effect. Possible fields of a rule:
- **id** - rule id (optional, if not defined policy is hidden);
- **target** - target expression (optional);
- **condition** - any boolean expression which together with target defines if rule is applicable to the request (optional, if not defined rule is applicable when target matches);
- **effect** - **Deny** or **Permit**;
- **obligations** - set of obligations.

For example a rule with all fields:
```yaml
# Rule with all its fields
...
id: "Test Rule"
target: # x == "test"
- equal:
  - attr: x
  - val:
      type: string
      content: "test"
condition: # not (c contains 192.0.2.1 or b)
  not:
  - or:
    - contains:
      - attr: c
      - val:
          type: address
          content: "192.0.2.1"
    - attr: b
effect: Permit
obligations:
- a:
   val:
     type: address
     content: "192.0.2.1"

```

### Target
Any particular policy set or policy or rule is applicable only if request matches its target. Target is a list of **any** expressions. **Any** expression is a list of **all** expressions and **all** expression is a list of match expression. Match expression is a boolean expression of two arguments. One of arguments should be a request attribute and other should be a immediate value. Only two functions can represent match expression **equal** and **contains**. If list of match expressions for particular **all** expression contains single element **all** keyword can be dropped. Similarly if list of **all** expressions for particular **any** expression consists of one element **any** keyword can be dropped.

Request matches target when all **any** expressions match (if one or more of **any** expression doesn't match, target also doesn't match). **Any** expression matches request if one or more of its **all** expressions match the request (if all **all** expressions don't match, **any** expression doesn't match as well). And similarly to target **all** expression matches if all its inner expressions match as well. If during target evaluation error occurs the policy set, policy or rule effect becomes **indeterminate** (if rule effect is permit it is **indeterminateP** if deny - **indeterminateD** for policy and policy set kind of **indeterminate** depends on combining algorithm (see below).

Below an example of policy with different kinds of targets:
```yaml
# Target examples
...
# ((x == test and c contains address(192.0.2.1)) or
#  x == example) and
# (network(192.0.2.0/28) contains a or network(192.0.2.16/28) contains a)
target:
- any:
  - all:
    - equal:
      - attr: x
      - val:
          type: string
          content: "test"
    - contains:
      - attr: c
      - val:
          type: address
          content: 192.0.2.1
  - equal:
    - attr: x
    - val:
        type: string
        content: "example"
- any:
  - contains:
    - val:
        type: network
        content: 192.0.2.0/28
    - attr: a
  - contains:
    - val:
        type: network
        content: 192.0.2.16/28
    - attr: a
...
# (x == test or x == example) and (network(192.0.2.0/28) contains a or network(192.0.2.16/28) contains a)
target:
- any:
  - equal:
    - attr: x
    - val:
        type: string
        content: "test"
  - equal:
    - attr: x
    - val:
        type: string
        content: "example"
- any:
  - contains:
    - val:
        type: network
        content: 192.0.2.0/28
    - attr: a
  - contains:
    - val:
        type: network
        content: 192.0.2.16/28
    - attr: a
...
# x == test and network(192.0.2.0/24) contains a
target:
- equal:
  - attr: x
  - val:
      type: string
      content: "test"
- contains:
  - val:
      type: network
      content: 192.0.2.0/24
  - attr: a
...
# x == test
target:
- equal:
  - attr: x
  - val:
      type: string
      content: "test"
```

### Condition
Condition is rule field which can be any boolean expression (for example see above "Rule with all its fields"). Following functions available to make such expression:
- **equal** - expects two arguments, where the result is true if the arguments are equal. The two arguments can be:
    - both strings
    - both integers, in which case, the comparison will be performed using integer arithmetic.
    - both floats, in which case, the comparison will be performed using floating point arithmetic.
    - a float and an integer, in which case, the integer is promoted to float before floating point arithmetic is applied.
- **greater** - expects two arguments, where the result is true if the first argument is greater than the second. The two arguments can be:
    - both integers, in which case, the comparison will be performed using integer arithmetic.
    - both floats, in which case, the comparison will be performed using floating point arithmetic.
    - a float and an integer, in which case, the integer is promoted to float before floating point arithmetic is applied.
- **contains**:
  - string contains substring - expects two string arguments first is a string to search in and second is a substring to search for;
  - network constains address;
  - set of strings contains string;
  - set of networks contains address;
  - set of domains contains doman;
- **not** - boolean not (expects boolean as its single argument);
- **and**, **or** - boolean and and or (expect booleans as its arguments (requires at least one).

In any expression attribute can be referred with **attr** keyword and immediate value with **val** keyword (see below). There is special **selector** expression which is described below.

### Immediate Value
Immediate value can be refered with **val** keyword and has fields:
- **type** - value type (any of available types);
- **content** - value data.
For **obligations** only data itself requires as type can be derived from attribute definition.

Example of policy with all possible values:
```yaml
# All values example
...
# String
val:
  type: string
  content: test
...
obligations:
- s:
   val:
     type: string
     content: example
...
# Address
val:
  type: address
  content: 192.0.2.1
...
obligations:
- a:
   val:
     type: address
     content: 192.0.2.2
...
# Network
val:
  type: network
  content: 192.0.2.0/28
...
obligations:
- c:
   val:
     type: network
     content: 192.0.2.16/28
...
# Domain
val:
  type: domain
  content: example.com
...
obligations:
- d:
   val:
     type: domain
     content: example.com
...
# Set of Strings
val:
  type: set of strings
  content:
  - first
  - second
...
obligations:
- ss:
    val:
      type: set of strings
      content:
      - first
      - second
...
# Set of Networks
val:
  type: set of networks
  content:
  - 192.0.2.0/28
  - 192.0.2.16/28
...
obligations:
- sn:
    val:
      type: set of networks
      content:
      - 192.0.2.0/28
      - 192.0.2.16/28
...
# Set of Domains
val:
  type: set of domains
  content:
  - example.com
  - example.net
...
obligations:
- sd:
    val:
      type: set of domains
      content:
      - example.com
      - example.net
...
# List of Strings
val:
  type: list of strings
  content:
  - first
  - second
...
obligations:
- ls:
    val:
      type: list of strings
      content:
      - first
      - second
...
# Custom Flags
types:
  colors:
    meta: flags
    flags: [red, green, blue]

...
val:
  type: colors
  content: [blue, green]
```

### Selector
Selector expression is an expression to access additionally supplied data. Selector uses **uri** field to locate source of such data. Currently only "local" URI schema is supported which defines local selector.

Local selector uses local content data (see below) and has following fields:
- **uri** - URI of local content ("local:&lt;content-id&gt;/&lt;content-item-id&gt;);
- **path** - defines path to data in local content (optional, if not set selector extracts immediate value from content item). Path represents a list of expressions. It should match to content item keys (see below). Selector calculates path expressions one by one and extracts value from next mapping step of content item until reaches desired value;
- **type** - type of data in local content (any of available types).

Example of local selector:
```yaml
# Selector example
...
selector:
  uri: "local:content/domain-addresses"
  path:
  - val:
      type: string
      content: good
  - attr: d
  type: set of networks
```

Content for the example:
```json
{
  "id": "content",
  "items": {
    "domain-addresses": {
      "keys": ["string", "domain"],
      "type": "set of networks",
      "data": {
        "good": {
          "example.com": ["192.0.2.16/28", "192.0.2.32/28"],
          "test.com": ["192.0.2.48/28", "192.0.2.64/28"]
        },
        "bad": {
          "example.com": ["2001:db8:1000::/40", "2001:db8:2000::/40"],
          "test.com": ["2001:db8:3000::/40", "2001:db8:4000::/40"]
        }
      }
    }
  }
}
```

When result of mapping is a flags value selector matches flag names by order of definition. Flags defined in policies YAST or JSON must have the same number of names as content have. For example if policy defines names "one", "two", "three" for some flags type and content defines names "red", "green" and "blue" for a map referred by selector URI, selector returns "one" for "red", "two" for "green" and "three" for "blue".

### Numerical Expression
A numerical expression is constructed using the following numerical functions:
- **add** - accepts two arguments, where the result is the sum of the arguments
- **subtract** - accepts two arguments, where the result is first argument subtracted by the second argument
- **multiply** - accepts two arguments, where the result is the product of the arguments
- **divide** - accepts two arguments, where the result is the first argument divided by the second argument
- **range** - accepts three arguments - the first two of which specify a range, and the third a value to be compared against the range. The arguments are:
  - **min** - the minimum value of the range
  - **max** - the maximum value of the range
  - **val** - the value to compare against the range

  The result of the **range** function is a string, which can be one of the following:
  - **Below** - if **val** is less than **min**
  - **Above** - if **val** is greater than **max**
  - **Within** - if it's not one of the above

Each of the above numerical functions performs the respective numerical operation using either integer or floating point arithmetic according to the following rule:
- if all the arguments are integers, the operation is performed using integer arithmetic and the result returned as an integer. (In the case of **range**, the result is a string)
- if all the arguments are floats, the operation is performed using floating point arithmetic and the result returned as a float. (In the case of **range**, the result is a string)
- if the arguments include a combination of integers and floats, the integers are first promoted to floats and then operation is performed using floating point arithmetic. The result is returned as a float. (In the case of **range**, the result is a string)

### Other functions
There are several other functions available:
- **list of strings** - converts its argument to list of strings. It accepts set of strings, list of strings and flags. In case of set of strings the function returns list of strings sorted in order maintained by set (set keeps order of initial value definition). List of strings returned by the function as is. And for flags it returns list of names for flags which are set (keeping order of names from flags type definition).
- **concat** - concatenates all given arguments to single list of strings. The function treats MissingValueError in special way. If at least one argument returns some data, any MissingValueError is ignored. But when all arguments return the error, **concat** returns the error as well. It accepts strings, lists of strings, sets of strings and flags as arguments. **concat** handles lists of strings, sets of strings and flags the same way as function **list of strings**.
- **try** - returns result of first expression which calculated with no error. If all arguments calculated with error it throws the last one. It accepts expressions of any types but all of them must be of the same type (which becomes type of function result).

### Local Content
Local content is a set of content **items** (see example above). It's identified by **id** field which can be any string with no slash character (`/`). Each content item also has id (key of "items" JSON object) and following fields:
- **keys** - list of types of nested maps (optional, if not present data should contain immediate value of type);
- **type** - any built-in type (or flags type defintion or name for domain map);
- **data** - list of nested maps with keys of mentioned types or immediate value of given type.

Local content supports string map (key type "string"), domain map (key type "domain") and network map (key type "network" or "address"). Selector expectes string expression as path item for string map, domain - for domain map and address or network - for network map ("address" expression is allowed even if content key is "network" and vice verse).

Any map supports also mapping to flags type. To create such map user needs to define flags in type field. Being defined a flags type can be used by name within the content and its updates. Type definition goes to **type** field of content item but it's represented by JSON object instead of string. The object should have following fields:
- **meta** - string "flags" (the only supported metatype for now);
- **name** - type name (can be used late instead of the definition);
- **flags** - list of flag names (up to 64).

For example:
```json
{
  "id": "content",
  "items": {
    "example-domains": {
      "keys": ["domain"],
      "type": {
        "meta": "flags",
        "name": "tags",
        "flags": ["red", "green", "blue"]
      },
      "data": {
        "example.com": ["red", "green"],
        "example.net": ["red", "blue"],
        "example.org": ["blue", "green"]
      }
    },
    "test-domains": {
      "keys": ["domain"],
      "type": "tags",
      "data": {
        "test.com": ["red", "green"],
        "test.net": ["red", "blue"],
        "test.org": ["blue", "green"]
      }
    }
  }
}
```

### Policy and Rule Combining Algorithms
Policy and rule combinig algorithms define how to use child policies or rules of given policy set or policy and how to combine their effects, statuses and obligations. Themis supports following algorithms:
- **FirstApplicableEffect** - evaluates child policies or rules one by one until meets any other than **NotApplicable** effect (see details below);
- **DenyOverrides** - evaluates child policies or rules one by one until meets **Deny** effect;
- **Mapper** - evaluates map expression and uses result to find child policy or rule to evaluate.

For any algorithm if effect of children evaluation is **Deny** or **Permit** policy or policy set adds its obligation to what it got from children.

#### First Applicable Effect
The algorithm iterates child policies or rules and evaluates them one by one. User should not relay on any particular order (however for now it goes sequentionaly from first to the last). If any child evaluation effect is not **NotApplicable** the effect becomes overall policy or policy set result.

#### DenyOverrides
The algorithm as well evaluates child policies or rules one by one. User should not relay on any particular order (however for now it goes sequentionaly from first to the last). If any effect is **Deny** the effect becomes overall policy or policy set result and any other evaluation results are dropped. Other effects are combined as following:

| Effects | Result |
| --- | --- |
| at least one **IndeterminateDP** or at least one **IndeterminateD** with at least one **Permit** or at least one **IndeterminateP** and any **NotApplicable** | **IndeterminateDP** |
| at least one **IndeterminateD** and any **NotApplicable** | **IndeterminateD** |
| at least one **Permit** and any **IndeterminateP** or **NotApplicable** | **Permit** |
| at least one **IndeterminateP** and any **NotApplicable** | **IndeterminateP** |
| only **NotApplicable** | **NotApplicable** |

In case of any **Indeterminate** result all statuses are combined together.

#### Mapper
The algorithm is capable to select particular child policy or rule with no evaluation other children one by one. It has some parameters:
- **id** - always "mapper" for the algorithm;
- **map** - expression to get resulting policy or rule id (it can be string, set of strings or list of strings expression);
- **default** - policy or rule id to evaluate if result of map expression doesn't match any child id (optional, if absent mapper policy effect is **Indeterminate**);
- **error** - policy or rule id to evaluate if map expression can't be evaluated (optional, if absent mapper policy effect is **Indeterminate**);
- **alg** - nested algorithm to use if map expression result is set of strings or list of strings (required if map is set of strings or list of strings expression, otherwise ignored). Other mapper can be used here (but for it **default** and **error** fields are ignored);
- **order** - chooses order in which policies or rules selected by map expression are passed to nested algorithm. The order can be "External" (default) and "Internal". With external order policies or rules follow order of ids returned by map expression. Internal order stands for order of the policies or rules as they appear in parent policy set or policy.

Any hidden child policy or rule is ignored by mapper algorithm.

For example policies:
```yaml
# Mapper example
...
alg:
  id: Mapper
  map:
    attr: p
  default: DenyPolicy
  error: ErrorPolicy
...
alg:
  id: Mapper
  map:
    selector:
      uri: "local:content/domain-policies"
      path:
      - attr: d
      type: list of strings
  default: DenyRule
  alg: FirstApplicableEffect
  order: Internal
```

and content for them:
```json
{
  "id": "content",
  "items": {
    "domain-policies": {
      "keys": ["domain"],
      "type": "list of strings",
      "data": {
        "example.com": ["PermitCom", "DenyCom"],
        "example.net": ["PermitNet", "DenyNet"]
      }
    }
  }
}
```

If result type of **map** is a flags type its flag names treated as id of policy to run. If flags value has several flags set they are ordered according of order in type definiton and passed to nested combining algorithm.

# PDPServer
PDP server allows to run and control PDP. Additionally the server provides endpoint for healthcheck and supports OpenZipkin tracing. Started with no options pdpservers gets no initial policies and content. Policies and content in the case should be provided by control interface. Option `-p` provides initial policy for the server from given YAML file. Option `-j` provides content (it can be specified several times). For example (`-v 3` sets maximal log level):
```
$ pdpserver -v 3 -p policy.yaml -j mapper.json -j content.json
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=policy.yaml
INFO[0000] Parsing policy                                policy=policy.yaml
INFO[0000] Opening content                               content=mapper.json
INFO[0000] Parsing content                               content=mapper.json
INFO[0000] Opening content                               content=content.json
INFO[0000] Parsing content                               content=content.json
INFO[0000] Opening service port                          address="0.0.0.0:5555"
INFO[0000] Opening control port                          address="0.0.0.0:5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
```
Other pdpserver options:
- `-c` - listen for policies on given address:port (default "0.0.0.0:5554");
- `-health` - health check endpoint;
- `-l` - listen for decision requests on given address:port (default "0.0.0.0:5555");
- `-pprof` - performance profiler endpoint (see go tool pprof);
- `-t` - OpenZipkin tracing endpoint;
- `-v` - log verbosity (0 - error, 1 - warn (default), 2 - info, 3 - debug).

## Requests
To make decision requests there are 3 options. Create client from scratch which implements protocol defined by `proto/service.proto`, use golang client package `themis\pep` to implement client application (see `contrib/coredns/policy`) and for debug use simple PEPCLI client. To use PEPCLI client create requests YAML file for example:
```yaml
attributes:
  s: string
  a: address

requests:
- s: Local Test
  a: 127.0.0.1

- s: Example
  a: 192.0.2.1
```
start PDP server with some policy (use for example "All permit policy" above):
```
$ pdpserver -v 3 -p all-permit-policy.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=all-permit-policy.yaml
INFO[0000] Parsing policy                                policy=all-permit-policy.yaml
INFO[0000] Opening service port                          address="0.0.0.0:5555"
INFO[0000] Opening control port                          address="0.0.0.0:5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
```
in the other terminal run PEPCLI:
```
$ pepcli -i requests.yaml test
Got 2 requests. Sending...
- effect: PERMIT
  reason: "Ok"

- effect: PERMIT
  reason: "Ok"
```
PEPCLI sends two requests which are listed in the file ang gets all permitted as instructed by the policy. In PDP's terminal you can see respective logs:
```
...
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
INFO[0089] Validating context
DEBU[0089] Request context                               context=attributes:
- a.(Address): 127.0.0.1
- s.(String): "Local Test"
INFO[0089] Returning response
DEBU[0089] Response                                      effect=PERMIT obligation=no attributes reason=Ok
INFO[0089] Validating context
DEBU[0089] Request context                               context=attributes:
- s.(String): "Example"
- a.(Address): 192.0.2.1
INFO[0089] Returning response
DEBU[0089] Response                                      effect=PERMIT obligation=no attributes reason=Ok
...
```

## Policies and content uploading and updating
PDP Server accepts control requests to upload and update policies or content. Themis user can implement her own client from scratch using protocol definition from `proto/control.proto` or using golang package `themis/pdpctrl-client`. To make control requests for debug purpose Themis provides PAPCLI tool.

Update is a list of commands each contains three fields:
- **op** - add or delete;
- **path** - list of ids;
- **entity** - contains an entity to add as child to entity defined by path.

PDP supporst **add** and **delete** commands. In case of **add** path should point to parent item and **entity** should contain appropriate child. For example if it is policy update and **path** points to policy set **entity** can be policy or other policy set. If **path** points to policy then only rule can be accepted as **entity**.

For example in one terminal start PDP server with no policies:
```
$ pdpserver -v 3
INFO[0000] Starting PDP server
INFO[0000] Opening service port                          address="0.0.0.0:5555"
INFO[0000] Opening control port                          address="0.0.0.0:5554"
INFO[0000] Creating service protocol handler
INFO[0000] Creating control protocol handler
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
```
Then upload policy with PAPCLI ("All permit policy"):
```
$ papcli -s 127.0.0.1:5554 -p all-permit-policy.yaml
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```
PDP got the data:
```
...
INFO[0000] Serving decision requests
INFO[0000] Serving control requests
INFO[0004] Got new control request
INFO[0004] Got new data stream
INFO[0004] Got apply command
INFO[0004] New policy has been applied                   id=1
...
```

PDP Server doesn't accept updates to policy with no tag so upload other policy and set tag to it to make update later (here policy similar to "permit if x is test" is used but with ids (because commands add and delete doen't see hidden policies or rules):
```yaml
# Permit if x is "test" otherwise Not Applicable
attributes:
  x: string

policies:
  id: Root
  alg: FirstApplicableEffect
  target:
  - equal:
    - attr: x
    - val:
        type: string
        content: "test"
  rules:
  - id: First Rule
    effect: Permit
```
Run PAPCLI with the policy and initial tag (option `-vt` the tag should be correct UUID):
```
$ papcli -s 127.0.0.1:5554 -p permit-test-x-policy.yaml -vt 823f79f2-0001-4eb2-9ba0-2a8c1b284443
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```
PDP Server accepts the policy:
```
...
INFO[0016] Got new control request
INFO[0016] Got new data stream
INFO[0016] Got apply command
INFO[0016] New policy has been applied                   id=1 tag=823f79f2-0001-4eb2-9ba0-2a8c1b284443
...
```

Then policy can be updated (with following update which removes "First Rule" and adds other one):
```yaml
- op: add
  path:
  - Root
  entity:
    id: Permit Rule With Obligation
    effect: Permit
    obligations:
    - x: example

- op: delete
  path:
  - Root
  - First Rule
```
Run PAPCLI with the update (you need to specify previous tag with option `-vf` and new tag with option `-vt`, when both options present PDP server considers data as update and checks if `-vf` tag matches to tag current tag of updated to maintain update consistency):
```
$ papcli -s 127.0.0.1:5554 -p permit-test-x-policy-update.yaml -vf 823f79f2-0001-4eb2-9ba0-2a8c1b284443 -vt 93a17ce2-788d-476f-bd11-a5580a2f35f3
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```
PDP accepts the update:
```
...
INFO[0373] Got new control request
INFO[0373] Got new data stream
DEBU[0373] Policy update                                 update=policy update: 823f79f2-0001-4eb2-9ba0-2a8c1b284443 - 93a17ce2-788d-476f-bd11-a5580a2f35f3
commands:
- Add Path ("Root")
- Delete Path ("Root"/"First Rule")
INFO[0373] Got apply command
INFO[0373] Policy update has been applied                curr-tag=93a17ce2-788d-476f-bd11-a5580a2f35f3 id=3 prev-tag=823f79f2-0001-4eb2-9ba0-2a8c1b284443
...
```
Consider content update. For example use "Selector example" policy. Start PDP with the policy:
```
$ pdpserver -v 3 -p selector-examle.yaml
INFO[0000] Starting PDP server
INFO[0000] Loading policy                                policy=selector-examle.yaml
INFO[0000] Parsing policy                                policy=selector-examle.yaml
INFO[0000] Opening service port                          address="0.0.0.0:5555"
INFO[0000] Opening control port                          address="0.0.0.0:5554"
INFO[0000] Creating service protocol handler
INFO[0000] Serving decision requests
INFO[0000] Creating control protocol handler
INFO[0000] Serving control requests
```
Then upload content with some tag (to be able to update it):
```
$ papcli -s 127.0.0.1:5554 -j selector-examle.json -vt 823f79f2-0001-4eb2-9ba0-2a8c1b284443
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```
PDP server accepts upload:
```
...
INFO[0265] Got new control request
INFO[0265] Got new data stream
INFO[0265] Got apply command
INFO[0265] New content has been applied                  id=1 tag=823f79f2-0001-4eb2-9ba0-2a8c1b284443
...
```
Now lets move IPv4 addresses from "good" to "bad" map and IPv6 from
"bad" to "good" for "example.com":
```json
[
  {
    "op": "delete",
    "path": ["domain-addresses", "good", "example.com"]
  },
  {
    "op": "add",
    "path": ["domain-addresses", "good", "example.com"],
    "entity": {
      "type": "set of networks",
      "data": ["2001:db8:1000::/40", "2001:db8:2000::/40"]
    }
  },
  {
    "op": "delete",
    "path": ["domain-addresses", "bad", "example.com"]
  },
  {
    "op": "add",
    "path": ["domain-addresses", "bad", "example.com"],
    "entity": {
      "type": "set of networks",
      "data": ["192.0.2.16/28", "192.0.2.32/28"]
    }
  }
]
```
Note that update's entities doesn't contain **keys** field as **data** is immediate value (and has no any mappings). If update add some entity with mapping entity should have a **keys** filed. For example:
```json
[
  {
    "op": "delete",
    "path": ["domain-addresses", "good"]
  },
  {
    "op": "add",
    "path": ["domain-addresses", "good"],
    "entity": {
      "type": "set of networks",
      "keys": ["domain"],
      "data": {
        "example.com": ["2001:db8:1000::/40", "2001:db8:2000::/40"],
        "test.com": ["2001:db8:3000::/40", "2001:db8:4000::/40"]
      }
    }
  }
]
```

Run PAPCLI with content update file:
```
$ papcli -s 127.0.0.1:5554 -id content -j selector-examle-update.json -vf 823f79f2-0001-4eb2-9ba0-2a8c1b284443 -vt 93a17ce2-788d-476f-bd11-a5580a2f35f3
INFO[0000] Requesting data upload to PDP servers...
INFO[0000] Uploading data to PDP servers...
```

Check PDP logs:
```
...
INFO[2190] Got new control request
INFO[2190] Got new data stream
DEBU[2190] Content update                                update=content update: 823f79f2-0001-4eb2-9ba0-2a8c1b284443 - 93a17ce2-788d-476f-bd11-a5580a2f35f3
content: "content"
commands:
- Delete Path ("domain-addresses"/"good"/"example.com")
- Add Path ("domain-addresses"/"good"/"example.com")
- Delete Path ("domain-addresses"/"bad"/"example.com")
- Add Path ("domain-addresses"/"bad"/"example.com")
INFO[2190] Got apply command
INFO[2190] Content update has been applied               cid=content curr-tag=93a17ce2-788d-476f-bd11-a5580a2f35f3 id=5
...
```

Contents with different ids and policies can be updated independently and in parallel.

# References
**[XACML-V3.0]** *eXtensible Access Control Markup Language (XACML) Version 3.0.* 22 January 2013. OASIS Standard. http://docs.oasis-open.org/xacml/3.0/xacml-3.0-core-spec-os-en.html.

