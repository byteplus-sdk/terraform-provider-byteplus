package iam_access_key_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_access_key"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamAccessKeyCreateConfig = `
resource "byteplus_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "byteplus_iam_access_key" "foo" {
	user_name = "${byteplus_iam_user.foo.user_name}"
    secret_file = "./sk"
    status = "active"
}
`

func TestAccByteplusIamAccessKeyResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_access_key.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.ByteplusIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamAccessKeyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "active"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
		},
	})
}

const testAccByteplusIamAccessKeyUpdateConfig = `
resource "byteplus_iam_user" "foo" {
  	user_name = "acc-test-user"
  	description = "acc-test"
  	display_name = "name"
}

resource "byteplus_iam_access_key" "foo" {
	user_name = "${byteplus_iam_user.foo.user_name}"
    secret_file = "./sk"
    status = "inactive"
}
`

func TestAccByteplusIamAccessKeyResource_Update(t *testing.T) {
	resourceName := "byteplus_iam_access_key.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_access_key.ByteplusIamAccessKeyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamAccessKeyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "active"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
			{
				Config: testAccByteplusIamAccessKeyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "user_name", "acc-test-user"),
					resource.TestCheckResourceAttr(acc.ResourceId, "secret_file", "./sk"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "inactive"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "pgp_key"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "encrypted_secret"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "key_fingerprint"),
				),
			},
			{
				Config:             testAccByteplusIamAccessKeyUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
