package dnat_entry_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/dnat_entry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusDnatEntryCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_nat_gateway" "foo" {
	vpc_id = "${byteplus_vpc.foo.id}"
    subnet_id = "${byteplus_subnet.foo.id}"
	spec = "Small"
	nat_gateway_name = "acc-test-ng"
	description = "acc-test"
	billing_type = "PostPaid"
	project_name = "default"
	tags {
		key = "k1"
		value = "v1"
	}
}

resource "byteplus_eip_address" "foo" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "byteplus_eip_associate" "foo" {
	allocation_id = "${byteplus_eip_address.foo.id}"
	instance_id = "${byteplus_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "byteplus_dnat_entry" "foo" {
	dnat_entry_name = "acc-test-dnat-entry"
    external_ip = "${byteplus_eip_address.foo.eip_address}"
    external_port = 80
    internal_ip = "172.16.0.10"
    internal_port = 80
    nat_gateway_id = "${byteplus_nat_gateway.foo.id}"
    protocol = "tcp"
	depends_on = [byteplus_eip_associate.foo]
}
`

func TestAccByteplusDnatEntryResource_Basic(t *testing.T) {
	resourceName := "byteplus_dnat_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &dnat_entry.ByteplusDnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusDnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccByteplusDnatEntryUpdateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_nat_gateway" "foo" {
	vpc_id = "${byteplus_vpc.foo.id}"
    subnet_id = "${byteplus_subnet.foo.id}"
	spec = "Small"
	nat_gateway_name = "acc-test-ng"
	description = "acc-test"
	billing_type = "PostPaid"
	project_name = "default"
	tags {
		key = "k1"
		value = "v1"
	}
}

resource "byteplus_eip_address" "foo" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "byteplus_eip_associate" "foo" {
	allocation_id = "${byteplus_eip_address.foo.id}"
	instance_id = "${byteplus_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "byteplus_eip_address" "foo1" {
	name = "acc-test-eip"
    description = "acc-test"
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    isp = "BGP"
}

resource "byteplus_eip_associate" "foo1" {
	allocation_id = "${byteplus_eip_address.foo1.id}"
	instance_id = "${byteplus_nat_gateway.foo.id}"
	instance_type = "Nat"
}

resource "byteplus_dnat_entry" "foo" {
	dnat_entry_name = "acc-test-dnat-entry-new"
    external_ip = "${byteplus_eip_address.foo1.eip_address}"
    external_port = 90
    internal_ip = "172.16.0.17"
    internal_port = 90
    nat_gateway_id = "${byteplus_nat_gateway.foo.id}"
    protocol = "udp"
	depends_on = [byteplus_eip_associate.foo1]
}
`

func TestAccByteplusDnatEntryResource_Update(t *testing.T) {
	resourceName := "byteplus_dnat_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &dnat_entry.ByteplusDnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusDnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
				),
			},
			{
				Config: testAccByteplusDnatEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "dnat_entry_name", "acc-test-dnat-entry-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "external_port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_ip", "172.16.0.17"),
					resource.TestCheckResourceAttr(acc.ResourceId, "internal_port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "udp"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "external_ip"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "dnat_entry_id"),
				),
			},
			{
				Config:             testAccByteplusDnatEntryUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
