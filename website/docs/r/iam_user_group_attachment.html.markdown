---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_user_group_attachment"
sidebar_current: "docs-byteplus-resource-iam_user_group_attachment"
description: |-
  Provides a resource to manage iam user group attachment
---
# byteplus_iam_user_group_attachment
Provides a resource to manage iam user group attachment
## Example Usage
```hcl
resource "byteplus_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc test"
  display_name = "name"
}

resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description     = "acc-test"
  display_name    = "acctest"
}

resource "byteplus_iam_user_group_attachment" "foo" {
  user_group_name = byteplus_iam_user_group.foo.user_group_name
  user_name       = byteplus_iam_user.foo.user_name
}
```
## Argument Reference
The following arguments are supported:
* `user_group_name` - (Required, ForceNew) The name of the user group.
* `user_name` - (Required, ForceNew) The name of the user.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
IamUserGroupAttachment can be imported using the id, e.g.
```
$ terraform import byteplus_iam_user_group_attachment.default user_group_id:user_id
```

