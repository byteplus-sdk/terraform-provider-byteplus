package cen_service_route_entry_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_service_route_entry"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusCenServiceRouteEntryCreateConfig = `
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
`

func TestAccByteplusCenServiceRouteEntryResource_Basic(t *testing.T) {
	resourceName := "byteplus_cen_service_route_entry.foo"

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
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCenServiceRouteEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "destination_cidr_block", "100.64.0.0/11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_mode", "Custom"),
					resource.TestCheckResourceAttr(acc.ResourceId, "service_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_to_instances.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "publish_to_instances.*", map[string]string{
						"instance_type":      "VPC",
						"instance_region_id": "cn-beijing",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cen_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "service_vpc_id"),
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

const testAccByteplusCenServiceRouteEntryUpdateConfig = `
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
  description            = "acc-test-new"
  publish_mode           = "LocalDCGW"
}
`

func TestAccByteplusCenServiceRouteEntryResource_Update(t *testing.T) {
	resourceName := "byteplus_cen_service_route_entry.foo"

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
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusCenServiceRouteEntryCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "destination_cidr_block", "100.64.0.0/11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_mode", "Custom"),
					resource.TestCheckResourceAttr(acc.ResourceId, "service_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_to_instances.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "publish_to_instances.*", map[string]string{
						"instance_type":      "VPC",
						"instance_region_id": "cn-beijing",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cen_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "service_vpc_id"),
				),
			},
			{
				Config: testAccByteplusCenServiceRouteEntryUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "destination_cidr_block", "100.64.0.0/11"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_mode", "LocalDCGW"),
					resource.TestCheckResourceAttr(acc.ResourceId, "service_region_id", "cn-beijing"),
					resource.TestCheckResourceAttr(acc.ResourceId, "publish_to_instances.#", "0"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cen_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "service_vpc_id"),
				),
			},
			{
				Config:             testAccByteplusCenServiceRouteEntryUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
