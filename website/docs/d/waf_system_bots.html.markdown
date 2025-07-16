---
subcategory: "WAF"
layout: "byteplus"
page_title: "Byteplus: byteplus_waf_system_bots"
sidebar_current: "docs-byteplus-datasource-waf_system_bots"
description: |-
  Use this data source to query detailed information of waf system bots
---
# byteplus_waf_system_bots
Use this data source to query detailed information of waf system bots
## Example Usage
```hcl
data "byteplus_waf_system_bots" "foo" {
  host = "www.tf-test.com"
}
```
## Argument Reference
The following arguments are supported:
* `host` - (Required) Domain name information.
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `data` - Host the Bot configuration information.
    * `action` - The execution action of the Bot.
    * `bot_type` - The name of Bot.
    * `description` - The description of Bot.
    * `enable` - Whether to enable Bot.
    * `rule_tag` - The rule ID corresponding to Bot.
* `total_count` - The total count of query.


