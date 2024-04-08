package vpc_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/vpc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccVpcDatasourceConfig = `
resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_vpc" "foo1" {
  vpc_name   = "acc-test-vpc1"
  cidr_block = "172.16.0.0/16"
}

data "byteplus_vpcs" "foo"{
  ids = ["${byteplus_vpc.foo1.id}", "${byteplus_vpc.foo.id}"]
}
`

func TestAccByteplusVpcDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vpcs.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpc.ByteplusVpcService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "vpcs.#", "2"),
				),
			},
		},
	})
}
