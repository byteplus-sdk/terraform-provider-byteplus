resource "byteplus_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc test"
  display_name = "name"
}
data "byteplus_iam_users" "foo" {
  user_names = [byteplus_iam_user.foo.user_name]
}