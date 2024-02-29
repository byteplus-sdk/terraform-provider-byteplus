package zone_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/zone"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusZonesDatasourceConfig = `
data "byteplus_zones" "foo"{
    ids = ["cn-beijing-a"]
}
`

func TestAccByteplusZonesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_zones.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &zone.ByteplusZoneService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusZonesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "zones.#", "1"),
				),
			},
		},
	})
}
