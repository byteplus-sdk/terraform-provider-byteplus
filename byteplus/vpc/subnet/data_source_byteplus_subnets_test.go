package subnet_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/subnet"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccSubnetDatasourceConfig = `
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

resource "byteplus_subnet" "foo1" {
  subnet_name = "acc-test-subnet1"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_subnets" "foo"{
  ids = ["${byteplus_subnet.foo.id}", "${byteplus_subnet.foo1.id}"]
}
`

func TestAccByteplusSubnetDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_subnets.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.ByteplusSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "subnets.#", "2"),
				),
			},
		},
	})
}
