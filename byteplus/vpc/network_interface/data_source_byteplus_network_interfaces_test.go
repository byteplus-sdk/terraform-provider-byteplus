package network_interface_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_interface"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusNetworkInterfacesDatasourceConfig = `
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-sg"
  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_network_interface" "foo" {
  network_interface_name = "acc-test-eni-${count.index}"
  subnet_id = "${byteplus_subnet.foo.id}"
  security_group_ids = ["${byteplus_security_group.foo.id}"]
  count = 3
}

data "byteplus_network_interfaces" "foo"{
    ids = byteplus_network_interface.foo[*].id
}
`

func TestAccByteplusNetworkInterfacesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_network_interfaces.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_interface.ByteplusNetworkInterfaceService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusNetworkInterfacesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "network_interfaces.#", "3"),
				),
			},
		},
	})
}
