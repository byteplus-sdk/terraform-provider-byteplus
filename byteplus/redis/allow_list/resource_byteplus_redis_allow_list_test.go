package allow_list_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/redis/allow_list"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRedisAllowListCreateConfig = `
resource "byteplus_redis_allow_list" "foo" {
    allow_list = ["192.168.0.0/24"]
    allow_list_name = "acc-test-allowlist"
}
`

const testAccByteplusRedisAllowListUpdateConfig = `
resource "byteplus_redis_allow_list" "foo" {
    allow_list = ["192.168.0.0/24", "192.168.1.0/24"]
    allow_list_desc = "acctest"
    allow_list_name = "acc-test-allowlist1"
}
`

func TestAccByteplusRedisAllowListResource_Basic(t *testing.T) {
	resourceName := "byteplus_redis_allow_list.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allow_list.NewRedisAllowListService(client)
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
				Config: testAccByteplusRedisAllowListCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist"),
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

func TestAccByteplusRedisAllowListResource_Update(t *testing.T) {
	resourceName := "byteplus_redis_allow_list.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allow_list.NewRedisAllowListService(client)
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
				Config: testAccByteplusRedisAllowListCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist"),
				),
			},
			{
				Config: testAccByteplusRedisAllowListUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_desc", "acctest"),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_name", "acc-test-allowlist1"),
				),
			},
			{
				Config:             testAccByteplusRedisAllowListUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
