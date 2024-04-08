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

resource "byteplus_nat_gateway" "foo" {
  vpc_id           = byteplus_vpc.foo.id
  subnet_id        = byteplus_subnet.foo.id
  spec             = "Small"
  nat_gateway_name = "acc-test-ng"
  description      = "acc-test"
  billing_type     = "PostPaid"
  project_name     = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_eip_address" "foo" {
  name         = "acc-test-eip"
  description  = "acc-test"
  bandwidth    = 1
  billing_type = "PostPaidByBandwidth"
  isp          = "BGP"
}

resource "byteplus_eip_associate" "foo" {
  allocation_id = byteplus_eip_address.foo.id
  instance_id   = byteplus_nat_gateway.foo.id
  instance_type = "Nat"
}

resource "byteplus_snat_entry" "foo" {
  snat_entry_name = "acc-test-snat-entry"
  nat_gateway_id  = byteplus_nat_gateway.foo.id
  eip_id          = byteplus_eip_address.foo.id
  source_cidr     = "172.16.0.0/24"
  depends_on      = ["byteplus_eip_associate.foo"]
}
