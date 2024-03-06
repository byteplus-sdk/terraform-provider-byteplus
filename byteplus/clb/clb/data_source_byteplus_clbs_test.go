package clb_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/clb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusClbsDatasourceConfig = `
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
  	description = "acc-test-demo"
  	load_balancer_name = "acc-test-clb-${count.index}"
	load_balancer_billing_type = "PostPaid"
  	eip_billing_config {
    	isp = "BGP"
    	eip_billing_type = "PostPaidByBandwidth"
    	bandwidth = 1
  	}
	tags {
		key = "k1"
		value = "v1"
	}
	count = 3
}

data "byteplus_clbs" "foo"{
    ids = byteplus_clb.foo[*].id
}
`

func TestAccByteplusClbsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_clbs.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &clb.ByteplusClbService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "clbs.#", "3"),
				),
			},
		},
	})
}
