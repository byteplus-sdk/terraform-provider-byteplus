---
subcategory: "CLB"
layout: "byteplus"
page_title: "Byteplus: byteplus_server_group"
sidebar_current: "docs-byteplus-resource-server_group"
description: |-
  Provides a resource to manage server group
---
# byteplus_server_group
Provides a resource to manage server group
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}
resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_clb" "foo" {
  type               = "public"
  subnet_id          = byteplus_subnet.foo.id
  load_balancer_spec = "small_1"
  description        = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp              = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id  = byteplus_clb.foo.id
  server_group_name = "acc-test-create"
  description       = "hello demo11"
}
```
## Argument Reference
The following arguments are supported:
* `load_balancer_id` - (Required, ForceNew) The ID of the Clb.
* `address_ip_version` - (Optional, ForceNew) The address ip version of the ServerGroup. Valid values: `ipv4`, `ipv6`. Default is `ipv4`.
* `description` - (Optional) The description of ServerGroup.
* `server_group_id` - (Optional) The ID of the ServerGroup.
* `server_group_name` - (Optional) The name of the ServerGroup.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
ServerGroup can be imported using the id, e.g.
```
$ terraform import byteplus_server_group.default rsp-273yv0kir1vk07fap8tt9jtwg
```

