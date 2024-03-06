package server_group_server_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/server_group_server"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusServerGroupServerCreateConfig = `
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

resource "byteplus_security_group" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
}

data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
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
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = "${byteplus_server_group.foo.id}"
  instance_id = "${byteplus_ecs_instance.foo.id}"
  type = "ecs"
  weight = 100
  port = 80
  description = "This is a acc test server"
}

`

const testAccByteplusServerGroupServerUpdateConfig = `
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

resource "byteplus_security_group" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
}

data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
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
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = "${byteplus_server_group.foo.id}"
  instance_id = "${byteplus_ecs_instance.foo.id}"
  type = "ecs"
  weight = 80
  port = 90
  description = "This is a acc test server 2"
}

`

func TestAccByteplusServerGroupServerResource_Basic(t *testing.T) {
	resourceName := "byteplus_server_group_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group_server.ByteplusServerGroupServerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupServerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "This is a acc test server"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "weight", "100"),
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

func TestAccByteplusServerGroupServerResource_Update(t *testing.T) {
	resourceName := "byteplus_server_group_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group_server.ByteplusServerGroupServerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupServerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "This is a acc test server"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "weight", "100"),
				),
			},
			{
				Config: testAccByteplusServerGroupServerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "This is a acc test server 2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "weight", "80"),
				),
			},
			{
				Config:             testAccByteplusServerGroupServerUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusServerGroupServerCreateConfigIpv6 = `
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

resource "byteplus_security_group" "foo" {
  vpc_id = byteplus_vpc.vpc_ipv6.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_ecs_instance" "foo" {
  image_id = "image-ycjwwciuzy5pkh54xx8f"
  instance_type = "ecs.c3i.large"
  instance_name = "acc-test-ecs-ipv6"
  password = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type = "ESSD_PL0"
  system_volume_size = 40
  subnet_id = byteplus_subnet.subnet_ipv6.id
  security_group_ids = [byteplus_security_group.foo.id]
  ipv6_address_count = 2
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = byteplus_server_group.foo.id
  instance_id = byteplus_ecs_instance.foo.id
  type = "ecs"
  weight = 100
  port = 80
  description = "This is a acc test server"
}
`

func TestAccByteplusServerGroupServerResource_CreateIpv6(t *testing.T) {
	resourceName := "byteplus_server_group_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group_server.ByteplusServerGroupServerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupServerCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "This is a acc test server"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "weight", "100"),
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

const testAccByteplusServerGroupServerCreateConfigEni = `
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

resource "byteplus_security_group" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
}

data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
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
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description = "acc-test"
  subnet_id = byteplus_subnet.foo.id
  security_group_ids = [byteplus_security_group.foo.id]
}

resource "byteplus_network_interface_attach" "foo" {
  instance_id = byteplus_ecs_instance.foo.id
  network_interface_id = byteplus_network_interface.foo.id
}

resource "byteplus_server_group_server" "foo" {
  server_group_id = byteplus_server_group.foo.id
  instance_id = byteplus_network_interface.foo.id
  type = "eni"
  weight = 100
  port = 80
  description = "This is a acc test server"
  depends_on = [byteplus_network_interface_attach.foo]
}
`

func TestAccByteplusServerGroupServerResource_CreateEni(t *testing.T) {
	resourceName := "byteplus_server_group_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group_server.ByteplusServerGroupServerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupServerCreateConfigEni,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "This is a acc test server"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "weight", "100"),
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
