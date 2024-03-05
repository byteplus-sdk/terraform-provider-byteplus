package eip_address_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEipAddressesDatasourceConfig = `
resource "byteplus_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
}
data "byteplus_eip_addresses" "foo"{
    ids = ["${byteplus_eip_address.foo.id}"]
}
`

func TestAccByteplusEipAddressesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_eip_addresses.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.ByteplusEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEipAddressesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "addresses.#", "1"),
				),
			},
		},
	})
}
