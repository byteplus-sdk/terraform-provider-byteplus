package subnet_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/subnet"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccSubnetForCreate = `
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
`

const testAccSubnetForUpdate = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet_new"
  description = "tfdesc"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  vpc_id = "${byteplus_vpc.foo.id}"
}
`

func TestAccByteplusSubnetResource_Basic(t *testing.T) {
	resourceName := "byteplus_subnet.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.ByteplusSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
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

func TestAccByteplusSubnetResource_Update(t *testing.T) {
	resourceName := "byteplus_subnet.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &subnet.ByteplusSubnetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccSubnetForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config: testAccSubnetForUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "tfdesc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "subnet_name", "acc-test-subnet_new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cidr_block", "172.16.0.0/24"),
					// compute status
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
				),
			},
			{
				Config:             testAccSubnetForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
