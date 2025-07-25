---
subcategory: "KAFKA"
layout: "byteplus"
page_title: "Byteplus: byteplus_kafka_public_address"
sidebar_current: "docs-byteplus-resource-kafka_public_address"
description: |-
  Provides a resource to manage kafka public address
---
# byteplus_kafka_public_address
Provides a resource to manage kafka public address
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

resource "byteplus_kafka_instance" "foo" {
  instance_name        = "acc-test-kafka"
  instance_description = "tf-test"
  version              = "2.2.2"
  compute_spec         = "kafka.20xrate.hw"
  subnet_id            = byteplus_subnet.foo.id
  user_name            = "tf-user"
  user_password        = "tf-pass!@q1"
  charge_type          = "PostPaid"
  storage_space        = 300
  partition_number     = 350
  project_name         = "default"
  tags {
    key   = "k1"
    value = "v1"
  }

  parameters {
    parameter_name  = "MessageMaxByte"
    parameter_value = "12"
  }
  parameters {
    parameter_name  = "LogRetentionHours"
    parameter_value = "70"
  }
}

resource "byteplus_eip_address" "foo" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "BGP"
  name         = "acc-test-eip"
  description  = "tf-test"
  project_name = "default"
}

resource "byteplus_kafka_public_address" "foo" {
  instance_id = byteplus_kafka_instance.foo.id
  eip_id      = byteplus_eip_address.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `eip_id` - (Required, ForceNew) The id of eip.
* `instance_id` - (Required, ForceNew) The id of kafka instance.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `endpoint_type` - The endpoint type of instance.
* `network_type` - The network type of instance.
* `public_endpoint` - The public endpoint of instance.


## Import
KafkaPublicAddress can be imported using the instance_id:eip_id, e.g.
```
$ terraform import byteplus_kafka_public_address.default instance_id:eip_id
```

