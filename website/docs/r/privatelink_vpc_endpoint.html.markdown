---
subcategory: "PRIVATELINK"
layout: "byteplus"
page_title: "Byteplus: byteplus_privatelink_vpc_endpoint"
sidebar_current: "docs-byteplus-resource-privatelink_vpc_endpoint"
description: |-
  Provides a resource to manage privatelink vpc endpoint
---
# byteplus_privatelink_vpc_endpoint
Provides a resource to manage privatelink vpc endpoint
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
}
```
## Argument Reference
The following arguments are supported:
* `security_group_ids` - (Required) The security group ids of vpc endpoint. It is recommended to bind security groups using the 'security_group_ids' field in this resource instead of using `byteplus_privatelink_security_group`.
For operations that jointly use this resource and `byteplus_privatelink_security_group`, use lifecycle ignore_changes to suppress changes to the 'security_group_ids' field.
* `service_id` - (Required, ForceNew) The id of vpc endpoint service.
* `description` - (Optional) The description of vpc endpoint.
* `endpoint_name` - (Optional) The name of vpc endpoint.
* `service_name` - (Optional, ForceNew) The name of vpc endpoint service.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `business_status` - Whether the vpc endpoint is locked.
* `connection_status` - The connection  status of vpc endpoint.
* `creation_time` - The create time of vpc endpoint.
* `deleted_time` - The delete time of vpc endpoint.
* `endpoint_domain` - The domain of vpc endpoint.
* `endpoint_type` - The type of vpc endpoint.
* `status` - The status of vpc endpoint.
* `update_time` - The update time of vpc endpoint.
* `vpc_id` - The vpc id of vpc endpoint.


## Import
VpcEndpoint can be imported using the id, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint.default ep-3rel74u229dz45zsk2i6l****
```

