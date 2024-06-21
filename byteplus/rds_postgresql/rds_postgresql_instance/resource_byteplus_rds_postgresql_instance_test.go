package rds_postgresql_instance_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsPostgresqlInstanceCreateConfig = `
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
`

const testAccByteplusRdsPostgresqlInstanceUpdateConfig = `
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
    storage_space          = 20
    subnet_id = byteplus_subnet.foo.id
    instance_name          = "acc-test-2"
    charge_info {
        charge_type = "PostPaid"
    }
    project_name = "default"
    tags {
        key   = "tfk2"
        value = "tfv2"
    }
    parameters {
        name  = "auto_explain.log_analyze"
        value = "on"
    }
    parameters {
        name  = "auto_explain.log_format"
        value = "xml"
    }
}
`

func TestAccByteplusRdsPostgresqlInstanceResource_Basic(t *testing.T) {
	resourceName := "byteplus_rds_postgresql_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_postgresql_instance.NewRdsPostgresqlInstanceService(client)
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
				Config: testAccByteplusRdsPostgresqlInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "PostgreSQL_12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.1c2g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tfk1",
						"value": "tfv1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_analyze",
						"value": "off",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_format",
						"value": "text",
					}),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameters"},
			},
		},
	})
}

func TestAccByteplusRdsPostgresqlInstanceResource_Update(t *testing.T) {
	resourceName := "byteplus_rds_postgresql_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_postgresql_instance.NewRdsPostgresqlInstanceService(client)
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
				Config: testAccByteplusRdsPostgresqlInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "PostgreSQL_12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.1c2g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "40"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tfk1",
						"value": "tfv1",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_analyze",
						"value": "off",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_format",
						"value": "text",
					}),
				),
			},
			{
				Config: testAccByteplusRdsPostgresqlInstanceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "PostgreSQL_12"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test-2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.postgres.1c2g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "20"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "tfk2",
						"value": "tfv2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_analyze",
						"value": "on",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"name":  "auto_explain.log_format",
						"value": "xml",
					}),
				),
			},
			{
				Config:             testAccByteplusRdsPostgresqlInstanceUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
