---
subcategory: "MONGODB"
layout: "byteplus"
page_title: "Byteplus: byteplus_mongodb_instance_parameter"
sidebar_current: "docs-byteplus-resource-mongodb_instance_parameter"
description: |-
  Provides a resource to manage mongodb instance parameter
---
# byteplus_mongodb_instance_parameter
Provides a resource to manage mongodb instance parameter
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
  mongos_node_number     = 32
  shard_number           = 3
  storage_space_gb       = 20
  subnet_id              = byteplus_subnet.foo.id
  zone_id                = data.byteplus_zones.foo.zones[0].id
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_mongodb_instance_parameter" "foo" {
  instance_id     = byteplus_mongodb_instance.foo.id
  parameter_name  = "cursorTimeoutMillis"
  parameter_role  = "Node"
  parameter_value = "600111"
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The instance ID.
* `parameter_name` - (Required, ForceNew) The name of parameter. The parameter resource can only be added or modified, deleting this resource will not actually execute any operation.
* `parameter_role` - (Required, ForceNew) The node type to which the parameter belongs. The value range is as follows: Node, Shard, ConfigServer, Mongos.
* `parameter_value` - (Required) The value of parameter.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
mongodb parameter can be imported using the param:instanceId:parameterName:parameterRole, e.g.
```
$ terraform import byteplus_mongodb_instance_parameter.default param:mongo-replica-e405f8e2****:connPoolMaxConnsPerHost
```

