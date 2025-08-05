package allow_list_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/allow_list"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusMongodbAllowListsDatasourceConfig = `
resource "byteplus_mongodb_allow_list" "foo" {
    allow_list_name="acc-test"
    allow_list_desc="acc-test"
    allow_list_type="IPv4"
    allow_list="10.1.1.3,10.2.3.0/24,10.1.1.1"
}

data "byteplus_mongodb_allow_lists" "foo"{
    allow_list_ids = [byteplus_mongodb_allow_list.foo.id]
    region_id = "cn-beijing"
}
`

func TestAccByteplusMongodbAllowListsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_mongodb_allow_lists.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allow_list.NewMongoDBAllowListService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusMongodbAllowListsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_lists.#", "1"),
				),
			},
		},
	})
}
