data "byteplus_zones" "foo" {
}
resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_clb" "foo" {
  type               = "public"
  subnet_id          = byteplus_subnet.foo.id
  load_balancer_spec = "small_1"
  description        = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp              = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth        = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id  = byteplus_clb.foo.id
  server_group_name = "acc-test-create"
  description       = "hello demo11"
}

resource "byteplus_listener" "foo" {
  load_balancer_id = byteplus_clb.foo.id
  listener_name    = "acc-test-listener"
  protocol         = "HTTP"
  port             = 90
  server_group_id  = byteplus_server_group.foo.id
  health_check {
    enabled              = "on"
    interval             = 10
    timeout              = 3
    healthy_threshold    = 5
    un_healthy_threshold = 2
    domain               = "byteplus.com"
    http_code            = "http_2xx"
    method               = "GET"
    uri                  = "/"
  }
  enabled = "on"
}
resource "byteplus_clb_rule" "foo" {
  listener_id     = byteplus_listener.foo.id
  server_group_id = byteplus_server_group.foo.id
  domain          = "test-byte123.com"
  url             = "/tftest"
}