---
subcategory: "PRIVATELINK"
layout: "byteplus"
page_title: "Byteplus: byteplus_privatelink_vpc_endpoints"
sidebar_current: "docs-byteplus-datasource-privatelink_vpc_endpoints"
description: |-
  Use this data source to query detailed information of privatelink vpc endpoints
---
# byteplus_privatelink_vpc_endpoints
Use this data source to query detailed information of privatelink vpc endpoints
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = byteplus_vpc.foo.id
}

resource "byteplus_clb" "foo" {
  type                       = "public"
  subnet_id                  = byteplus_subnet.foo.id
  load_balancer_spec         = "small_1"
  description                = "acc-test-demo"
  load_balancer_name         = "acc-test-clb"
  load_balancer_billing_type = "PostPaid"
  eip_billing_config {
    isp              = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_privatelink_vpc_endpoint_service" "foo" {
  resources {
    resource_id   = byteplus_clb.foo.id
    resource_type = "CLB"
  }
  description         = "acc-test"
  auto_accept_enabled = true
}

resource "byteplus_privatelink_vpc_endpoint" "foo" {
  security_group_ids = [byteplus_security_group.foo.id]
  service_id         = byteplus_privatelink_vpc_endpoint_service.foo.id
  endpoint_name      = "acc-test-ep"
  description        = "acc-test"
  count              = 2
}

data "byteplus_privatelink_vpc_endpoints" "foo" {
  ids = byteplus_privatelink_vpc_endpoint.foo[*].id
}
```
## Argument Reference
The following arguments are supported:
* `endpoint_name` - (Optional) The name of vpc endpoint.
* `ids` - (Optional) The IDs of vpc endpoint.
* `name_regex` - (Optional) A Name Regex of vpc endpoint.
* `output_file` - (Optional) File name where to save data source results.
* `service_name` - (Optional) The name of vpc endpoint service.
* `status` - (Optional) The status of vpc endpoint. Valid values: `Creating`, `Pending`, `Available`, `Deleting`, `Inactive`.
* `vpc_id` - (Optional) The vpc id of vpc endpoint.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - Returns the total amount of the data list.
* `vpc_endpoints` - The collection of query.
    * `business_status` - Whether the vpc endpoint is locked.
    * `connection_status` - The connection  status of vpc endpoint.
    * `creation_time` - The create time of vpc endpoint.
    * `deleted_time` - The delete time of vpc endpoint.
    * `description` - The description of vpc endpoint.
    * `endpoint_domain` - The domain of vpc endpoint.
    * `endpoint_id` - The Id of vpc endpoint.
    * `endpoint_name` - The name of vpc endpoint.
    * `endpoint_type` - The type of vpc endpoint.
    * `id` - The Id of vpc endpoint.
    * `service_id` - The Id of vpc endpoint service.
    * `service_name` - The name of vpc endpoint service.
    * `status` - The status of vpc endpoint.
    * `update_time` - The update time of vpc endpoint.
    * `vpc_id` - The vpc id of vpc endpoint.


