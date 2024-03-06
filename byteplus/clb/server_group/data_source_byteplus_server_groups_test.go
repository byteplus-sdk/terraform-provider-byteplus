package server_group_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/server_group"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusServerGroupsDatasourceConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

data "byteplus_server_groups" "foo"{
    ids = ["${byteplus_server_group.foo.id}"]
}
`

func TestAccByteplusServerGroupsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_server_groups.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &server_group.ByteplusServerGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusServerGroupsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "groups.#", "1"),
				),
			},
		},
	})
}
