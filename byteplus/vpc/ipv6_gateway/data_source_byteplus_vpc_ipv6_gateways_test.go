package ipv6_gateway_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccIpv6GatewayConfig = `
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

	data "byteplus_vpc_ipv6_gateways" "foo" {
		ids = ["${byteplus_vpc_ipv6_gateway.foo.id}"]
	}
`

func TestAccByteplusIpv6GatewayDataSource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vpc_ipv6_gateways.foo"
	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ipv6_gateway.ByteplusIpv6GatewayService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccIpv6GatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_gateways.#", "1"),
				),
			},
		},
	})
}
