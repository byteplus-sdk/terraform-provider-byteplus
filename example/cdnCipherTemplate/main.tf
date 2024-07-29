resource "byteplus_cdn_cipher_template" "foo" {
  title = "tf-test"
  message = "test for tf"
  project = "test"
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