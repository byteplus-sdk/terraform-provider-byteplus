package network_acl_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_acl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccNetworkAclConfig = `
	data "byteplus_zones" "foo"{
	}

	resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
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

	data "byteplus_network_acls" "foo"{
	  ids = [byteplus_network_acl.foo.id]
	}
`

func TestAccByteplusNetworkAclDataSource_Basic(t *testing.T) {
	resourceName := "data.byteplus_network_acls.foo"
	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &network_acl.ByteplusNetworkAclService{},
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkAclConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "network_acls.#", "1"),
				),
			},
		},
	})
}
