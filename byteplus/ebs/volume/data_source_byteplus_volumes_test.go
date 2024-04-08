package volume_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ebs/volume"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVolumesDatasourceConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_volume" "foo" {
	volume_name = "acc-test-volume-${count.index}"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 60
    zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
	count = 3
}

data "byteplus_volumes" "foo"{
    ids = byteplus_volume.foo[*].id
}
`

func TestAccByteplusVolumesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_volumes.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.ByteplusVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVolumesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.#", "3"),
				),
			},
		},
	})
}
