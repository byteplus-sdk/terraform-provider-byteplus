---
subcategory: "PRIVATELINK"
layout: "byteplus"
page_title: "Byteplus: byteplus_privatelink_vpc_endpoint_service_permissions"
sidebar_current: "docs-byteplus-datasource-privatelink_vpc_endpoint_service_permissions"
description: |-
  Use this data source to query detailed information of privatelink vpc endpoint service permissions
---
# byteplus_privatelink_vpc_endpoint_service_permissions
Use this data source to query detailed information of privatelink vpc endpoint service permissions
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

resource "byteplus_privatelink_vpc_endpoint_service_permission" "foo" {
  service_id        = byteplus_privatelink_vpc_endpoint_service.foo.id
  permit_account_id = "210000000"
}

data "byteplus_privatelink_vpc_endpoint_service_permissions" "foo" {
  permit_account_id = byteplus_privatelink_vpc_endpoint_service_permission.foo.permit_account_id
  service_id        = byteplus_privatelink_vpc_endpoint_service.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `service_id` - (Required) The Id of service.
* `output_file` - (Optional) File name where to save data source results.
* `permit_account_id` - (Optional) The Id of permit account.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `permissions` - The collection of query.
    * `permit_account_id` - The permit account id.
* `total_count` - Returns the total amount of the data list.


