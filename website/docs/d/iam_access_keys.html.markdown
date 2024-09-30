---
subcategory: "IAM"
layout: "byteplus"
page_title: "Byteplus: byteplus_iam_access_keys"
sidebar_current: "docs-byteplus-datasource-iam_access_keys"
description: |-
  Use this data source to query detailed information of iam access keys
---
# byteplus_iam_access_keys
Use this data source to query detailed information of iam access keys
## Example Usage
```hcl
data "byteplus_iam_access_keys" "foo" {
}
```
## Argument Reference
The following arguments are supported:
* `name_regex` - (Optional) A Name Regex of IAM.
* `output_file` - (Optional) File name where to save data source results.
* `user_name` - (Optional) The user names.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `access_key_metadata` - The collection of access keys.
    * `access_key_id` - The user access key id.
    * `create_date` - The user access key create date.
    * `status` - The user access key status.
    * `update_date` - The user access key update date.
    * `user_name` - The user name.
* `total_count` - The total count of user query.


