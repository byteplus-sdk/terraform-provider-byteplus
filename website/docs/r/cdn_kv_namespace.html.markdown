---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_kv_namespace"
sidebar_current: "docs-byteplus-resource-cdn_kv_namespace"
description: |-
  Provides a resource to manage cdn kv namespace
---
# byteplus_cdn_kv_namespace
Provides a resource to manage cdn kv namespace
## Example Usage
```hcl
resource "byteplus_cdn_kv_namespace" "foo" {
  namespace    = "acc-test-kv-namespace"
  description  = "tf-test"
  project_name = "default"
}
```
## Argument Reference
The following arguments are supported:
* `namespace` - (Required) Set a recognizable name for the namespace. The input requirements are as follows:
Length should be between 2 and 64 characters.
It can only contain English letters, numbers, hyphens (-), and underscores (_).
* `description` - (Optional) Set a description for the namespace. The input requirements are as follows:
Any characters are allowed.
The length should not exceed 80 characters.
* `project_name` - (Optional) The name of the project to which the namespace belongs, defaulting to `default`.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CdnKvNamespace can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_kv_namespace.default resource_id
```

