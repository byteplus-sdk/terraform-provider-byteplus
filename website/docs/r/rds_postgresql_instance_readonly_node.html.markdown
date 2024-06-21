---
subcategory: "RDS_POSTGRESQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_postgresql_instance_readonly_node"
sidebar_current: "docs-byteplus-resource-rds_postgresql_instance_readonly_node"
description: |-
  Provides a resource to manage rds postgresql instance readonly node
---
# byteplus_rds_postgresql_instance_readonly_node
Provides a resource to manage rds postgresql instance readonly node
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-project1"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-subnet-test-2"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}


resource "byteplus_rds_postgresql_instance" "foo" {
  db_engine_version = "PostgreSQL_12"
  node_spec         = "rds.postgres.1c2g"
  primary_zone_id   = data.byteplus_zones.foo.zones[0].id
  secondary_zone_id = data.byteplus_zones.foo.zones[0].id
  storage_space     = 40
  subnet_id         = byteplus_subnet.foo.id
  instance_name     = "acc-test-1"
  charge_info {
    charge_type = "PostPaid"
  }
  project_name = "default"
  tags {
    key   = "tfk1"
    value = "tfv1"
  }
  parameters {
    name  = "auto_explain.log_analyze"
    value = "off"
  }
  parameters {
    name  = "auto_explain.log_format"
    value = "text"
  }
}

resource "byteplus_rds_postgresql_instance_readonly_node" "foo" {
  instance_id = byteplus_rds_postgresql_instance.foo.id
  node_spec   = "rds.postgres.1c2g"
  zone_id     = data.byteplus_zones.foo.zones[0].id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The RDS PostgreSQL instance id of the readonly node.
* `node_spec` - (Required) The specification of readonly node.
* `zone_id` - (Required, ForceNew) The available zone of readonly node.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `node_id` - The id of the readonly node.


## Import
RdsPostgresqlInstanceReadonlyNode can be imported using the instance_id:node_id, e.g.
```
$ terraform import byteplus_rds_postgresql_instance_readonly_node.default postgres-21a3333b****:postgres-ca7b7019****
```

