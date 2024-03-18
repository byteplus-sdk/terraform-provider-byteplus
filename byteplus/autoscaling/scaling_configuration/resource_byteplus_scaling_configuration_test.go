package scaling_configuration_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_configuration"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingConfigurationCreateConfig = `
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
`

const testAccByteplusScalingConfigurationUpdateConfig = `
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
    instance_name = "acc-test-instance-new"
    instance_types = ["ecs.g1.large"]
    password = "93f0cb0614Aab12"
    scaling_configuration_name = "acc-test-scaling-config-new"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
    security_group_ids = ["${byteplus_security_group.foo.id}"]
	volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}

	volumes {
    	volume_type = "ESSD_PL0"
    	size = 100
    	delete_with_instance = true
  	}

	tags {
    	key = "k1"
    	value = "v1"
  	}
	ipv6_address_count = 0
}
`

func TestAccByteplusScalingConfigurationResource_Basic(t *testing.T) {
	resourceName := "byteplus_scaling_configuration.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_configuration.ByteplusScalingConfigurationService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingConfigurationCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-instance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_types.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_configuration_name", "acc-test-scaling-config"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.delete_with_instance", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func TestAccByteplusScalingConfigurationResource_Update(t *testing.T) {
	resourceName := "byteplus_scaling_configuration.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_configuration.ByteplusScalingConfigurationService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingConfigurationCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-instance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_types.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_configuration_name", "acc-test-scaling-config"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.delete_with_instance", "true"),
				),
			},
			{
				Config: testAccByteplusScalingConfigurationUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-instance-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_types.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_configuration_name", "acc-test-scaling-config-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.1.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.1.size", "100"),
					resource.TestCheckResourceAttr(acc.ResourceId, "volumes.1.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_address_count", "0"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config:             testAccByteplusScalingConfigurationUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
