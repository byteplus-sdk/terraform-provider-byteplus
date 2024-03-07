package vpn_gateway_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVpnGatewaysDatasourceConfig = `
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

resource "byteplus_vpn_gateway" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  subnet_id = "${byteplus_subnet.foo.id}"
  bandwidth = 20
  vpn_gateway_name = "acc-test"
  description = "acc-test"
  period = 2
  project_name = "default"
}


data "byteplus_vpn_gateways" "foo"{
    ids = ["${byteplus_vpn_gateway.foo.id}"]
}
`

func TestAccByteplusVpnGatewaysDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vpn_gateways.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpn_gateway.ByteplusVpnGatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVpnGatewaysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "vpn_gateways.#", "1"),
				),
			},
		},
	})
}
