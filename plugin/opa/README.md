# opa

*opa* - enables OPA to be used as a CoreDNS _firewall_ policy engine.

## Syntax

```
opa ENGINE-NAME {
    endpoint URL
    tls CERT KEY CACERT
    fields FIELD [FIELD...]
}
```

* **ENGINE-NAME** is the name of the policy engine, used by the firewall
  plugin to uniquely identify the instance. Each instance of _opa_ in
  the Corefile must have a unique **ENGINE-NAME**.

* `endpoint` defines the OPA endpoint **URL**.  It should include the
  full path to the rule.

* `tls` **CERT** **KEY** **CACERT** are the TLS cert, key and the CA
  cert file names for the OPA connection.
   
* `fields` lists the fields that will be sent to OPA when evaluating the
  policy for a DNS request/response. Fields available are the same as in
  *firewall* plugin expresions: *metadata* from other plugins, and data
  from the request/response ("type", "name", "proto", "client_ip", etc).
  See the *firewall* README for a list. If this option is omitted, the
  following fields are sent: "client_ip", "name", "rcode", "response_ip"


## Firewall Policy Engine

This plugin is not a standalone plugin.  It must be used in conjunction
with the _firewall_ plugin to function. For this plugin to be active,
the _firewall_ plugin must reference it in a rule.  See the "Policy
Engine Plugins" section of the _firewall_ plugin README for more
information.

## Writing the OPA Policy

This plugin assumes that the rule referenced in the `endpoint` URL will
evaluate to one of following values:
* "allow" - allows the dns request/response to proceed as normal
* "refuse" - sends a REFUSED response to the client
* "block" - sends a NXDOMAIN response to the client
* "drop" - sends no response to the client

When writing a rules in OPA, all `fields` are available as input.

## Examples

Point to a local OPA instance using a rule named `action` in the `dns`
package.

~~~ txt
. {
  opa myengine {
        endpoint http://127.0.0.1:8181/v1/data/dns/action
  }

  firewall query {
    opa myengine
  }
}
~~~

The following is an example OPA package. It implements a simple name
blacklist, while whitelisting a client subnet. The rule "action" allows
any request from clients with an IP address in the `1.2.3.0/24` network.
For all other clients it  blocks `a.example.com.`, refuses
`b.example.com`, and drops requests for `x.example.com`. It allows all
other requests.

~~~ rego
package dns
  
import input.name
import input.client_ip

default action = "allow"

# Action Priority
action = "allow" {
  allow
} else = "refuse" {
  refuse
} else = "block" {
  block
} else = "drop" {
  drop
}

block { name == "a.example.com." }

refuse { name == "b.example.com." }

drop { name == "x.example.com." }

allow { net.cidr_contains("1.2.3.0/24", client_ip) }
~~~