package security_group_rule_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/security_group_rule"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccSecurityGroupDatasourceConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
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
  cidr_ip           = "172.16.0.0/24"
}

data "byteplus_security_group_rules" "foo"{
  security_group_id = "${byteplus_security_group.foo.id}"
  direction = "${byteplus_security_group_rule.foo.direction}"
  cidr_ip = "${byteplus_security_group_rule.foo.cidr_ip}"
}
`

func TestAccByteplusSecurityGroupRulesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_security_group_rules.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group_rule.ByteplusSecurityGroupRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityGroupDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_rules.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_rules.0.direction", "egress"),
				),
			},
		},
	})
}
