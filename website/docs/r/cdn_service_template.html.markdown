---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_service_template"
sidebar_current: "docs-byteplus-resource-cdn_service_template"
description: |-
  Provides a resource to manage cdn service template
---
# byteplus_cdn_service_template
Provides a resource to manage cdn service template
## Example Usage
```hcl
resource "byteplus_cdn_service_template" "foo" {
  title   = "tf-test2"
  message = "test2"
  project = "test"
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
```
## Argument Reference
The following arguments are supported:
* `service_template_config` - (Required) The service template configuration. Please convert the configuration module structure into json and pass it into a string. You must specify the Origin module. The OriginProtocol parameter, and other domain configuration modules are optional. For detailed parameter introduction, please refer to `https://docs.byteplus.com/en/docs/byteplus-cdn/reference-updateservicetemplate`.
* `title` - (Required) Indicates the name of the encryption policy you want to create. The name must not exceed 100 characters.
* `lock_template` - (Optional) Whether to lock the template. If you set this field to true, then the template will be locked. Please note that the template cannot be modified or unlocked after it is locked. When you want to use this template to create a domain name, please lock the template first. The default value is false.
* `message` - (Optional) Indicates the description of the encryption policy, which must not exceed 120 characters.
* `project` - (Optional) Indicates the project to which this encryption policy belongs. The default value of the parameter is default, indicating the Default project.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CdnServiceTemplate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_service_template.default resource_id
```

