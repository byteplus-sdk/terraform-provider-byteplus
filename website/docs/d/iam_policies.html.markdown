---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_policies"
sidebar_current: "docs-byteplus-datasource-iam_policies"
description: |-
  Use this data source to query detailed information of iam policies
---
# byteplus_iam_policies
Use this data source to query detailed information of iam policies
## Example Usage
```hcl
resource "byteplus_iam_policy" "foo" {
  policy_name     = "acc-test-policy"
  description     = "acc-test"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

data "byteplus_iam_policies" "foo" {
  query = byteplus_iam_policy.foo.description
  #  user_name = "user-test"
  #  role_name = "test-role"
}
```
## Argument Reference
The following arguments are supported:
* `name_regex` - (Optional) A Name Regex of Policy.
* `output_file` - (Optional) File name where to save data source results.
* `query` - (Optional) Query policies, support policy name or description.
* `role_name` - (Optional) The name of the IAM role.
* `scope` - (Optional) The scope of the Policy.
* `status` - (Optional) The status of policy.
* `user_name` - (Optional) The name of the IAM user.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `policies` - The collection of Policy query.
    * `create_date` - The create time of the Policy.
    * `description` - The description of the Policy.
    * `id` - The ID of the Policy.
    * `policy_document` - The document of the Policy.
    * `policy_name` - The name of the Policy.
    * `policy_trn` - The resource name of the Policy.
    * `policy_type` - The type of the Policy.
    * `role_attach_date` - The role attach time of the Policy.The data show only query with role_name.
    * `role_name` - The name of the IAM role.The data show only query with role_name.
    * `update_date` - The update time of the Policy.
    * `user_attach_date` - The user attach time of the Policy.The data show only query with user_name.
    * `user_name` - The name of the IAM user.The data show only query with user_name.
* `total_count` - The total count of Policy query.


