package acl_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusAclsDatasourceConfig = `
resource "byteplus_acl" "foo" {
	acl_name = "acc-test-acl-${count.index}"
	description = "acc-test-demo"
	project_name = "default"
	acl_entries {
    	entry = "172.20.1.0/24"
    	description = "e1"
  	}
	count = 3
}

data "byteplus_acls" "foo"{
    ids = byteplus_acl.foo[*].id
}
`

func TestAccByteplusAclsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_acls.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.ByteplusAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAclsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "acls.#", "3"),
				),
			},
		},
	})
}
