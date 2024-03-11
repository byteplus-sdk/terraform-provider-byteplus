package default_node_pool_test

import (
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/default_node_pool"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

const testAccByteplusVkeDefaultNodePoolCreateConfig = `
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
        resource_public_access_default_enabled = false
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

resource "byteplus_vke_default_node_pool" "foo" {
    cluster_id = byteplus_vke_cluster.foo.id
    node_config {
        security {
            login {
                password = "amw4WTdVcTRJVVFsUXpVTw=="
            }
            security_group_ids = [byteplus_security_group.foo.id]
            security_strategies = ["Hids"]
        }
        initialize_script = "ISMvYmluL2Jhc2gKZWNobyAx"

    }
    kubernetes_config {
        labels {
            key   = "tf-key1"
            value = "tf-value1"
        }
        labels {
            key   = "tf-key2"
            value = "tf-value2"
        }
        taints {
            key = "tf-key3"
            value = "tf-value3"
            effect = "NoSchedule"
        }
        taints {
            key = "tf-key4"
            value = "tf-value4"
            effect = "NoSchedule"
        }
        cordon = true
    }
    tags {
        key = "tf-k1"
        value = "tf-v1"
    }
}

`

const testAccByteplusVkeDefaultNodePoolUpdateConfig = `
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
        resource_public_access_default_enabled = false
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

resource "byteplus_vke_default_node_pool" "foo" {
    cluster_id = byteplus_vke_cluster.foo.id
    node_config {
        security {
            login {
                password = "UHdkMTIzNDU2"
            }
            security_group_ids = [byteplus_security_group.foo.id]
            security_strategies = ["Hids"]
        }
        initialize_script = "ISMvYmluL2Jhc2gKZWNobyAx"

    }
    kubernetes_config {
        labels {
            key   = "tf-key1"
            value = "tf-value1"
        }
        labels {
            key   = "tf-key2"
            value = "tf-value2"
        }
		labels {
            key   = "tf-key3"
            value = "tf-value3"
        }
        taints {
            key = "tf-key3"
            value = "tf-value3"
            effect = "NoSchedule"
        }
        taints {
            key = "tf-key4"
            value = "tf-value4"
            effect = "NoSchedule"
        }
		taints {
            key = "tf-key5"
            value = "tf-value5"
            effect = "NoSchedule"
        }
        cordon = true
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

func TestAccByteplusVkeDefaultNodePoolResource_Basic(t *testing.T) {
	resourceName := "byteplus_vke_default_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool.NewDefaultNodePoolService(client)
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
				Config: testAccByteplusVkeDefaultNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "amw4WTdVcTRJVVFsUXpVTw=="),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.security.0.login.0.password", "is_import"},
			},
		},
	})
}

func TestAccByteplusVkeDefaultNodePoolResource_Update(t *testing.T) {
	resourceName := "byteplus_vke_default_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return default_node_pool.NewDefaultNodePoolService(client)
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
				Config: testAccByteplusVkeDefaultNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "amw4WTdVcTRJVVFsUXpVTw=="),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
				),
			},
			{
				Config: testAccByteplusVkeDefaultNodePoolUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "3"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key1",
						"value": "tf-value1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key2",
						"value": "tf-value2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "tf-key3",
						"value": "tf-value3",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "tf-key3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "tf-value3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "tf-key4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "tf-value4"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.key", "tf-key5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.value", "tf-value5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.2.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ISMvYmluL2Jhc2gKZWNobyAx"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k1",
						"value": "tf-v1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tf-k2",
						"value": "tf-v2",
					}),
				),
			},
			{
				Config:             testAccByteplusVkeDefaultNodePoolUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
