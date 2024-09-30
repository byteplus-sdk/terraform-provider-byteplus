package iam_user_policy_attachment_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_policy_attachment"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusIamUserPolicyAttachmentCreateConfig = `
resource "byteplus_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
}
resource "byteplus_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}
resource "byteplus_iam_user_policy_attachment" "foo" {
    policy_name = byteplus_iam_policy.foo.policy_name
    policy_type = "Custom"
    user_name = byteplus_iam_user.foo.user_name
}
`

func TestAccByteplusIamUserPolicyAttachmentResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_user_policy_attachment.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_user_policy_attachment.ByteplusIamUserPolicyAttachmentService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamUserPolicyAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_type", "Custom"),
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
