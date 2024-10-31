---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_domain"
sidebar_current: "docs-byteplus-resource-cdn_domain"
description: |-
  Provides a resource to manage cdn domain
---
# byteplus_cdn_domain
Provides a resource to manage cdn domain
## Example Usage
```hcl
resource "byteplus_cdn_service_template" "foo" {
  title   = "tf-test2"
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
              Address      = "10.10.10.10"
              HttpPort     = "80"
              HttpsPort    = "443"
              InstanceType = "ip"
              OriginType   = "primary"
              Weight       = "1"
            }
          ]
        }
      }]
      OriginHost     = ""
      OriginProtocol = "http"
      OriginHost     = ""
    }
  )
}

resource "byteplus_cdn_cipher_template" "foo" {
  title         = "tf-test"
  message       = "test for tf"
  project       = ""
  lock_template = true
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

resource "byteplus_cdn_domain" "foo" {
  domain              = "tf-test.com"
  service_template_id = byteplus_cdn_service_template.foo.id
  https_switch        = "off"
  cipher_template_id  = byteplus_cdn_cipher_template.foo.id
  project             = ""
  service_region      = "outside_chinese_mainland"
  tags {
    key   = "k1"
    value = "v1"
  }
  tags {
    key   = "k2"
    value = "v2"
  }
}
```
## Argument Reference
The following arguments are supported:
* `domain` - (Required, ForceNew) Indicates a domain name you want to add. The domain name you add must meet all of the following conditions: Length does not exceed 100 characters. Cannot contain uppercase letters. Does not include any of these suffixes: zjgslb.com, yangyi19.com, volcgslb.com, veew-alb-cn1.com, bplgslb.com, bplslb.com, ttgslb.com. When you bind your domain name with a delivery policy, the origin address specified in the policy must not be the same as your domain name.
* `service_template_id` - (Required) Indicates a delivery policy to be bound with this domain name. You can use DescribeTemplates to obtain the ID of the delivery policy you want to bind.
* `cert_id` - (Optional) Indicates the ID of a certificate. This certificate is stored in the BytePlus Certificate Center and will be associated with the domain name. If HTTPSSwitch is on, this parameter is required. Before using this API, you need to grant CDN access to the Certificate Center, then upload your certificate to the BytePlus Certificate Center to obtain the ID of the certificate. It is recommended to authorize CDN access to the Certificate Center using the primary account. You can use ListCertInfo to obtain the ID of the certificate you want to associate. If HTTPSSwitch is off, this parameter does not take effect.
* `cipher_template_id` - (Optional) Indicates an encryption policy to be bound with this domain name. You can use DescribeTemplates to obtain the ID of the encryption policy you want to bind. If this parameter is not specified, it means that the domain name will not be bound to any encryption policy at present.
* `https_switch` - (Optional) Indicates whether to enable "HTTPS Encryption Service" for this domain name. This parameter can take the following values: on: Indicates to enable this service. off: Indicates not to enable this service. The default value of this parameter is off.
* `project` - (Optional) Indicates the project to which this domain name belongs, with the default value being default.
* `service_region` - (Optional) Indicates the service region enabled for this domain name. This parameter can take the following values: outside_chinese_mainland: Indicates "Global (excluding Chinese Mainland)". chinese_mainland: Indicates "Chinese Mainland". global: Indicates "Global". The default value of this parameter is outside_chinese_mainland. Note that chinese_mainland or global are not available by default. To make the two service regions available, please submit a ticket. Also, since both regions include Chinese Mainland, you must complete the following additional actions: Perform real-name authentication for your BytePlus account. Perform ICP filing for your domain name.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `status` - Indicates the status of the domain name. This parameter can be: online: Indicates the status is Enabled. offline: Indicates the status is Disabled. configuring: Indicates the status is Configuring.


## Import
CdnDomain can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_domain.default resource_id
```

