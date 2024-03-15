package cen_bandwidth_package_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_bandwidth_package"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusCenBandwidthPackagesDatasourceConfig = `
resource "byteplus_cen_bandwidth_package" "foo" {
  local_geographic_region_set_id = "China"
  peer_geographic_region_set_id  = "China"
  bandwidth                      = 2
  cen_bandwidth_package_name     = "acc-test-cen-bp-${count.index}"
  description                    = "acc-test"
  billing_type                   = "PrePaid"
  period_unit                    = "Month"
  period                         = 1
  project_name                   = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
  count = 2
}

data "byteplus_cen_bandwidth_packages" "foo"{
    ids = byteplus_cen_bandwidth_package.foo[*].id
}
`

func TestAccByteplusCenBandwidthPackagesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_cen_bandwidth_packages.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return cen_bandwidth_package.NewCenBandwidthPackageService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCenBandwidthPackagesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth_packages.#", "2"),
				),
			},
		},
	})
}
