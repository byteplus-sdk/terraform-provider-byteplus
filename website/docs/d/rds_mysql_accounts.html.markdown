---
subcategory: "RDS_MYSQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_mysql_accounts"
sidebar_current: "docs-byteplus-datasource-rds_mysql_accounts"
description: |-
  Use this data source to query detailed information of rds mysql accounts
---
# byteplus_rds_mysql_accounts
Use this data source to query detailed information of rds mysql accounts
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

resource "byteplus_rds_mysql_instance" "foo" {
  instance_name          = "acc-test-rds-mysql"
  db_engine_version      = "MySQL_5_7"
  node_spec              = "rds.mysql.1c2g"
  primary_zone_id        = data.byteplus_zones.foo.zones[0].id
  secondary_zone_id      = data.byteplus_zones.foo.zones[0].id
  storage_space          = 80
  subnet_id              = byteplus_subnet.foo.id
  lower_case_table_names = "1"
  charge_info {
    charge_type = "PostPaid"
  }
  parameters {
    parameter_name  = "auto_increment_increment"
    parameter_value = "2"
  }
  parameters {
    parameter_name  = "auto_increment_offset"
    parameter_value = "4"
  }
}

resource "byteplus_rds_mysql_database" "foo" {
  db_name     = "acc-test-db"
  instance_id = byteplus_rds_mysql_instance.foo.id
}

resource "byteplus_rds_mysql_account" "foo" {
  account_name     = "acc-test-account"
  account_password = "93f0cb0614Aab12"
  account_type     = "Normal"
  instance_id      = byteplus_rds_mysql_instance.foo.id
  account_privileges {
    db_name                  = byteplus_rds_mysql_database.foo.db_name
    account_privilege        = "Custom"
    account_privilege_detail = "SELECT,INSERT"
  }
}

data "byteplus_rds_mysql_accounts" "foo" {
  instance_id  = byteplus_rds_mysql_instance.foo.id
  account_name = byteplus_rds_mysql_account.foo.account_name
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required) The id of the RDS instance.
* `account_name` - (Optional) The name of the database account. This field supports fuzzy query.
* `name_regex` - (Optional) A Name Regex of database account.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `accounts` - The collection of RDS instance account query.
    * `account_name` - The name of the database account.
    * `account_privileges` - The privilege detail list of RDS mysql instance account.
        * `account_privilege_detail` - The privilege detail of the account.
        * `account_privilege` - The privilege type of the account.
        * `db_name` - The name of database.
    * `account_status` - The status of the database account.
    * `account_type` - The type of the database account.
* `total_count` - The total count of database account query.


