package scaling_configuration_attachment_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_configuration_attachment"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingConfigurationAttachmentCreateConfig = `
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

resource "byteplus_scaling_group" "foo" {
  scaling_group_name = "acc-test-scaling-group"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 1
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 10
}

resource "byteplus_scaling_configuration" "foo" {
    image_id = "${data.byteplus_images.foo.images[0].image_id}"
    instance_name = "acc-test-instance"
    instance_types = ["ecs.g1.large"]
    password = "93f0cb0614Aab12"
    scaling_configuration_name = "acc-test-scaling-config"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
    security_group_ids = ["${byteplus_security_group.foo.id}"]
	volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
}

resource "byteplus_scaling_configuration_attachment" "foo" {
    scaling_configuration_id = "${byteplus_scaling_configuration.foo.id}"
}
`

const testAccByteplusScalingConfigurationAttachmentUpdateConfig = `
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

resource "byteplus_ecs_key_pair" "foo" {
    description = "acc-test-2"
    key_pair_name = "acc-test-key-pair-name"
}

resource "byteplus_ecs_launch_template" "foo" {
    description = "acc-test-desc"
    eip_bandwidth = 200
    eip_billing_type = "PostPaidByBandwidth"
    eip_isp = "BGP"
    host_name = "acc-hostname"
    image_id = "${data.byteplus_images.foo.images[0].image_id}"
    instance_charge_type = "PostPaid"
    instance_name = "acc-instance-name"
    instance_type_id = "ecs.g1.large"
    key_pair_name = "${byteplus_ecs_key_pair.foo.key_pair_name}"
    launch_template_name = "acc-test-template"
    network_interfaces {
        subnet_id = "${byteplus_subnet.foo.id}"
        security_group_ids = ["${byteplus_security_group.foo.id}"]
    }
	volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
}

resource "byteplus_scaling_group" "foo" {
  scaling_group_name = "acc-test-scaling-group"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 1
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 10
  launch_template_id = "${byteplus_ecs_launch_template.foo.id}"
  launch_template_version = "Default"
}

resource "byteplus_scaling_configuration" "foo" {
    image_id = "${data.byteplus_images.foo.images[0].image_id}"
    instance_name = "acc-test-instance"
    instance_types = ["ecs.g1.large"]
    password = "93f0cb0614Aab12"
    scaling_configuration_name = "acc-test-scaling-config"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
    security_group_ids = ["${byteplus_security_group.foo.id}"]
	volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
}

resource "byteplus_scaling_configuration_attachment" "foo" {
    scaling_configuration_id = "${byteplus_scaling_configuration.foo.id}"
}
`

func TestAccByteplusScalingConfigurationAttachmentResource_Basic(t *testing.T) {
	resourceName := "byteplus_scaling_configuration_attachment.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_configuration_attachment.ByteplusScalingConfigurationAttachmentService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingConfigurationAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccByteplusScalingConfigurationAttachmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
				),
			},
		},
	})
}
