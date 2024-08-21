resource "byteplus_cdn_certificate" "foo" {
  certificate = "-----BEGIN CERTIFICATE----- *** -----END CERTIFICATE-----"
  private_key = "-----BEGIN PRIVATE KEY----- *** -----END PRIVATE KEY-----"
  desc        = "tf-test"
  repeatable  = true
}
