---
subcategory: "EBS"
layout: "byteplus"
page_title: "Byteplus: byteplus_volume_attach"
sidebar_current: "docs-byteplus-resource-volume_attach"
description: |-
  Provides a resource to manage volume attach
---
# byteplus_volume_attach
Provides a resource to manage volume attach
## Example Usage
```hcl
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

resource "byteplus_security_group" "foo" {
  security_group_name = "acc-test-security-group"
  vpc_id              = byteplus_vpc.foo.id
}

data "byteplus_images" "foo" {
  os_type          = "Linux"
  visibility       = "public"
  instance_type_id = "ecs.g1.large"
}

resource "byteplus_ecs_instance" "foo" {
  instance_name        = "acc-test-ecs"
  description          = "acc-test"
  host_name            = "tf-acc-test"
  image_id             = data.byteplus_images.foo.images[0].image_id
  instance_type        = "ecs.g1.large"
  password             = "93f0cb0614Aab12"
  instance_charge_type = "PostPaid"
  system_volume_type   = "ESSD_PL0"
  system_volume_size   = 40
  subnet_id            = byteplus_subnet.foo.id
  security_group_ids   = [byteplus_security_group.foo.id]
  project_name         = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_volume" "foo" {
  volume_name        = "acc-test-volume"
  volume_type        = "ESSD_PL0"
  description        = "acc-test"
  kind               = "data"
  size               = 40
  zone_id            = data.byteplus_zones.foo.zones[0].id
  volume_charge_type = "PostPaid"
  project_name       = "default"
}

resource "byteplus_volume_attach" "foo" {
  instance_id = byteplus_ecs_instance.foo.id
  volume_id   = byteplus_volume.foo.id
}
```
## Argument Reference
The following arguments are supported:
* `instance_id` - (Required, ForceNew) The Id of Instance.
* `volume_id` - (Required, ForceNew) The Id of Volume.
* `delete_with_instance` - (Optional) Delete Volume with Attached Instance.It is not recommended to use this field. If used, please ensure that the value of this field is consistent with the value of `delete_with_instance` in byteplus_volume.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `created_at` - Creation time of Volume.
* `status` - Status of Volume.
* `updated_at` - Update time of Volume.


## Import
VolumeAttach can be imported using the id, e.g.
```
$ terraform import byteplus_volume_attach.default vol-abc12345:i-abc12345
```

