---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_edge_function"
sidebar_current: "docs-byteplus-resource-cdn_edge_function"
description: |-
  Provides a resource to manage cdn edge function
---
# byteplus_cdn_edge_function
Provides a resource to manage cdn edge function
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
```
## Argument Reference
The following arguments are supported:
* `name` - (Required) The name of the edge function.
* `canary_countries` - (Optional) The array of countries where the canary cluster is located.
* `envs` - (Optional) The environment variables of the edge function.
* `project_name` - (Optional) The name of the project to which the edge function belongs, defaulting to `default`.
* `remark` - (Optional) The remark of the edge function.
* `source_code` - (Optional) Code content. The input requirements are as follows: 
Not empty.
Value after base64 encoding.

The `envs` object supports the following:

* `key` - (Required) The key of the environment variable.
* `value` - (Required) The value of the environment variable.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `account_identity` - The account id of the edge function.
* `create_time` - The create time of the edge function. Displayed in UNIX timestamp format.
* `creator` - The creator of the edge function.
* `domain` - The domain name bound to the edge function.
* `status` - The status of the edge function.
* `update_time` - The update time of the edge function. Displayed in UNIX timestamp format.
* `user_identity` - The user id of the edge function.


## Import
CdnEdgeFunction can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_edge_function.default resource_id
```

