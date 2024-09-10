---
subcategory: "CEN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cen_bandwidth_package_associate"
sidebar_current: "docs-byteplus-resource-cen_bandwidth_package_associate"
description: |-
  Provides a resource to manage cen bandwidth package associate
---
# byteplus_cen_bandwidth_package_associate
Provides a resource to manage cen bandwidth package associate
## Example Usage
```hcl
resource "byteplus_cen" "foo" {
  cen_name     = "acc-test-cen"
  description  = "acc-test"
  project_name = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_cen_bandwidth_package" "foo" {
  local_geographic_region_set_id = "China"
  peer_geographic_region_set_id  = "China"
  bandwidth                      = 2
  cen_bandwidth_package_name     = "acc-test-cen-bp"
  description                    = "acc-test"
  billing_type                   = "PrePaid"
  period_unit                    = "Month"
  period                         = 1
  project_name                   = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_cen_bandwidth_package_associate" "foo" {
  cen_bandwidth_package_id = byteplus_cen_bandwidth_package.foo.id
  cen_id                   = byteplus_cen.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `cen_bandwidth_package_id` - (Required, ForceNew) The ID of the cen bandwidth package.
* `cen_id` - (Required, ForceNew) The ID of the cen.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Cen bandwidth package associate can be imported using the CenBandwidthPackageId:CenId, e.g.
```
$ terraform import byteplus_cen_bandwidth_package_associate.default cbp-4c2zaavbvh5fx****:cen-7qthudw0ll6jmc****
```

