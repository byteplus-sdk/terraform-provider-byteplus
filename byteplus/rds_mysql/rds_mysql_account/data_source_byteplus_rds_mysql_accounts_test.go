package rds_mysql_account_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_account"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsMysqlAccountsDatasourceConfig = `
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

resource "byteplus_rds_mysql_instance" "foo" {
	instance_name = "acc-test-rds-mysql"
  	db_engine_version = "MySQL_5_7"
  	node_spec = "rds.mysql.1c2g"
  	primary_zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  	secondary_zone_id = "${data.byteplus_zones.foo.zones[0].id}"
  	storage_space = 80
  	subnet_id = "${byteplus_subnet.foo.id}"
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
}

resource "byteplus_rds_mysql_database" "foo" {
    db_name = "acc-test-db"
    instance_id = "${byteplus_rds_mysql_instance.foo.id}"
}

resource "byteplus_rds_mysql_account" "foo" {
    account_name = "acc-test-account"
    account_password = "93f0cb0614Aab12"
    account_type = "Normal"
    instance_id = "${byteplus_rds_mysql_instance.foo.id}"
	account_privileges {
		db_name = "${byteplus_rds_mysql_database.foo.db_name}"
		account_privilege = "Custom"
		account_privilege_detail = "SELECT,INSERT"
	}
}

data "byteplus_rds_mysql_accounts" "foo"{
    instance_id = "${byteplus_rds_mysql_instance.foo.id}"
	account_name = "${byteplus_rds_mysql_account.foo.account_name}"
}
`

func TestAccByteplusRdsMysqlAccountsDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_rds_mysql_accounts.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return rds_mysql_account.NewRdsMysqlAccountService(client)
		},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusRdsMysqlAccountsDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "accounts.#", "1"),
				),
			},
		},
	})
}
