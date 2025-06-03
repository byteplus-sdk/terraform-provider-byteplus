---
subcategory: "RDS_MYSQL"
layout: "byteplus"
page_title: "Byteplus: byteplus_rds_mysql_regions"
sidebar_current: "docs-byteplus-datasource-rds_mysql_regions"
description: |-
  Use this data source to query detailed information of rds mysql regions
---
# byteplus_rds_mysql_regions
Use this data source to query detailed information of rds mysql regions
## Example Usage
```hcl
data "byteplus_rds_mysql_regions" "foo" {
}
```
## Argument Reference
The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `regions` - The collection of region query.
    * `region_id` - The id of the region.
    * `region_name` - The name of region.
* `total_count` - The total count of region query.


