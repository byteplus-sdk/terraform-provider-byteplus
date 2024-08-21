---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_edge_function_associate"
sidebar_current: "docs-byteplus-resource-cdn_edge_function_associate"
description: |-
  Provides a resource to manage cdn edge function associate
---
# byteplus_cdn_edge_function_associate
Provides a resource to manage cdn edge function associate
## Example Usage
```hcl
resource "byteplus_cdn_edge_function" "foo" {
  name         = "acc-test-function"
  remark       = "tf-test"
  project_name = "default"
  source_code  = base64encode("hello world")
  envs {
    key   = "k1"
    value = "v1"
  }
  canary_countries = ["China", "Japan", "United Kingdom"]
}

resource "byteplus_cdn_edge_function_associate" "foo" {
  function_id = byteplus_cdn_edge_function.foo.id
  domain      = "tf.com"
}
```
## Argument Reference
The following arguments are supported:
* `domain` - (Required, ForceNew) The domain name which you wish to bind with the function.
* `function_id` - (Required, ForceNew) The id of the function for which you want to bind to domain.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CdnEdgeFunctionAssociate can be imported using the function_id:domain, e.g.
```
$ terraform import byteplus_cdn_edge_function_associate.default function_id:domain
```

