package bandwidth_package_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/bandwidth_package/bandwidth_package"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusBandwidthPackagesDatasourceConfig = `
resource "byteplus_bandwidth_package" "foo" {
  bandwidth_package_name    = "acc-test-bp"
  billing_type              = "PostPaidByBandwidth"
  isp                       = "BGP"
  description               = "acc-test"
  bandwidth                 = 2
  protocol                  = "IPv4"
  security_protection_types = ["AntiDDoS_Enhanced"]
  tags {
    key   = "k1"
    value = "v1"
  }
  count = 2
}

data "byteplus_bandwidth_packages" "foo" {
  ids = byteplus_bandwidth_package.foo[*].id
}
`

func TestAccByteplusBandwidthPackagesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_bandwidth_packages.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return bandwidth_package.NewBandwidthPackageService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusBandwidthPackagesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "packages.#", "2"),
				),
			},
		},
	})
}
