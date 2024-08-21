---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_cron_job_state"
sidebar_current: "docs-byteplus-resource-cdn_cron_job_state"
description: |-
  Provides a resource to manage cdn cron job state
---
# byteplus_cdn_cron_job_state
Provides a resource to manage cdn cron job state
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

resource "byteplus_cdn_cron_job_state" "foo" {
  function_id = byteplus_cdn_edge_function.foo.id
  job_name    = byteplus_cdn_cron_job.foo.job_name
  action      = "Start"
}
```
## Argument Reference
The following arguments are supported:
* `action` - (Required) Start or Stop of corn job, the value can be `Start` or `Stop`. 
If the target status of the action is consistent with the current status of the corn job, the action will not actually be executed.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `function_id` - (Required, ForceNew) The id of the function.
* `job_name` - (Required, ForceNew) The name of the cron job.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `job_state` - The status of the cron job.


## Import
CdnCronJobState can be imported using the state:function_id:job_name, e.g.
```
$ terraform import byteplus_cdn_cron_job_state.default state:function_id:job_name
```

