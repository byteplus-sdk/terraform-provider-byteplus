---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_webhook"
sidebar_current: "docs-byteplus-resource-cloud_monitor_webhook"
description: |-
  Provides a resource to manage cloud monitor webhook
---
# byteplus_cloud_monitor_webhook
Provides a resource to manage cloud monitor webhook
## Example Usage
```hcl
resource "byteplus_cloud_monitor_webhook" "foo1" {
  name = "acc-test-webhook-"
  type = "custom"
  url  = "http://alert.volc.com/callback"
}
```
## Argument Reference
The following arguments are supported:
* `name` - (Required) The name of the webhook.

Length limit must not exceed 512 bytes.
The name can be repeated.
* `type` - (Required) Type of the webhook.

custom: custom webhook
wecom: WeChat webhook
lark: Lark webhook
dingtalk: DingTalk webhook.
* `url` - (Required) The address of the webhook.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `created_at` - The creation time of the webhook.
* `event_rule_ids` - Event rule IDs.
* `rule_ids` - Alarm strategy IDs.
* `updated_at` - The update time of the webhook.


## Import
CloudMonitorWebhook can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_webhook.default resource_id
```

