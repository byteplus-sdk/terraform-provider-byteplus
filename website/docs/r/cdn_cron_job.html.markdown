---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_cron_job"
sidebar_current: "docs-byteplus-resource-cdn_cron_job"
description: |-
  Provides a resource to manage cdn cron job
---
# byteplus_cdn_cron_job
Provides a resource to manage cdn cron job
## Example Usage
```hcl
resource "byteplus_cdn_edge_function" "foo" {
  name         = "acc-test-function"
  remark       = "tf-test"
  project_name = "default"
  source_code  = base64encode("hello world")
  envs {
    key   = "k1"
    value = "v1"
  }
  canary_countries = ["China", "Japan", "United Kingdom"]
}

resource "byteplus_cdn_cron_job" "foo" {
  function_id     = byteplus_cdn_edge_function.foo.id
  job_name        = "acc-test-cron-job"
  description     = "tf-test"
  cron_type       = 1
  cron_expression = "0 17 10 * *"
  parameter       = "test"
}
```
## Argument Reference
The following arguments are supported:
* `cron_expression` - (Required) The execution expression. The expression must meet the following requirements:
Supports cron expressions (does not support second-level triggers).
* `cron_type` - (Required) The schedule type of the cron job. Possible values:
1: Global schedule.
2: Single point schedule.
* `function_id` - (Required, ForceNew) The id of the function.
* `job_name` - (Required, ForceNew) The name of the cron job. The name must meet the following requirements:
Each cron job name for a function must be unique
Length should not exceed 128 characters.
* `description` - (Optional) The description of the cron job.
* `parameter` - (Optional) The execution parameter of the cron job.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `create_time` - The creation time of the cron job. Displayed in UNIX timestamp format.
* `job_state` - The status of the cron job.
* `update_time` - The update time of the cron job. Displayed in UNIX timestamp format.


## Import
CdnCronJob can be imported using the function_id:job_name, e.g.
```
$ terraform import byteplus_cdn_cron_job.default function_id:job_name
```

