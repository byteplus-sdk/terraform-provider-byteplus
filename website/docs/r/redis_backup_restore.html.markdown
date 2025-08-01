---
subcategory: "REDIS"
layout: "byteplus"
page_title: "Byteplus: byteplus_redis_backup_restore"
sidebar_current: "docs-byteplus-resource-redis_backup_restore"
description: |-
  Provides a resource to manage redis backup restore
---
# byteplus_redis_backup_restore
Provides a resource to manage redis backup restore
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

resource "byteplus_redis_backup" "foo" {
  instance_id = byteplus_redis_instance.foo.id
}

resource "byteplus_redis_backup_restore" "foo" {
  instance_id = byteplus_redis_instance.foo.id
  time_point  = byteplus_redis_backup.foo.end_time
  backup_type = "Full"
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) Id of instance.
* `backup_point_id` - (Optional) Backup ID, used to specify the backups to be used when restoring by the backup set. When choosing to restore by backup set (i.e., BackupType is Full), this parameter is required. Use lifecycle and ignore_changes in import.
* `backup_type` - (Optional, ForceNew) The type of backup. The value can be Full or Inc.
* `time_point` - (Optional) Time point of backup, for example: 2021-11-09T06:07:26Z. Use lifecycle and ignore_changes in import.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Redis Backup Restore can be imported using the restore:instanceId, e.g.
```
$ terraform import byteplus_redis_backup_restore.default restore:redis-asdljioeixxxx
```

