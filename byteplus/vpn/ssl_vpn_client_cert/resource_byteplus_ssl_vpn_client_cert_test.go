package ssl_vpn_client_cert_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/ssl_vpn_client_cert"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusSslVpnClientCertCreateConfig = `
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

resource "byteplus_ssl_vpn_client_cert" "foo" {
  ssl_vpn_server_id = byteplus_ssl_vpn_server.foo.id
  ssl_vpn_client_cert_name = "acc-test-client-cert"
  description = "acc-test"
}
`

func TestAccByteplusSslVpnClientCertResource_Basic(t *testing.T) {
	resourceName := "byteplus_ssl_vpn_client_cert.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_vpn_client_cert.NewSslVpnClientCertService(client)
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
				Config: testAccByteplusSslVpnClientCertCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_client_cert_name", "acc-test-client-cert"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "certificate_status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "open_vpn_client_config"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_vpn_server_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ca_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_key"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "creation_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "expired_time"),
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

const testAccByteplusSslVpnClientCertUpdateConfig = `
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

resource "byteplus_ssl_vpn_client_cert" "foo" {
  ssl_vpn_server_id = byteplus_ssl_vpn_server.foo.id
  ssl_vpn_client_cert_name = "acc-test-client-cert-new"
  description = "acc-test-new"
}
`

func TestAccByteplusSslVpnClientCertResource_Update(t *testing.T) {
	resourceName := "byteplus_ssl_vpn_client_cert.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_vpn_client_cert.NewSslVpnClientCertService(client)
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
				Config: testAccByteplusSslVpnClientCertCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_client_cert_name", "acc-test-client-cert"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "certificate_status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "open_vpn_client_config"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_vpn_server_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ca_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_key"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "creation_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "expired_time"),
				),
			},
			{
				Config: testAccByteplusSslVpnClientCertUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_vpn_client_cert_name", "acc-test-client-cert-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "status", "Available"),
					resource.TestCheckResourceAttr(acc.ResourceId, "certificate_status", "Available"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "open_vpn_client_config"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_vpn_server_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ca_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_certificate"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "client_key"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "creation_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "expired_time"),
				),
			},
			{
				Config:             testAccByteplusSslVpnClientCertUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
