---
subcategory: "PRIVATELINK"
layout: "byteplus"
page_title: "Byteplus: byteplus_privatelink_vpc_endpoint_service_permission"
sidebar_current: "docs-byteplus-resource-privatelink_vpc_endpoint_service_permission"
description: |-
  Provides a resource to manage privatelink vpc endpoint service permission
---
# byteplus_privatelink_vpc_endpoint_service_permission
Provides a resource to manage privatelink vpc endpoint service permission
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
```
## Argument Reference
The following arguments are supported:
* `permit_account_id` - (Required, ForceNew) The id of account.
* `service_id` - (Required, ForceNew) The id of service.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
VpcEndpointServicePermission can be imported using the serviceId:permitAccountId, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint_service_permission.default epsvc-2fe630gurkl37k5gfuy33****:2100000000
```

