resource "byteplus_vpc" "foo" {
  vpc_name   = "acc-test-project1"
  cidr_block = "172.16.0.0/16"
}

resource "byteplus_subnet" "foo" {
  subnet_name = "acc-subnet-test-2"
  cidr_block  = "172.16.0.0/24"
  zone_id     = "cn-beijing-a"
  vpc_id      = byteplus_vpc.foo.id
}

resource "byteplus_security_group" "foo" {
  vpc_id              = byteplus_vpc.foo.id
  security_group_name = "acc-test-security-group2"
}

resource "byteplus_ecs_instance" "foo" {
  image_id             = "image-ybqi99s7yq8rx7mnk44b"
  instance_type        = "ecs.g1ie.large"
  instance_name        = "acc-test-ecs-name2"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
  lifecycle {
    ignore_changes = [security_group_ids, instance_name]
  }
}

resource "byteplus_ecs_instance" "foo2" {
  image_id             = "image-ybqi99s7yq8rx7mnk44b"
  instance_type        = "ecs.g1ie.large"
  instance_name        = "acc-test-ecs-name2"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
  lifecycle {
    ignore_changes = [security_group_ids, instance_name]
  }
}

resource "byteplus_vke_cluster" "foo" {
  name                      = "acc-test-1"
  description               = "created by terraform"
  delete_protection_enabled = false
  cluster_config {
    subnet_ids                       = [byteplus_subnet.foo.id]
    api_server_public_access_enabled = true
    api_server_public_access_config {
      public_access_network_config {
        billing_type = "PostPaidByBandwidth"
        bandwidth    = 1
      }
    }
    resource_public_access_default_enabled = true
  }
  pods_config {
    pod_network_mode = "VpcCniShared"
    vpc_cni_config {
      subnet_ids = [byteplus_subnet.foo.id]
    }
  }
  services_config {
    service_cidrsv4 = ["172.30.0.0/18"]
  }
  tags {
    key   = "tf-k1"
    value = "tf-v1"
  }
}

resource "byteplus_vke_default_node_pool" "foo" {
  cluster_id = byteplus_vke_cluster.foo.id
  node_config {
    security {
      login {
        password = "amw4WTdVcTRJVVFsUXpVTw=="
      }
      security_group_ids  = [byteplus_security_group.foo.id]
      security_strategies = ["Hids"]
    }
    initialize_script = "ISMvYmluL2Jhc2gKZWNobyAx"

  }
  kubernetes_config {
    labels {
      key   = "tf-key1"
      value = "tf-value1"
    }
    labels {
      key   = "tf-key2"
      value = "tf-value2"
    }
    taints {
      key    = "tf-key3"
      value  = "tf-value3"
      effect = "NoSchedule"
    }
    taints {
      key    = "tf-key4"
      value  = "tf-value4"
      effect = "NoSchedule"
    }
    cordon = true
  }
  tags {
    key   = "tf-k1"
    value = "tf-v1"
  }
}

resource "byteplus_vke_default_node_pool_batch_attach" "foo" {
  cluster_id           = byteplus_vke_cluster.foo.id
  default_node_pool_id = byteplus_vke_default_node_pool.foo.id
  instances {
    instance_id                          = byteplus_ecs_instance.foo.id
    keep_instance_name                   = true
    additional_container_storage_enabled = false
  }
  instances {
    instance_id                          = byteplus_ecs_instance.foo2.id
    keep_instance_name                   = true
    additional_container_storage_enabled = false
  }
  kubernetes_config {
    labels {
      key   = "tf-key1"
      value = "tf-value1"
    }
    labels {
      key   = "tf-key2"
      value = "tf-value2"
    }
    taints {
      key    = "tf-key3"
      value  = "tf-value3"
      effect = "NoSchedule"
    }
    taints {
      key    = "tf-key4"
      value  = "tf-value4"
      effect = "NoSchedule"
    }
    cordon = true
  }
}