# firewall

## Name

*firewall* - enables filtering of query and response based on expressions.

## Description

The firewall plugin defines a list of rules that trigger workflow action on the DNS query or its response.
A rule list is an ordered set of rules that are evaluated in sequence.
Rules can be an expression rule, or a policy engine rule. 
An expression rule has two parts: an action and an expression. When the rule is evaluated,
first the expression is evaluated.
- If the expression evaluates to `true` the action is performed on the query and the rule list evaluation ceases.
- If the expression does not evaluates to `true` then next rule in sequence is evaluated.

The firewall plugin can also refer to other policy engines to determine the action to take.

## Syntax

~~~ txt
firewall DIRECTION {
    ACTION EXPRESSION
    POLICY-PLUGIN ENGINE-NAME
}
~~~~

* **DIRECTION** indicates if the _rule list_ will be applied to queries or responses. It can be `query` or `response`.

* **ACTION** defines the workflow action to take if the **EXPRESSION** evaluates to `true`.
Available actions:
  - `allow` : continue the DNS resolution process (i.e. continue to the next plugin in the chain)
  - `refuse` : interrupt the DNS resolution, reply with REFUSE code
  - `block` : interrupt the DNS resolution, reply with NXDOMAIN code
  - `drop` : interrupt the DNS resolution, do not reply (client will time out)

  An action must be followed by an **EXPRESSION**, which defines the boolean expression for the rule.  Expression uses 
  a [c-like expression format](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) where the variables are either
  the `metadata` of CoreDNS or the fields of a DNS query/response.  Common operators apply.

  Expression Examples:
  * `client_ip == '10.0.0.20'`
  * `type == 'AAAA'`
  * `type IN ('AAAA', 'A', 'TXT')`
  * `type IN ('AAAA', 'A') && name =~ 'google.com'`
  * `[mac/address] =~ '.*:FF:.*'`

  NOTE: because of the `/` separator included in a label of metadata, those labels must be enclosed in
  brackets `[...]`.

  The following names are supported for querying information in expressions

  * `type`: type of the request (A, AAAA, TXT, ..)
  * `name`: name of the request (the domain requested)
  * `class`: class of the request (IN, CS, CH, ...)
  * `proto`: protocol used (tcp or udp)
  * `client_ip`: client's IP address, for IPv6 addresses these are enclosed in brackets: `[::1]`
  * `size`: request size in bytes
  * `port`: client's port
  * `rcode`: response CODE (NOERROR, NXDOMAIN, SERVFAIL, ...)
  * `rsize`: raw (uncompressed), response size (a client may receive a smaller response)
  * `>rflags`: response flags, each set flag will be displayed, e.g. "aa, tc". This includes the qr
    bit as well
  * `>bufsize`: the EDNS0 buffer size advertised in the query
  * `>do`: is the EDNS0 DO (DNSSEC OK) bit set in the query
  * `>id`: query ID
  * `>opcode`: query OPCODE
  * `server_ip`: server's IP address, for IPv6 addresses these are enclosed in brackets: `[::1]`
  * `server_port` : client's port
  * `response_ip` : the IP returned in the first A or AAAA record of the Answer section

* **POLICY-PLUGIN** : is the name of another plugin that implements a firewall policy engine. 
  **ENGINE-NAME** is the name of an engine defined in your Corefile. Requests/responses will be evaluated by
  that plugin policy engine to determine the action.

## Policy Engine Plugins

In addition to using the built-in action/expression syntax, the _firewall_ plugin can use a policy engine plugin
to evaluate policy.

To use a policy engine plugin, you'll need to compile plugin into CoreDNS, the declare the plugin in in your
Corefile, and reference the plugin as an action of a firewall rule.  See the "Using a Policy Engine Plugin" example below.

When authoring a new policy engine plugin, the plugin must implement the `Engineer` interface defined in firewall/policy.

This reposiory includes two Policy Engine Plugins:
* *themis* - enables Infoblox's Themis policy engine yo be used as CoreDNS firewall policy engine
* *opa* - enables OPA to be used as a CoreDNS firewall policy engine.

## External Plugin

*Firewall* and other associated policy plugins in this repository are *external* plugins, which means it they are not included in CoreDNS releases.
To use the plugins in this repository, you'll need to build a CoreDNS image with the plugins you want to add included in the plugin list. In a nutshell you'll need to:
* Clone <https://github.com/coredns/coredns>
* Add this plugin to [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg) per instructions therein.  Order in this file is important, it dictates the order of plugin execution when processing a query.  The firewall plugin should execute before plugins that generate responses.
* `make -f Makefile.release DOCKER=your-docker-repo release`
* `make -f Makefile.release DOCKER=your-docker-repo docker`
* `make -f Makefile.release DOCKER=your-docker-repo docker-push`

## Examples

### Client IP ACL
This example shows how to use *firewall* to create a basic client IP based ACL. Here `10.120.1.11` is expressly REFUSED.
Other clients in `10.120.0.0/24` and `10.120.1.0/24` are allowed.  All other clients are not responded to.

~~~ corefile
. {
   firewall query {
      refuse client_ip == '10.120.1.11'
      allow client_ip =~ '10\.120\.0\..*'
      allow client_ip =~ '10\.120\.1\..*'
      drop true
   }
}
~~~

### Domain Name Policy
Allow all queries for example.com.
Allow A and AAAA queries for google.com.
NXDOMAIN all other queries.

~~~ corefile
. {
   firewall query {
      allow name =~ 'example.com'
      allow name =~ 'google.com' && (type == 'A' || type == 'AAAA')
      block true
   }
}
~~~

### EDNS0 Metadata Policy
This example uses the *metadata_edns0* plugin to define labels `group_id` and `client_id` with values extracted from EDNS0.
The firewall rules use those metadata to REFUSE any query without a group_id of `123456789` or client_id of `ABCDEF`.

~~~ corefile
example.org {
   metadata
   metadata_edns0 {
      group_id edns0 0xffed bytes
      client_id edns0 0xffee bytes
   }
   firewall query {
      refuse [metadata_edns0/client_id] != 'ABCDEF'
      refuse [metadata_edns0/group_id] != '123456789'
      allow true
   }
}
~~~

### Kubernetes Metadata Multi-Tenancy Policy
This example shows how *firewall* could be useful in a Kubernetes multi-tenancy application. It uses the *kubernetes*
plugin metadata to prevent Pods in certain namespaces from looking up Services in other namespaces.
Specifically, if the requesting Pod is in a namespace beginning with 'tenant-', it permits only lookups to
Services that are in the same namespace or in the 'default' namespace. Note here that `pods verified` is
required in *kubernetes* plugin to enable the `[kubernetes/client-namespace]` metadata variable.

~~~ corefile
cluster.local {
   metadata
   kubernetes {
      pods verified
   }
   firewall query {
      allow [kubernetes/client-namespace] !~ '^tenant-'
      allow [kubernetes/namespace] == [kubernetes/client-namespace]
      allow [kubernetes/namespace] == 'default'
      block true
   }
}
~~~

### Using a Policy Engine Plugin

The following example illustrates how the a policy engine plugin (*themis* in this example) can be used by the *firewall* plugin.
Note that the *themis* plugin options are not defined here, and are replaced by `...`.

~~~
. {
  firewall query {
    themis myengine
  }
   
  themis myengine {
    ...
  }
}
~~~
