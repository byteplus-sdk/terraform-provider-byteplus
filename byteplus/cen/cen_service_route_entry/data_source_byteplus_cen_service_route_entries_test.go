package cen_service_route_entry_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_service_route_entry"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusCenServiceRouteEntriesDatasourceConfig = `
resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
  count      = 3
}

resource "byteplus_cen" "foo" {
  cen_name     = "acc-test-cen"
  description  = "acc-test"
  project_name = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_cen_attach_instance" "foo" {
  cen_id             = byteplus_cen.foo.id
  instance_id        = byteplus_vpc.foo[count.index].id
  instance_region_id = "cn-beijing"
  instance_type      = "VPC"
  count              = 3
}

resource "byteplus_cen_service_route_entry" "foo" {
  cen_id                 = byteplus_cen.foo.id
  destination_cidr_block = "100.64.0.0/11"
  service_region_id      = "cn-beijing"
  service_vpc_id         = byteplus_cen_attach_instance.foo[0].instance_id
  description            = "acc-test"
  publish_mode           = "Custom"
  publish_to_instances {
    instance_region_id = "cn-beijing"
    instance_type      = "VPC"
    instance_id        = byteplus_cen_attach_instance.foo[1].instance_id
  }
  publish_to_instances {
    instance_region_id = "cn-beijing"
    instance_type      = "VPC"
    instance_id        = byteplus_cen_attach_instance.foo[2].instance_id
  }
}

data "byteplus_cen_service_route_entries" "foo"{
  cen_id                 = byteplus_cen.foo.id
  destination_cidr_block = byteplus_cen_service_route_entry.foo.destination_cidr_block
}
`

func TestAccByteplusCenServiceRouteEntriesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_cen_service_route_entries.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return cen_service_route_entry.NewCenServiceRouteEntryService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCenServiceRouteEntriesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "service_route_entries.#", "1"),
				),
			},
		},
	})
}
