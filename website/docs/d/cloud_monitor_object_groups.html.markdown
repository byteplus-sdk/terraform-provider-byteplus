---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_object_groups"
sidebar_current: "docs-byteplus-datasource-cloud_monitor_object_groups"
description: |-
  Use this data source to query detailed information of cloud monitor object groups
---
# byteplus_cloud_monitor_object_groups
Use this data source to query detailed information of cloud monitor object groups
## Example Usage
```hcl
data "byteplus_cloud_monitor_object_groups" "foo" {
  ids = ["189968995094908****", "189969233297575****"]
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of cloud monitor object group ids.
* `name_regex` - (Optional) A Name Regex of Resource.
* `name` - (Optional) The keyword of the object group names. Fuzzy match is supported.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `object_groups` - The collection of query.
    * `alert_template_id` - The alarm template ID associated with the resource group.
    * `alert_template_name` - The alarm template name associated with the resource group.
    * `created_at` - The creation time of the resource group.
    * `id` - Resource group ID.
    * `name` - Resource group name.
    * `objects` - List of cloud product resources under the resource group.
        * `dimensions` - Collection of cloud product resource IDs.
            * `key` - Key for retrieving metrics.
            * `value` - Value corresponding to the Key.
        * `id` - Resource grouping ID.
        * `namespace` - The product space to which the cloud product belongs in cloud monitoring.
        * `region` - Availability zone associated with the cloud product under the current resource.
        * `type` - Type of resource collection.
    * `updated_at` - The update time of the resource group.
* `total_count` - The total count of query.


