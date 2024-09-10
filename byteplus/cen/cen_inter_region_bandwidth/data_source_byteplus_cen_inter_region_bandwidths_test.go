package cen_inter_region_bandwidth_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_inter_region_bandwidth"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusCenInterRegionBandwidthsDatasourceConfig = `
resource "byteplus_cen" "foo" {
  cen_name     = "acc-test-cen"
  description  = "acc-test"
  project_name = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_cen_bandwidth_package" "foo" {
  local_geographic_region_set_id = "China"
  peer_geographic_region_set_id  = "China"
  bandwidth                      = 5
  cen_bandwidth_package_name     = "acc-test-cen-bp"
  description                    = "acc-test"
  billing_type                   = "PrePaid"
  period_unit                    = "Month"
  period                         = 1
  project_name                   = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_cen_bandwidth_package_associate" "foo" {
  cen_bandwidth_package_id = byteplus_cen_bandwidth_package.foo.id
  cen_id                   = byteplus_cen.foo.id
}

resource "byteplus_cen_inter_region_bandwidth" "foo" {
  cen_id          = byteplus_cen.foo.id
  local_region_id = "cn-beijing"
  peer_region_id  = "cn-shanghai"
  bandwidth       = 2
  depends_on      = [byteplus_cen_bandwidth_package_associate.foo]
}

data "byteplus_cen_inter_region_bandwidths" "foo"{
  ids = [byteplus_cen_inter_region_bandwidth.foo.id]
}
`

func TestAccByteplusCenInterRegionBandwidthsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_cen_inter_region_bandwidths.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return cen_inter_region_bandwidth.NewCenInterRegionBandwidthService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCenInterRegionBandwidthsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "inter_region_bandwidths.#", "1"),
				),
			},
		},
	})
}
