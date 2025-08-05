package endpoint_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/endpoint"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusMongodbEndpointsDatasourceConfig1 = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
    vpc_name     = "acc-test-vpc"
    cidr_block   = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-test-subnet"
    cidr_block  = "172.16.0.0/24"
    zone_id     = data.byteplus_zones.foo.zones[0].id
    vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_mongodb_instance" "foo"{
    db_engine_version = "MongoDB_4_2"
    instance_type = "ShardedCluster"
    super_account_password = "@acc-test-123"
    node_spec = "mongo.shard.1c2g"
    mongos_node_spec = "mongo.mongos.1c2g"
    instance_name = "acc-test-mongo-shard"
    charge_type = "PostPaid"
    project_name = "default"
    mongos_node_number = 2
    shard_number = 2
    storage_space_gb = 20
    subnet_id = byteplus_subnet.foo.id
    zone_id = data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_eip_address" "foo" {
    billing_type = "PostPaidByBandwidth"
    bandwidth = 1
    isp = "ChinaUnicom"
    name = "acc-eip-${count.index}"
    description = "acc-test"
    project_name = "default"
    count = 2
}

resource "byteplus_mongodb_endpoint" "foo_public"{
    instance_id = byteplus_mongodb_instance.foo.id
    network_type = "Public"
	object_id = byteplus_mongodb_instance.foo.mongos_id
	mongos_node_ids = [byteplus_mongodb_instance.foo.mongos[0].mongos_node_id, byteplus_mongodb_instance.foo.mongos[1].mongos_node_id]
	eip_ids = byteplus_eip_address.foo[*].id
}

resource "byteplus_mongodb_endpoint" "foo_private"{
    instance_id = byteplus_mongodb_instance.foo.id
    network_type = "Private"
	object_id = byteplus_mongodb_instance.foo.config_servers_id
}
`

const testAccByteplusMongodbEndpointsDatasourceConfig2 = `
data "byteplus_zones" "foo"{
}

resource "byteplus_vpc" "foo" {
    vpc_name     = "acc-test-vpc"
    cidr_block   = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-test-subnet"
    cidr_block  = "172.16.0.0/24"
    zone_id     = data.byteplus_zones.foo.zones[0].id
    vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_mongodb_instance" "foo"{
    db_engine_version = "MongoDB_4_2"
    instance_type = "ShardedCluster"
    super_account_password = "@acc-test-123"
    node_spec = "mongo.shard.1c2g"
    mongos_node_spec = "mongo.mongos.1c2g"
    instance_name = "acc-test-mongo-shard"
    charge_type = "PostPaid"
    project_name = "default"
    mongos_node_number = 2
    shard_number = 2
    storage_space_gb = 20
    subnet_id = byteplus_subnet.foo.id
    zone_id = data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_eip_address" "foo" {
    billing_type = "PostPaidByBandwidth"
    bandwidth = 1
    isp = "ChinaUnicom"
    name = "acc-eip-${count.index}"
    description = "acc-test"
    project_name = "default"
    count = 2
}

resource "byteplus_mongodb_endpoint" "foo_public"{
    instance_id = byteplus_mongodb_instance.foo.id
    network_type = "Public"
	object_id = byteplus_mongodb_instance.foo.mongos_id
	mongos_node_ids = [byteplus_mongodb_instance.foo.mongos[0].mongos_node_id, byteplus_mongodb_instance.foo.mongos[1].mongos_node_id]
	eip_ids = byteplus_eip_address.foo[*].id
}

resource "byteplus_mongodb_endpoint" "foo_private"{
    instance_id = byteplus_mongodb_instance.foo.id
    network_type = "Private"
	object_id = byteplus_mongodb_instance.foo.config_servers_id
}

data "byteplus_mongodb_endpoints" "foo"{
    instance_id = byteplus_mongodb_instance.foo.id
}
`

func TestAccByteplusMongodbEndpointsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_mongodb_endpoints.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return endpoint.NewMongoDBEndpointService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusMongodbEndpointsDatasourceConfig1,
			},
			{
				Config: testAccByteplusMongodbEndpointsDatasourceConfig2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "endpoints.#", "3"),
				),
			},
		},
	})
}
