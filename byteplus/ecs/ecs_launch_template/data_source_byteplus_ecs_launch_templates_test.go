package ecs_launch_template_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_launch_template"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsLaunchTemplatesDatasourceConfig = `
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

data "byteplus_ecs_launch_templates" "foo"{
    ids = ["${byteplus_ecs_launch_template.foo.id}"]
}
`

func TestAccByteplusEcsLaunchTemplatesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_ecs_launch_templates.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_launch_template.ByteplusEcsLaunchTemplateService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsLaunchTemplatesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "launch_templates.#", "1"),
				),
			},
		},
	})
}
