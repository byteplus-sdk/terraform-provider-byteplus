package ecs_deployment_set_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsDeploymentSetsDatasourceConfig = `
resource "byteplus_ecs_deployment_set" "foo" {
    deployment_set_name = "acc-test-ecs-ds-${count.index}"
	description = "acc-test"
    granularity = "switch"
    strategy = "Availability"
	count = 3
}

data "byteplus_ecs_deployment_sets" "foo"{
    granularity = "switch"
    ids = byteplus_ecs_deployment_set.foo[*].id
}
`

func TestAccByteplusEcsDeploymentSetsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_ecs_deployment_sets.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_deployment_set.ByteplusEcsDeploymentSetService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsDeploymentSetsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "deployment_sets.#", "3"),
				),
			},
		},
	})
}
