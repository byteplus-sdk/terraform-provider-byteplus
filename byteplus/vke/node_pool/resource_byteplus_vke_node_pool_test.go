package node_pool_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/node_pool"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusVkeNodePoolCreateConfig = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_Basic(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.security.0.login.0.password"},
			},
		},
	})
}

const testAccByteplusVkeNodePoolUpdateConfig = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool-new"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k2"
        value = "node-pool-v2"
    }
	tags {
        key = "node-pool-k3"
        value = "node-pool-v3"
    }
}
`

func TestAccByteplusVkeNodePoolResource_Update(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k2",
						"value": "node-pool-v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k3",
						"value": "node-pool-v3",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolUpdateAutoScalingConfig = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 1
		max_replicas = 20
		desired_replicas = 2
		priority = 20
        subnet_policy = "Priority"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_UpdateAutoScalingConfig(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdateAutoScalingConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "20"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "20"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "Priority"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolUpdateAutoScalingConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolUpdateNodeConfig = `
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

resource "byteplus_subnet" "foo2" {
  	subnet_name = "acc-test-subnet2"
  	cidr_block = "172.16.2.0/24"
  	zone_id = "${data.byteplus_zones.foo.zones[0].id}"
	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

resource "byteplus_security_group" "foo2" {
  	security_group_name = "acc-test-security-group2"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.large"]
        subnet_ids = ["${byteplus_subnet.foo.id}", "${byteplus_subnet.foo2.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB2a2Uh"
		security {
            login {
                 password = "UHdkMTIzNDU2Nzg="
            }
			security_strategies = []
            security_group_ids = ["${byteplus_security_group.foo.id}", "${byteplus_security_group.foo2.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test-new"
        ecs_tags {
            key = "ecs_k2"
            value = "ecs_v2"
        }
		ecs_tags {
            key = "ecs_k3"
            value = "ecs_v3"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_UpdateNodeConfig(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdateNodeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.large"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB2a2Uh"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2Nzg="),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k2",
						"value": "ecs_v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k3",
						"value": "ecs_v3",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolUpdateNodeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolUpdateKubernetesConfig1 = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label2"
            value = "value2"
        }
		labels {
            key   = "label3"
            value = "value3"
        }
        taints {
            key   = "taint-key/node-type-1"
            value = "taint-value-1"
			effect = "PreferNoSchedule"
        }
		taints {
            key   = "taint-key/node-type-2"
            value = "taint-value-2"
			effect = "PreferNoSchedule"
        }
        cordon = false
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

const testAccByteplusVkeNodePoolUpdateKubernetesConfig2 = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_UpdateKubernetesConfig(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdateKubernetesConfig1,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label2",
						"value": "value2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label3",
						"value": "value3",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.key", "taint-key/node-type-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.value", "taint-value-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.1.effect", "PreferNoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdateKubernetesConfig2,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolUpdateKubernetesConfig2,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolCreatePrePaidConfig = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = false
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PrePaid"
		period = 2
		auto_renew = false
		auto_renew_period = 1
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_CreatePrePaid(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreatePrePaidConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.period", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew_period", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_config.0.security.0.login.0.password"},
			},
		},
	})
}

const testAccByteplusVkeNodePoolUpdatePrePaidConfig = `
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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = false
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "60"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf1"
        }
		data_volumes {
            type = "ESSD_PL0"
            size = "60"
			mount_point = "/tf2"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PrePaid"
		period = 3
		auto_renew = true
		auto_renew_period = 6
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}
`

func TestAccByteplusVkeNodePoolResource_UpdatePrePaidConfig(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreatePrePaidConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.period", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew_period", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolUpdatePrePaidConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PrePaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.period", "3"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.auto_renew_period", "6"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolUpdatePrePaidConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolSystemVolumeAndDataVolumeConfig = `

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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "80"
		}
        data_volumes {
            type = "ESSD_PL0"
            size = "80"
			mount_point = "/tf1"
        }
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = true
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}

`

// 修改系统盘（修改）和数据盘（修改&减少）
func TestAccByteplusVkeNodePoolResource_VolumeUpdate(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolSystemVolumeAndDataVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolSystemVolumeAndDataVolumeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}

const testAccByteplusVkeNodePoolDeleteDataVolumeConfig = `

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

resource "byteplus_security_group" "foo" {
  	security_group_name = "acc-test-security-group"
  	vpc_id = "${byteplus_vpc.foo.id}"
}

data "byteplus_images" "foo" {
  name_regex = "veLinux 1.0 CentOS兼容版 64位"
}

resource "byteplus_vke_cluster" "foo" {
    name = "acc-test-cluster"
    description = "created by terraform"
    delete_protection_enabled = false
    cluster_config {
        subnet_ids = ["${byteplus_subnet.foo.id}"]
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
            subnet_ids = ["${byteplus_subnet.foo.id}"]
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

resource "byteplus_vke_node_pool" "foo" {
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool"
	auto_scaling {
        enabled = true
		min_replicas = 0
		max_replicas = 5
		desired_replicas = 0
		priority = 5
        subnet_policy = "ZoneBalance"
    }
	node_config {
		instance_type_ids = ["ecs.g1ie.xlarge"]
        subnet_ids = ["${byteplus_subnet.foo.id}"]
		image_id = [for image in data.byteplus_images.foo.images : image.image_id if image.image_name == "veLinux 1.0 CentOS兼容版 64位"][0]
		system_volume {
			type = "ESSD_PL0"
            size = "80"
		}
		initialize_script = "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"
		security {
            login {
                 password = "UHdkMTIzNDU2"
            }
			security_strategies = ["Hids"]
            security_group_ids = ["${byteplus_security_group.foo.id}"]
        }
		additional_container_storage_enabled = false
        instance_charge_type = "PostPaid"
        name_prefix = "acc-test"
        ecs_tags {
            key = "ecs_k1"
            value = "ecs_v1"
        }
	}
	kubernetes_config {
        labels {
            key   = "label1"
            value = "value1"
        }
        taints {
            key   = "taint-key/node-type"
            value = "taint-value"
			effect = "NoSchedule"
        }
        cordon = true
    }
	tags {
        key = "node-pool-k1"
        value = "node-pool-v1"
    }
}

`

// 清空数据盘
func TestAccByteplusVkeNodePoolResource_ClearVolumeUpdate(t *testing.T) {
	resourceName := "byteplus_vke_node_pool.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return node_pool.NewNodePoolService(client)
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
				Config: testAccByteplusVkeNodePoolCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.0.mount_point", "/tf1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.size", "60"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.1.mount_point", "/tf2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config: testAccByteplusVkeNodePoolDeleteDataVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "name", "acc-test-node-pool"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.enabled", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.min_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.max_replicas", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.desired_replicas", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.priority", "5"),
					resource.TestCheckResourceAttr(acc.ResourceId, "auto_scaling.0.subnet_policy", "ZoneBalance"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.labels.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "kubernetes_config.0.labels.*", map[string]string{
						"key":   "label1",
						"value": "value1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.key", "taint-key/node-type"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.value", "taint-value"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.taints.0.effect", "NoSchedule"),
					resource.TestCheckResourceAttr(acc.ResourceId, "kubernetes_config.0.cordon", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_type_ids.0", "ecs.g1ie.xlarge"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.subnet_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.initialize_script", "ZWNobyBoZWxsbyB0ZXJyYWZvcm0h"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.additional_container_storage_enabled", "false"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.instance_charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.name_prefix", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.type", "ESSD_PL0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.system_volume.0.size", "80"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.data_volumes.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.security_strategies.0", "Hids"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.security.0.login.0.password", "UHdkMTIzNDU2"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "node_config.0.image_id"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_config.0.ecs_tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "node_config.0.ecs_tags.*", map[string]string{
						"key":   "ecs_k1",
						"value": "ecs_v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "node-pool-k1",
						"value": "node-pool-v1",
					}),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "cluster_id"),
				),
			},
			{
				Config:             testAccByteplusVkeNodePoolDeleteDataVolumeConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
