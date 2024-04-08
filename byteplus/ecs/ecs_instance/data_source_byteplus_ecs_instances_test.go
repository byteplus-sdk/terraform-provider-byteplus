package ecs_instance_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsInstancesDatasourceConfig = `
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
 	instance_name = "acc-test-ecs-${count.index}"
	description = "acc-test"
	host_name = "tf-acc-test"
  	image_id = "${data.byteplus_images.foo.images[0].image_id}"
  	instance_type = "ecs.g1.large"
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
	project_name = "default"
	tags {
    	key = "k1"
    	value = "v1"
  	}
	count = 2
}

data "byteplus_ecs_instances" "foo" {
  ids = byteplus_ecs_instance.foo[*].id
}
`

func TestAccByteplusEcsInstancesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_ecs_instances.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_instance.ByteplusEcsService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsInstancesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "instances.#", "2"),
				),
			},
		},
	})
}
