package scaling_group_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_group"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccScalingGroupForCreate = `
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

resource "byteplus_scaling_group" "foo" {
  scaling_group_name = "acc-test-scaling-group"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 1
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 10
  scaling_mode = "recycle"
}
`

const testAccScalingGroupForUpdate = `
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

resource "byteplus_scaling_group" "foo" {
  scaling_group_name = "acc-test-scaling-group-new"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 10
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 30
  scaling_mode = "recycle"
  tags {
    key = "k2"
    value = "v2"
  }

  tags {
    key = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusScalingGroupResource_Basic(t *testing.T) {
	resourceName := "byteplus_scaling_group.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_group.ByteplusScalingGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccScalingGroupForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_group_name", "acc-test-scaling-group"),
					resource.TestCheckResourceAttr(acc.ResourceId, "multi_az_policy", "BALANCE"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_terminate_policy", "OldestInstance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_instance_number", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "default_cooldown", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_mode", "recycle"),
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

func TestAccByteplusScalingGroupResource_Update(t *testing.T) {
	resourceName := "byteplus_scaling_group.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_group.ByteplusScalingGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccScalingGroupForCreate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_group_name", "acc-test-scaling-group"),
					resource.TestCheckResourceAttr(acc.ResourceId, "multi_az_policy", "BALANCE"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_terminate_policy", "OldestInstance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_instance_number", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "default_cooldown", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_mode", "recycle"),
				),
			},
			{
				Config: testAccScalingGroupForUpdate,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_group_name", "acc-test-scaling-group-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "multi_az_policy", "BALANCE"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_terminate_policy", "OldestInstance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_instance_number", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "default_cooldown", "30"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_mode", "recycle"),
				),
			},
			{
				Config:             testAccScalingGroupForUpdate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
