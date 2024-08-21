---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_kvs"
sidebar_current: "docs-byteplus-datasource-cdn_kvs"
description: |-
  Use this data source to query detailed information of cdn kvs
---
# byteplus_cdn_kvs
Use this data source to query detailed information of cdn kvs
## Example Usage
```hcl
data "byteplus_cdn_kvs" "foo" {
  namespace_id = "4723722642589338688"
  namespace    = "acc-test-kv-namespace"
}
```
## Argument Reference
The following arguments are supported:
* `namespace_id` - (Required) The id of the kv namespace.
* `namespace` - (Required) The name of the kv namespace.
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `namespace_keys` - The collection of query.
    * `create_time` - The creation time of the kv namespace key. Displayed in UNIX timestamp format.
    * `ddl` - Data expiration time. After the data expires, the Value in the Key will be inaccessible.
Displayed in UNIX timestamp format.
0: Permanent storage.
    * `key_status` - The status of the kv namespace key.
    * `key` - The key of the kv namespace key.
    * `namespace_id` - The id of the kv namespace key.
    * `namespace` - The name of the kv namespace key.
    * `update_time` - The update time of the kv namespace key. Displayed in UNIX timestamp format.
    * `value` - The value of the kv namespace key.
* `total_count` - The total count of query.


