resource "byteplus_waf_cdn_domain" "foo" {
  domain_name = "www.tf-test.com"
  project_follow = 1
  tls_enable =  1
  tls_fields_config {
    headers_config {
      enable = 1
    }
  }
}
