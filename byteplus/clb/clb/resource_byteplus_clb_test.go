package clb_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/clb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusClbCreateConfig = `
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
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb"
	load_balancer_billing_type = "PostPaid"
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
}
`

func TestAccByteplusClbResource_Basic(t *testing.T) {
	resourceName := "byteplus_clb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
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

const testAccByteplusClbUpdateBasicAttributeConfig = `
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
  	load_balancer_spec = "small_2"
  	description = "acc-test-demo-new"
  	load_balancer_name = "acc-test-clb-new"
	load_balancer_billing_type = "PostPaid"
	modification_protection_status = "ConsoleProtection"
	modification_protection_reason = "reason"
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
}
`

func TestAccByteplusClbResource_UpdateBasicAttribute(t *testing.T) {
	resourceName := "byteplus_clb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
				),
			},
			{
				Config: testAccByteplusClbUpdateBasicAttributeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", "ConsoleProtection"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", "reason"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
				),
			},
			{
				Config:             testAccByteplusClbUpdateBasicAttributeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusClbUpdateBillingTypeConfig1 = `
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
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb"
	load_balancer_billing_type = "PrePaid"
	period = 1
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
}
`

const testAccByteplusClbUpdateBillingTypeConfig2 = `
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
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb"
	load_balancer_billing_type = "PrePaid"
	period = 2
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PrePaid"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
}
`

func TestAccByteplusClbResource_UpdateBillingType(t *testing.T) {
	resourceName := "byteplus_clb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
				),
			},
			{
				Config:             testAccByteplusClbUpdateBillingTypeConfig1,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "period", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "renew_type", "ManualRenew"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PrePaid",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
				),
			},
			{
				Config: testAccByteplusClbUpdateBillingTypeConfig2,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "period", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "renew_type", "ManualRenew"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PrePaid",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
				),
			},
			{
				Config:             testAccByteplusClbUpdateBillingTypeConfig2,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusClbUpdateTagsConfig = `
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
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb"
	load_balancer_billing_type = "PostPaid"
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
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

func TestAccByteplusClbResource_UpdateTags(t *testing.T) {
	resourceName := "byteplus_clb.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
				),
			},
			{
				Config: testAccByteplusClbUpdateTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "public"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "eip_billing_config.*", map[string]string{
						"isp":              "BGP",
						"eip_billing_type": "PostPaidByBandwidth",
						"bandwidth":        "1",
					}),
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
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eip_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
				),
			},
			{
				Config:             testAccByteplusClbUpdateTagsConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusClbCreateConfigIpv6 = `
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
`

func TestAccByteplusClbResource_CreateIpv6(t *testing.T) {
	resourceName := "byteplus_clb.private_clb_ipv6"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbCreateConfigIpv6,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_name", "acc-test-clb-ipv6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_spec", "small_1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "load_balancer_billing_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "type", "private"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "address_ip_version", "DualStack"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_address", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_reason", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "modification_protection_status", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "vpc_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "eni_ipv6_address"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "master_zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "region_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "renew_type"),
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
