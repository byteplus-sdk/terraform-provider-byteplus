---
subcategory: "CLB"
layout: "byteplus"
page_title: "Byteplus: byteplus_server_group_server"
sidebar_current: "docs-byteplus-resource-server_group_server"
description: |-
  Provides a resource to manage server group server
---
# byteplus_server_group_server
Provides a resource to manage server group server
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

resource "byteplus_security_group" "foo" {
  vpc_id              = byteplus_vpc.foo.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_ecs_instance" "foo" {
  image_id             = "image-ycjwwciuzy5pkh54xx8f"
  instance_type        = "ecs.c3i.large"
  instance_name        = "acc-test-ecs-name"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = byteplus_server_group.foo.id
  instance_id     = byteplus_ecs_instance.foo.id
  type            = "ecs"
  weight          = 100
  port            = 80
  description     = "This is a acc test server"
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of ecs instance or the network card bound to ecs instance.
* `port` - (Required) The port receiving request.
* `server_group_id` - (Required, ForceNew) The ID of the ServerGroup.
* `type` - (Required, ForceNew) The type of instance. Optional choice contains `ecs`, `eni`.
* `description` - (Optional) The description of the instance.
* `ip` - (Optional, ForceNew) The private ip of the instance.
* `weight` - (Optional) The weight of the instance, range in 0~100.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `server_id` - The server id of instance in ServerGroup.


## Import
ServerGroupServer can be imported using the id, e.g.
```
$ terraform import byteplus_server_group_server.default rsp-274xltv2*****8tlv3q3s:rs-3ciynux6i1x4w****rszh49sj
```

