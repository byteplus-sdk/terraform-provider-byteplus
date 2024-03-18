package scaling_configuration_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingConfigurationsDatasourceConfig = `
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
    count = 3
    image_id = "${data.byteplus_images.foo.images[0].image_id}"
    instance_name = "acc-test-instance"
    instance_types = ["ecs.g1.large"]
    password = "93f0cb0614Aab12"
    scaling_configuration_name = "acc-test-scaling-config-${count.index}"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
    security_group_ids = ["${byteplus_security_group.foo.id}"]
	volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
}

data "byteplus_scaling_configurations" "foo"{
    ids = byteplus_scaling_configuration.foo[*].id
}
`

func TestAccByteplusScalingConfigurationsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_scaling_configurations.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_configuration.ByteplusScalingConfigurationService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingConfigurationsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_configurations.#", "3"),
				),
			},
		},
	})
}
