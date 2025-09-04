---
subcategory: "KAFKA"
layout: "byteplus"
page_title: "Byteplus: byteplus_kafka_groups"
sidebar_current: "docs-byteplus-datasource-kafka_groups"
description: |-
  Use this data source to query detailed information of kafka groups
---
# byteplus_kafka_groups
Use this data source to query detailed information of kafka groups
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

resource "byteplus_kafka_group" "foo" {
  instance_id = byteplus_kafka_instance.foo.id
  group_id    = "acc-test-group"
  description = "tf-test"
}

data "byteplus_kafka_groups" "default" {
  instance_id = byteplus_kafka_group.foo.instance_id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required) The instance id of kafka group.
* `group_id` - (Optional) The id of kafka group, support fuzzy matching.
* `name_regex` - (Optional) A Name Regex of kafka group.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `groups` - The collection of query.
    * `group_id` - The id of kafka group.
    * `state` - The state of kafka group.
* `total_count` - The total count of query.


