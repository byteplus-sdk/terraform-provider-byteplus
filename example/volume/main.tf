data "byteplus_zones" "foo" {
}

resource "byteplus_volume" "PostVolume" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 40
  zone_id            = data.byteplus_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
}
