# firewall

## Name

*firewall* - enables filtering of query and response based on expressions.

## Description

The firewall plugin defines a rule list of expressions that trigger workflow action on the DNS query or its response.
A rule list is an ordered set of rules that are evaluated in sequence.
Each rule has two parts: an action and an expression. When the rule is evaluated,
first the expression is evaluated.
- If the expression evaluates to `true` the action is performed on the query and the rule list evaluation ceases.
- If the expression does not evaluates to `true` then next rule in sequence is evaluated.


## Syntax

~~~ txt
firewall DIRECTION {
    ACTION EXPRESSION
    ACTION EXPRESSION
    ...
}
~~~~

* **DIRECTION** indicates if the _rule list_ will be applied to queries or responses. It can be `query` or `response`.

* **ACTION** defines the workflow action to take if the **EXPRESSION** evaluates to `true`.
Available actions:
  - `allow` : continue the DNS resolution process (i.e. continue to the next plugin in the chain)
  - `refuse` : interrupt the DNS resolution, reply with REFUSE code
  - `block` : interrupt the DNS resolution, reply with NXDOMAIN code
  - `drop` : interrupt the DNS resolution, do not reply (client will time out)

* **EXPRESSION** defines the boolen expression for the rule.  Expression is a [go-like language](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) where the variables are either the `metadata` of CoreDNS
either a list of names associated with the DNS query/response information.
Usual operators applies.

  Examples:
  * `client_ip == '10.0.0.20'`
  * `type == 'AAAA'`
  * `type IN ('AAAA', 'A', 'TXT')`
  * `type IN ('AAAA', 'A') && name =~ 'google.com'`
  * `[mac/address] =~ '.*:FF:.*'`

  NOTE: because of the `/` separator included in a label of metadata, those labels must be enclosed on 
  bracket [...] for acorrect evaluation by the expression engine

  The following names are supported for querying information

  * `type`: type of the request (A, AAAA, TXT, ..)
  * `name`: name of the request (the domain requested)
  * `class`: class of the request (IN, CS, CH, ...)
  * `proto`: protocol used (tcp or udp)
  * `remote`: client's IP address, for IPv6 addresses these are enclosed in brackets: `[::1]`
  * `size`: request size in bytes
  * `port`: client's port
  * `duration`: response duration
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

## External Plugin

*Firewall* and other associated policy plugins in this repository are *external* plugins, which means it they are not included in CoreDNS releases.  To use the plugins in this repository, you'll need to build a CoreDNS image with the plugins you want to add included in the plugin list. In a nutshell you'll need to:
* Clone https://github.com/coredns/coredns
* Add this plugin to [plugin.cfg](https://github.com/coredns/coredns/blob/master/plugin.cfg) per instructions therein.
* `make -f Makefile.release DOCKER=your-docker-repo release`
* `make -f Makefile.release DOCKER=your-docker-repo docker`
* `make -f Makefile.release DOCKER=your-docker-repo docker-push`

## Examples

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

### Custom Metadata Policy
This example uses the *metadata* plugin to define labels `group_id` and `client_id` with values extracted from EDNS0. The firewall rules use those metadata to REFUSE any query without a group_id of `123456789` or client_id of `ABCDEF`.

~~~ corefile
example.org {
   metadata {
      group_id edns0 0xffed bytes
      client_id edns0 0xffee bytes
   }
   firewall query {
      refuse [metadata/client_id] != 'ABCDEF'
      refuse [metadata/group_id] != '123456789'
      block true
   }
}
~~~

