---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_user_group"
sidebar_current: "docs-byteplus-resource-iam_user_group"
description: |-
  Provides a resource to manage iam user group
---
# byteplus_iam_user_group
Provides a resource to manage iam user group
## Example Usage
```hcl
resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description     = "acc-test"
  display_name    = "acctest"
}
```
## Argument Reference
The following arguments are supported:
* `user_group_name` - (Required, ForceNew) The name of the user group.
* `description` - (Optional) The description of the user group.
* `display_name` - (Optional) The display name of the user group.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
IamUserGroup can be imported using the id, e.g.
```
$ terraform import byteplus_iam_user_group.default user_group_name
```

