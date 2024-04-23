resource "byteplus_eip_address" "foo" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "BGP"
  name         = "acc-test-eip"
  description  = "acc-test"
  project_name = "default"
}

resource "byteplus_bandwidth_package" "ipv4" {
  bandwidth_package_name = "acc-test-bp"
  billing_type           = "PostPaidByBandwidth"
  isp                    = "BGP"
  description            = "acc-test"
  bandwidth              = 2
  protocol               = "IPv4"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_bandwidth_package_attachment" "ipv4" {
  allocation_id        = byteplus_eip_address.foo.id
  bandwidth_package_id = byteplus_bandwidth_package.ipv4.id
}

data "byteplus_zones" "foo" {
}

data "byteplus_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g1.large"
}

resource "byteplus_vpc" "foo" {
  vpc_name    = "acc-test-vpc"
  cidr_block  = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "byteplus_subnet" "foo" {
  subnet_name     = "acc-test-subnet"
  cidr_block      = "172.16.0.0/24"
  zone_id         = data.byteplus_zones.foo.zones[0].id
  vpc_id          = byteplus_vpc.foo.id
  ipv6_cidr_block = 1
}

resource "byteplus_security_group" "foo" {
  vpc_id              = byteplus_vpc.foo.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_vpc_ipv6_gateway" "foo" {
  vpc_id      = byteplus_vpc.foo.id
  name        = "acc-test-1"
  description = "test"
}

resource "byteplus_ecs_instance" "foo" {
  image_id             = data.byteplus_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  instance_name        = "acc-test-ecs-name"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
  ipv6_address_count   = 1
}

data "byteplus_vpc_ipv6_addresses" "foo" {
  associated_instance_id = byteplus_ecs_instance.foo.id
}

resource "byteplus_vpc_ipv6_address_bandwidth" "foo" {
  ipv6_address = data.byteplus_vpc_ipv6_addresses.foo.ipv6_addresses.0.ipv6_address
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 5
}

resource "byteplus_bandwidth_package" "ipv6" {
  bandwidth_package_name = "acc-test-bp"
  billing_type           = "PostPaidByBandwidth"
  isp                    = "BGP"
  description            = "acc-test"
  bandwidth              = 2
  protocol               = "IPv6"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_bandwidth_package_attachment" "foo" {
  allocation_id        = byteplus_vpc_ipv6_address_bandwidth.foo.id
  bandwidth_package_id = byteplus_bandwidth_package.ipv6.id
}
