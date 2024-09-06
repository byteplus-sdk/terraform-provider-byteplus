resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "cn-beijing-a"
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id           = byteplus_vpc.foo.id
  subnet_id        = byteplus_subnet.foo.id
  bandwidth        = 50
  vpn_gateway_name = "acc-test1"
  description      = "acc-test1"
  period           = 7
  project_name     = "default"
}