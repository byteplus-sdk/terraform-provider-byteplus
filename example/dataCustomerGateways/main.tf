resource "byteplus_customer_gateway" "foo" {
  ip_address            = "192.0.1.3"
  customer_gateway_name = "acc-test"
  description           = "acc-test"
  project_name          = "default"
}
data "byteplus_customer_gateways" "foo" {
  ids = [byteplus_customer_gateway.foo.id]
}