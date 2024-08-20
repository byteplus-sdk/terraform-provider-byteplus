resource "byteplus_cdn_edge_function" "foo" {
  name         = "acc-test-function"
  remark       = "tf-test"
  project_name = "default"
  source_code  = base64encode("hello world")
  envs {
    key   = "k1"
    value = "v1"
  }
  canary_countries = ["China", "Japan", "United Kingdom"]
}

resource "byteplus_cdn_cron_job" "foo" {
  function_id     = byteplus_cdn_edge_function.foo.id
  job_name        = "acc-test-cron-job"
  description     = "tf-test"
  cron_type       = 1
  cron_expression = "0 17 10 * *"
  parameter       = "test"
}
