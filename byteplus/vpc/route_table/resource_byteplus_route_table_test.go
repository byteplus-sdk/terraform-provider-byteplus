package route_table_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_table"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccRouteTableForCreate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_route_table" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
}
`

func TestAccByteplusRouteTableResource_Basic(t *testing.T) {
	resourceName := "byteplus_route_table.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.ByteplusRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table"),
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

const testAccRouteTableForUpdate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_route_table" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table-new"
  description = "tfdesc"
}
`

func TestAccByteplusRouteTableResource_Update(t *testing.T) {
	resourceName := "byteplus_route_table.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table.ByteplusRouteTableService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table"),
				),
			},
			{
				Config: testAccRouteTableForUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "route_table_name", "acc-test-route-table-new"),
				),
			},
			{
				Config:             testAccRouteTableForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
