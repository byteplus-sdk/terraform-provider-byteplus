package ecs_available_resource_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_available_resource"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusAvailableResourcesDatasourceConfig = `
data "byteplus_ecs_available_resources" "foo"{
    destination_resource = "InstanceType"
}
`

func TestAccByteplusAvailableResourcesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_ecs_available_resources.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ecs_available_resource.NewEcsAvailableResourceService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAvailableResourcesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "available_zones.#", "3"),
				),
			},
		},
	})
}
