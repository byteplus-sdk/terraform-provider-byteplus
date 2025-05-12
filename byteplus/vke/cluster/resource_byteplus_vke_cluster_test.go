package cluster_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/cluster"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVkeClusterCreateConfig = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
    vpc_name = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block = "172.16.0.0/24"
    zone_id = data.byteplus_zones.foo.zones[0].id
    vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_security_group" "foo" {
    vpc_id = byteplus_vpc.foo.id
    security_group_name = "acc-test-security-group2"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-1"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = [byteplus_subnet.foo.id]
        api_server_public_access_enabled = true
        api_server_public_access_config {
            public_access_network_config {
                billing_type = "PostPaidByBandwidth"
                bandwidth = 1
            }
        }
        resource_public_access_default_enabled = true
    }
    pods_config {
        pod_network_mode = "VpcCniShared"
        vpc_cni_config {
            subnet_ids = [byteplus_subnet.foo.id]
        }
    }
    services_config {
        service_cidrsv4 = ["172.30.0.0/18"]
    }
    tags {
        key = "tf-k1"
        value = "tf-v1"
    }
}

`

const testAccByteplusVkeClusterUpdateConfig = `
resource "byteplus_vpc" "foo" {
    vpc_name = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block = "172.16.0.0/24"
    zone_id = "cn-beijing-a"
    vpc_id = byteplus_vpc.foo.id
}

resource "byteplus_security_group" "foo" {
    vpc_id = byteplus_vpc.foo.id
    security_group_name = "acc-test-security-group2"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-2"
    description = "created by terraform update"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = [byteplus_subnet.foo.id]
        api_server_public_access_enabled = false
        api_server_public_access_config {
            public_access_network_config {
                billing_type = "PostPaidByBandwidth"
                bandwidth = 2
            }
        }
        resource_public_access_default_enabled = true
    }
    pods_config {
        pod_network_mode = "VpcCniShared"
        vpc_cni_config {
            subnet_ids = [byteplus_subnet.foo.id]
        }
    }
    services_config {
        service_cidrsv4 = ["172.30.0.0/18"]
    }
    tags {
        key = "tf-k1"
        value = "tf-v1"
    }
    tags {
        key = "tf-k2"
        value = "tf-v2"
    }
}

`

func TestAccByteplusVkeClusterResource_Basic(t *testing.T) {
	resourceName := "byteplus_vke_cluster.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &cluster.ByteplusVkeClusterService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVkeClusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "1"),
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

func TestAccByteplusVkeClusterResource_Update(t *testing.T) {
	resourceName := "byteplus_vke_cluster.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &cluster.ByteplusVkeClusterService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVkeClusterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "1"),
				),
			},
			{
				Config: testAccByteplusVkeClusterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "delete_protection_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "created by terraform update"),
					resource.TestCheckResourceAttr(acc.ResourceId, "logging_config.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "pods_config.0.pod_network_mode", "VpcCniShared"),
					resource.TestCheckResourceAttr(acc.ResourceId, "services_config.#", "1"),
					byteplus.TestCheckTypeSetElemAttr(acc.ResourceId, "services_config.0.service_cidrsv4.*", "172.30.0.0/18"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k2",
						"value": "tf-v2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth", "0"),
				),
			},
			{
				Config:             testAccByteplusVkeClusterUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
