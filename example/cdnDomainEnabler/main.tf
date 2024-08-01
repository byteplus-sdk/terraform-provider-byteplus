resource "byteplus_cdn_service_template" "foo" {
  title = "tf-test2"
  message = "test2"
  project = ""
  # 是否发布模版
  lock_template = true
  service_template_config = jsonencode(
    {
      OriginIpv6 = "followclient"
      ConditionalOrigin = {
        OriginRules = []
      }
      Origin = [{
        OriginAction = {
          OriginLines = [
            {
              Address = "10.10.10.10"
              HttpPort = "80"
              HttpsPort = "443"
              InstanceType = "ip"
              OriginType = "primary"
              Weight = "1"
            }
          ]
        }
      }]
      OriginHost = ""
      OriginProtocol = "http"
      OriginHost = ""
    }
  )
}

resource "byteplus_cdn_cipher_template" "foo" {
  title = "tf-test"
  message = "test for tf"
  project = ""
  lock_template = true
  https {
    disable_http = false
    forced_redirect {
      enable_forced_redirect = true
      status_code = "302"
    }
    http2 = false
    ocsp = false
    tls_version = ["tlsv1.1", "tlsv1.2", "tlsv1.3"]
    hsts {
      subdomain = "exclude"
      switch = true
      ttl = 3600
    }
  }
  #  http_forced_redirect {
  #    enable_forced_redirect = false
  #    status_code = "301"
  #  }
  quic {
    switch = false
  }
}

resource "byteplus_cdn_domain" "foo" {
  domain = "tf.byte-test.com"
  service_template_id = byteplus_cdn_service_template.foo.id
  https_switch = "on"
  cert_id = "cert-"
  cipher_template_id = byteplus_cdn_cipher_template.foo.id
  project = ""
  service_region = "outside_chinese_mainland"
}

resource "byteplus_cdn_domain_enabler" "foo" {
  domain = byteplus_cdn_domain.foo.id
}