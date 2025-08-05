resource "byteplus_mongodb_allow_list" "foo" {
  allow_list_name = "acc-test"
  allow_list_desc = "acc-test"
  allow_list_type = "IPv4"
  allow_list      = "10.1.1.3,10.2.3.0/24,10.1.1.1"
}

data "byteplus_mongodb_allow_lists" "foo" {
  allow_list_ids = [byteplus_mongodb_allow_list.foo.id]
  region_id      = "cn-beijing"
}