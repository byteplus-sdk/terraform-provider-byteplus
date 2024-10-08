---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_login_profile"
sidebar_current: "docs-byteplus-resource-iam_login_profile"
description: |-
  Provides a resource to manage iam login profile
---
# byteplus_iam_login_profile
Provides a resource to manage iam login profile
## Example Usage
```hcl
resource "byteplus_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc-test"
  display_name = "name"
}

resource "byteplus_iam_login_profile" "foo" {
  user_name               = byteplus_iam_user.foo.user_name
  password                = "93f0cb0614Aab12"
  login_allowed           = true
  password_reset_required = false
}
```
## Argument Reference
The following arguments are supported:
* `password` - (Required) The password.
* `user_name` - (Required, ForceNew) The user name.
* `login_allowed` - (Optional) The flag of login allowed.
* `password_reset_required` - (Optional) Is required reset password when next time login in.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Login profile can be imported using the UserName, e.g.
```
$ terraform import byteplus_iam_login_profile.default user_name
```

