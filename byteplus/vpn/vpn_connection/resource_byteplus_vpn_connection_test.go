package vpn_connection_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_connection"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVpnConnectionCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  subnet_id = "${byteplus_subnet.foo.id}"
  bandwidth = 20
  vpn_gateway_name = "acc-test"
  description = "acc-test"
  period = 2
  project_name = "default"
}

resource "byteplus_customer_gateway" "foo" {
  ip_address = "192.0.1.3"
  customer_gateway_name = "acc-test"
  description = "acc-test"
  project_name = "default"
}

resource "byteplus_vpn_connection" "foo" {
  vpn_connection_name = "acc-tf-test"
  description = "acc-tf-test"
  vpn_gateway_id = "${byteplus_vpn_gateway.foo.id}"
  customer_gateway_id = "${byteplus_customer_gateway.foo.id}"
  local_subnet = ["192.168.0.0/22"]
  remote_subnet = ["192.161.0.0/20"]
  dpd_action = "none"
  nat_traversal = true
  ike_config_psk = "acctest@!3"
  ike_config_version = "ikev1"
  ike_config_mode = "main"
  ike_config_enc_alg = "aes"
  ike_config_auth_alg = "md5"
  ike_config_dh_group = "group2"
  ike_config_lifetime = 9000
  ike_config_local_id = "acc_test"
  ike_config_remote_id = "acc_test"
  ipsec_config_enc_alg = "aes"
  ipsec_config_auth_alg = "sha256"
  ipsec_config_dh_group = "group2"
  ipsec_config_lifetime = 9000
  project_name = "default"
  log_enabled = false
}

`

const testAccByteplusVpnConnectionUpdateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
	  vpc_name   = "acc-test-vpc"
	  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
	  subnet_name = "acc-test-subnet"
	  cidr_block = "172.16.0.0/24"
	  zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	  vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_vpn_gateway" "foo" {
  vpc_id = "${byteplus_vpc.foo.id}"
  subnet_id = "${byteplus_subnet.foo.id}"
  bandwidth = 20
  vpn_gateway_name = "acc-test"
  description = "acc-test"
  period = 2
  project_name = "default"
}

resource "byteplus_customer_gateway" "foo" {
  ip_address = "192.0.1.3"
  customer_gateway_name = "acc-test"
  description = "acc-test"
  project_name = "default"
}

resource "byteplus_vpn_connection" "foo" {
  vpn_connection_name = "acc-tf-test1"
  description = "acc-tf-test1"
  vpn_gateway_id = "${byteplus_vpn_gateway.foo.id}"
  customer_gateway_id = "${byteplus_customer_gateway.foo.id}"
  local_subnet = ["192.168.0.0/22"]
  remote_subnet = ["192.161.0.0/20"]
  dpd_action = "clear"
  nat_traversal = true
  ike_config_psk = "acctest@!31"
  ike_config_version = "ikev1"
  ike_config_mode = "main"
  ike_config_enc_alg = "aes"
  ike_config_auth_alg = "md5"
  ike_config_dh_group = "group2"
  ike_config_lifetime = 9000
  ike_config_local_id = "acc_test"
  ike_config_remote_id = "acc_test"
  ipsec_config_enc_alg = "aes"
  ipsec_config_auth_alg = "sha256"
  ipsec_config_dh_group = "group1"
  ipsec_config_lifetime = 10000
  project_name = "default"
  log_enabled = true
}

`

func TestAccByteplusVpnConnectionResource_Basic(t *testing.T) {
	resourceName := "byteplus_vpn_connection.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpn_connection.ByteplusVpnConnectionService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVpnConnectionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-tf-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "dpd_action", "none"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_auth_alg", "md5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_dh_group", "group2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_lifetime", "9000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_local_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_mode", "main"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_psk", "acctest@!3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_remote_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_version", "ikev1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_auth_alg", "sha256"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_dh_group", "group2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_lifetime", "9000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "log_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "nat_traversal", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpn_connection_name", "acc-tf-test"),
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

func TestAccByteplusVpnConnectionResource_Update(t *testing.T) {
	resourceName := "byteplus_vpn_connection.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpn_connection.ByteplusVpnConnectionService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVpnConnectionCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-tf-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "dpd_action", "none"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_auth_alg", "md5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_dh_group", "group2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_lifetime", "9000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_local_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_mode", "main"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_psk", "acctest@!3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_remote_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_version", "ikev1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_auth_alg", "sha256"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_dh_group", "group2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_lifetime", "9000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "log_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "nat_traversal", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpn_connection_name", "acc-tf-test"),
				),
			},
			{
				Config: testAccByteplusVpnConnectionUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-tf-test1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "dpd_action", "clear"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_auth_alg", "md5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_dh_group", "group2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_lifetime", "9000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_local_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_mode", "main"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_psk", "acctest@!31"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_remote_id", "acc_test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ike_config_version", "ikev1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_auth_alg", "sha256"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_dh_group", "group1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_enc_alg", "aes"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ipsec_config_lifetime", "10000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "log_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "nat_traversal", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "vpn_connection_name", "acc-tf-test1"),
				),
			},
			{
				Config:             testAccByteplusVpnConnectionUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
