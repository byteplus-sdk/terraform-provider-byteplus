package eip_associate_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_associate"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEipAssociateCreateConfig = `
	data "byteplus_zones" "foo"{
	}

	data "byteplus_images" "foo" {
	  os_type = "Linux"
	  visibility = "public"
	  instance_type_id = "ecs.g1.large"
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
	  vpc_id = "${byteplus_vpc.foo.id}"
	  security_group_name = "acc-test-security-group"
	}

	resource "byteplus_ecs_instance" "foo" {
	  image_id = "${data.byteplus_images.foo.images[0].image_id}"
	  instance_type = "ecs.g1.large"
	  instance_name = "acc-test-ecs-name"
	  password = "93f0cb0614Aab12"
	  instance_charge_type = "PostPaid"
	  system_volume_type = "ESSD_PL0"
	  system_volume_size = 40
	  subnet_id = byteplus_subnet.foo.id
	  security_group_ids = [byteplus_security_group.foo.id]
	}
	
	resource "byteplus_eip_address" "foo" {
    billing_type = "PostPaidByTraffic"
	}

	resource "byteplus_eip_associate" "foo" {
		allocation_id = "${byteplus_eip_address.foo.id}"
		instance_id = "${byteplus_ecs_instance.foo.id}"
		instance_type = "EcsInstance"
	}
`

func TestAccByteplusEipAssociateResource_Basic(t *testing.T) {
	resourceName := "byteplus_eip_associate.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &eip_associate.ByteplusEipAssociateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEipAssociateCreateConfig,
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
