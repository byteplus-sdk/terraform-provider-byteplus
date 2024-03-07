package customer_gateway_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/customer_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusCustomerGatewaysDatasourceConfig = `
resource "byteplus_customer_gateway" "foo" {
  ip_address = "192.0.1.3"
  customer_gateway_name = "acc-test"
  description = "acc-test"
  project_name = "default"
}
data "byteplus_customer_gateways" "foo"{
    ids = ["${byteplus_customer_gateway.foo.id}"]
}
`

func TestAccByteplusCustomerGatewaysDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_customer_gateways.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &customer_gateway.ByteplusCustomerGatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCustomerGatewaysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "customer_gateways.#", "1"),
				),
			},
		},
	})
}
