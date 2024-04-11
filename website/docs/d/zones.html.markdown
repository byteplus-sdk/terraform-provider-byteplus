---
subcategory: "ECS"
layout: "byteplus"
page_title: "Byteplus: byteplus_zones"
sidebar_current: "docs-byteplus-datasource-zones"
description: |-
  Use this data source to query detailed information of zones
---
# byteplus_zones
Use this data source to query detailed information of zones
## Example Usage
```hcl
data "byteplus_zones" "default" {
  ids = ["cn-beijing-a"]
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of zone ids.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of zone query.
* `zones` - The collection of zone query.
    * `id` - The id of the zone.
    * `zone_id` - The id of the zone.


