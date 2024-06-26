---
subcategory: "CLB"
layout: "byteplus"
page_title: "Byteplus: byteplus_server_group_servers"
sidebar_current: "docs-byteplus-datasource-server_group_servers"
description: |-
  Use this data source to query detailed information of server group servers
---
# byteplus_server_group_servers
Use this data source to query detailed information of server group servers
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

data "byteplus_server_group_servers" "foo" {
  ids             = [element(split(":", byteplus_server_group_server.foo.id), length(split(":", byteplus_server_group_server.foo.id)) - 1)]
  server_group_id = byteplus_server_group.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `server_group_id` - (Required) The ID of the ServerGroup.
* `ids` - (Optional) The list of ServerGroupServer IDs.
* `name_regex` - (Optional) A Name Regex of ServerGroupServer.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `servers` - The server list of ServerGroup.
    * `description` - The description of the instance.
    * `id` - The server id of instance in ServerGroup.
    * `instance_id` - The ID of ecs instance or the network card bound to ecs instance.
    * `ip` - The private ip of the instance.
    * `port` - The port receiving request.
    * `server_id` - The server id of instance in ServerGroup.
    * `type` - The type of instance. Optional choice contains `ecs`, `eni`.
    * `weight` - The weight of the instance.
* `total_count` - The total count of ServerGroupServer query.


