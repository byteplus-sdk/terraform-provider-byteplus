package ssl_state_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/ssl_state"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusMongodbSslStateCreateConfig = `
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
    instance_type="ReplicaSet"
    super_account_password="@acc-test-123"
    node_spec="mongo.2c4g"
    mongos_node_spec="mongo.mongos.2c4g"
    instance_name="acc-test-mongo-replica"
    charge_type="PostPaid"
    project_name = "default"
    mongos_node_number = 2
    shard_number=3
    storage_space_gb=20
    subnet_id=byteplus_subnet.foo.id
    zone_id= data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_mongodb_ssl_state" "foo" {
    instance_id = byteplus_mongodb_instance.foo.id
}
`

func TestAccByteplusMongodbSslStateResource_Basic(t *testing.T) {
	resourceName := "byteplus_mongodb_ssl_state.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_state.NewMongoDBSSLStateService(client)
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
				Config: testAccByteplusMongodbSslStateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "is_valid", "true"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_expired_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ssl_action"),
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

const testAccByteplusMongodbSslStateUpdateConfig = `
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
    instance_type="ReplicaSet"
    super_account_password="@acc-test-123"
    node_spec="mongo.2c4g"
    mongos_node_spec="mongo.mongos.2c4g"
    instance_name="acc-test-mongo-replica"
    charge_type="PostPaid"
    project_name = "default"
    mongos_node_number = 2
    shard_number=3
    storage_space_gb=20
    subnet_id=byteplus_subnet.foo.id
    zone_id= data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_mongodb_ssl_state" "foo" {
    instance_id = byteplus_mongodb_instance.foo.id
    ssl_action = "Update"
}
`

func TestAccByteplusMongodbSslStateResource_Update(t *testing.T) {
	resourceName := "byteplus_mongodb_ssl_state.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return ssl_state.NewMongoDBSSLStateService(client)
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
				Config: testAccByteplusMongodbSslStateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "is_valid", "true"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_expired_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
					resource.TestCheckNoResourceAttr(acc.ResourceId, "ssl_action"),
				),
			},
			{
				Config: testAccByteplusMongodbSslStateUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_action", "Update"),
					resource.TestCheckResourceAttr(acc.ResourceId, "ssl_enable", "true"),
					resource.TestCheckResourceAttr(acc.ResourceId, "is_valid", "true"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "ssl_expired_time"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
				),
			},
			{
				Config:             testAccByteplusMongodbSslStateUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
