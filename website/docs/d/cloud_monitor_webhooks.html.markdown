---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_webhooks"
sidebar_current: "docs-byteplus-datasource-cloud_monitor_webhooks"
description: |-
  Use this data source to query detailed information of cloud monitor webhooks
---
# byteplus_cloud_monitor_webhooks
Use this data source to query detailed information of cloud monitor webhooks
## Example Usage
```hcl
data "byteplus_cloud_monitor_webhooks" "foo" {
  ids = ["189968992116123****"]
}
```
## Argument Reference
The following arguments are supported:
* `event_rule_id` - (Optional) Event Rule ID.
* `ids` - (Optional) A list of webhook IDs.
* `name_regex` - (Optional) A Name Regex of Resource.
* `name` - (Optional) Webhook name, fuzzy search by name.
* `output_file` - (Optional) File name where to save data source results.
* `rule_id` - (Optional) Alarm strategy ID.
* `type` - (Optional) Type of the webhook.

custom: Custom webhook
wecom: WeChat webhook
lark: Lark webhook
dingtalk: DingTalk webhook.
* `url` - (Optional) The address of the webhook.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `total_count` - The total count of query.
* `webhooks` - The collection of query.
    * `created_at` - The creation time of the webhook.
    * `event_rule_ids` - Event rule IDs.
    * `id` - The id of the webhook.
    * `name` - The name of the webhook.
    * `rule_ids` - Alarm strategy IDs.
    * `type` - Type of the webhook.
    * `updated_at` - The update time of the webhook.
    * `url` - The address of the webhook.


