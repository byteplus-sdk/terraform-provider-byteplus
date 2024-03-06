package acl_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusAclCreateConfig = `
resource "byteplus_acl" "foo" {
	acl_name = "acc-test-acl"
	description = "acc-test-demo"
	project_name = "default"
	acl_entries {
    	entry = "172.20.1.0/24"
    	description = "e1"
  	}
}
`

func TestAccByteplusAclResource_Basic(t *testing.T) {
	resourceName := "byteplus_acl.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.ByteplusAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAclCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.1.0/24",
						"description": "e1",
					}),
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

const testAccByteplusAclUpdateConfig = `
resource "byteplus_acl" "foo" {
    acl_name = "acc-test-acl-new"
    description = "acc-test-demo-new"
    project_name = "default"
	acl_entries {
    	entry = "172.20.2.0/24"
    	description = "e2"
  	}
	acl_entries {
    	entry = "172.20.3.0/24"
    	description = "e3"
  	}
}
`

func TestAccByteplusAclResource_Update(t *testing.T) {
	resourceName := "byteplus_acl.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &acl.ByteplusAclService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAclCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.1.0/24",
						"description": "e1",
					}),
				),
			},
			{
				Config: testAccByteplusAclUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_name", "acc-test-acl-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-demo-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_time"),
					resource.TestCheckResourceAttr(acc.ResourceId, "acl_entries.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.2.0/24",
						"description": "e2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "acl_entries.*", map[string]string{
						"entry":       "172.20.3.0/24",
						"description": "e3",
					}),
				),
			},
			{
				Config:             testAccByteplusAclUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
