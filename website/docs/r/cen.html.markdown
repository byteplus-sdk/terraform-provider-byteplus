---
subcategory: "CEN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cen"
sidebar_current: "docs-byteplus-resource-cen"
description: |-
  Provides a resource to manage cen
---
# byteplus_cen
Provides a resource to manage cen
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
```
## Argument Reference
The following arguments are supported:
* `cen_name` - (Optional) The name of the cen.
* `description` - (Optional) The description of the cen.
* `project_name` - (Optional) The ProjectName of the cen instance.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `account_id` - The account ID of the cen.
* `cen_bandwidth_package_ids` - A list of bandwidth package IDs of the cen.
* `cen_id` - The ID of the cen.
* `creation_time` - The create time of the cen.
* `status` - The status of the cen.
* `update_time` - The update time of the cen.


## Import
Cen can be imported using the id, e.g.
```
$ terraform import byteplus_cen.default cen-7qthudw0ll6jmc****
```

