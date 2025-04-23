---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_object_group"
sidebar_current: "docs-byteplus-resource-cloud_monitor_object_group"
description: |-
  Provides a resource to manage cloud monitor object group
---
# byteplus_cloud_monitor_object_group
Provides a resource to manage cloud monitor object group
## Example Usage
```hcl
resource "byteplus_eip_address" "foo" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "BGP"
  name         = "acc-eip"
  description  = "acc-test"
  project_name = "default"
}

resource "byteplus_volume" "foo" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 20
  zone_id            = "ap-southeast-1a"
  volume_charge_type = "PostPaid"
}

resource "byteplus_cloud_monitor_object_group" "foo" {
  name = "acc_test_object_group"
  objects {
    namespace = "VCM_EIP"
    region    = ["ap-southeast-1"]
    dimensions {
      key   = "ResourceID"
      value = [byteplus_eip_address.foo.id]
    }
  }
  objects {
    namespace = "VCM_EBS"
    region    = ["ap-southeast-1"]
    dimensions {
      key   = "ResourceID"
      value = [byteplus_volume.foo.id]
    }
  }
}
```
## Argument Reference
The following arguments are supported:
* `name` - (Required) The name of resource group.

Can only contain Chinese, English, or underscores
The length is limited to 1-64 characters.
* `objects` - (Required) Need to group the list of cloud product resources, the maximum length of the list is 100.

The `dimensions` object supports the following:

* `key` - (Required) Key for retrieving metrics.
* `value` - (Required) Value corresponding to the Key.

The `objects` object supports the following:

* `dimensions` - (Required) Collection of cloud product resource IDs.
* `namespace` - (Required) The product space to which the cloud product belongs in cloud monitoring.
* `region` - (Required) Availability zone associated with the cloud product under the current resource. Only one region id can be specified currently.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CloudMonitorObjectGroup can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_object_group.default resource_id
```

