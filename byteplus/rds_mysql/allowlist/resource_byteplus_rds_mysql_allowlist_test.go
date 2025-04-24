package allowlist_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/allowlist"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsMysqlAllowlistCreateConfig = `
resource "byteplus_rds_mysql_allowlist" "foo" {
    allow_list_name = "acc-test-allowlist"
	allow_list_desc = "acc-test"
	allow_list_type = "IPv4"
	allow_list = ["192.168.0.0/24", "192.168.1.0/24"]
}
`

func TestAccByteplusRdsMysqlAllowlistResource_Basic(t *testing.T) {
	resourceName := "byteplus_rds_mysql_allowlist.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allowlist.NewRdsMysqlAllowListService(client)
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
				Config: testAccByteplusRdsMysqlAllowlistCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_type", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "2"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.0.0/24"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "allow_list_id"),
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

const testAccByteplusRdsMysqlAllowlistUpdateConfig = `
resource "byteplus_rds_mysql_allowlist" "foo" {
    allow_list_name = "acc-test-allowlist-new"
	allow_list_desc = "acc-test-new"
	allow_list_type = "IPv4"
	allow_list = ["192.168.0.0/24", "192.168.3.0/24", "192.168.4.0/24"]
}
`

func TestAccByteplusRdsMysqlAllowlistResource_Update(t *testing.T) {
	resourceName := "byteplus_rds_mysql_allowlist.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allowlist.NewRdsMysqlAllowListService(client)
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
				Config: testAccByteplusRdsMysqlAllowlistCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_type", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "2"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.0.0/24"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.1.0/24"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "allow_list_id"),
				),
			},
			{
				Config: testAccByteplusRdsMysqlAllowlistUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_type", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "3"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.0.0/24"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.3.0/24"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "allow_list.*", "192.168.4.0/24"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "allow_list_id"),
				),
			},
			{
				Config:             testAccByteplusRdsMysqlAllowlistUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
