package ipv6_gateway_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccVpcIpv6GatewayCreate = `
	data "byteplus_zones" "foo"{
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
`

const testAccVpcIpv6GatewayUpdate = `
	data "byteplus_zones" "foo"{
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
	  name = "acc-test-2"
	  description = "test update"
	}
`

func TestAccVpcIpv6GatewayResource_Basic(t *testing.T) {
	resourceName := "byteplus_vpc_ipv6_gateway.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_gateway.ByteplusIpv6GatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6GatewayCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "test"),
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

func TestAccVpcIpv6GatewayResource_Update(t *testing.T) {
	resourceName := "byteplus_vpc_ipv6_gateway.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_gateway.ByteplusIpv6GatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcIpv6GatewayCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "test"),
				),
			},
			{
				Config: testAccVpcIpv6GatewayUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "test update"),
				),
			},
			{
				Config:             testAccVpcIpv6GatewayUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
