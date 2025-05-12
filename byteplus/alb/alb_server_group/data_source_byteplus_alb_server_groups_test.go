package alb_server_group_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb_server_group"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusAlbServerGroupsDatasourceConfig = `
resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_alb_server_group" "foo" {
  vpc_id = byteplus_vpc.foo.id
  server_group_name = "acc-test-server-group-${count.index}"
  description = "acc-test"
  server_group_type = "instance"
  scheduler = "sh"
  project_name = "default"
  health_check {
    enabled = "on"
    interval = 3
    timeout = 3
    method = "GET"
  }
  sticky_session_config {
    sticky_session_enabled = "on"
    sticky_session_type = "insert"
    cookie_timeout = "1100"
  }
  count = 3
}

data "byteplus_alb_server_groups" "foo"{
    ids = byteplus_alb_server_group.foo[*].id
}
`

func TestAccByteplusAlbServerGroupsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_alb_server_groups.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return alb_server_group.NewAlbServerGroupService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAlbServerGroupsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "server_groups.#", "3"),
				),
			},
		},
	})
}
