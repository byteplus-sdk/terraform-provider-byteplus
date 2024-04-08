package ecs_key_pair_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_key_pair"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsKeyPairsDatasourceConfig = `
resource "byteplus_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description ="acc-test"
}
data "byteplus_ecs_key_pairs" "foo"{
    key_pair_name = "${byteplus_ecs_key_pair.foo.key_pair_name}"
}
`

func TestAccByteplusEcsKeyPairsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_ecs_key_pairs.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.ByteplusEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsKeyPairsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pairs.#", "1"),
				),
			},
		},
	})
}
