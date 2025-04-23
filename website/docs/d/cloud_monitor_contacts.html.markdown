---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_contacts"
sidebar_current: "docs-byteplus-datasource-cloud_monitor_contacts"
description: |-
  Use this data source to query detailed information of cloud monitor contacts
---
# byteplus_cloud_monitor_contacts
Use this data source to query detailed information of cloud monitor contacts
## Example Usage
```hcl
data "byteplus_cloud_monitor_contacts" "foo" {
  ids = ["17******516", "1712**********0"]
}
```
## Argument Reference
The following arguments are supported:
* `email` - (Optional) The email of the cloud monitor contact. This field support fuzzy query.
* `ids` - (Optional) A list of Contact IDs.
* `name` - (Optional) The keyword of contact names. Fuzzy match is supported.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `contacts` - The collection of query.
    * `email` - The email of contact.
    * `id` - The ID of contact.
    * `name` - The name of contact.
* `total_count` - The total count of query.


