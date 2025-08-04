package instance_parameter_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/instance_parameter"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusMongodbInstanceParameterCreateConfig = `
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
    mongos_node_number = 32
    shard_number=3
    storage_space_gb=20
    subnet_id=byteplus_subnet.foo.id
    zone_id= data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_mongodb_instance_parameter" "foo" {
    instance_id = byteplus_mongodb_instance.foo.id
    parameter_name = "cursorTimeoutMillis"
    parameter_role = "Node"
    parameter_value = "600001"
}
`

const testAccByteplusMongodbInstanceParameterUpdateConfig = `
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
    mongos_node_number = 32
    shard_number=3
    storage_space_gb=20
    subnet_id=byteplus_subnet.foo.id
    zone_id= data.byteplus_zones.foo.zones[0].id
    tags {
        key = "k1"
        value = "v1"
    }
}

resource "byteplus_mongodb_instance_parameter" "foo" {
    instance_id = byteplus_mongodb_instance.foo.id
    parameter_name = "cursorTimeoutMillis"
    parameter_role = "Node"
    parameter_value = "600111"
}
`

func TestAccByteplusMongodbInstanceParameterResource_Basic(t *testing.T) {
	resourceName := "byteplus_mongodb_instance_parameter.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return instance_parameter.NewMongoDBInstanceParameterService(client)
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
				Config: testAccByteplusMongodbInstanceParameterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_name", "cursorTimeoutMillis"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_role", "Node"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_value", "600001"),
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

func TestAccByteplusMongodbInstanceParameterResource_Update(t *testing.T) {
	resourceName := "byteplus_mongodb_instance_parameter.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return instance_parameter.NewMongoDBInstanceParameterService(client)
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
				Config: testAccByteplusMongodbInstanceParameterCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_name", "cursorTimeoutMillis"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_role", "Node"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_value", "600001"),
				),
			},
			{
				Config: testAccByteplusMongodbInstanceParameterUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_name", "cursorTimeoutMillis"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_role", "Node"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameter_value", "600111"),
				),
			},
			{
				Config:             testAccByteplusMongodbInstanceParameterUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
