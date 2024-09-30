package iam_user_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusIamUserCreateConfig = `
resource "byteplus_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
}
`

const testAccByteplusIamUserUpdateConfig = `
resource "byteplus_iam_user" "foo" {
    description = "acc test update"
    display_name = "name2"
    email = "xxx@163.com"
    user_name = "acc-test-user2"
}
`

func TestAccByteplusIamUserResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_user.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user.ByteplusIamUserService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamUserCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
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

func TestAccByteplusIamUserResource_Update(t *testing.T) {
	resourceName := "byteplus_iam_user.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user.ByteplusIamUserService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamUserCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				Config: testAccByteplusIamUserUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc test update"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "name2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "email", "xxx@163.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "mobile_phone", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user2"),
				),
			},
			{
				Config:             testAccByteplusIamUserUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
