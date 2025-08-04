data "byteplus_zones" "foo" {
}

resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-vpc"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-test-subnet"
  cidr_block  = "172.16.0.0/24"
  zone_id     = data.byteplus_zones.foo.zones[0].id
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_mongodb_instance" "foo" {
  zone_ids               = [data.byteplus_zones.foo.zones[0].id]
  db_engine_version      = "MongoDB_4_2"
  instance_type          = "ReplicaSet"
  node_spec              = "mongo.2c4g"
  storage_space_gb       = 20
  subnet_id              = byteplus_subnet.foo.id
  instance_name          = "acc-test-mongodb-replica"
  charge_type            = "PostPaid"
  super_account_password = "93f0cb0614Aab12"
  project_name           = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_mongodb_account" "foo" {
  instance_id      = byteplus_mongodb_instance.foo.id
  account_name     = "acc-test-mongodb-account"
  auth_db          = "admin"
  account_password = "93f0cb0614Aab12"
  account_desc     = "acc-test"
  account_privileges {
    db_name    = "admin"
    role_names = ["userAdmin", "clusterMonitor"]
  }
  account_privileges {
    db_name    = "config"
    role_names = ["read"]
  }
  account_privileges {
    db_name    = "local"
    role_names = ["read"]
  }
}
