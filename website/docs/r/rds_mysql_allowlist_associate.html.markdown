---
subcategory: "RDS_MYSQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_mysql_allowlist_associate"
sidebar_current: "docs-byteplus-resource-rds_mysql_allowlist_associate"
description: |-
  Provides a resource to manage rds mysql allowlist associate
---
# byteplus_rds_mysql_allowlist_associate
Provides a resource to manage rds mysql allowlist associate
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

resource "byteplus_rds_mysql_allowlist" "foo" {
  allow_list_name = "acc-test-allowlist"
  allow_list_desc = "acc-test"
  allow_list_type = "IPv4"
  allow_list      = ["192.168.0.0/24", "192.168.1.0/24"]
}

resource "byteplus_rds_mysql_allowlist_associate" "foo" {
  allow_list_id = byteplus_rds_mysql_allowlist.foo.id
  instance_id   = byteplus_rds_mysql_instance.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `allow_list_id` - (Required, ForceNew) The id of the allow list.
* `instance_id` - (Required, ForceNew) The id of the mysql instance.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
RDS AllowList Associate can be imported using the instance id and allow list id, e.g.
```
$ terraform import byteplus_rds_mysql_allowlist_associate.default rds-mysql-h441603c68aaa:acl-d1fd76693bd54e658912e7337d5b****
```

