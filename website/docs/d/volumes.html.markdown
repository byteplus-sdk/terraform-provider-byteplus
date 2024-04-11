---
subcategory: "EBS"
layout: "byteplus"
page_title: "Byteplus: byteplus_volumes"
sidebar_current: "docs-byteplus-datasource-volumes"
description: |-
  Use this data source to query detailed information of volumes
---
# byteplus_volumes
Use this data source to query detailed information of volumes
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}

resource "byteplus_volume" "foo" {
  volume_name        = "acc-test-volume-${count.index}"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 60
  zone_id            = data.byteplus_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
  count              = 3
}

data "byteplus_volumes" "foo" {
  ids = byteplus_volume.foo[*].id
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of Volume IDs.
* `instance_id` - (Optional) The Id of instance.
* `kind` - (Optional) The Kind of Volume.
* `name_regex` - (Optional) A Name Regex of Volume.
* `output_file` - (Optional) File name where to save data source results.
* `volume_name` - (Optional) The name of Volume.
* `volume_status` - (Optional) The Status of Volume, the value can be `available` or `attaching` or `attached` or `detaching` or `creating` or `deleting` or `error` or `extending`.
* `volume_type` - (Optional) The type of Volume.
* `zone_id` - (Optional) The Id of Zone.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of Volume query.
* `volumes` - The collection of Volume query.


