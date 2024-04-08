package zone_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/zone"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusClbZonesDatasourceConfig = `
data "byteplus_clb_zones" "foo"{
}
`

func TestAccByteplusClbZonesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_clb_zones.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &zone.ByteplusClbZoneService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbZonesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "master_zones.#", "1"),
				),
			},
		},
	})
}
