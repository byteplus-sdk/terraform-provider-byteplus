---
subcategory: "REDIS"
layout: "byteplus"
page_title: "Byteplus: byteplus_redis_account"
sidebar_current: "docs-byteplus-resource-redis_account"
description: |-
  Provides a resource to manage redis account
---
# byteplus_redis_account
Provides a resource to manage redis account
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

resource "byteplus_redis_instance" "foo" {
  zone_ids            = [data.byteplus_zones.foo.zones[0].id]
  instance_name       = "acc-test-tf-redis"
  sharded_cluster     = 1
  password            = "1qaz!QAZ12"
  node_number         = 2
  shard_capacity      = 1024
  shard_number        = 2
  engine_version      = "5.0"
  subnet_id           = byteplus_subnet.foo.id
  deletion_protection = "disabled"
  vpc_auth_mode       = "close"
  charge_type         = "PostPaid"
  port                = 6381
  project_name        = "default"
}

resource "byteplus_redis_account" "foo" {
  account_name = "acc_test_account"
  instance_id  = byteplus_redis_instance.foo.id
  password     = "Password@@"
  role_name    = "ReadOnly"
}
```
## Argument Reference
The following arguments are supported:
* `account_name` - (Required, ForceNew) Redis account name.
* `instance_id` - (Required, ForceNew) The ID of the Redis instance.
* `password` - (Required) The password of the redis account. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `role_name` - (Required) Role type, the valid value can be `Administrator`, `ReadWrite`, `ReadOnly`, `NotDangerous`.
* `description` - (Optional) The description of the redis account.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Redis account can be imported using the instanceId:accountName, e.g.
```
$ terraform import byteplus_redis_account.default redis-42b38c769c4b:test
```

