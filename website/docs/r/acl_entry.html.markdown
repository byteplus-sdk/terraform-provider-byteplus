---
subcategory: "CLB"
layout: "byteplus"
page_title: "Byteplus: byteplus_acl_entry"
sidebar_current: "docs-byteplus-resource-acl_entry"
description: |-
  Provides a resource to manage acl entry
---
# byteplus_acl_entry
Provides a resource to manage acl entry
## Example Usage
```hcl
resource "byteplus_acl" "foo" {
  acl_name    = "tf-test-3"
  description = "tf-test"
}

resource "byteplus_acl_entry" "foo" {
  acl_id      = byteplus_acl.foo.id
  description = "tf acl entry desc demo"
  entry       = "192.2.2.1/32"
}
```
## Argument Reference
The following arguments are supported:
* `acl_id` - (Required, ForceNew) The ID of Acl.
* `entry` - (Required, ForceNew) The content of the AclEntry.
* `description` - (Optional, ForceNew) The description of the AclEntry.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
AclEntry can be imported using the id, e.g.
```
$ terraform import byteplus_acl_entry.default ID is a string concatenated with colons(AclId:Entry)
```

