package acl_entry_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl_entry"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusAclEntryCreateConfig = `
resource "byteplus_acl" "foo" {
	acl_name = "acc-test-acl"
	description = "acc-test-demo"
	project_name = "default"
}

resource "byteplus_acl_entry" "foo" {
    acl_id = "${byteplus_acl.foo.id}"
    entry = "172.20.1.0/24"
	description = "entry"
}
`

func TestAccByteplusAclEntryResource_Basic(t *testing.T) {
	resourceName := "byteplus_acl_entry.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl_entry.ByteplusAclEntryService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAclEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "entry", "172.20.1.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "entry"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "acl_id"),
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
