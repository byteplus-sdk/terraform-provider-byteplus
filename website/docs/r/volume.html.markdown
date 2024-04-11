---
subcategory: "EBS"
layout: "byteplus"
page_title: "Byteplus: byteplus_volume"
sidebar_current: "docs-byteplus-resource-volume"
description: |-
  Provides a resource to manage volume
---
# byteplus_volume
Provides a resource to manage volume
## Notice
When Destroy this resource,If the resource charge type is PrePaid,Please unsubscribe the resource 
in  [BytePlus Console](https://console.byteplus.com/home),when complete console operation,yon can
use 'terraform state rm ${resourceId}' to remove.
## Example Usage
```hcl
data "byteplus_zones" "foo" {
}

resource "byteplus_volume" "PostVolume" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 40
  zone_id            = data.byteplus_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
}
```
## Argument Reference
The following arguments are supported:
* `kind` - (Required, ForceNew) The kind of Volume, the value is `data`.
* `size` - (Required) The size of Volume.
* `volume_name` - (Required) The name of Volume.
* `volume_type` - (Required, ForceNew) The type of Volume, the value is `PTSSD` or `ESSD_PL0` or `ESSD_PL1` or `ESSD_PL2` or `ESSD_FlexPL`.
* `zone_id` - (Required, ForceNew) The id of the Zone.
* `delete_with_instance` - (Optional) Delete Volume with Attached Instance.
* `description` - (Optional) The description of the Volume.
* `instance_id` - (Optional, ForceNew) The ID of the instance to which the created volume is automatically attached. Please note this field needs to ask the system administrator to apply for a whitelist.
When use this field to attach ecs instance, the attached volume cannot be deleted by terraform, please use `terraform state rm byteplus_volume.resource_name` command to remove it from terraform state file and management.
* `project_name` - (Optional) The ProjectName of the Volume.
* `volume_charge_type` - (Optional) The charge type of the Volume, the value is `PostPaid` or `PrePaid`. The `PrePaid` volume cannot be detached. Cannot convert `PrePaid` volume to `PostPaid`.Please note that `PrePaid` type needs to ask the system administrator to apply for a whitelist.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `created_at` - Creation time of Volume.
* `status` - Status of Volume.
* `trade_status` - Status of Trade.


## Import
Volume can be imported using the id, e.g.
```
$ terraform import byteplus_volume.default vol-mizl7m1kqccg5smt1bdpijuj
```

