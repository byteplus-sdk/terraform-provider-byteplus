resource "byteplus_waf_cdn_domain" "foo" {
  domain = "tf-test.com"
  project_follow = 1
  tls_enable =  1
  tls_fields_config {
    headers_config {
      enable = 1
    }
  }
  project_name = "default"
}
