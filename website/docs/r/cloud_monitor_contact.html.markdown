---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_contact"
sidebar_current: "docs-byteplus-resource-cloud_monitor_contact"
description: |-
  Provides a resource to manage cloud monitor contact
---
# byteplus_cloud_monitor_contact
Provides a resource to manage cloud monitor contact
## Example Usage
```hcl
resource "byteplus_cloud_monitor_contact" "default" {
  name  = "acc-test-contact"
  email = "test.com"
}
```
## Argument Reference
The following arguments are supported:
* `email` - (Required) The email of contact.
* `name` - (Required) The name of contact.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CloudMonitor Contact can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_contact.default 145258255725730****
```

