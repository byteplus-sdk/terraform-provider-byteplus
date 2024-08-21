---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_edge_functions"
sidebar_current: "docs-byteplus-datasource-cdn_edge_functions"
description: |-
  Use this data source to query detailed information of cdn edge functions
---
# byteplus_cdn_edge_functions
Use this data source to query detailed information of cdn edge functions
## Example Usage
```hcl
data "byteplus_cdn_edge_functions" "foo" {
  status = 100
}
```
## Argument Reference
The following arguments are supported:
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.
* `status` - (Optional) The status of the function. 
100: running. 
400: unassociated. 
500: configuring.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `edge_functions` - The collection of query.
    * `account_identity` - The account id of the edge function.
    * `continent_cluster` - The canary cluster info of the edge function.
        * `cluster_type` - The type of the cluster.
        * `continent` - The continent where the cluster is located.
        * `country` - The country where the cluster is located.
        * `traffics` - The versions of the function deployed on this cluster.
    * `create_time` - The create time of the edge function. Displayed in UNIX timestamp format.
    * `creator` - The creator of the edge function.
    * `domain` - The domain name bound to the edge function.
    * `envs` - The environment variables of the edge function.
        * `key` - The key of the environment variable.
        * `value` - The value of the environment variable.
    * `function_id` - The id of the edge function.
    * `id` - The id of the edge function.
    * `name` - The name of the edge function.
    * `project_name` - The name of the project to which the edge function belongs.
    * `remark` - The remark of the edge function.
    * `source_code` - The latest code content of the edge function. The code is transformed into a Base64-encoded format.
    * `status` - The status of the edge function.
    * `update_time` - The update time of the edge function. Displayed in UNIX timestamp format.
    * `user_identity` - The user id of the edge function.
* `total_count` - The total count of query.


