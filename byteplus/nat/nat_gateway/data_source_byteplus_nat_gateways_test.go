package nat_gateway_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/nat_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusNatGatewaysDatasourceConfig = `
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

resource "byteplus_nat_gateway" "foo" {
	vpc_id = "${byteplus_vpc.foo.id}"
    subnet_id = "${byteplus_subnet.foo.id}"
	spec = "Small"
	nat_gateway_name = "acc-test-ng-${count.index}"
	description = "acc-test"
	billing_type = "PostPaid"
	project_name = "default"
	tags {
		key = "k1"
		value = "v1"
	}
	count =3 
}

data "byteplus_nat_gateways" "foo"{
    ids = byteplus_nat_gateway.foo[*].id
}
`

func TestAccByteplusNatGatewaysDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_nat_gateways.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &nat_gateway.ByteplusNatGatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNatGatewaysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "nat_gateways.#", "3"),
				),
			},
		},
	})
}
