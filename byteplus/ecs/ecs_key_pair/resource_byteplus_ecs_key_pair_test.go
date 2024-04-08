package ecs_key_pair_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_key_pair"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusEcsKeyPairCreateConfig = `
resource "byteplus_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description ="acc-test"
}
`

const testAccByteplusEcsKeyPairUpdateConfig = `
resource "byteplus_ecs_key_pair" "foo" {
    description = "acc-test-2"
    key_pair_name = "acc-test-key-name"
}
`

func TestAccByteplusEcsKeyPairResource_Basic(t *testing.T) {
	resourceName := "byteplus_ecs_key_pair.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.ByteplusEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsKeyPairCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
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

func TestAccByteplusEcsKeyPairResource_Update(t *testing.T) {
	resourceName := "byteplus_ecs_key_pair.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &ecs_key_pair.ByteplusEcsKeyPairService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusEcsKeyPairCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
				),
			},
			{
				Config: testAccByteplusEcsKeyPairUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "key_pair_name", "acc-test-key-name"),
				),
			},
			{
				Config:             testAccByteplusEcsKeyPairUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
