data "byteplus_cr_namespaces" "foo" {
  registry = "tf-1"
  names    = ["namespace-*"]
}