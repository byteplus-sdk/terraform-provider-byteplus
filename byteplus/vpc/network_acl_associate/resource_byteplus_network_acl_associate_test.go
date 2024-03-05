package network_acl_associate_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_acl_associate"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccNetworkAclAssociateConfig = `
	data "byteplus_zones" "foo"{
	}

	resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
	}

	resource "byteplus_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/16"
	  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	  vpc_id = "${byteplus_vpc.foo.id}"
	}		

	resource "byteplus_network_acl" "foo" {
	  vpc_id = "${byteplus_vpc.foo.id}"
	  network_acl_name = "acc-test-acl"
	  ingress_acl_entries {
		network_acl_entry_name = "acc-ingress1"
		policy = "accept"
		protocol = "all"
		source_cidr_ip = "192.168.0.0/24"
	  }
	  egress_acl_entries {
		network_acl_entry_name = "acc-egress2"
		policy = "accept"
		protocol = "all"
		destination_cidr_ip = "192.168.0.0/16"
	  }
	}

	resource "byteplus_network_acl_associate" "foo" {
	  network_acl_id = "${byteplus_network_acl.foo.id}"
	  resource_id = "${byteplus_subnet.foo.id}"
	}

`

func TestAccByteplusNetworkAclAssociateResource_Basic(t *testing.T) {
	resourceName := "byteplus_network_acl_associate.foo"
	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_acl_associate.ByteplusNetworkAclAssociateService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclAssociateConfig,
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
