data "byteplus_cdn_certificates" "foo" {
  configured_domain = ["byteplus-demo.byte-test.com"]
  name              = "*.byte-test.com"
  fuzzy_match       = true
  status            = "running"
}
