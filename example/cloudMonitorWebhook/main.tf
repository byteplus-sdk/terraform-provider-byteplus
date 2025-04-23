resource "byteplus_cloud_monitor_webhook" "foo1" {
  name = "acc-test-webhook-"
  type = "custom"
  url  = "http://alert.volc.com/callback"
}
