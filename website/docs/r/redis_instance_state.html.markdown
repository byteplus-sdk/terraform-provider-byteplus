---
subcategory: "REDIS"
layout: "byteplus"
page_title: "Byteplus: byteplus_redis_instance_state"
sidebar_current: "docs-byteplus-resource-redis_instance_state"
description: |-
  Provides a resource to manage redis instance state
---
# byteplus_redis_instance_state
Provides a resource to manage redis instance state
## Example Usage
```hcl
resource "byteplus_redis_instance_state" "foo" {
  action      = "Restart"
  instance_id = "redis-cnlficlt4974swtbz"
}
```
## Argument Reference
The following arguments are supported:
* `action` - (Required, ForceNew) Instance Action, the value can be `Restart`.
* `instance_id` - (Required, ForceNew) Id of Instance.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Redis State Instance can be imported using the id, e.g.
```
$ terraform import byteplus_redis_instance_state.default state:redis-mizl7m1kqccg5smt1bdpijuj
```

