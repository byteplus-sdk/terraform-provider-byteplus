package rds_postgresql_schema_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_schema"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsPostgresqlSchemasDatasourceConfig = `
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

resource "byteplus_rds_postgresql_database" "foo" {
    db_name     = "acc-test"
    instance_id = byteplus_rds_postgresql_instance.foo.id
    c_type      = "C"
    collate     = "zh_CN.utf8"
}

resource "byteplus_rds_postgresql_account" "foo" {
    account_name       = "acc-test-account"
    account_password   = "9wc@********12"
    account_type       = "Normal"
    instance_id        = byteplus_rds_postgresql_instance.foo.id
    account_privileges = "Inherit,Login,CreateRole,CreateDB"
}

resource "byteplus_rds_postgresql_account" "foo1" {
    account_name       = "acc-test-account1"
    account_password   = "9wc@*******12"
    account_type       = "Normal"
    instance_id        = byteplus_rds_postgresql_instance.foo.id
    account_privileges = "Inherit,Login,CreateRole,CreateDB"
}

resource "byteplus_rds_postgresql_schema" "foo" {
    db_name = byteplus_rds_postgresql_database.foo.db_name
    instance_id = byteplus_rds_postgresql_instance.foo.id
    owner = byteplus_rds_postgresql_account.foo.account_name
    schema_name = "acc-test-schema"
}

data "byteplus_rds_postgresql_schemas" "foo"{
    db_name = byteplus_rds_postgresql_schema.foo.db_name
    instance_id = byteplus_rds_postgresql_instance.foo.id
}
`

func TestAccByteplusRdsPostgresqlSchemasDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_rds_postgresql_schemas.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_postgresql_schema.NewRdsPostgresqlSchemaService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusRdsPostgresqlSchemasDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "schemas.#", "1"),
				),
			},
		},
	})
}
