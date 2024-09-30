package iam_user_group_policy_attachment_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group_policy_attachment"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamUserGroupPolicyAttachmentCreateConfig = `
resource "byteplus_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description = "acc-test"
  display_name = "acc-test"
}

resource "byteplus_iam_user_group_policy_attachment" "foo" {
    policy_name = byteplus_iam_policy.foo.policy_name
    policy_type = "Custom"
    user_group_name = byteplus_iam_user_group.foo.user_group_name
}
`

func TestAccByteplusIamUserGroupPolicyAttachmentResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_user_group_policy_attachment.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return iam_user_group_policy_attachment.NewIamUserGroupPolicyAttachmentService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamUserGroupPolicyAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_type", "Custom"),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_group_name", "acc-test-group"),
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
