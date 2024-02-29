package volume_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ebs/volume"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVolumeCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_volume" "foo" {
	volume_name = "acc-test-volume"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 40
    zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
}
`

func TestAccByteplusVolumeResource_Basic(t *testing.T) {
	resourceName := "byteplus_volume.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.ByteplusVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
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

const testAccByteplusVolumeUpdateBasicAttributeConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_volume" "foo" {
	volume_name = "acc-test-volume-new"
    volume_type = "ESSD_PL0"
	description = "acc-test-new"
    kind = "data"
    size = 40
    zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
	delete_with_instance = true
}
`

func TestAccByteplusVolumeResource_UpdateBasicAttribute(t *testing.T) {
	resourceName := "byteplus_volume.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.ByteplusVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config: testAccByteplusVolumeUpdateBasicAttributeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config:             testAccByteplusVolumeUpdateBasicAttributeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVolumeUpdateVolumeSizeConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_volume" "foo" {
	volume_name = "acc-test-volume"
    volume_type = "ESSD_PL0"
	description = "acc-test"
    kind = "data"
    size = 60
    zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	volume_charge_type = "PostPaid"
	project_name = "default"
}
`

func TestAccByteplusVolumeResource_UpdateVolumeSize(t *testing.T) {
	resourceName := "byteplus_volume.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume.ByteplusVolumeService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVolumeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config: testAccByteplusVolumeUpdateVolumeSizeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_name", "acc-test-volume"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kind", "data"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volume_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_id", ""),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trade_status"),
				),
			},
			{
				Config:             testAccByteplusVolumeUpdateVolumeSizeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
