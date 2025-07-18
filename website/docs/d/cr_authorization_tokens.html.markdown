---
subcategory: "CR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cr_authorization_tokens"
sidebar_current: "docs-byteplus-datasource-cr_authorization_tokens"
description: |-
  Use this data source to query detailed information of cr authorization tokens
---
# byteplus_cr_authorization_tokens
Use this data source to query detailed information of cr authorization tokens
## Example Usage
```hcl
data "byteplus_cr_authorization_tokens" "foo" {
  registry = "tf-1"
}
```
## Argument Reference
The following arguments are supported:
* `registry` - (Required) The cr instance name want to query.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `tokens` - The collection of users.
    * `expire_time` - The expiration time of the temporary access token.
    * `token` - The Temporary access token.
    * `username` - The username for login repository instance.
* `total_count` - The total count of instance query.


