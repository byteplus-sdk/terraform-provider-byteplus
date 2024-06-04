package scaling_group_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_group"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccScalingGroupDatasourceConfig = `
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
  count = 3
  scaling_group_name = "acc-test-scaling-group-${count.index}"
  subnet_ids = ["${byteplus_subnet.foo.id}"]
  multi_az_policy = "BALANCE"
  desire_instance_number = 0
  min_instance_number = 0
  max_instance_number = 10
  instance_terminate_policy = "OldestInstance"
  default_cooldown = 30

  tags {
    key = "k2"
    value = "v2"
  }

  tags {
    key = "k1"
    value = "v1"
  }
}

data "byteplus_scaling_groups" "foo"{
  ids = byteplus_scaling_group.foo[*].id
}
`

func TestAccByteplusScalingGroupDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_scaling_groups.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &scaling_group.ByteplusScalingGroupService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccScalingGroupDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "scaling_groups.#", "3"),
				),
			},
		},
	})
}
