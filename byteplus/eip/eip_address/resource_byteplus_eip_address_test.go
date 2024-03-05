package eip_address_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEipAddressCreateConfig = `
resource "byteplus_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
}
`

const testAccByteplusEipAddressUpdateConfig = `
resource "byteplus_eip_address" "foo" {
    bandwidth = 1
    billing_type = "PostPaidByBandwidth"
    description = "acc-test"
    isp = "BGP"
    name = "acc-test-eip"
}
`

func TestAccByteplusEipAddressResource_Basic(t *testing.T) {
	resourceName := "byteplus_eip_address.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.ByteplusEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEipAddressCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByTraffic"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
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

func TestAccByteplusEipAddressResource_Update(t *testing.T) {
	resourceName := "byteplus_eip_address.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_address.ByteplusEipAddressService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEipAddressCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByTraffic"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccByteplusEipAddressUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-eip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccByteplusEipAddressUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
