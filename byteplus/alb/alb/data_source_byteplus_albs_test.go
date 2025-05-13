package alb_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/alb/alb"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusAlbsDatasourceConfig = `
data "byteplus_alb_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "subnet_1" {
  subnet_name = "acc-test-subnet-1"
  cidr_block = "172.16.1.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_subnet" "subnet_2" {
  subnet_name = "acc-test-subnet-2"
  cidr_block = "172.16.2.0/24"
  zone_id = data.byteplus_alb_zones.foo.zones[1].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_alb" "foo" {
  address_ip_version = "IPv4"
  type = "private"
  load_balancer_name = "acc-test-alb-private-${count.index}"
  description = "acc-test"
  subnet_ids = [byteplus_subnet.subnet_1.id, byteplus_subnet.subnet_2.id]
  project_name = "default"
  delete_protection = "off"
  tags {
    key = "k1"
    value = "v1"
  }
  count = 3
}

data "byteplus_albs" "foo" {
  ids = byteplus_alb.foo[*].id
}
`

func TestAccByteplusAlbsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_albs.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return alb.NewAlbService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusAlbsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "albs.#", "3"),
				),
			},
		},
	})
}
