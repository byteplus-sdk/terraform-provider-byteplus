package iam_user_group_attachment_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group_attachment"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamUserGroupAttachmentCreateConfig = `
resource "byteplus_iam_user" "foo" {
  user_name = "acc-test-user"
  description = "acc test"
  display_name = "name"
}

resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description = "acc-test"
  display_name = "acctest"
}

resource "byteplus_iam_user_group_attachment" "foo" {
    user_group_name = byteplus_iam_user_group.foo.user_group_name
    user_name = byteplus_iam_user.foo.user_name
}
`

func TestAccByteplusIamUserGroupAttachmentResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_user_group_attachment.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return iam_user_group_attachment.NewIamUserGroupAttachmentService(client)
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
				Config: testAccByteplusIamUserGroupAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_group_name", "acc-test-group"),
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
