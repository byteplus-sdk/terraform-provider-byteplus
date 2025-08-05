---
subcategory: "MONGODB"
layout: "byteplus"
page_title: "Byteplus: byteplus_mongodb_ssl_state"
sidebar_current: "docs-byteplus-resource-mongodb_ssl_state"
description: |-
  Provides a resource to manage mongodb ssl state
---
# byteplus_mongodb_ssl_state
Provides a resource to manage mongodb ssl state
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

resource "byteplus_mongodb_instance" "foo" {
  db_engine_version      = "MongoDB_4_2"
  instance_type          = "ReplicaSet"
  super_account_password = "@acc-test-123"
  node_spec              = "mongo.2c4g"
  mongos_node_spec       = "mongo.mongos.2c4g"
  instance_name          = "acc-test-mongo-replica"
  charge_type            = "PostPaid"
  project_name           = "default"
  mongos_node_number     = 2
  shard_number           = 3
  storage_space_gb       = 20
  subnet_id              = byteplus_subnet.foo.id
  zone_id                = data.byteplus_zones.foo.zones[0].id
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_mongodb_ssl_state" "foo" {
  instance_id = byteplus_mongodb_instance.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The ID of mongodb instance.
* `ssl_action` - (Optional) The action of ssl, valid value contains `Update`. Set `ssl_action` to `Update` will will trigger an SSL update operation when executing `terraform apply`.When the current time is less than 30 days from the `ssl_expired_time`, executing `terraform apply` will automatically renew the SSL.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `is_valid` - Whetehr SSL is valid.
* `ssl_enable` - Whether SSL is enabled.
* `ssl_expired_time` - The expire time of SSL.


## Import
mongodb ssl state can be imported using the ssl:instanceId, e.g.
```
$ terraform import byteplus_mongodb_ssl_state.default ssl:mongo-shard-d050db19xxx
```

