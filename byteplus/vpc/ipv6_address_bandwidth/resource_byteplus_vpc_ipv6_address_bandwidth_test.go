package ipv6_address_bandwidth_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_address_bandwidth"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccVpcIpv6AddressBandwidthCreate = `
	data "byteplus_zones" "foo"{
	}

	data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
	}

	resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
	  enable_ipv6 = true
	}

	resource "byteplus_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	  vpc_id = "${byteplus_vpc.foo.id}"
	  ipv6_cidr_block = 1
	}

	resource "byteplus_security_group" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
	}

	resource "byteplus_vpc_ipv6_gateway" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  name = "acc-test-1"
	  description = "test"
	}

	resource "byteplus_ecs_instance" "foo" {
	  image_id = "${data.byteplus_images.foo.images[0].image_id}"
	  instance_type = "ecs.g1.large"
	  instance_name = "acc-test-ecs-name"
	  password = "93f0cb0614Aab12"
	  instance_charge_type = "PostPaid"
	  system_volume_type = "ESSD_PL0"
	  system_volume_size = 40
	  subnet_id = byteplus_subnet.foo.id
	  security_group_ids = [byteplus_security_group.foo.id]
	  ipv6_address_count = 1
	}

	data "byteplus_vpc_ipv6_addresses" "foo"{
	  associated_instance_id = "${byteplus_ecs_instance.foo.id}"
	}

	resource "byteplus_vpc_ipv6_address_bandwidth" "foo" {
	  ipv6_address = data.byteplus_vpc_ipv6_addresses.foo.ipv6_addresses.0.ipv6_address
	  billing_type = "PostPaidByBandwidth"
	  bandwidth = 5
	}
`

const testAccVpcIpv6AddressBandwidthUpdate = `
	data "byteplus_zones" "foo"{
	}

	data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
	}

	resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
	  enable_ipv6 = true
	}

	resource "byteplus_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	  vpc_id = "${byteplus_vpc.foo.id}"
	  ipv6_cidr_block = 1
	}

	resource "byteplus_security_group" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
	}
	
	resource "byteplus_vpc_ipv6_gateway" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  name = "acc-test-1"
	  description = "test"
	}
	
	resource "byteplus_ecs_instance" "foo" {
	  image_id = "${data.byteplus_images.foo.images[0].image_id}"
	  instance_type = "ecs.g1.large"
	  instance_name = "acc-test-ecs-name"
	  password = "93f0cb0614Aab12"
	  instance_charge_type = "PostPaid"
	  system_volume_type = "ESSD_PL0"
	  system_volume_size = 40
	  subnet_id = byteplus_subnet.foo.id
	  security_group_ids = [byteplus_security_group.foo.id]
	  ipv6_address_count = 1
	}

	data "byteplus_vpc_ipv6_addresses" "foo"{
	  associated_instance_id = "${byteplus_ecs_instance.foo.id}"
	}

	resource "byteplus_vpc_ipv6_address_bandwidth" "foo" {
	  ipv6_address = data.byteplus_vpc_ipv6_addresses.foo.ipv6_addresses.0.ipv6_address
	  billing_type = "PostPaidByBandwidth"
	  bandwidth = 10
	}
`

func TestAccVpcIpv6AddressBandwidthResource_Basic(t *testing.T) {
	resourceName := "byteplus_vpc_ipv6_address_bandwidth.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_address_bandwidth.ByteplusIpv6AddressBandwidthService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6AddressBandwidthCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "5"),
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

func TestAccVpcIpv6AddressBandwidthResource_Update(t *testing.T) {
	resourceName := "byteplus_vpc_ipv6_address_bandwidth.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_address_bandwidth.ByteplusIpv6AddressBandwidthService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6AddressBandwidthCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "5"),
				),
			},
			{
				Config: testAccVpcIpv6AddressBandwidthUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "10"),
				),
			},
			{
				Config:             testAccVpcIpv6AddressBandwidthUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
