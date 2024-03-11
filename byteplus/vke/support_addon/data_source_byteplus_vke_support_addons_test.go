package support_addon_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/support_addon"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVkeSupportAddonsDatasourceConfig = `
data "byteplus_vke_support_addons" "foo"{
}
`

func TestAccByteplusVkeSupportAddonsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vke_support_addons.foo"

	_ = &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &support_addon.ByteplusVkeSupportAddonService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVkeSupportAddonsDatasourceConfig,
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}
