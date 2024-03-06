package server_group_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/server_group"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusServerGroupCreateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

`

const testAccByteplusServerGroupUpdateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create1"
  description = "acc hello demo11"
}

`

func TestAccByteplusServerGroupResource_Basic(t *testing.T) {
	resourceName := "byteplus_server_group.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.ByteplusServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create"),
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

func TestAccByteplusServerGroupResource_Update(t *testing.T) {
	resourceName := "byteplus_server_group.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.ByteplusServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create"),
				),
			},
			{
				Config: testAccByteplusServerGroupUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc hello demo11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-create1"),
				),
			},
			{
				Config:             testAccByteplusServerGroupUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusServerGroupCreateConfigIpv6 = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "vpc_ipv6" {
  vpc_name = "acc-test-vpc-ipv6"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "byteplus_subnet" "subnet_ipv6" {
  subnet_name = "acc-test-subnet-ipv6"
  cidr_block = "172.16.0.0/24"
  zone_id = data.byteplus_zones.foo.zones[1].id
  vpc_id = byteplus_vpc.vpc_ipv6.id
  ipv6_cidr_block = 1
}

resource "byteplus_clb" "private_clb_ipv6" {
  type = "private"
  subnet_id = byteplus_subnet.subnet_ipv6.id
  load_balancer_name = "acc-test-clb-ipv6"
  load_balancer_spec = "small_1"
  description = "acc-test-demo"
  project_name = "default"
  address_ip_version = "DualStack"
  tags {
    key = "k1"
    value = "v1"
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.private_clb_ipv6.id}"
  server_group_name = "acc-test-sg-ipv6"
  description = "acc-test"
  address_ip_version = "ipv6"
}
`

func TestAccByteplusServerGroupResource_CreateIpv6(t *testing.T) {
	resourceName := "byteplus_server_group.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.ByteplusServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "server_group_name", "acc-test-sg-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "ipv6"),
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
