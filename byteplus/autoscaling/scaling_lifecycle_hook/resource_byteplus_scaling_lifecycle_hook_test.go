package scaling_lifecycle_hook_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_lifecycle_hook"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingLifecycleHookCreateConfig = `
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
  scaling_group_name = "acc-test-scaling-group-lifecycle"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 1
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 10
}

resource "byteplus_scaling_lifecycle_hook" "foo" {
    lifecycle_hook_name = "acc-test-lifecycle"
    lifecycle_hook_policy = "CONTINUE"
    lifecycle_hook_timeout = 30
    lifecycle_hook_type = "SCALE_IN"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
}
`

const testAccByteplusScalingLifecycleHookUpdateConfig = `
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

resource "byteplus_ecs_command" "foo" {
  name = "acc-test-command"
  description = "tf"
  working_dir = "/home"
  username = "root"
  timeout = 100
  command_content = "IyEvYmluL2Jhc2gKCgplY2hvICJvcGVyYXRpb24gc3VjY2VzcyEi"
}

resource "byteplus_scaling_group" "foo" {
  scaling_group_name = "acc-test-scaling-group-lifecycle"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 1
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 10
}

resource "byteplus_scaling_lifecycle_hook" "foo" {
    lifecycle_hook_name = "acc-test-lifecycle"
    lifecycle_hook_policy = "ROLLBACK"
    lifecycle_hook_timeout = 300
    lifecycle_hook_type = "SCALE_OUT"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
	lifecycle_command {
    command_id = byteplus_ecs_command.foo.id
    parameters = "{}"
  }
}
`

func TestAccByteplusScalingLifecycleHookResource_Basic(t *testing.T) {
	resourceName := "byteplus_scaling_lifecycle_hook.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_lifecycle_hook.ByteplusScalingLifecycleHookService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingLifecycleHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_name", "acc-test-lifecycle"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_policy", "CONTINUE"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_timeout", "30"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_type", "SCALE_IN"),
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

func TestAccByteplusScalingLifecycleHookResource_Update(t *testing.T) {
	resourceName := "byteplus_scaling_lifecycle_hook.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_lifecycle_hook.ByteplusScalingLifecycleHookService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingLifecycleHookCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_name", "acc-test-lifecycle"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_policy", "CONTINUE"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_timeout", "30"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_type", "SCALE_IN"),
				),
			},
			{
				Config: testAccByteplusScalingLifecycleHookUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_name", "acc-test-lifecycle"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_policy", "ROLLBACK"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_timeout", "300"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hook_type", "SCALE_OUT"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_command.#", "1"),
				),
			},
			{
				Config:             testAccByteplusScalingLifecycleHookUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
