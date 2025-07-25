---
subcategory: "KAFKA"
layout: "byteplus"
page_title: "Byteplus: byteplus_kafka_topic_partitions"
sidebar_current: "docs-byteplus-datasource-kafka_topic_partitions"
description: |-
  Use this data source to query detailed information of kafka topic partitions
---
# byteplus_kafka_topic_partitions
Use this data source to query detailed information of kafka topic partitions
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

resource "byteplus_kafka_sasl_user" "foo" {
  user_name     = "acc-test-user"
  instance_id   = byteplus_kafka_instance.foo.id
  user_password = "suqsnis123!"
  description   = "tf-test"
  all_authority = true
  password_type = "Scram"
}

resource "byteplus_kafka_topic" "foo" {
  topic_name       = "acc-test-topic"
  instance_id      = byteplus_kafka_instance.foo.id
  description      = "tf-test"
  partition_number = 15
  replica_number   = 3

  parameters {
    min_insync_replica_number = 2
    message_max_byte          = 10
    log_retention_hours       = 96
  }

  all_authority = false
  access_policies {
    user_name     = byteplus_kafka_sasl_user.foo.user_name
    access_policy = "Pub"
  }
}

data "byteplus_kafka_topic_partitions" "default" {
  instance_id   = byteplus_kafka_instance.foo.id
  topic_name    = byteplus_kafka_topic.foo.topic_name
  partition_ids = [1, 2]
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required) The id of kafka instance.
* `topic_name` - (Required) The name of kafka topic.
* `output_file` - (Optional) File name where to save data source results.
* `partition_ids` - (Optional) The index number of partition.
* `under_insync_only` - (Optional) Whether to only query the list of partitions that have out-of-sync replicas, the default value is false.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `partitions` - The collection of query.
    * `end_offset` - The end offset of partition leader.
    * `insync_replicas` - The insync replica info.
    * `leader` - The leader info of partition.
    * `message_count` - The count of message.
    * `partition_id` - The index number of partition.
    * `replicas` - The replica info.
    * `start_offset` - The start offset of partition leader.
    * `under_insync_replicas` - The under insync replica info.
* `total_count` - The total count of query.


