package iam_user_group_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamUserGroupsDatasourceConfig = `
resource "byteplus_iam_user_group" "foo" {
  user_group_name = "acc-test-group"
  description = "acc-test"
  display_name = "acc-test"
}

data "byteplus_iam_user_groups" "foo"{
    query = "acc-test-group"
}
`

func TestAccByteplusIamUserGroupsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_iam_user_groups.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return iam_user_group.NewIamUserGroupService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamUserGroupsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(acc.ResourceId, "user_groups"),
				),
			},
		},
	})
}
