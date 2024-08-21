---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_cipher_template"
sidebar_current: "docs-byteplus-resource-cdn_cipher_template"
description: |-
  Provides a resource to manage cdn cipher template
---
# byteplus_cdn_cipher_template
Provides a resource to manage cdn cipher template
## Example Usage
```hcl
resource "byteplus_cdn_cipher_template" "foo" {
  title   = "tf-test"
  message = "test for tf"
  project = "test"
  https {
    disable_http = false
    forced_redirect {
      enable_forced_redirect = true
      status_code            = "302"
    }
    http2       = false
    ocsp        = false
    tls_version = ["tlsv1.1", "tlsv1.2", "tlsv1.3"]
    hsts {
      subdomain = "exclude"
      switch    = true
      ttl       = 3600
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
```
## Argument Reference
The following arguments are supported:
* `title` - (Required) Indicates the name of the encryption policy you want to create. The name must not exceed 100 characters.
* `http_forced_redirect` - (Optional) Indicates the configuration module for the forced redirection from HTTPS to HTTP. This feature is disabled by default.
* `https` - (Optional) Indicates the configuration module for the HTTPS encryption service.
* `lock_template` - (Optional) Whether to lock the template. If you set this field to true, then the template will be locked. Please note that the template cannot be modified or unlocked after it is locked. When you want to use this template to create a domain name, please lock the template first. The default value is false.
* `message` - (Optional) Indicates the description of the encryption policy, which must not exceed 120 characters.
* `project` - (Optional) Indicates the project to which this encryption policy belongs. The default value of the parameter is default, indicating the Default project.
* `quic` - (Optional) Indicates the QUIC configuration module. This feature is disabled by default.

The `forced_redirect` object supports the following:

* `enable_forced_redirect` - (Required) Indicates the switch for the Forced Redirect configuration. This parameter can take the following values: true: Indicates to enable Forced Redirect. false: Indicates to disable Forced Redirect.
* `status_code` - (Required) Indicates the status code returned to the client by the CDN when forced redirect occurs. This parameter can take the following values: 301: Indicates that the returned status code is 301. 302: Indicates that the returned status code is 302. The default value for this parameter is 301.

The `hsts` object supports the following:

* `subdomain` - (Optional) Indicates whether the HSTS configuration should also be applied to the subdomains of the domain name. This parameter can take the following values: include: Indicates that HSTS settings apply to subdomains. exclude: Indicates that HSTS settings do not apply to subdomains. The default value for this parameter is exclude.
* `switch` - (Optional) Indicates whether to enable HSTS. This parameter can take the following values: true: Indicates to enable HSTS. false: Indicates to disable HSTS. The default value for this parameter is false.
* `ttl` - (Optional) Indicates the expiration time for the Strict-Transport-Security response header in the browser cache, in seconds. If Switch is true, this parameter is required. The value range for this parameter is 0 - 31,536,000 seconds, where 31,536,000 seconds represents 365 days. If the value of this parameter is 0, it is equivalent to disabling the HSTS settings.

The `http_forced_redirect` object supports the following:

* `enable_forced_redirect` - (Required) Indicates whether to enable the forced redirection from HTTPS. This parameter can take the following values: true: Indicates to enable the forced redirection from HTTPS. Once enabled, the content delivery network will respond with StatusCode to inform the browser to send an HTTPS request when it receives an HTTP request from a user. false: Indicates to disable the forced redirection from HTTPS.
* `status_code` - (Required) Indicates the status code returned by the content delivery network when forced redirection from HTTPS occurs. The default value for this parameter is 301.

The `https` object supports the following:

* `disable_http` - (Optional) Indicates whether the CDN accepts HTTP user requests. This parameter can take the following values: true: Indicates that it does not accept. If an HTTP request is received, the CDN will reject the request. false: Indicates that it accepts. The default value for this parameter is false.
* `forced_redirect` - (Optional) Indicates the configuration for the mandatory redirection from HTTP to HTTPS. This feature is disabled by default.
* `hsts` - (Optional) Indicates the HSTS (HTTP Strict Transport Security) configuration module. This feature is disabled by default.
* `http2` - (Optional) Indicates the switch for HTTP/2 configuration. This parameter can take the following values: true: Indicates to enable HTTP/2. false: Indicates to disable HTTP/2. The default value for this parameter is true.
* `ocsp` - (Optional) Indicates whether to enable OCSP Stapling. This parameter can take the following values: true: Indicates to enable OCSP Stapling. false: Indicates to disable OCSP Stapling. The default value for this parameter is false.
* `tls_version` - (Optional) Indicates a list that specifies the TLS versions supported by the domain name. This parameter can take the following values: tlsv1.0: Indicates TLS 1.0. tlsv1.1: Indicates TLS 1.1. tlsv1.2: Indicates TLS 1.2. tlsv1.3: Indicates TLS 1.3. The default value for this parameter is ["tlsv1.1", "tlsv1.2", "tlsv1.3"].

The `quic` object supports the following:

* `switch` - (Required) Indicates whether to enable QUIC. This parameter can take the following values: true: Indicates to enable QUIC. false: Indicates to disable QUIC.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CdnCipherTemplate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_cipher_template.default resource_id
```

