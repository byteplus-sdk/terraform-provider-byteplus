resource "byteplus_iam_user" "foo" {
  user_name    = "acc-test-user"
  description  = "acc-test"
  display_name = "name"
}

resource "byteplus_iam_access_key" "foo" {
  user_name   = byteplus_iam_user.foo.user_name
  secret_file = "./sk"
  status      = "active"
  #  pgp_key = "keybase:some_person_that_exists"
}
