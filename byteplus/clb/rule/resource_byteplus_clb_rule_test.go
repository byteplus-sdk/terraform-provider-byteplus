package rule_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/rule"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusClbRuleCreateConfig = `
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
resource "byteplus_clb_rule" "foo" {
  listener_id = "${byteplus_listener.foo.id}"
  server_group_id = "${byteplus_server_group.foo.id}"
  domain = "test-byte123.com"
  url = "/yyyy"
}
`

const testAccByteplusClbRuleUpdateConfig = `
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
resource "byteplus_clb_rule" "foo" {
  listener_id = "${byteplus_listener.foo.id}"
  server_group_id = "${byteplus_server_group.foo.id}"
  domain = "acc-test-byte123.com"
  url = "/accyyyy"
}
`

func TestAccByteplusClbRuleResource_Basic(t *testing.T) {
	resourceName := "byteplus_clb_rule.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &rule.ByteplusRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbRuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "test-byte123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/yyyy"),
				),
			},
		},
	})
}

func TestAccByteplusClbRuleResource_Update(t *testing.T) {
	resourceName := "byteplus_clb_rule.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &rule.ByteplusRuleService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusClbRuleCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "test-byte123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/yyyy"),
				),
			},
			{
				Config: testAccByteplusClbRuleUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", ""),
					resource.TestCheckResourceAttr(acc.ResourceId, "domain", "acc-test-byte123.com"),
					resource.TestCheckResourceAttr(acc.ResourceId, "url", "/accyyyy"),
				),
			},
			{
				Config:             testAccByteplusClbRuleUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
