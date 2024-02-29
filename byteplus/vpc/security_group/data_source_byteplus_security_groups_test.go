package security_group_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/security_group"
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
  count = 3
}

data "byteplus_security_groups" "foo"{
  ids = ["${byteplus_security_group.foo[0].id}", "${byteplus_security_group.foo[1].id}", "${byteplus_security_group.foo[2].id}"]
}
`

func TestAccByteplusSecurityGroupDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_security_groups.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &security_group.ByteplusSecurityGroupService{},
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
					resource.TestCheckResourceAttr(acc.ResourceId, "security_groups.#", "3"),
				),
			},
		},
	})
}
