package rds_mysql_instance_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsMysqlInstanceCreateConfig = `
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

resource "byteplus_rds_mysql_instance" "foo" {
  db_engine_version = "MySQL_5_7"
  node_spec = "rds.mysql.1c2g"
  primary_zone_id = data.byteplus_zones.foo.zones[0].id
  secondary_zone_id = data.byteplus_zones.foo.zones[0].id
  storage_space = 80
  subnet_id = byteplus_subnet.foo.id
  instance_name = "acc-test"
  lower_case_table_names = "1"

  charge_info {
    charge_type = "PostPaid"
  }

  parameters {
    parameter_name = "auto_increment_increment"
    parameter_value = "2"
  }
  parameters {
    parameter_name = "auto_increment_offset"
    parameter_value = "4"
  }

  project_name = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}
`

const testAccByteplusRdsMysqlInstanceUpdateConfig = `
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

resource "byteplus_rds_mysql_allowlist" "foo" {
    allow_list_name = "acc-test-allowlist"
	allow_list_desc = "acc-test"
	allow_list_type = "IPv4"
	allow_list = ["192.168.0.0/24", "192.168.1.0/24"]
}

resource "byteplus_rds_mysql_allowlist" "foo1" {
    allow_list_name = "acc-test-allowlist1"
	allow_list_desc = "acc-test1"
	allow_list_type = "IPv4"
	allow_list = ["192.168.0.0/24", "192.168.1.0/24"]
}

resource "byteplus_rds_mysql_instance" "foo" {
  db_engine_version = "MySQL_5_7"
  node_spec = "rds.mysql.2c4g"
  primary_zone_id = data.byteplus_zones.foo.zones[0].id
  secondary_zone_id = data.byteplus_zones.foo.zones[0].id
  storage_space = 100
  subnet_id = byteplus_subnet.foo.id
  instance_name = "acc-test1"
  lower_case_table_names = "1"

  charge_info {
    charge_type = "PostPaid"
  }

  allow_list_ids = [byteplus_rds_mysql_allowlist.foo.id, byteplus_rds_mysql_allowlist.foo1.id]

  parameters {
    parameter_name = "auto_increment_increment"
    parameter_value = "4"
  }
  parameters {
    parameter_name = "auto_increment_offset"
    parameter_value = "8"
  }
  parameters {
    parameter_name = "innodb_thread_concurrency"
    parameter_value = "0"
  }

  project_name = "default"
  tags {
    key   = "k2"
    value = "v2"
  }
  tags {
    key   = "k3"
    value = "v3"
  }
}
`

func TestAccByteplusRdsMysqlInstanceResource_Basic(t *testing.T) {
	resourceName := "byteplus_rds_mysql_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_mysql_instance.NewRdsMysqlInstanceService(client)
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
				Config: testAccByteplusRdsMysqlInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_ids.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "MySQL_5_7"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lower_case_table_names", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.mysql.1c2g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_increment",
						"parameter_value": "2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_offset",
						"parameter_value": "4",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "80"),
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

func TestAccByteplusRdsMysqlInstanceResource_Update(t *testing.T) {
	resourceName := "byteplus_rds_mysql_instance.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_mysql_instance.NewRdsMysqlInstanceService(client)
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
				Config: testAccByteplusRdsMysqlInstanceCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_ids.#", "0"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "MySQL_5_7"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lower_case_table_names", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.mysql.1c2g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "1"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k1",
						"value": "v1",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_increment",
						"parameter_value": "2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_offset",
						"parameter_value": "4",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "80"),
				),
			},
			{
				Config: testAccByteplusRdsMysqlInstanceUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "allow_list_ids.#", "2"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.#", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "charge_info.0.charge_type", "PostPaid"),
					resource.TestCheckResourceAttr(acc.ResourceId, "db_engine_version", "MySQL_5_7"),
					resource.TestCheckResourceAttr(acc.ResourceId, "instance_name", "acc-test1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "lower_case_table_names", "1"),
					resource.TestCheckResourceAttr(acc.ResourceId, "node_spec", "rds.mysql.2c4g"),
					resource.TestCheckResourceAttr(acc.ResourceId, "project_name", "default"),
					resource.TestCheckResourceAttr(acc.ResourceId, "tags.#", "2"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k2",
						"value": "v2",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "tags.*", map[string]string{
						"key":   "k3",
						"value": "v3",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "parameters.#", "3"),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_increment",
						"parameter_value": "4",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "auto_increment_offset",
						"parameter_value": "8",
					}),
					byteplus.TestCheckTypeSetElemNestedAttrs(acc.ResourceId, "parameters.*", map[string]string{
						"parameter_name":  "innodb_thread_concurrency",
						"parameter_value": "0",
					}),
					resource.TestCheckResourceAttr(acc.ResourceId, "storage_space", "100"),
				),
			},
			{
				Config:             testAccByteplusRdsMysqlInstanceUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
