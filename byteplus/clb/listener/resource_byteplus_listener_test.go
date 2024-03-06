package listener_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/listener"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusListenerCreateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "byteplus_listener" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  listener_name = "acc-test-listener"
  protocol = "HTTP"
  port = 90
  server_group_id = "${byteplus_server_group.foo.id}"
  health_check {
    enabled = "on"
    interval = 10
    timeout = 3
    healthy_threshold = 5
    un_healthy_threshold = 2
    domain = "byteplus.com"
    http_code = "http_2xx"
    method = "GET"
    uri = "/"
  }
  enabled = "on"
}

`

const testAccByteplusListenerUpdateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "byteplus_listener" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  listener_name = "acc-test-listener1"
  description = "hello demo"
  protocol = "HTTP"
  port = 90
  server_group_id = "${byteplus_server_group.foo.id}"
  health_check {
    enabled = "on"
    interval = 10
    timeout = 3
    healthy_threshold = 5
    un_healthy_threshold = 2
    domain = "byteplus.com"
    http_code = "http_2xx"
    method = "GET"
    uri = "/"
  }
  enabled = "off"
}

`

func TestAccByteplusListenerResource_Basic(t *testing.T) {
	resourceName := "byteplus_listener.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &listener.ByteplusListenerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusListenerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "HTTP"),
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

func TestAccByteplusListenerResource_Update(t *testing.T) {
	resourceName := "byteplus_listener.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &listener.ByteplusListenerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusListenerCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "HTTP"),
				),
			},
			{
				Config: testAccByteplusListenerUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "hello demo"),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "off"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "HTTP"),
				),
			},
			{
				Config:             testAccByteplusListenerUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusListenerTCPCreateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "byteplus_listener" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  listener_name = "acc-test-listener"
  protocol = "TCP"
  port = 90
  server_group_id = "${byteplus_server_group.foo.id}"
  enabled = "on"
  bandwidth = 2
  proxy_protocol_type = "standard"
  persistence_type = "source_ip"
  persistence_timeout = 100
  connection_drain_enabled = "on"
  connection_drain_timeout = 100
}

`

func TestAccByteplusListenerResource_BasicTCP(t *testing.T) {
	resourceName := "byteplus_listener.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &listener.ByteplusListenerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusListenerTCPCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "TCP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "proxy_protocol_type", "standard"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_type", "source_ip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_timeout", "100"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_timeout", "100"),
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

const testAccByteplusListenerTCPUpdateConfig = `
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

resource "byteplus_clb" "foo" {
  type = "public"
  subnet_id = "${byteplus_subnet.foo.id}"
  load_balancer_spec = "small_1"
  description = "acc0Demo"
  load_balancer_name = "acc-test-create"
  eip_billing_config {
    isp = "BGP"
    eip_billing_type = "PostPaidByBandwidth"
    bandwidth = 1
  }
}

resource "byteplus_server_group" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  server_group_name = "acc-test-create"
  description = "hello demo11"
}

resource "byteplus_listener" "foo" {
  load_balancer_id = "${byteplus_clb.foo.id}"
  listener_name = "acc-test-listener"
  protocol = "TCP"
  port = 90
  server_group_id = "${byteplus_server_group.foo.id}"
  enabled = "on"
  bandwidth = 2
  proxy_protocol_type = "standard"
  persistence_type = "source_ip"
  persistence_timeout = 1000
  connection_drain_enabled = "on"
  connection_drain_timeout = 200
}

`

func TestAccByteplusListenerResource_UpdateTCP(t *testing.T) {
	resourceName := "byteplus_listener.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &listener.ByteplusListenerService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusListenerTCPCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "TCP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "proxy_protocol_type", "standard"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_type", "source_ip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_timeout", "100"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_timeout", "100"),
				),
			},
			{
				Config: testAccByteplusListenerTCPUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "health_check.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "listener_name", "acc-test-listener"),
					resource.TestCheckResourceAttr(acc.ResourceId, "port", "90"),
					resource.TestCheckResourceAttr(acc.ResourceId, "protocol", "TCP"),
					resource.TestCheckResourceAttr(acc.ResourceId, "bandwidth", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "proxy_protocol_type", "standard"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_type", "source_ip"),
					resource.TestCheckResourceAttr(acc.ResourceId, "persistence_timeout", "1000"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_enabled", "on"),
					resource.TestCheckResourceAttr(acc.ResourceId, "connection_drain_timeout", "200"),
				),
			},
			{
				Config:             testAccByteplusListenerTCPUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
