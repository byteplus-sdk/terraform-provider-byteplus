---
subcategory: "KAFKA"
layout: "byteplus"
page_title: "Byteplus: byteplus_kafka_consumed_topics"
sidebar_current: "docs-byteplus-datasource-kafka_consumed_topics"
description: |-
  Use this data source to query detailed information of kafka consumed topics
---
# byteplus_kafka_consumed_topics
Use this data source to query detailed information of kafka consumed topics
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

data "byteplus_kafka_consumed_topics" "default" {
  instance_id = byteplus_kafka_instance.foo.id
  group_id    = byteplus_kafka_group.foo.group_id
  topic_name  = byteplus_kafka_topic.foo.topic_name
}
```
## Argument Reference
The following arguments are supported:
* `group_id` - (Required) The id of kafka group.
* `instance_id` - (Required) The id of kafka instance.
* `output_file` - (Optional) File name where to save data source results.
* `topic_name` - (Optional) The name of kafka topic. This field supports fuzzy query.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `consumed_topics` - The collection of query.
    * `accumulation` - The total amount of message accumulation in this topic for the consumer group.
    * `topic_name` - The name of kafka topic.
* `total_count` - The total count of query.


