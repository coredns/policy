# themis

*themis* - implements Infoblox's Themis policy engine as CoreDNS _firewall_ policy engine.

## Syntax

```
themis ENGINE-NAME {
    pdp POLICY-FILE CONTENT [CONTENT...]
    endpoint PDP [PDP...]
    attr NAME LABEL [DSTTYPE]
    debug_query_suffix SUFFIX
    debug_id ID
    metrics METRIC [METRIC...]
    streams COUNT [BALANCE]
    connection_timeout
    transfer ATTR [ATTR...]
    dnstap
    connection_timeout
    log
    max_request_size [[auto] SIZE]
    max_response_attributes auto | COUNT
    cache [TTL [SIZE]]
}
```

* **ENGINE-NAME** is the name of the policy engine, used by the firewall plugin to uniquely identify the instance.
  Each instance of _themis_ in the Corefile must have a unique **ENGINE-NAME**.

* `pdp` defines themis policy and content files for local policy evaluation

* `endpoint` defines a list themis **PDP** addresses for remote policy evaluation

* `attr` is used for assigning labels into PDP attributes. `attr` may be defined multiple times.
  **DSTTYPE** allowed values depends on Themis PDP implementation, e.g. string (default), domain, address.

* `debug_query_suffix` enables debug query feature. **SUFFIX** must end with a dot. 

* `debug_id` is used to assist debugging. **ID** is a unique id that can be used to help determine
  which CoreDNS instance created a response.

* `metrics` defines a list of prometheus metrics to report.

* `streams` **COUNT** sets the number of gRPC streams for PDP connections.
  **BALANCE** can be `round-robin` (default), or `hot-spot`.

* `transfer` defines the set of attributes from domain validation response tha
  should be inserted into IP validation request.

* `dnstap` defines attributes to be included in the extra field of DNStap message if received
  from the PDP.

* `connection_timeout` sets the timeout for query validation when no PDP servers are available.
  A negative value or `no` means wait forever, the default behavior. A timeout of `0` causes
  validation to fail instantly if there are no PDP servers. The option works only if gRPC streams are
  greater than 0.

* `log` enables logging of the PDP request and response

* `max_request_size` sets maximum buffer size in bytes for serialized request. Setting the limit
  too high will make the plugin to allocate too much memory. Setting the limit too small can lead
  to buffer overflow errors during validation. If `auto` is set, the plugin allocates the required
  amount of bytes for each request. If both `auto` and **SIZE** are set, **SIZE** is only used for
  cache allocations and not for limiting the request buffer. 

* `max_response_attributes` sets the maximum number of attributes expected from the PDP. If the value
  is `auto`, the plugin automatically allocates the necessary number of attributes for each PDP response.

* `cache` enables decision cache. **TTL** default value is 10 minutes. **SIZE** limits the memory cache 
  will use to given number of megabytes. If **SIZE** isn't provided cache can grow indeterminately.

## Firewall Policy Engine

This plugin is not a standalone plugin.  It must be used in conjunction with the _firewall_ plugin to function.
For this plugin to be active, the _firewall_ plugin must reference it in a rule.  See the "Policy Engine Plugins"
section of the _firewall_ plugin README for more information.

## Examples

In the Corefile below, edns0 options with code 0xffee is split into two values - client_id (first 16 bytes)
and group_id (last 16 bytes). Edns0 options less than 32 bytes in size will not assign a client_id or group_id.

~~~ txt
. {
  themis myengine {
    endpoint 10.0.0.7:1234 10.0.0.8:1234
    attr client_id request/client_id
    attr group_id request/group_id
    attr uid request/another_id
    attr source_ip request/source_ip address 
    attr client_name request/client_name
    debug_query_suffix debug.
    debug_id instance_1
    streams 100
    transfer gid uid
    log
  }

  firewall query {
    themis myengine
  }
}
~~~

