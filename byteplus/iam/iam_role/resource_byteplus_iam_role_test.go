package iam_role_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_role"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamRoleCreateConfig = `
resource "byteplus_iam_role" "foo" {
	role_name = "acc-test-role"
    display_name = "acc-test"
	description = "acc-test"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"
	max_session_duration = 3600
}
`

func TestAccByteplusIamRoleResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_role.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.ByteplusIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamRoleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3600"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
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

const testAccByteplusIamRoleUpdateConfig = `
resource "byteplus_iam_role" "foo" {
    role_name = "acc-test-role-new"
    display_name = "acc-test-new"
	description = "acc-test-new"
    trust_policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"
	max_session_duration = 3700
}
`

func TestAccByteplusIamRoleResource_Update(t *testing.T) {
	resourceName := "byteplus_iam_role.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_role.ByteplusIamRoleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamRoleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3600"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"auto_scaling\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
				),
			},
			{
				Config: testAccByteplusIamRoleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "display_name", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "max_session_duration", "3700"),
					resource.TestCheckResourceAttr(acc.ResourceId, "role_name", "acc-test-role-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "trust_policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"sts:AssumeRole\"],\"Principal\":{\"Service\":[\"ecs\"]}}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "trn"),
				),
			},
			{
				Config:             testAccByteplusIamRoleUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
