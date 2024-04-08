package route_table_associate_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_table_associate"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccRouteTableAssociateForCreate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc-attach"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet-attach"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_subnet" "foo1" {
  subnet_name = "acc-test-subnet-attach1"
  cidr_block = "172.16.16.0/20"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_subnet" "foo2" {
  subnet_name = "acc-test-subnet-attach2"
  cidr_block = "172.16.6.0/23"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_subnet" "foo3" {
  subnet_name = "acc-test-subnet-attach3"
  cidr_block = "172.16.14.0/26"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_route_table" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach"
}

resource "byteplus_route_table" "foo1" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach1"
}

resource "byteplus_route_table" "foo2" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach2"
}

resource "byteplus_route_table" "foo3" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table-attach3"
}

resource "byteplus_route_table_associate" "foo" {
  route_table_id = "${byteplus_route_table.foo.id}"
  subnet_id = "${byteplus_subnet.foo.id}"
}

resource "byteplus_route_table_associate" "foo1" {
  route_table_id = "${byteplus_route_table.foo1.id}"
  subnet_id = "${byteplus_subnet.foo1.id}"
}

resource "byteplus_route_table_associate" "foo2" {
  route_table_id = "${byteplus_route_table.foo2.id}"
  subnet_id = "${byteplus_subnet.foo2.id}"
}

resource "byteplus_route_table_associate" "foo3" {
  route_table_id = "${byteplus_route_table.foo3.id}"
  subnet_id = "${byteplus_subnet.foo3.id}"
}

`

func TestAccByteplusRouteTableAssociateResource_Basic(t *testing.T) {
	resourceName := "byteplus_route_table_associate.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_table_associate.ByteplusRouteTableAssociateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTableAssociateForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
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
