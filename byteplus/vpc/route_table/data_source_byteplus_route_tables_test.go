package route_table_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_table"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccRouteTableDatasourceConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_route_table" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
  count = 3
}

data "byteplus_route_tables" "foo" {
  ids = ["${byteplus_route_table.foo[0].id}", "${byteplus_route_table.foo[1].id}", "${byteplus_route_table.foo[2].id}"]
}
`

func TestAccByteplusRouteTableDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_route_tables.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.ByteplusRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "route_tables.#", "3"),
				),
			},
		},
	})
}
