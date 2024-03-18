package scaling_lifecycle_hook_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_lifecycle_hook"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingLifecycleHooksDatasourceConfig = `
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
    count = 3
    lifecycle_hook_name = "acc-test-lifecycle-${count.index}"
    lifecycle_hook_policy = "CONTINUE"
    lifecycle_hook_timeout = 30
    lifecycle_hook_type = "SCALE_IN"
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
}

data "byteplus_scaling_lifecycle_hooks" "foo"{
    ids = byteplus_scaling_lifecycle_hook.foo[*].lifecycle_hook_id
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
}
`

func TestAccByteplusScalingLifecycleHooksDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_scaling_lifecycle_hooks.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_lifecycle_hook.ByteplusScalingLifecycleHookService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingLifecycleHooksDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "lifecycle_hooks.#", "3"),
				),
			},
		},
	})
}
