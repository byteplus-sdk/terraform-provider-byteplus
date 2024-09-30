resource "byteplus_iam_role" "role" {
  role_name             = "TerraformTestRole"
  display_name          = "terraform role"
  trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
  description           = "created by terraform"
  max_session_duration  = 43200
}

resource "byteplus_iam_policy" "policy" {
  policy_name     = "TerraformResourceTest1"
  description     = "created by terraform 1"
  policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "byteplus_iam_role_policy_attachment" "foo" {
  role_name   = byteplus_iam_role.role.id
  policy_name = byteplus_iam_policy.policy.id
  policy_type = byteplus_iam_policy.policy.policy_type
}