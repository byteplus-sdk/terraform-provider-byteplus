resource "byteplus_eip_address" "foo" {
  billing_type = "PostPaidByBandwidth"
  bandwidth    = 1
  isp          = "BGP"
  name         = "acc-eip"
  description  = "acc-test"
  project_name = "default"
}

resource "byteplus_volume" "foo" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 20
  zone_id            = "ap-southeast-1a"
  volume_charge_type = "PostPaid"
}

resource "byteplus_cloud_monitor_object_group" "foo" {
  name = "acc_test_object_group"
  objects {
    namespace = "VCM_EIP"
    region    = ["ap-southeast-1"]
    dimensions {
      key   = "ResourceID"
      value = [byteplus_eip_address.foo.id]
    }
  }
  objects {
    namespace = "VCM_EBS"
    region    = ["ap-southeast-1"]
    dimensions {
      key   = "ResourceID"
      value = [byteplus_volume.foo.id]
    }
  }
}
