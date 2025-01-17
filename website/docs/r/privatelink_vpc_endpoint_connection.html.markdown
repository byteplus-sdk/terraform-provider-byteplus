---
subcategory: "PRIVATELINK"
layout: "byteplus"
page_title: "Byteplus: byteplus_privatelink_vpc_endpoint_connection"
sidebar_current: "docs-byteplus-resource-privatelink_vpc_endpoint_connection"
description: |-
  Provides a resource to manage privatelink vpc endpoint connection
---
# byteplus_privatelink_vpc_endpoint_connection
Provides a resource to manage privatelink vpc endpoint connection
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
  description = "acc-test"
}

resource "byteplus_privatelink_vpc_endpoint" "foo" {
  security_group_ids = [byteplus_security_group.foo.id]
  service_id         = byteplus_privatelink_vpc_endpoint_service.foo.id
  endpoint_name      = "acc-test-ep"
  description        = "acc-test"
}

resource "byteplus_privatelink_vpc_endpoint_connection" "foo" {
  endpoint_id = byteplus_privatelink_vpc_endpoint.foo.id
  service_id  = byteplus_privatelink_vpc_endpoint_service.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `endpoint_id` - (Required, ForceNew) The id of the endpoint.
* `service_id` - (Required, ForceNew) The id of the security group.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `connection_status` - The status of the connection.
* `creation_time` - The create time of the connection.
* `endpoint_owner_account_id` - The account id of the vpc endpoint.
* `endpoint_vpc_id` - The vpc id of the vpc endpoint.
* `update_time` - The update time of the connection.
* `zones` - The available zones.
    * `network_interface_id` - The id of the network interface.
    * `network_interface_ip` - The ip address of the network interface.
    * `resource_id` - The id of the resource.
    * `subnet_id` - The id of the subnet.
    * `zone_domain` - The domain of the zone.
    * `zone_id` - The id of the zone.
    * `zone_status` - The status of the zone.


## Import
PrivateLink Vpc Endpoint Connection Service can be imported using the endpoint id and service id, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint_connection.default ep-3rel74u229dz45zsk2i6l69qa:epsvc-2byz5mykk9y4g2dx0efs4aqz3
```

