---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_rule"
sidebar_current: "docs-byteplus-resource-cloud_monitor_rule"
description: |-
  Provides a resource to manage cloud monitor rule
---
# byteplus_cloud_monitor_rule
Provides a resource to manage cloud monitor rule
## Example Usage
```hcl
resource "byteplus_cloud_monitor_contact" "foo1" {
  name  = "acc-test-contact-1"
  email = "test1@163.com"
}

resource "byteplus_cloud_monitor_contact" "foo2" {
  name  = "acc-test-contact-2"
  email = "test2@163.com"
}

resource "byteplus_cloud_monitor_contact_group" "foo" {
  name             = "acc-test-contact-group-new"
  description      = "tf-test-new"
  contacts_id_list = [byteplus_cloud_monitor_contact.foo1.id, byteplus_cloud_monitor_contact.foo2.id]
}

resource "byteplus_cloud_monitor_webhook" "foo" {
  name = "acc-test-webhook"
  type = "custom"
  url  = "http://alert.volc.com/callback"
}

# use level conditions
resource "byteplus_cloud_monitor_rule" "foo1" {
  rule_name        = "acc-test-rule-level-conditions"
  description      = "acc-test"
  namespace        = "VCM_ECS"
  sub_namespace    = "Storage"
  enable_state     = "disable"
  evaluation_count = 5
  effect_start_at  = "00:15"
  effect_end_at    = "22:55"
  silence_time     = 5
  alert_methods    = ["Email", "Webhook"]
  #  web_hook = "http://alert.volc.com/callback"
  webhook_ids         = [byteplus_cloud_monitor_webhook.foo.id]
  contact_group_ids   = [byteplus_cloud_monitor_contact_group.foo.id]
  multiple_conditions = true
  condition_operator  = "||"
  notify_mode         = "rule"
  regions             = ["ap-southeast-1"]
  original_dimensions {
    key   = "ResourceID"
    value = ["*"]
  }
  level_conditions {
    level = "warning"
    conditions {
      metric_name         = "DiskUsageAvail"
      metric_unit         = "Megabytes"
      statistics          = "avg"
      comparison_operator = ">"
      threshold           = "100"
    }
    conditions {
      metric_name         = "DiskUsageUtilization"
      metric_unit         = "Percent"
      statistics          = "avg"
      comparison_operator = ">"
      threshold           = "90"
    }
  }
  level_conditions {
    level = "critical"
    conditions {
      metric_name         = "DiskUsageAvail"
      metric_unit         = "Megabytes"
      statistics          = "avg"
      comparison_operator = ">"
      threshold           = "100"
    }
    conditions {
      metric_name         = "DiskUsageUtilization"
      metric_unit         = "Percent"
      statistics          = "avg"
      comparison_operator = ">"
      threshold           = "90"
    }
  }
  recovery_notify {
    enable = true
  }
  no_data {
    enable           = true
    evaluation_count = 5
  }
  project_name = "default"
}
```
## Argument Reference
The following arguments are supported:
* `alert_methods` - (Required) The notification method of the cloud monitor rule. Valid values: `Email`, `Webhook`.
* `effect_end_at` - (Required) The effect end time of the cloud monitor rule. The expression is `HH:MM`.
* `effect_start_at` - (Required) The effect start time of the cloud monitor rule. The expression is `HH:MM`.
* `enable_state` - (Required) Whether to enable the cloud monitor rule. Valid values: `enable`, `disable`.
* `evaluation_count` - (Required) The number of successive times for which the threshold is reached before the alarm is triggered. Unit in minutes. Supports configurations of 1, 3, 5, 10, 15, 30, 60, and 120.
* `namespace` - (Required, ForceNew) The namespace of the cloud monitor rule.
* `original_dimensions` - (Required) The original dimensions of the cloud monitor rule.
* `regions` - (Required, ForceNew) The region ids of the cloud monitor rule.
* `rule_name` - (Required) The name of the cloud monitor rule.
* `silence_time` - (Required) The interval at which alarms are sent. Unit in minutes. Valid values: 5, 30, 60, 180, 360, 720, 1440.
* `sub_namespace` - (Required, ForceNew) The sub namespace of the cloud monitor rule.
* `condition_operator` - (Optional) Specifies whether the alarm is triggered only when the conditions on multiple metrics are met. Valid values: `&&`, `||`.
* `conditions` - (Optional) The conditions that trigger the alarm.
Specify an array that contains a maximum of 10 metric math expressions. One of `conditions` and `level_conditions` must be specified.
* `contact_group_ids` - (Optional) The contact group ids of the cloud monitor rule. When the alert method is `Email`,, This field must be specified.
* `description` - (Optional) The description of the cloud monitor rule.
* `level_conditions` - (Optional) The level conditions that trigger the alarm. One of `conditions` and `level_conditions` must be specified.
* `level` - (Optional) The severity level of the cloud monitor rule. Valid values: `critical`, `warning`, `notice`. One of `level` and `level_conditions` must be specified.
* `multiple_conditions` - (Optional) Whether to use multiple metrics in the cloud monitor rule.
* `no_data` - (Optional) No-data alarm.
* `notification_id` - (Optional) The notification id of the cloud monitor rule.
* `notify_mode` - (Optional) Alarm sending aggregation strategy.

rule(default): aggregation by rule.
resource: aggregation by rule and resource.
* `project_name` - (Optional) The project name of the cloud monitor rule.
* `recovery_notify` - (Optional) The recovery notify of the cloud monitor rule.
* `web_hook` - (Optional) The webhook URL that is used when an alarm is triggered. When the alert method is `Webhook`, one of `web_hook` and `webhook_ids` must be specified.
* `webhook_ids` - (Optional) The web hook id list of the cloud monitor rule. When the alert method is `Webhook`, one of `web_hook` and `webhook_ids` must be specified.

The `conditions` object supports the following:

* `comparison_operator` - (Required) The comparison operation of the cloud monitor rule. Valid values: `>`, `>=`, `<`, `<=`, `!=`, `=`.
* `metric_name` - (Required) The metric name of the cloud monitor rule.
* `metric_unit` - (Required) The metric unit of the cloud monitor rule.
* `statistics` - (Required) The statistics of the cloud monitor rule. Valid values: `avg`, `max`, `min`.
* `threshold` - (Required) The threshold of the cloud monitor rule.

The `level_conditions` object supports the following:

* `level` - (Required) The severity level of the cloud monitor rule. Valid values: `critical`, `warning`, `notice`.
* `conditions` - (Optional) The conditions that trigger the alarm.
Specify an array that contains a maximum of 10 metric math expressions.

The `no_data` object supports the following:

* `enable` - (Optional) Specifies whether to enable no-data alarm. The default value is false.
* `evaluation_count` - (Optional) No data alarm triggering threshold. When `enable` is set to true, `evaluation_count` is mandatory. The range of values is integers between 3 and 20.

The `original_dimensions` object supports the following:

* `key` - (Required) The key of the dimension.
* `value` - (Required) The value of the dimension. If you want to specify all possible values of the dimension, set the value to an asterisk ( * ).

The `recovery_notify` object supports the following:

* `enable` - (Optional) Specifies whether alarm recovery notifications are sent.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `alert_state` - The alert state of the cloud monitor rule.
* `created_at` - The created time of the cloud monitor rule.
* `updated_at` - The updated time of the cloud monitor rule.


## Import
CloudMonitorRule can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_rule.default 174284623567451****
```

