package volume_attach_test

import (
	"regexp"
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ebs/volume_attach"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusVolumeAttachCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
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

resource "byteplus_volume_attach" "foo" {
    instance_id = "${byteplus_ecs_instance.foo.id}"
    volume_id = "${byteplus_volume.foo.id}"
}
`

func TestAccByteplusVolumeAttachResource_Basic(t *testing.T) {
	resourceName := "byteplus_volume_attach.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume_attach.ByteplusVolumeAttachService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVolumeAttachCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "attached"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "false"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "volume_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "updated_at"),
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

const testAccByteplusVolumeAttachDeleteWithInstanceConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	vpc_name   = "acc-test-vpc"
  	cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  	subnet_name = "acc-test-subnet"
  	cidr_block = "172.16.0.0/24"
  	zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
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
	delete_with_instance = true
}

resource "byteplus_volume_attach" "foo" {
    instance_id = "${byteplus_ecs_instance.foo.id}"
    volume_id = "${byteplus_volume.foo.id}"
	delete_with_instance = true
}
`

func TestAccByteplusVolumeAttachResource_DeleteWithInstance(t *testing.T) {
	resourceName := "byteplus_volume_attach.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &volume_attach.ByteplusVolumeAttachService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config:      testAccByteplusVolumeAttachDeleteWithInstanceConfig,
				ExpectError: regexp.MustCompile("^After applying this step, the plan was not empty(.|\n)*UPDATE: byteplus_volume\\.foo(.|\n)*delete_with_instance: \"false\" => \"true\""),
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "attached"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_with_instance", "true"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "volume_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "created_at"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "updated_at"),
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
