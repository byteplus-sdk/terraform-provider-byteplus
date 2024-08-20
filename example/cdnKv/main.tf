resource "byteplus_cdn_kv_namespace" "foo" {
  namespace    = "acc-test-kv-namespace"
  description  = "tf-test"
  project_name = "default"
}

resource "byteplus_cdn_kv" "foo" {
  namespace_id = byteplus_cdn_kv_namespace.foo.id
  namespace    = byteplus_cdn_kv_namespace.foo.namespace
  key          = "acc-test-key"
  value        = base64encode("tf-test")
  ttl          = 1000
}
