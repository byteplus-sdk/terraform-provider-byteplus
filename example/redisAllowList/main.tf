resource "byteplus_redis_allow_list" "foo" {
  allow_list_name = "acc_test_tf_allowlist_create"
  allow_list      = ["0.0.0.0/0", "192.168.0.0/24", "192.168.1.1", "192.168.2.22"]
  allow_list_desc = "acctftestallowlist"
}
