resource "byteplus_cr_registry" "foo" {
  name    = "acc-test-cr-registry"
  project = "default"
}

resource "byteplus_cr_endpoint" "foo" {
  registry = byteplus_cr_registry.foo.id
  enabled  = true
}

resource "byteplus_cr_endpoint_acl_policy" "foo" {
  registry    = byteplus_cr_endpoint.foo.registry
  type        = "Public"
  entry       = "192.168.0.${count.index}"
  description = "test-${count.index}"
  count       = 3
}
