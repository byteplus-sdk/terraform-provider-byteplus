package security_group_rule_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/security_group_rule"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccSecurityGroupRuleForCreate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "byteplus_security_group" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  security_group_name = "acc-test-security-group"
}

resource "byteplus_security_group_rule" "foo" {
  direction         = "egress"
  security_group_id = "${byteplus_security_group.foo.id}"
  protocol          = "tcp"
  port_start        = 8000
  port_end          = 9003
  cidr_ip           = "2406:d440:10d:ff00::/64"
}
`

const testAccSecurityGroupRuleForUpdate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
  enable_ipv6 = true
}

resource "byteplus_security_group" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  security_group_name = "acc-test-security-group"
}

resource "byteplus_security_group_rule" "foo" {
  direction         = "egress"
  security_group_id = "${byteplus_security_group.foo.id}"
  protocol          = "tcp"
  port_start        = 8000
  port_end          = 9003
  cidr_ip           = "2406:d440:10d:ff00::/64"
  description       = "tfdesc"
}
`

func TestAccByteplusSecurityGroupRuleResource_Basic(t *testing.T) {
	resourceName := "byteplus_security_group_rule.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.ByteplusSecurityGroupRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "2406:d440:10d:ff00::/64"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSecurityGroupRuleForCreate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccByteplusSubnetResource_Update(t *testing.T) {
	resourceName := "byteplus_security_group_rule.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.ByteplusSecurityGroupRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupRuleForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "2406:d440:10d:ff00::/64"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccSecurityGroupRuleForUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "direction", "egress"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "tcp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_start", "8000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port_end", "9003"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_ip", "2406:d440:10d:ff00::/64"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSecurityGroupRuleForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
