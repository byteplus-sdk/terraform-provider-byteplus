package customer_gateway_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/customer_gateway"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusCustomerGatewayCreateConfig = `
resource "byteplus_customer_gateway" "foo" {
  ip_address = "192.0.1.3"
  customer_gateway_name = "acc-test"
  description = "acc-test"
  project_name = "default"
}
`

const testAccByteplusCustomerGatewayUpdateConfig = `
resource "byteplus_customer_gateway" "foo" {
    customer_gateway_name = "acc-test1"
    description = "acc-test1"
    ip_address = "192.0.1.3"
    project_name = "default"
}
`

func TestAccByteplusCustomerGatewayResource_Basic(t *testing.T) {
	resourceName := "byteplus_customer_gateway.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &customer_gateway.ByteplusCustomerGatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCustomerGatewayCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "customer_gateway_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "192.0.1.3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
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

func TestAccByteplusCustomerGatewayResource_Update(t *testing.T) {
	resourceName := "byteplus_customer_gateway.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &customer_gateway.ByteplusCustomerGatewayService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCustomerGatewayCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "customer_gateway_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "192.0.1.3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
				),
			},
			{
				Config: testAccByteplusCustomerGatewayUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "customer_gateway_name", "acc-test1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ip_address", "192.0.1.3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
				),
			},
			{
				Config:             testAccByteplusCustomerGatewayUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
