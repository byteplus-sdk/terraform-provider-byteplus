---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_kv"
sidebar_current: "docs-byteplus-resource-cdn_kv"
description: |-
  Provides a resource to manage cdn kv
---
# byteplus_cdn_kv
Provides a resource to manage cdn kv
## Example Usage
```hcl
resource "byteplus_cdn_kv_namespace" "foo" {
  namespace    = "acc-test-kv-namespace"
  description  = "tf-test"
  project_name = "default"
}

resource "byteplus_cdn_kv" "foo" {
  namespace_id = byteplus_cdn_kv_namespace.foo.id
  namespace    = byteplus_cdn_kv_namespace.foo.namespace
  key          = "acc-test-key"
  value        = base64encode("tf-test")
  ttl          = 1000
}
```
## Argument Reference
The following arguments are supported:
* `key` - (Required, ForceNew) The key of the kv namespace.
* `namespace_id` - (Required, ForceNew) The id of the kv namespace.
* `namespace` - (Required, ForceNew) The name of the kv namespace.
* `value` - (Required) The value of the kv namespace key. Single Value upload data does not exceed 128KB. This field must be encrypted with base64.
* `ttl` - (Optional) Set the data storage time. Unit: second. After the data expires, the Value in the Key will be inaccessible.
If this parameter is not specified or the parameter value is 0, it is stored permanently by default.
The storage time cannot be less than 60s.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `create_time` - The creation time of the kv namespace key. Displayed in UNIX timestamp format.
* `ddl` - Data expiration time. After the data expires, the Value in the Key will be inaccessible.
Displayed in UNIX timestamp format.
0: Permanent storage.
* `key_status` - The status of the kv namespace key.
* `update_time` - The update time of the kv namespace key. Displayed in UNIX timestamp format.


## Import
CdnKv can be imported using the namespace_id:namespace:key, e.g.
```
$ terraform import byteplus_cdn_kv.default namespace_id:namespace:key
```

