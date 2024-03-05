resource "byteplus_eip_address" "foo" {
  billing_type = "PostPaidByTraffic"
}
data "byteplus_eip_addresses" "foo" {
  ids = [byteplus_eip_address.foo.id]
}