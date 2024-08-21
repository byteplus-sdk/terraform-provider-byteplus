---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_cron_jobs"
sidebar_current: "docs-byteplus-datasource-cdn_cron_jobs"
description: |-
  Use this data source to query detailed information of cdn cron jobs
---
# byteplus_cdn_cron_jobs
Use this data source to query detailed information of cdn cron jobs
## Example Usage
```hcl
data "byteplus_cdn_cron_jobs" "foo" {
  function_id = "8f06f8db8d6b4bcdb979db68273f****"
}
```
## Argument Reference
The following arguments are supported:
* `function_id` - (Required) The id of the function.
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `cron_jobs` - The collection of query.
    * `create_time` - The creation time of the cron job. Displayed in UNIX timestamp format.
    * `cron_expression` - The cron expression of the cron job.
    * `cron_type` - The type of the cron job.
    * `description` - The description of the cron job.
    * `job_name` - The name of the cron job.
    * `job_state` - The status of the cron job.
    * `parameter` - The parameter of the cron job.
    * `update_time` - The update time of the cron job. Displayed in UNIX timestamp format.
* `total_count` - The total count of query.


