package scaling_policy_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusScalingPoliciesDatasourceConfig = `
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
}

resource "byteplus_scaling_policy" "foo" {
  count = 3
  active = false
  scaling_group_id = "${byteplus_scaling_group.foo.id}"
  scaling_policy_name = "acc-tf-sg-policy-test-${count.index}"
  scaling_policy_type = "Alarm"
  adjustment_type = "QuantityChangeInCapacity"
  adjustment_value = 100
  cooldown = 10
  alarm_policy_rule_type = "Static"
  alarm_policy_evaluation_count = 1
  alarm_policy_condition_metric_name = "Instance_CpuBusy_Avg"
  alarm_policy_condition_metric_unit = "Percent"
  alarm_policy_condition_comparison_operator = "="
  alarm_policy_condition_threshold = 100
}

data "byteplus_scaling_policies" "foo"{
    scaling_group_id = "${byteplus_scaling_group.foo.id}"
	ids = byteplus_scaling_policy.foo[*].scaling_policy_id
}
`

func TestAccByteplusScalingPoliciesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_scaling_policies.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_policy.ByteplusScalingPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusScalingPoliciesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_policies.#", "3"),
				),
			},
		},
	})
}
