---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_contact_group"
sidebar_current: "docs-byteplus-resource-cloud_monitor_contact_group"
description: |-
  Provides a resource to manage cloud monitor contact group
---
# byteplus_cloud_monitor_contact_group
Provides a resource to manage cloud monitor contact group
## Example Usage
```hcl
resource "byteplus_cloud_monitor_contact" "foo1" {
  name  = "acc-test-contact-1"
  email = "test1@163.com"
}

resource "byteplus_cloud_monitor_contact" "foo2" {
  name  = "acc-test-contact-2"
  email = "test2@163.com"
}

resource "byteplus_cloud_monitor_contact_group" "foo" {
  name             = "acc-test-contact-group-new"
  description      = "tf-test-new"
  contacts_id_list = [byteplus_cloud_monitor_contact.foo1.id, byteplus_cloud_monitor_contact.foo2.id]
}
```
## Argument Reference
The following arguments are supported:
* `name` - (Required) The name of the contact group.
* `contacts_id_list` - (Optional) When creating a contact group, contacts should be added with their contact ID. The maximum number of IDs allowed is 100, meaning that the maximum number of members in a single contact group is 100.
* `description` - (Optional) The description of the contact group.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
CloudMonitorContactGroup can be imported using the id, e.g.
```
$ terraform import byteplus_cloud_monitor_contact_group.default resource_id
```

