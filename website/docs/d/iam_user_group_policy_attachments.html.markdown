---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_user_group_policy_attachments"
sidebar_current: "docs-byteplus-datasource-iam_user_group_policy_attachments"
description: |-
  Use this data source to query detailed information of iam user group policy attachments
---
# byteplus_iam_user_group_policy_attachments
Use this data source to query detailed information of iam user group policy attachments
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

data "byteplus_iam_user_group_policy_attachments" "foo" {
  user_group_name = byteplus_iam_user_group_policy_attachment.foo.user_group_name
}
```
## Argument Reference
The following arguments are supported:
* `user_group_name` - (Required) A name of user group.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `policies` - The collection of query.
    * `attach_date` - Attached time.
    * `description` - The description.
    * `policy_name` - Name of the policy.
    * `policy_trn` - Resource name of the strategy.
    * `policy_type` - The type of the policy.
* `total_count` - The total count of query.


