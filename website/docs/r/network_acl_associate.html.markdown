---
subcategory: "VPC"
layout: "byteplus"
page_title: "Byteplus: byteplus_network_acl_associate"
sidebar_current: "docs-byteplus-resource-network_acl_associate"
description: |-
  Provides a resource to manage network acl associate
---
# byteplus_network_acl_associate
Provides a resource to manage network acl associate
## Example Usage
```hcl
resource "byteplus_network_acl" "foo" {
  vpc_id           = "vpc-ru0wv9alfoxsu3nuld85rpp"
  network_acl_name = "tf-test-acl"
}

resource "byteplus_network_acl_associate" "foo1" {
  network_acl_id = byteplus_network_acl.foo.id
  resource_id    = "subnet-637jxq81u5mon3gd6ivc7rj"
}
```
## Argument Reference
The following arguments are supported:
* `network_acl_id` - (Required, ForceNew) The id of Network Acl.
* `resource_id` - (Required, ForceNew) The resource id of Network Acl.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
NetworkAcl associate can be imported using the network_acl_id:resource_id, e.g.
```
$ terraform import byteplus_network_acl_associate.default nacl-172leak37mi9s4d1w33pswqkh:subnet-637jxq81u5mon3gd6ivc7rj
```

