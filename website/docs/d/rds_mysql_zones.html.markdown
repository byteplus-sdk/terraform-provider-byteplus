---
subcategory: "RDS_MYSQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_mysql_zones"
sidebar_current: "docs-byteplus-datasource-rds_mysql_zones"
description: |-
  Use this data source to query detailed information of rds mysql zones
---
# byteplus_rds_mysql_zones
Use this data source to query detailed information of rds mysql zones
## Example Usage
```hcl
data "byteplus_rds_mysql_zones" "foo" {
}
```
## Argument Reference
The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of query.
* `zones` - The collection of zone query.
    * `description` - The description of the zone.
    * `id` - The id of the zone.
    * `zone_id` - The id of the zone.
    * `zone_name` - The name of the zone.


