---
subcategory: "RDS_POSTGRESQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_postgresql_schemas"
sidebar_current: "docs-byteplus-datasource-rds_postgresql_schemas"
description: |-
  Use this data source to query detailed information of rds postgresql schemas
---
# byteplus_rds_postgresql_schemas
Use this data source to query detailed information of rds postgresql schemas
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

resource "byteplus_rds_postgresql_database" "foo" {
  db_name     = "acc-test"
  instance_id = byteplus_rds_postgresql_instance.foo.id
  c_type      = "C"
  collate     = "zh_CN.utf8"
}

resource "byteplus_rds_postgresql_account" "foo" {
  account_name       = "acc-test-account"
  account_password   = "9wc@********12"
  account_type       = "Normal"
  instance_id        = byteplus_rds_postgresql_instance.foo.id
  account_privileges = "Inherit,Login,CreateRole,CreateDB"
}

resource "byteplus_rds_postgresql_account" "foo1" {
  account_name       = "acc-test-account1"
  account_password   = "9wc@*******12"
  account_type       = "Normal"
  instance_id        = byteplus_rds_postgresql_instance.foo.id
  account_privileges = "Inherit,Login,CreateRole,CreateDB"
}

resource "byteplus_rds_postgresql_schema" "foo" {
  db_name     = byteplus_rds_postgresql_database.foo.db_name
  instance_id = byteplus_rds_postgresql_instance.foo.id
  owner       = byteplus_rds_postgresql_account.foo.account_name
  schema_name = "acc-test-schema"
}

data "byteplus_rds_postgresql_schemas" "foo" {
  db_name     = byteplus_rds_postgresql_schema.foo.db_name
  instance_id = byteplus_rds_postgresql_instance.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required) The id of the instance.
* `db_name` - (Optional) The name of the database.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `schemas` - The collection of query.
    * `db_name` - The name of the database.
    * `owner` - The owner of the schema.
    * `schema_name` - The name of the schema.
* `total_count` - The total count of query.


