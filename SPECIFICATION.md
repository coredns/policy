# Policy ... or Firewall 

Adding firewall capability to CoreDNS.

## History
Infoblox support on open source Policy Engine named Themis.
This project is available here : https://github.com/infobloxopen/themis

For Infoblox purposes, this Engine is tied with CoreDNS by a plugin that filter DNS Queries (and responses) by sending Policy queries to Themis
This implemetnation is production in Infoblox SAAS offering.

We would like to take advantage of this experience to define a "Firewall" plugin in CoreDNS that would be an enhancement of current implementation:

- more easy to use
- not limited to only one Policy Engine : CNCF host Open Policy Agent which is another engine
- propose easy immediate language for filtering, directly in the Corefile
- leverage Metadata plugin and ability for any plugin to contribute to the properties usable for filtering
- have, eventually, the ability to migrate Infoblox's plugin to this new Firewall plugin.


**Related enhancement or features already raised in CoreDNS:** 
- https://github.com/coredns/coredns/issues/178
- https://github.com/coredns/coredns/pull/1337



## Requirement

Expectation is to be able to filter the incoming DNS Queries, and the DNS Reply of this queries.

Filtering would be based on properties of the Querie/Reply in conjunction with Environment.
Examples:
- allow incoming queries based on the IP of the client, or the domain targeted
- allow incoming queries based on suspicious activity related to the domain targeted
- allow incoming queries based on the both the source IP, but also a kind of authorization on domains targeted
- same as above but with a selecting the authorization policy based on parameter provided in the query (edns0 value)
- refuse the reply if the IP returned is blacklisted
- etc ...

### Policy Engines

Several components could be used as an evaluation of the policy:  
 1. a simple Expression language : domain == 'google.com' and IP in-cidr(172.28.0.0/24)
 2. a policy engine available in CNCF projects : OPA - it suppose that 
 3. another opensource policy engine Themis, already used in conjunction with CoreDNS:  
-- either as a remote policy Engine with its own policy configuration  
-- either as an embedded poliy engine, in CoreDNS pod, with policy files provided as configuration  
 4. others ... could be another expression language more close to another business activity.

### Environment information

All these Engines would need to evaluate the policy against properties of the query or of the environment of execution:
- client IP
- domain targeted
- class of query
- returned IP
- authorization labels attached to the pod of the client that issued the query (for k8s)
- internal client-id encoded in the EDNS0 part of the query, by the DNS client
- etc ... 

We can split these information in 2 sets:

1. the ones attached to the request itself
2. the ones that can be collected on environement on execution by any plugin of CoreDNS (labels on the POD, process env variable, etc ..)

### Enhancement

The same Policy Engines could be used in other thing that filtering the query.
For instance, define what IP to return in a A/AAAA query based on some GEO logic, or other logic based on the query (client identification)

In that case, the same Engines, same Environment information could be reused but with another plugin that is not policy, 
but, in current case, a smart load balancer for the reply.


## Design proposition

In order to be generic enough to accept existing and future uses cases around, 
I propose to take advantage of the plugin organization of CoreDNS and split this whole feature in several plugins.

A. Separate the Policy Engine and configuration from their usage as a filtering policy  
1. the generic Policy plugin: define the filters, where to apply (query/reply) - a set of ordered of boolean rules that reply (accept/refuse) 
2. a plugin per Engine of Policy (expression, Open Policy Agent, Themis, others ...) - application of the plugin will create an instance of Engine with it's own configuration
  
B. Take advantage of the "Metadata" to carry the environement information to be evaluated in the rules
1. a Request plugin can implement the Provider interface and provide all metadata tied to the DNS query, including extraction of EDNS0 field values
2. Kubernetes can be extended to add the "label" information attached to POD or SERVICE
3. List to be continued depending the needs

NOTE: if we want to create later on a "SmartLB" plugin based on Policy Engine, we will create a dependancy with Policy because of the declaration if Engine interfaces and registration.
Which means this part, should, at one point, go into the internal plugins of CoreDNS.


### Sample of Corefile

#### Make a filter based on query type and domain targeted

**It is expected that the firewall plugin would be more often use in that simple way, for an easy filtering on expressions based on metadata**

```
.:53 {
   metadata
   request
   expression rule
   firewall . query {
      rule permit in-cidr([request/client_ip], 172.23.0.0/16) and ([request/qtype] == 'A')
      rule permit request/qtype == 'TXT'
      rule deny true
   }
} 
```

I suppose we will make a short-cut and have the Expression engine built directly in the firewall plugin.
Therefore any expression will be available immediately without requiring the plign "Expression" (or similar name)

In that case, we can simplify the above Corefile with:

```
.:53 {
   metadata
   request
   firewall . query {
      permit in-cidr([request/client_ip], 172.23.0.0/16) and (request/qtype == 'A')
      permit request/qtype == 'TXT'
      deny true
   }
} 
```

and specifying the 'permit' and 'deny' keywords as reserved and used for the expression built-in plugin


#### Make a filter based on client_id encoded in EDNS0 and filter reply on IP

Plugin request is xpected to be able to provide metadata from property of the DNS query
It will be able to extract values from the EDNS0 records.
'hex' is an encoding.

```
.:53 {
   metadata
   request {
      edns0 0xffed client_id      
      edns0 0xffee group_id hex 16 0 16      
   }
   firewall . query {
      permit request/client_id ~= 'infoblox'
      deny true
   }
   firewall . response {
      permit inc-cidr(request/returned_ip, 10.0.0.0/8)
      deny true
   }
} 
```

#### Make a filter based on external engine instances

```
.:53 {
   request {
      edns0 0xffed client_id      
      edns0 0xffee group_id      
   }
   themis permission-client {
      endpoint 10.0.0.7:1234, 10.0.0.8:1234
      attr client-id string request/client_id
      streams 100
   }
   themis permission-domain {
      file /home/root/policy-file-domain 
      attr domain string request/qname
   }
   themis permission-ip {
      file /home/root/policy-file-ip 
      attr ip string request/returned_ip
      transfer gid uid
   }
   kubernetes { ... 
   }
   opa opa-permission-ip {
   }
   firewall . query {
      themis permission-client
      themis permission-domain
      permit request/client_id ~= 'infoblox'      
      deny true

      # sample for kubernetes label related rules      
      permit k8s/clientpod.security > level4
      
      # k8s could be an Polcy Engine that filter based on RBAC
      kubernetes rbac
   }
   firewall . response {
      themis permission-domain
      deny true
   }
   load-balancer {
       
   }
} 

```

NOTE: It would be very time consuming to use 3 separated instance of Themis for the filtering. 
It is implemented here only for capability show use case

   
NOTE: Themis plugin would allow 2 instantiation types : one for embedded engine (same binary) and use xml files for policy. 
The second one is a remote engine running in its own container, with a communication over gRPC.

Open Policy Agent plugin is not known yet, but it is supposed to be a remote OPA Engine with REST API. 

### Features that are not yet in the design

Option passthrough : a way to to avoid filters
=> can be implemented using permit rules based on domain

Option log : need to check what is logged, when 
Dnstap logging in the queries. That is in the current PEP. Should it be "themis" ?  

Option debug : a mechanism to debug the Engine evaluations.



## Details specifications

### Request plugin

Should be regular plugin for CoreDNS
Could be the metadata plugin. 
Nothing specific.
Will need to reuse the EDNS0 syntax for extraction 

~~~
request {
      edns0 0xffed client_id      
      edns0 0xffee group_id hex 16 0 16
      edns0 <filed-id> <label> 
      ends0 <filed-id> <label> <encoded-format> <params of format ...>          
}
~~~

so far, only 'hex' format is supported with params <length>  <start> <end>


currently supported metadata is based on the variables used in REWRITE section. these are:

	queryName  = "qname"
	queryType  = "qtype"
	clientIP   = "client_ip"
	clientPort = "client_port"
	protocol   = "protocol"
	serverIP   = "server_ip"
	serverPort = "server_port"
 

### Firewall Plugin

#### Ruleset and point of application

It defines:
- a point of application of the filter. two application points are allowed : query or response
- each point of application is optional (can be omitted)
- at each point, the Firewall works like a Firewall ruleset:
- defined on ordered set of rules.
-- each rule is apply in order and return a result : undefined, allow, log, redirect, refuse, block, drop
NOTE: this list is coming from existing plugin used with Themis. It may need to be reviewd. 
for instance is 'redirect' a consistent option here ?  
-- if allow or block, the result of the filter is known and applied
-- if undefined, the net rule is evaluated
-- at the end of the ruleset, if the evaluation is still undefined, then the default operation apply : block

Depending the result of the ruleset the following operation happen:

- refuse : return SERVFAIL
- allow - just let next handler serve the query
- log - add a log and let next handler serve the query
- redirect - an IP is returned with evaluation to send the query to
- refuse - return a REFUSE to client
- block - return a NXDOMAIN to client
- drop - return a NOERROR with no info.

#### Interaction with Engine Instances

Firewall plugin scans all available Engine Instances for registration.
Each Engine Instance is expected to implement an interface that allow:

- build a rule from the initial expression (one before starting DNS Service)
- to build a context for rule execution (once per DNS query)  
- execute the rule (once per rule associated with this Engine Instance) 

#### Default Engine included on Firewall

For simplification of syntax and easyness, the Expression Engine will be embedded directly in the firewall plugin
it will have no name.
The rule syntax for this Engine will be:

**permit** <boolean expression>
**deny** <boolean expression>

a <boolean expression> is an expression that can evaluate to something we can finally result into a 'true' or 'false'.
in order to simpliy the expression, usual meaning of tranformation into boolean will ne implemented.
eg. "" is false, 0 is false, etc ...

This Engine will be based on the go expression evaluator : https://github.com/Knetic/govaluate

It may need some adaptation:
- accept variable names with / (so we can use directly labels of metadata as variables of expression). 
Right now we need to enclose a label into bracket to be correcly evaluated : [request/qtype] == "A"
- add some buil-in function or operator : in-cidr 
- resolve the evaluation of variables that are not provided as metadata : It should be considered as existing parameters with empty values.

NOTE: see here operator accepted : https://github.com/Knetic/govaluate/blob/master/MANUAL.md.
It includes regexp comparators, IN membership.

The language can be extended by adding functions. 

WARNING: looking more closely in the repo of this evaluator, I found:
- a PR is pending for adding IP related support (see https://github.com/Knetic/govaluate/pull/101) 
- there is no merge activity on this repo since almost one year .. !  

### Open Policy Agent - a Policy Engine hosted by CNCF

see : https://www.openpolicyagent.org/

### Themis  Open source Policy Engine, supported by Infoblox

repository is available here : https://github.com/infobloxopen/themis
Themis is already integrated with CoreDNS and used in that context by Infoblox.

We take advantage of this experience to build here a more generic firewall plugin.
We would like that current usage of Themis within Infoblox can eventually migrate to this new firewall plugin.
It will requires a re-compilation of the Corefile, but we expect that all options will be available.

As of today, the plugin defined in infobloxopen, includes:
- transformation of request, edns0 information into variables suitable for Themis (attributes)
- call of a unique rule for query and response that send evaluation to a remote Themis engine
- allow specific option on plugin like : log and fallthrough 


a sample of stanza for this plugin looks like:

~~~ txt
policy {
    endpoint 10.0.0.7:1234, 10.0.0.8:1234
    edns0 0xffee client_id hex string 32 0 32
    edns0 0xffee group_id hex string 32 16 32
    edns0 0xffef uid // equal edns0 0xffef uid hex string
    edns0 0xffea source_ip ip address
    edns0 0xffeb client_name bytes string
    debug_query_suffix debug.
    debug_id instance_1
    streams 100
    transfer gid uid
    passthrough mycompanyname.com. mycompanyname.org.
    log
}
~~~


