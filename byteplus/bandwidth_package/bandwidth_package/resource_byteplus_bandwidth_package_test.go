package bandwidth_package_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/bandwidth_package/bandwidth_package"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusBandwidthPackageCreateConfig = `
resource "byteplus_bandwidth_package" "foo" {
  bandwidth_package_name    = "acc-test-bp"
  billing_type              = "PostPaidByBandwidth"
  isp                       = "BGP"
  description               = "acc-test"
  bandwidth                 = 10
  protocol                  = "IPv4"
  security_protection_types = ["AntiDDoS_Enhanced"]
  tags {
    key   = "k1"
    value = "v1"
  }
}
`

func TestAccByteplusBandwidthPackageResource_Basic(t *testing.T) {
	resourceName := "byteplus_bandwidth_package.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return bandwidth_package.NewBandwidthPackageService(client)
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
				Config: testAccByteplusBandwidthPackageCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth_package_name", "acc-test-bp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_protection_types.0", "AntiDDoS_Enhanced"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
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

const testAccByteplusBandwidthPackageUpdateConfig = `
resource "byteplus_bandwidth_package" "foo" {
  bandwidth_package_name    = "acc-test-bp-new"
  billing_type              = "PostPaidByBandwidth"
  isp                       = "BGP"
  description               = "acc-test-new"
  bandwidth                 = 5
  protocol                  = "IPv4"
  security_protection_types = ["AntiDDoS_Enhanced"]
  project_name              = "default"
  tags {
    key   = "k2"
    value = "v2"
  }
  tags {
    key   = "k3"
    value = "v3"
  }
}
`

func TestAccByteplusBandwidthPackageResource_Update(t *testing.T) {
	resourceName := "byteplus_bandwidth_package.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return bandwidth_package.NewBandwidthPackageService(client)
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
				Config: testAccByteplusBandwidthPackageCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth_package_name", "acc-test-bp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_protection_types.0", "AntiDDoS_Enhanced"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				Config: testAccByteplusBandwidthPackageUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth_package_name", "acc-test-bp-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PostPaidByBandwidth"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_protection_types.0", "AntiDDoS_Enhanced"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k3",
						"value": "v3",
					}),
				),
			},
			{
				Config:             testAccByteplusBandwidthPackageUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusBandwidthPackagePrePaidCreateConfig = `
resource "byteplus_bandwidth_package" "foo" {
  bandwidth_package_name    = "acc-test-bp"
  billing_type              = "PrePaid"
  period                    = "2"
  isp                       = "BGP"
  description               = "acc-test"
  bandwidth                 = 5
  protocol                  = "IPv4"
  security_protection_types = ["AntiDDoS_Enhanced"]
  tags {
    key   = "k1"
    value = "v1"
  }
  count = 3
}
`

func TestAccByteplusBandwidthPackageResource_PrePaid(t *testing.T) {
	resourceName := "byteplus_bandwidth_package.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return bandwidth_package.NewBandwidthPackageService(client)
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
				Config: testAccByteplusBandwidthPackagePrePaidCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "10"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth_package_name", "acc-test-bp"),
					resource.TestCheckResourceAttr(acc.ResourceId, "billing_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "isp", "BGP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "IPv4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "security_protection_types.0", "AntiDDoS_Enhanced"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}
