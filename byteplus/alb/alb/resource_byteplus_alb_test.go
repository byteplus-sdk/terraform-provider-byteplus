package alb_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusAlbCreateConfig = `
data "byteplus_alb_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "subnet_1" {
  subnet_name = "acc-test-subnet-1"
  cidr_block = "172.16.1.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_subnet" "subnet_2" {
  subnet_name = "acc-test-subnet-2"
  cidr_block = "172.16.2.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[1].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_alb" "foo" {
  address_ip_version = "IPv4"
  type = "private"
  load_balancer_name = "acc-test-alb-private"
  description = "acc-test"
  subnet_ids = [byteplus_subnet.subnet_1.id, byteplus_subnet.subnet_2.id]
  project_name = "default"
  delete_protection = "off"
  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusAlbResource_Basic(t *testing.T) {
	resourceName := "byteplus_alb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return alb.NewAlbService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAlbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-alb-private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection", "off"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "zone_mappings.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "dns_name"),
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

const testAccByteplusAlbUpdateConfig = `
data "byteplus_alb_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "subnet_1" {
  subnet_name = "acc-test-subnet-1"
  cidr_block = "172.16.1.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_subnet" "subnet_2" {
  subnet_name = "acc-test-subnet-2"
  cidr_block = "172.16.2.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[1].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_alb" "foo" {
  address_ip_version = "IPv4"
  type = "private"
  load_balancer_name = "acc-test-alb-private-new"
  description = "acc-test-new"
  subnet_ids = [byteplus_subnet.subnet_1.id, byteplus_subnet.subnet_2.id]
  project_name = "default"
  delete_protection = "off"
  tags {
    key = "k2"
    value = "v2"
  }
  tags {
    key = "k3"
    value = "v3"
  }
}
`

func TestAccByteplusAlbResource_Update(t *testing.T) {
	resourceName := "byteplus_alb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return alb.NewAlbService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAlbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-alb-private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection", "off"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "zone_mappings.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "dns_name"),
				),
			},
			{
				Config: testAccByteplusAlbUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-alb-private-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection", "off"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "zone_mappings.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k3",
						"value": "v3",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "dns_name"),
				),
			},
			{
				Config:             testAccByteplusAlbUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusAlbCreateDualStackConfig = `
data "byteplus_alb_zones" "foo"{
}

resource "byteplus_vpc" "vpc_ipv6" {
  vpc_name = "acc-test-vpc-ipv6"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "byteplus_subnet" "subnet_ipv6_1" {
  subnet_name = "acc-test-subnet-ipv6-1"
  cidr_block = "172.16.1.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.vpc_ipv6.id
  ipv6_cidr_block = 1
}

resource "byteplus_subnet" "subnet_ipv6_2" {
  subnet_name = "acc-test-subnet-ipv6-2"
  cidr_block = "172.16.2.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[1].id
  vpc_id = byteplus_vpc.vpc_ipv6.id
  ipv6_cidr_block = 2
}

resource "byteplus_vpc_ipv6_gateway" "ipv6_gateway" {
  vpc_id = byteplus_vpc.vpc_ipv6.id
  name = "acc-test-ipv6-gateway"
}

resource "byteplus_alb" "foo" {
  address_ip_version = "DualStack"
  type = "public"
  load_balancer_name = "acc-test-alb-public"
  description = "acc-test"
  subnet_ids = [byteplus_subnet.subnet_ipv6_1.id, byteplus_subnet.subnet_ipv6_2.id]
  project_name = "default"
  delete_protection = "off"

  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
  ipv6_eip_billing_config {
    isp = "BGP"
    billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }

  tags {
    key = "k1"
    value = "v1"
  }
  depends_on = [byteplus_vpc_ipv6_gateway.ipv6_gateway]
}
`

func TestAccByteplusAlbResource_CreateDualStack(t *testing.T) {
	resourceName := "byteplus_alb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return alb.NewAlbService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAlbCreateDualStackConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "DualStack"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-alb-public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection", "off"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "zone_mappings.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "ipv6_eip_billing_config.*", map[string]string{
						"isp":          "BGP",
						"billing_type": "PostPaidByBandwidth",
						"bandwidth":    "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "dns_name"),
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
