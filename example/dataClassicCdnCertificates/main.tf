resource "byteplus_classic_cdn_certificate" "foo" {
  certificate = ""
  private_key = ""
  desc        = "tf-test"
  source      = "cdn_cert_hosting"
}

data "byteplus_classic_cdn_certificates" "foo" {
  source = byteplus_classic_cdn_certificate.foo.source
}