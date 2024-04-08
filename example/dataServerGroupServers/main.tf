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

resource "byteplus_security_group" "foo" {
  vpc_id              = byteplus_vpc.foo.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_ecs_instance" "foo" {
  image_id             = "image-ycjwwciuzy5pkh54xx8f"
  instance_type        = "ecs.c3i.large"
  instance_name        = "acc-test-ecs-name"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = byteplus_server_group.foo.id
  instance_id     = byteplus_ecs_instance.foo.id
  type            = "ecs"
  weight          = 100
  port            = 80
  description     = "This is a acc test server"
}

data "byteplus_server_group_servers" "foo" {
  ids             = [element(split(":", byteplus_server_group_server.foo.id), length(split(":", byteplus_server_group_server.foo.id)) - 1)]
  server_group_id = byteplus_server_group.foo.id
}