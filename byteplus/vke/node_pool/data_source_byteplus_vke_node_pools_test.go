package node_pool_test

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/node_pool"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusVkeNodePoolsDatasourceConfig = `
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
        resource_public_access_default_enabled = true
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
	count = 3
	cluster_id = "${byteplus_vke_cluster.foo.id}"
	name = "acc-test-node-pool-${count.index}"
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

data "byteplus_vke_node_pools" "foo"{
    ids = byteplus_vke_node_pool.foo[*].id
}
`

func TestAccByteplusVkeNodePoolsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vke_node_pools.foo"

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
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVkeNodePoolsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "node_pools.#", "3"),
				),
			},
		},
	})
}
