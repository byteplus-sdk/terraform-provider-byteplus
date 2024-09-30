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