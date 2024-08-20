resource "byteplus_cdn_kv_namespace" "foo" {
  namespace    = "acc-test-kv-namespace"
  description  = "tf-test"
  project_name = "default"
}
