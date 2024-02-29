package ecs_deployment_set_associate_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set_associate"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsDeploymentSetAssociateCreateConfig = `
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
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
}

resource "byteplus_ecs_instance_state" "foo" {
  	instance_id = "${byteplus_ecs_instance.foo.id}"
  	action = "Stop"
  	stopped_mode = "KeepCharging"
}

resource "byteplus_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
}

resource "byteplus_ecs_deployment_set_associate" "foo" {
    deployment_set_id = "${byteplus_ecs_deployment_set.foo.id}"
    instance_id = "${byteplus_ecs_instance.foo.id}"
	depends_on = [byteplus_ecs_instance_state.foo]
}
`

func TestAccByteplusEcsDeploymentSetAssociateResource_Basic(t *testing.T) {
	resourceName := "byteplus_ecs_deployment_set_associate.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set_associate.ByteplusEcsDeploymentSetAssociateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsDeploymentSetAssociateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
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
