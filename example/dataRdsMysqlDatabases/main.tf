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

resource "byteplus_rds_mysql_instance" "foo" {
  db_engine_version      = "MySQL_5_7"
  node_spec              = "rds.mysql.1c2g"
  primary_zone_id        = data.byteplus_zones.foo.zones[0].id
  secondary_zone_id      = data.byteplus_zones.foo.zones[0].id
  storage_space          = 80
  subnet_id              = byteplus_subnet.foo.id
  instance_name          = "acc-test"
  lower_case_table_names = "1"

  charge_info {
    charge_type = "PostPaid"
  }

  parameters {
    parameter_name  = "auto_increment_increment"
    parameter_value = "2"
  }
  parameters {
    parameter_name  = "auto_increment_offset"
    parameter_value = "4"
  }
}

resource "byteplus_rds_mysql_database" "foo" {
  db_name     = "acc-test"
  instance_id = byteplus_rds_mysql_instance.foo.id
}
data "byteplus_rds_mysql_databases" "foo" {
  db_name     = "acc-test"
  instance_id = byteplus_rds_mysql_instance.foo.id
}