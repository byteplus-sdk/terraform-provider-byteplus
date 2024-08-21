---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_edge_function_publishes"
sidebar_current: "docs-byteplus-datasource-cdn_edge_function_publishes"
description: |-
  Use this data source to query detailed information of cdn edge function publishes
---
# byteplus_cdn_edge_function_publishes
Use this data source to query detailed information of cdn edge function publishes
## Example Usage
```hcl
data "byteplus_cdn_edge_function_publishes" "foo" {
  function_id = "8f06f8db8d6b4bcdb979db68273f****"
}
```
## Argument Reference
The following arguments are supported:
* `function_id` - (Required) The id of the function.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `tickets` - The collection of query.
    * `content` - The content of the release record.
    * `create_time` - The create time of the release record. Displayed in UNIX timestamp format.
    * `creator` - The creator of the release record.
    * `description` - The description of the release record.
    * `function_id` - The function id.
    * `ticket_id` - The release record id.
    * `update_time` - The update time of the release record. Displayed in UNIX timestamp format.
* `total_count` - The total count of query.


