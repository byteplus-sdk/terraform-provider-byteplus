package rds_postgresql_instance_readonly_node_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_instance_readonly_node"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsPostgresqlInstanceReadonlyNodeCreateConfig = `
data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
    vpc_name   = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block  = "172.16.0.0/24"
    zone_id     = data.byteplus_zones.foo.zones[0].id
    vpc_id      = byteplus_vpc.foo.id
}


resource "byteplus_rds_postgresql_instance" "foo" {
    db_engine_version = "PostgreSQL_12"
    node_spec = "rds.postgres.1c2g"
    primary_zone_id        = data.byteplus_zones.foo.zones[0].id
    secondary_zone_id      = data.byteplus_zones.foo.zones[0].id
    storage_space          = 40
    subnet_id = byteplus_subnet.foo.id
    instance_name          = "acc-test-1"
    charge_info {
        charge_type = "PostPaid"
    }
    project_name = "default"
    tags {
        key   = "tfk1"
        value = "tfv1"
    }
    parameters {
        name  = "auto_explain.log_analyze"
        value = "off"
    }
    parameters {
        name  = "auto_explain.log_format"
        value = "text"
    }
}

resource "byteplus_rds_postgresql_instance_readonly_node" "foo" {
    instance_id = byteplus_rds_postgresql_instance.foo.id
    node_spec = "rds.postgres.1c2g"
    zone_id = data.byteplus_zones.foo.zones[0].id
}
`

const testAccByteplusRdsPostgresqlInstanceReadonlyNodeUpdateConfig = `
data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
    vpc_name   = "acc-test-project1"
    cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
    subnet_name = "acc-subnet-test-2"
    cidr_block  = "172.16.0.0/24"
    zone_id     = data.byteplus_zones.foo.zones[0].id
    vpc_id      = byteplus_vpc.foo.id
}


resource "byteplus_rds_postgresql_instance" "foo" {
    db_engine_version = "PostgreSQL_12"
    node_spec = "rds.postgres.1c2g"
    primary_zone_id        = data.byteplus_zones.foo.zones[0].id
    secondary_zone_id      = data.byteplus_zones.foo.zones[0].id
    storage_space          = 40
    subnet_id = byteplus_subnet.foo.id
    instance_name          = "acc-test-1"
    charge_info {
        charge_type = "PostPaid"
    }
    project_name = "default"
    tags {
        key   = "tfk1"
        value = "tfv1"
    }
    parameters {
        name  = "auto_explain.log_analyze"
        value = "off"
    }
    parameters {
        name  = "auto_explain.log_format"
        value = "text"
    }
}

resource "byteplus_rds_postgresql_instance_readonly_node" "foo" {
    instance_id = byteplus_rds_postgresql_instance.foo.id
    node_spec = "rds.postgres.2c4g"
    zone_id = data.byteplus_zones.foo.zones[0].id
}
`

func TestAccByteplusRdsPostgresqlInstanceReadonlyNodeResource_Basic(t *testing.T) {
	resourceName := "byteplus_rds_postgresql_instance_readonly_node.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_postgresql_instance_readonly_node.NewRdsPostgresqlInstanceReadonlyNodeService(client)
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
				Config: testAccByteplusRdsPostgresqlInstanceReadonlyNodeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.1c2g"),
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

func TestAccByteplusRdsPostgresqlInstanceReadonlyNodeResource_Update(t *testing.T) {
	resourceName := "byteplus_rds_postgresql_instance_readonly_node.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_postgresql_instance_readonly_node.NewRdsPostgresqlInstanceReadonlyNodeService(client)
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
				Config: testAccByteplusRdsPostgresqlInstanceReadonlyNodeCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.1c2g"),
				),
			},
			{
				Config: testAccByteplusRdsPostgresqlInstanceReadonlyNodeUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.2c4g"),
				),
			},
			{
				Config:             testAccByteplusRdsPostgresqlInstanceReadonlyNodeUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
