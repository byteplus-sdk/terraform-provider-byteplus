---
subcategory: "REDIS"
layout: "byteplus"
page_title: "Byteplus: byteplus_redis_pitr_time_periods"
sidebar_current: "docs-byteplus-datasource-redis_pitr_time_periods"
description: |-
  Use this data source to query detailed information of redis pitr time periods
---
# byteplus_redis_pitr_time_periods
Use this data source to query detailed information of redis pitr time periods
## Example Usage
```hcl
data "byteplus_redis_pitr_time_windows" "default" {
  ids = ["redis-cnlficlt4974swtbz", "redis-cnlfq69d1y1tnguxz"]
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Required) The ids of the instances.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `periods` - The list of time windows.
    * `end_time` - Recoverable end time (UTC time) supported when restoring data by point in time.
    * `instance_id` - The instance id.
    * `start_time` - The recoverable start time (in UTC time) supported when restoring data by point in time.
* `total_count` - The total count of redis instances time window query.


