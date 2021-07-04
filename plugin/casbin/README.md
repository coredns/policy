# Casbin

*casbin* - enables Casbin to be used as a CoreDNS *firewall* policy engine.

## Syntax

```
casbin ENGINE-NAME {
    model /path/to/model
    policy /path/to/policy
}
```

* **ENGINE-NAME** is the name of the policy engine, used by the firewall
  plugin to uniquely identify the instance. Each instance of _opa_ in
  the Corefile must have a unique **ENGINE-NAME**.

* **model** & **policy** are concepts in casbin. More details, please refer to [casbin](https://casbin.org/docs/en/how-it-works)

## Firewall Policy Engine

This plugin is not a standalone plugin.  It must be used in conjunction
with the _firewall_ plugin to function. For this plugin to be active,
the _firewall_ plugin must reference it in a rule.  See the "Policy
Engine Plugins" section of the _firewall_ plugin README for more
information.

## Examples

`myengine` points to a Casbin instance.

```txt
. {
  casbin myengine {
      model ./examples/model.conf
      policy ./examples/policy.csv
  }

  firewall query {
      casbin myengine
  }
}
```

model:

```conf
[request_definition]
r = client_ip, name

[policy_definition]
p = client_ip, name, action

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.client_ip == p.client_ip && r.name == p.name
```

policy:

```csv
p, 10.240.0.1, example.org., allow
```