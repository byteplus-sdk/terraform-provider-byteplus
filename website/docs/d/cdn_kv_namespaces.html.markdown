---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_kv_namespaces"
sidebar_current: "docs-byteplus-datasource-cdn_kv_namespaces"
description: |-
  Use this data source to query detailed information of cdn kv namespaces
---
# byteplus_cdn_kv_namespaces
Use this data source to query detailed information of cdn kv namespaces
## Example Usage
```hcl
data "byteplus_cdn_kv_namespaces" "foo" {

}
```
## Argument Reference
The following arguments are supported:
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `kv_namespaces` - The collection of query.
    * `create_time` - The creation time of the kv namespace. Displayed in UNIX timestamp format.
    * `creator` - The creator of the kv namespace.
    * `description` - The description of the kv namespace.
    * `id` - The id of the kv namespace.
    * `namespace_id` - The id of the kv namespace.
    * `namespace` - The name of the kv namespace.
    * `project_name` - The name of the project to which the namespace belongs, defaulting to `default`.
    * `update_time` - The update time of the kv namespace. Displayed in UNIX timestamp format.
* `total_count` - The total count of query.


