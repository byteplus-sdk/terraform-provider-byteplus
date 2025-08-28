resource "byteplus_waf_cdn_domain" "foo" {
  domain = "xxxxxx.com"
  project_follow = 1
  tls_enable =  1
  tls_fields_config {
    headers_config {
      enable = 1
    }
  }
  auto_cc_enable = 0
  cc_enable = 0
  project_name = "default"
}
