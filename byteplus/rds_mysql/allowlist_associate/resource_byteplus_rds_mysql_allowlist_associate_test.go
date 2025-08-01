package allowlist_associate_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/allowlist_associate"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusRdsMysqlAllowlistAssociateCreateConfig = `
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

resource "byteplus_rds_mysql_allowlist" "foo" {
    allow_list_name = "acc-test-allowlist"
	allow_list_desc = "acc-test"
	allow_list_type = "IPv4"
	allow_list = ["192.168.0.0/24", "192.168.1.0/24"]
}

resource "byteplus_rds_mysql_allowlist_associate" "foo" {
    allow_list_id = "${byteplus_rds_mysql_allowlist.foo.id}"
    instance_id = "${byteplus_rds_mysql_instance.foo.id}"
}
`

func TestAccByteplusRdsMysqlAllowlistAssociateResource_Basic(t *testing.T) {
	resourceName := "byteplus_rds_mysql_allowlist_associate.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		SvcInitFunc: func(client *bp.SdkClient) bp.ResourceService {
			return allowlist_associate.NewRdsMysqlAllowListAssociateService(client)
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
				Config: testAccByteplusRdsMysqlAllowlistAssociateCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "allow_list_id"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "instance_id"),
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
