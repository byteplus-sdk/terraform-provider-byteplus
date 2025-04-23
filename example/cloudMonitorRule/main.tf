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
