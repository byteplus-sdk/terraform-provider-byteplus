package ecs_launch_template_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_launch_template"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsLaunchTemplateCreateConfig = `
resource "byteplus_ecs_launch_template" "foo" {
    launch_template_name = "acc-test-template"
}
`

const testAccByteplusEcsLaunchTemplateUpdateConfig = `
resource "byteplus_ecs_launch_template" "foo" {
    description = "acc-test-desc"
    eip_bandwidth = 1
    eip_billing_type = "PostPaidByBandwidth"
    eip_isp = "ChinaMobile"
    host_name = "acc-xx"
    hpc_cluster_id = "acc-xx"
    image_id = "acc-xx"
    instance_charge_type = "acc-xx"
    instance_name = "acc-xx"
    instance_type_id = "acc-xx"
    key_pair_name = "acc-xx"
    launch_template_name = "acc-test-template2"
}
`

func TestAccByteplusEcsLaunchTemplateResource_Basic(t *testing.T) {
	resourceName := "byteplus_ecs_launch_template.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_launch_template.ByteplusEcsLaunchTemplateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsLaunchTemplateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "hpc_cluster_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "image_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "launch_template_name", "acc-test-template"),
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

func TestAccByteplusEcsLaunchTemplateResource_Update(t *testing.T) {
	resourceName := "byteplus_ecs_launch_template.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_launch_template.ByteplusEcsLaunchTemplateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsLaunchTemplateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "hpc_cluster_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "image_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type_id", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "launch_template_name", "acc-test-template"),
				),
			},
			{
				Config: testAccByteplusEcsLaunchTemplateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-desc"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_bandwidth", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "eip_isp", "ChinaMobile"),
					resource.TestCheckResourceAttr(acc.ResourceId, "host_name", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "hpc_cluster_id", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "image_id", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_charge_type", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_type_id", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-xx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "launch_template_name", "acc-test-template2"),
				),
			},
			{
				Config:             testAccByteplusEcsLaunchTemplateUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
