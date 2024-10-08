---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_user_group_policy_attachment"
sidebar_current: "docs-byteplus-resource-iam_user_group_policy_attachment"
description: |-
  Provides a resource to manage iam user group policy attachment
---
# byteplus_iam_user_group_policy_attachment
Provides a resource to manage iam user group policy attachment
## Example Usage
```hcl
resource "byteplus_iam_policy" "foo" {
  policy_name     = "acc-test-policy"
  description     = "acc-test"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description     = "acc-test"
  display_name    = "acc-test"
}

resource "byteplus_iam_user_group_policy_attachment" "foo" {
  policy_name     = byteplus_iam_policy.foo.policy_name
  policy_type     = "Custom"
  user_group_name = byteplus_iam_user_group.foo.user_group_name
}
```
## Argument Reference
The following arguments are supported:
* `policy_name` - (Required, ForceNew) The policy name.
* `policy_type` - (Required, ForceNew) Strategy types, System strategy, Custom strategy.
* `user_group_name` - (Required, ForceNew) The user group name.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
IamUserGroupPolicyAttachment can be imported using the user group name and policy name, e.g.
```
$ terraform import byteplus_iam_user_group_policy_attachment.default userGroupName:policyName
```

