package network_interface_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_interface"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusNetworkInterfaceCreateConfig = `
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-sg"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description = "acc-test"
  subnet_id = "${byteplus_subnet.foo.id}"
  security_group_ids = ["${byteplus_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.2"]
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusNetworkInterfaceResource_Basic(t *testing.T) {
	resourceName := "byteplus_network_interface.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
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

const testAccByteplusNetworkInterfaceUpdateConfig1 = `
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-sg"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni-new"
  description = "acc-test-new"
  subnet_id = "${byteplus_subnet.foo.id}"
  security_group_ids = ["${byteplus_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.2"]
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
  tags {
    key = "k2"
    value = "v2"
  }
}
`

func TestAccByteplusNetworkInterfaceResource_Update1(t *testing.T) {
	resourceName := "byteplus_network_interface.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccByteplusNetworkInterfaceUpdateConfig1,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
				),
			},
			{
				Config:             testAccByteplusNetworkInterfaceUpdateConfig1,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusNetworkInterfaceUpdateConfig2 = `
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-sg"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni"
  description = "acc-test"
  subnet_id = "${byteplus_subnet.foo.id}"
  security_group_ids = ["${byteplus_security_group.foo.id}"]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  private_ip_address = ["172.16.0.3", "172.16.0.4"]
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusNetworkInterfaceResource_Update2(t *testing.T) {
	resourceName := "byteplus_network_interface.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfaceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccByteplusNetworkInterfaceUpdateConfig2,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "2"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.3"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "private_ip_address.*", "172.16.0.4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config:             testAccByteplusNetworkInterfaceUpdateConfig2,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusNetworkInterfaceCreateConfigIpv6 = `
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

resource "byteplus_security_group" "foo" {
  vpc_id = byteplus_vpc.vpc_ipv6.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni-ipv6"
  description = "acc-test"
  subnet_id = byteplus_subnet.subnet_ipv6.id
  security_group_ids = [byteplus_security_group.foo.id]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  ipv6_address_count = 2
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusNetworkInterfaceResource_CreateIpv6(t *testing.T) {
	resourceName := "byteplus_network_interface.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfaceCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_address_count", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
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

const testAccByteplusNetworkInterfaceUpdateConfigIpv6 = `
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

resource "byteplus_security_group" "foo" {
  vpc_id = byteplus_vpc.vpc_ipv6.id
  security_group_name = "acc-test-security-group"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni-ipv6"
  description = "acc-test"
  subnet_id = byteplus_subnet.subnet_ipv6.id
  security_group_ids = [byteplus_security_group.foo.id]
  primary_ip_address = "172.16.0.253"
  port_security_enabled = false
  ipv6_address_count = 3
  project_name = "default"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusNetworkInterfaceResource_UpdateIpv6(t *testing.T) {
	resourceName := "byteplus_network_interface.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfaceCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_address_count", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccByteplusNetworkInterfaceUpdateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interface_name", "acc-test-eni-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_security_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.253"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "private_ip_address.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_address_count", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config:             testAccByteplusNetworkInterfaceUpdateConfigIpv6,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
