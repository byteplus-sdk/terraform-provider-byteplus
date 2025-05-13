---
subcategory: "ALB"
layout: "byteplus"
page_title: "Byteplus: byteplus_alb_zones"
sidebar_current: "docs-byteplus-datasource-alb_zones"
description: |-
  Use this data source to query detailed information of alb zones
---
# byteplus_alb_zones
Use this data source to query detailed information of alb zones
## Example Usage
```hcl
data "byteplus_alb_zones" "default" {
}
```
## Argument Reference
The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of zone query.
* `zones` - The collection of zone query.
    * `id` - The id of the zone.
    * `zone_id` - The id of the zone.


