resource "byteplus_iam_user" "user" {
  user_name   = "TfTest"
  description = "test"
}

resource "byteplus_iam_policy" "policy" {
  policy_name     = "TerraformResourceTest1"
  description     = "created by terraform 1"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "byteplus_iam_user_policy_attachment" "foo" {
  user_name   = byteplus_iam_user.user.user_name
  policy_name = byteplus_iam_policy.policy.policy_name
  policy_type = byteplus_iam_policy.policy.policy_type
}