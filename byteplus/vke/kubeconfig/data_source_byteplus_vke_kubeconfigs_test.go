package kubeconfig_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/kubeconfig"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusVkeKubeconfigsDatasourceConfig = `
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

resource "byteplus_vke_kubeconfig" "foo1" {
    cluster_id = "${byteplus_vke_cluster.foo.id}"
    type = "Private"
	valid_duration = 2
}

resource "byteplus_vke_kubeconfig" "foo2" {
    cluster_id = "${byteplus_vke_cluster.foo.id}"
    type = "Public"
	valid_duration = 2
}

data "byteplus_vke_kubeconfigs" "foo"{
    ids = ["${byteplus_vke_kubeconfig.foo1.id}", "${byteplus_vke_kubeconfig.foo2.id}"]
}
`

func TestAccByteplusVkeKubeconfigsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_vke_kubeconfigs.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &kubeconfig.ByteplusVkeKubeconfigService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusVkeKubeconfigsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "kubeconfigs.#", "2"),
				),
			},
		},
	})
}
