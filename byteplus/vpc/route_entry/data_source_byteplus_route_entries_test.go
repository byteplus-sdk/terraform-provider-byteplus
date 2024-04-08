package route_entry_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_entry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccRouteEntryDatasourceConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc-rn"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet-rn"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_nat_gateway" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  subnet_id = "${byteplus_subnet.foo.id}"
  spec = "Small"
  nat_gateway_name = "acc-test-nat-rn"
}

resource "byteplus_route_table" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  route_table_name = "acc-test-route-table"
}

resource "byteplus_route_entry" "foo" {
  route_table_id = "${byteplus_route_table.foo.id}"
  destination_cidr_block = "172.16.1.0/24"
  next_hop_type = "NatGW"
  next_hop_id = "${byteplus_nat_gateway.foo.id}"
  route_entry_name = "acc-test-route-entry"
}

data "byteplus_route_entries" "foo" {
  route_table_id = "${byteplus_route_table.foo.id}"
  ids = ["${byteplus_route_entry.foo.route_entry_id}"]
}
`

func TestAccByteplusRouteEntryDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_route_entries.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &route_entry.ByteplusRouteEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteEntryDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "route_entries.#", "1"),
				),
			},
		},
	})
}
