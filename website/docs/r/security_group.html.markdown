---
subcategory: "VPC"
layout: "byteplus"
page_title: "Byteplus: byteplus_security_group"
sidebar_current: "docs-byteplus-resource-security_group"
description: |-
  Provides a resource to manage security group
---
# byteplus_security_group
Provides a resource to manage security group
## Example Usage
```hcl
resource "byteplus_security_group" "g1test1" {
  vpc_id       = "vpc-2feppmy1ugt1c59gp688n1fld"
  project_name = "default"
}
```
## Argument Reference
The following arguments are supported:
* `vpc_id` - (Required, ForceNew) Id of the VPC.
* `description` - (Optional) Description of SecurityGroup.
* `project_name` - (Optional) The ProjectName of SecurityGroup.
* `security_group_name` - (Optional) Name of SecurityGroup.
* `tags` - (Optional) Tags.

The `tags` object supports the following:

* `key` - (Required) The Key of Tags.
* `value` - (Required) The Value of Tags.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `creation_time` - Creation time of SecurityGroup.
* `status` - Status of SecurityGroup.


## Import
SecurityGroup can be imported using the id, e.g.
```
$ terraform import byteplus_security_group.default sg-273ycgql3ig3k7fap8t3dyvqx
```

