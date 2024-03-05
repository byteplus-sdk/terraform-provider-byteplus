package snat_entry_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/snat_entry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusSnatEntryCreateConfig = `
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

resource "byteplus_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry"
    nat_gateway_id = "${byteplus_nat_gateway.foo.id}"
	eip_id = "${byteplus_eip_address.foo.id}"
	subnet_id = "${byteplus_subnet.foo.id}"
	depends_on = [byteplus_eip_associate.foo]
}
`

func TestAccByteplusSnatEntryResource_Basic(t *testing.T) {
	resourceName := "byteplus_snat_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.ByteplusSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusSnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
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

const testAccByteplusSnatEntryCreateSourceCidrConfig = `
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

resource "byteplus_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry"
    nat_gateway_id = "${byteplus_nat_gateway.foo.id}"
	eip_id = "${byteplus_eip_address.foo.id}"
	source_cidr = "172.16.0.0/24"
	depends_on = [byteplus_eip_associate.foo]
}
`

func TestAccByteplusSnatEntryResource_SourceCidr(t *testing.T) {
	resourceName := "byteplus_snat_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.ByteplusSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusSnatEntryCreateSourceCidrConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "source_cidr", "172.16.0.0/24"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
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

const testAccByteplusSnatEntryUpdateConfig = `
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
	name = "acc-test-eip1"
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

resource "byteplus_snat_entry" "foo" {
	snat_entry_name = "acc-test-snat-entry-new"
    nat_gateway_id = "${byteplus_nat_gateway.foo.id}"
	eip_id = "${byteplus_eip_address.foo1.id}"
	subnet_id = "${byteplus_subnet.foo.id}"
	depends_on = [byteplus_eip_associate.foo1]
}
`

func TestAccByteplusSnatEntryResource_Update(t *testing.T) {
	resourceName := "byteplus_snat_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &snat_entry.ByteplusSnatEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusSnatEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
				),
			},
			{
				Config: testAccByteplusSnatEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "snat_entry_name", "acc-test-snat-entry-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "nat_gateway_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "source_cidr"),
				),
			},
			{
				Config:             testAccByteplusSnatEntryUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
