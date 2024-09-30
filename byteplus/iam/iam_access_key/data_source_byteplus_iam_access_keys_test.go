package iam_access_key_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_access_key"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamAccessKeysDatasourceConfig = `
data "byteplus_iam_access_keys" "foo"{
  user_name = "inner-user"
}
`

func TestAccByteplusIamAccessKeysDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_iam_access_keys.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.ByteplusIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamAccessKeysDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "access_key_metadata.#", "1"),
				),
			},
		},
	})
}
