---
subcategory: "VPC"
layout: "byteplus"
page_title: "Byteplus: byteplus_route_entry"
sidebar_current: "docs-byteplus-resource-route_entry"
description: |-
  Provides a resource to manage route entry
---
# byteplus_route_entry
Provides a resource to manage route entry
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc-rn"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet-rn"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_nat_gateway" "foo" {
  vpc_id           = byteplus_vpc.foo.id
  subnet_id        = byteplus_subnet.foo.id
  spec             = "Small"
  nat_gateway_name = "acc-test-nat-rn"
}

resource "byteplus_route_table" "foo" {
  vpc_id           = byteplus_vpc.foo.id
  route_table_name = "acc-test-route-table"
}

resource "byteplus_route_entry" "foo" {
  route_table_id         = byteplus_route_table.foo.id
  destination_cidr_block = "172.16.1.0/24"
  next_hop_type          = "NatGW"
  next_hop_id            = byteplus_nat_gateway.foo.id
  route_entry_name       = "acc-test-route-entry-new"
}
```
## Argument Reference
The following arguments are supported:
* `destination_cidr_block` - (Required, ForceNew) The destination CIDR block of the route entry.
* `next_hop_id` - (Required, ForceNew) The id of the next hop.
* `next_hop_type` - (Required, ForceNew) The type of the next hop, Optional choice contains `Instance`, `HaVip`, `NetworkInterface`, `NatGW`, `VpnGW`, `TransitRouter`.
* `route_table_id` - (Required, ForceNew) The id of the route table.
* `description` - (Optional) The description of the route entry.
* `route_entry_name` - (Optional) The name of the route entry.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `route_entry_id` - The id of the route entry.
* `status` - The description of the route entry.


## Import
Route entry can be imported using the route_table_id:route_entry_id, e.g.
```
$ terraform import byteplus_route_entry.default vtb-274e19skkuhog7fap8u4i8ird:rte-274e1g9ei4k5c7fap8sp974fq
```

