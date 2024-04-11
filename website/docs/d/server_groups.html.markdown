---
subcategory: "CLB"
layout: "byteplus"
page_title: "Byteplus: byteplus_server_groups"
sidebar_current: "docs-byteplus-datasource-server_groups"
description: |-
  Use this data source to query detailed information of server groups
---
# byteplus_server_groups
Use this data source to query detailed information of server groups
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

data "byteplus_server_groups" "foo" {
  ids = [byteplus_server_group.foo.id]
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of ServerGroup IDs.
* `load_balancer_id` - (Optional) The id of the Clb.
* `name_regex` - (Optional) A Name Regex of ServerGroup.
* `output_file` - (Optional) File name where to save data source results.
* `server_group_name` - (Optional) The name of the ServerGroup.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `groups` - The collection of ServerGroup query.
    * `address_ip_version` - The address ip version of the ServerGroup.
    * `create_time` - The create time of the ServerGroup.
    * `description` - The description of the ServerGroup.
    * `id` - The ID of the ServerGroup.
    * `server_group_id` - The ID of the ServerGroup.
    * `server_group_name` - The name of the ServerGroup.
    * `update_time` - The update time of the ServerGroup.
* `total_count` - The total count of ServerGroup query.


