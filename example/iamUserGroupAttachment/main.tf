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