package iam_login_profile_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_login_profile"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamLoginProfileCreateConfig = `
resource "byteplus_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "byteplus_iam_login_profile" "foo" {
    user_name = "${byteplus_iam_user.foo.user_name}"
  	password = "93f0cb0614Aab12"
  	login_allowed = true
	password_reset_required = false
}
`

func TestAccByteplusIamLoginProfileResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_login_profile.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_login_profile.ByteplusIamLoginProfileService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamLoginProfileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

const testAccByteplusIamLoginProfileUpdateConfig = `
resource "byteplus_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "byteplus_iam_login_profile" "foo" {
    user_name = "${byteplus_iam_user.foo.user_name}"
  	password = "93f0cb0614Aab12177"
  	login_allowed = false
	password_reset_required = true
}
`

func TestAccByteplusIamLoginProfileResource_Update(t *testing.T) {
	resourceName := "byteplus_iam_login_profile.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_login_profile.ByteplusIamLoginProfileService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamLoginProfileCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				Config: testAccByteplusIamLoginProfileUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "login_allowed", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password", "93f0cb0614Aab12177"),
					resource.TestCheckResourceAttr(acc.ResourceId, "password_reset_required", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
				),
			},
			{
				Config:             testAccByteplusIamLoginProfileUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
