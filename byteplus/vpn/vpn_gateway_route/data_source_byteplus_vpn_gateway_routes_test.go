package vpn_gateway_route_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_gateway_route"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVpnGatewayRoutesDatasourceConfig = `
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

resource "byteplus_vpn_gateway_route" "foo" {
  vpn_gateway_id = "${byteplus_vpn_gateway.foo.id}"
  destination_cidr_block = "192.168.0.0/20"
  next_hop_id = "${byteplus_vpn_connection.foo.id}"
}

data "byteplus_vpn_gateway_routes" "foo"{
    ids = ["${byteplus_vpn_gateway_route.foo.id}"]
}
`

func TestAccByteplusVpnGatewayRoutesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vpn_gateway_routes.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &vpn_gateway_route.ByteplusVpnGatewayRouteService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVpnGatewayRoutesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "vpn_gateway_routes.#", "1"),
				),
			},
		},
	})
}
