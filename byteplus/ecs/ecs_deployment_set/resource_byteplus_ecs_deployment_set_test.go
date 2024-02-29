package ecs_deployment_set_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsDeploymentSetCreateConfig = `
resource "byteplus_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
}
`

func TestAccByteplusEcsDeploymentSetResource_Basic(t *testing.T) {
	resourceName := "byteplus_ecs_deployment_set.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.ByteplusEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsDeploymentSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"description"},
			},
		},
	})
}

const testAccByteplusEcsDeploymentSetUpdateConfig = `
resource "byteplus_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds-new"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
}
`

func TestAccByteplusEcsDeploymentSetResource_Update(t *testing.T) {
	resourceName := "byteplus_ecs_deployment_set.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.ByteplusEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsDeploymentSetCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				Config: testAccByteplusEcsDeploymentSetUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_set_name", "acc-test-ecs-ds-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "granularity", "switch"),
					resource.TestCheckResourceAttr(acc.ResourceId, "strategy", "Availability"),
				),
			},
			{
				Config:             testAccByteplusEcsDeploymentSetUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
