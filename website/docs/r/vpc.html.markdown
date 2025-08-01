---
subcategory: "VPC"
layout: "byteplus"
page_title: "Byteplus: byteplus_vpc"
sidebar_current: "docs-byteplus-resource-vpc"
description: |-
  Provides a resource to manage vpc
---
# byteplus_vpc
Provides a resource to manage vpc
## Example Usage
```hcl
# query available zones in current region
data "byteplus_zones" "foo" {
}

# create vpc
resource "byteplus_vpc" "foo" {
  vpc_name     = "acc-test-vpc"
  cidr_block   = "172.16.0.0/16"
  dns_servers  = ["8.8.8.8", "114.114.114.114"]
  project_name = "default"
}

# create subnet
resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

# create security group
resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = byteplus_vpc.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `cidr_block` - (Required, ForceNew) A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).
* `description` - (Optional) The description of the VPC.
* `dns_servers` - (Optional) The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.
* `enable_ipv6` - (Optional) Specifies whether to enable the IPv6 CIDR block of the VPC.
* `ipv6_cidr_block` - (Optional) The IPv6 CIDR block of the VPC.
* `project_name` - (Optional) The ProjectName of the VPC.
* `tags` - (Optional) Tags.
* `vpc_name` - (Optional) The name of the VPC.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `account_id` - The account ID of VPC.
* `associate_cens` - The associate cen list of VPC.
    * `cen_id` - The ID of CEN.
    * `cen_owner_id` - The owner ID of CEN.
    * `cen_status` - The status of CEN.
* `auxiliary_cidr_blocks` - The auxiliary cidr block list of VPC.
* `creation_time` - Creation time of VPC.
* `nat_gateway_ids` - The nat gateway ID list of VPC.
* `route_table_ids` - The route table ID list of VPC.
* `security_group_ids` - The security group ID list of VPC.
* `status` - Status of VPC.
* `subnet_ids` - The subnet ID list of VPC.
* `update_time` - The update time of VPC.
* `vpc_id` - The ID of VPC.


## Import
VPC can be imported using the id, e.g.
```
$ terraform import byteplus_vpc.default vpc-mizl7m1kqccg5smt1bdpijuj
```

