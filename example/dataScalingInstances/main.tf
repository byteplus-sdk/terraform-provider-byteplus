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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = byteplus_vpc.foo.id
}


data "byteplus_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g1.large"
}

resource "byteplus_ecs_key_pair" "foo" {
  description   = "acc-test-2"
  key_pair_name = "acc-test-key-pair-name"
}

resource "byteplus_ecs_launch_template" "foo" {
  description          = "acc-test-desc"
  eip_bandwidth        = 200
  eip_billing_type     = "PostPaidByBandwidth"
  eip_isp              = "BGP"
  host_name            = "acc-hostname"
  image_id             = data.byteplus_images.foo.images[0].image_id
  instance_charge_type = "PostPaid"
  instance_name        = "acc-instance-name"
  instance_type_id     = "ecs.g1.large"
  key_pair_name        = byteplus_ecs_key_pair.foo.key_pair_name
  launch_template_name = "acc-test-template"
  network_interfaces {
    subnet_id          = byteplus_subnet.foo.id
    security_group_ids = [byteplus_security_group.foo.id]
  }
  volumes {
    volume_type          = "ESSD_PL0"
    size                 = 50
    delete_with_instance = true
  }
}

resource "byteplus_scaling_group" "foo" {
  scaling_group_name        = "acc-test-scaling-group"
  subnet_ids                = [byteplus_subnet.foo.id]
  multi_az_policy           = "BALANCE"
  desire_instance_number    = -1
  min_instance_number       = 0
  max_instance_number       = 10
  instance_terminate_policy = "OldestInstance"
  default_cooldown          = 10
  launch_template_id        = byteplus_ecs_launch_template.foo.id
  launch_template_version   = "Default"
}

resource "byteplus_scaling_group_enabler" "foo" {
  scaling_group_id = byteplus_scaling_group.foo.id
}

resource "byteplus_ecs_instance" "foo" {
  count                = 3
  instance_name        = "acc-test-ecs-${count.index}"
  description          = "acc-test"
  host_name            = "tf-acc-test"
  image_id             = data.byteplus_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
}

resource "byteplus_scaling_instance_attachment" "foo" {
  count            = length(byteplus_ecs_instance.foo)
  instance_id      = byteplus_ecs_instance.foo[count.index].id
  scaling_group_id = byteplus_scaling_group.foo.id
  entrusted        = true

  depends_on = [
    byteplus_scaling_group_enabler.foo
  ]
}

data "byteplus_scaling_instances" "foo" {
  scaling_group_id = byteplus_scaling_group.foo.id
  ids              = byteplus_scaling_instance_attachment.foo[*].instance_id
}