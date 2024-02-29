data "byteplus_zones" "foo" {
}

resource "byteplus_volume" "foo" {
  volume_name        = "acc-test-volume-${count.index}"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 60
  zone_id            = data.byteplus_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
  count              = 3
}

data "byteplus_volumes" "foo" {
  ids = byteplus_volume.foo[*].id
}
