package ssl_vpn_server_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/ssl_vpn_server"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusSslVpnServerCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "172.16.0.0/24"
  zone_id = data.byteplus_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id = byteplus_vpc.foo.id
  subnet_id = byteplus_subnet.foo.id
  bandwidth = 5
  vpn_gateway_name = "acc-test1"
  description = "acc-test1"
  period = 7
  project_name = "default"
  ssl_enabled = true
  ssl_max_connections = 5
}

resource "byteplus_ssl_vpn_server" "foo" {
  vpn_gateway_id = byteplus_vpn_gateway.foo.id
  local_subnets = [byteplus_subnet.foo.cidr_block]
  client_ip_pool = "172.16.2.0/24"
  ssl_vpn_server_name = "acc-test-ssl"
  description = "acc-test"
  protocol = "UDP"
  cipher = "AES-128-CBC"
  auth = "SHA1"
  compress = true
}
`

const testAccByteplusSslVpnServerUpdateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block = "172.16.0.0/24"
  zone_id = data.byteplus_zones.foo.zones[0].id
  vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id = byteplus_vpc.foo.id
  subnet_id = byteplus_subnet.foo.id
  bandwidth = 5
  vpn_gateway_name = "acc-test1"
  description = "acc-test1"
  period = 7
  project_name = "default"
  ssl_enabled = true
  ssl_max_connections = 5
}

resource "byteplus_ssl_vpn_server" "foo" {
  vpn_gateway_id = byteplus_vpn_gateway.foo.id
  local_subnets = [byteplus_subnet.foo.cidr_block]
  client_ip_pool = "172.16.2.0/24"
  ssl_vpn_server_name = "acc-test-ssl1"
  description = "acc-test1"
  protocol = "UDP"
  cipher = "AES-128-CBC"
  auth = "SHA1"
  compress = true
}
`

func TestAccByteplusSslVpnServerResource_Basic(t *testing.T) {
	resourceName := "byteplus_ssl_vpn_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_vpn_server.NewSslVpnServerService(client)
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
				Config: testAccByteplusSslVpnServerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "auth", "SHA1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(acc.ResourceId, "client_ip_pool", "172.16.2.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "compress", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "UDP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_server_name", "acc-test-ssl"),
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

func TestAccByteplusSslVpnServerResource_Update(t *testing.T) {
	resourceName := "byteplus_ssl_vpn_server.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_vpn_server.NewSslVpnServerService(client)
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
				Config: testAccByteplusSslVpnServerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "auth", "SHA1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(acc.ResourceId, "client_ip_pool", "172.16.2.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "compress", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "UDP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_server_name", "acc-test-ssl"),
				),
			},
			{
				Config: testAccByteplusSslVpnServerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "auth", "SHA1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cipher", "AES-128-CBC"),
					resource.TestCheckResourceAttr(acc.ResourceId, "client_ip_pool", "172.16.2.0/24"),
					resource.TestCheckResourceAttr(acc.ResourceId, "compress", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "local_subnets.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "UDP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_server_name", "acc-test-ssl1"),
				),
			},
			{
				Config:             testAccByteplusSslVpnServerUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
