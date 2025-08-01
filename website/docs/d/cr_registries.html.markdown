---
subcategory: "CR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cr_registries"
sidebar_current: "docs-byteplus-datasource-cr_registries"
description: |-
  Use this data source to query detailed information of cr registries
---
# byteplus_cr_registries
Use this data source to query detailed information of cr registries
## Example Usage
```hcl
data "byteplus_cr_registries" "foo" {
  # names=["liaoliuqing-prune-test"]
  # types=["Enterprise"]
  statuses {
    phase     = "Running"
    condition = "Ok"
  }
}
```
## Argument Reference
The following arguments are supported:
* `names` - (Optional) The list of registry names to query.
* `output_file` - (Optional) File name where to save data source results.
* `projects` - (Optional) The list of project names to query.
* `resource_tags` - (Optional) The tags of cr registry.
* `statuses` - (Optional) The list of registry statuses.
* `types` - (Optional) The list of registry types to query.

The `resource_tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `values` - (Required) The Value of Tags.

The `statuses` object supports the following:

* `condition` - (Optional) The condition of registry.
* `phase` - (Optional) The phase of status.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `registries` - The collection of registry query.
    * `charge_type` - The charge type of registry.
    * `create_time` - The creation time of registry.
    * `domains` - The domain of registry.
        * `domain` - The domain of registry.
        * `type` - The domain type of registry.
    * `name` - The name of registry.
    * `project` - The ProjectName of the cr registry.
    * `proxy_cache_enabled` - Whether to enable proxy cache.
    * `proxy_cache` - The proxy cache of registry. This field is valid when proxy_cache_enabled is true.
        * `endpoint` - The endpoint of proxy cache.
        * `skip_ssl_verify` - Whether to skip ssl verify.
        * `type` - The type of proxy cache. Valid values: `DockerHub`, `DockerRegistry`.
        * `username` - The username of proxy cache.
    * `resource_tags` - Tags.
        * `key` - The Key of Tags.
        * `value` - The Value of Tags.
    * `status` - The status of registry.
        * `conditions` - The condition of registry.
        * `phase` - The phase status of registry.
    * `type` - The type of registry.
    * `user_status` - The status of user.
    * `username` - The username of cr instance.
* `total_count` - The total count of registry query.


