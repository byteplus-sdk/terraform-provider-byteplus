package ecs_instance_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsInstanceCreateConfig = `
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
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Basic(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "security_enhancement_strategy"},
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateBasicAttributeConfig = `
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
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs-new"
	description = "acc-test-new"
	host_name = "tf-acc-test"
	user_data = "ZWNobyBoZWxsbyBlY3Mh"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12new"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_BasicAttribute(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateBasicAttributeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", "echo hello ecs!"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateBasicAttributeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateSecurityGroupConfig = `
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
  	security_group_name = "acc-test-security-group-${count.index}"
  	vpc_id = "${byteplus_vpc.foo.id}"
	count = 3
}

data "byteplus_images" "foo" {
  	os_type = "Linux"
  	visibility = "public"
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = byteplus_security_group.foo[*].id
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_SecurityGroup(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateSecurityGroupConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateSystemVolumeConfig = `
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
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 50
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_SystemVolume(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateSystemVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateSystemVolumeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateInstanceTypeConfig = `
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
  	instance_type_id = "ecs.g1ie.xlarge"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.xlarge"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_InstanceType(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateInstanceTypeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateInstanceTypeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateImageConfig = `
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
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[1].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_Image(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateImageConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateImageConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusEcsInstanceUpdateTagsConfig = `
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
  	instance_type_id = "ecs.g1ie.large"
}

resource "byteplus_ecs_instance" "foo" {
 	instance_name = "acc-test-ecs"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1ie.large"
  	password = "93f0cb0614Aab12"
  	instance_charge_type = "PostPaid"
  	system_volume_type = "ESSD_PL0"
  	system_volume_size = 40
	data_volumes {
    	volume_type = "ESSD_PL0"
    	size = 50
    	delete_with_instance = true
  	}
	subnet_id = "${byteplus_subnet.foo.id}"
	security_group_ids = ["${byteplus_security_group.foo.id}"]
	primary_ip_address = "172.16.0.120"
	project_name = "default"
	tags {
    	key = "k2"
    	value = "v2"
  	}
	tags {
    	key = "k3"
    	value = "v3"
  	}
}
`

func TestAccByteplusEcsInstanceResource_Update_Tags(t *testing.T) {
	resourceName := "byteplus_ecs_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config: testAccByteplusEcsInstanceUpdateTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-ecs"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "primary_ip_address", "172.16.0.120"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "RUNNING"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.size", "50"),
					resource.TestCheckResourceAttr(acc.ResourceId, "data_volumes.0.delete_with_instance", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "tf-acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipv6_addresses.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secondary_network_interfaces.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_enhancement_strategy", "Active"),
					resource.TestCheckResourceAttr(acc.ResourceId, "spot_strategy", "NoSpot"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_size", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "system_volume_type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_data", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k3",
						"value": "v3",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "zone_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "image_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "network_interface_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "subnet_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "system_volume_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "auto_renew_period"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "include_data_volumes"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ipv6_address_count"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "keep_image_credential"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "period"),
				),
			},
			{
				Config:             testAccByteplusEcsInstanceUpdateTagsConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
