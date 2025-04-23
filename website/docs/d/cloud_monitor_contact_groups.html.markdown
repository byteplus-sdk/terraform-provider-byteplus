---
subcategory: "CLOUD_MONITOR"
layout: "byteplus"
page_title: "Byteplus: byteplus_cloud_monitor_contact_groups"
sidebar_current: "docs-byteplus-datasource-cloud_monitor_contact_groups"
description: |-
  Use this data source to query detailed information of cloud monitor contact groups
---
# byteplus_cloud_monitor_contact_groups
Use this data source to query detailed information of cloud monitor contact groups
## Example Usage
```hcl
data "byteplus_cloud_monitor_contact_groups" "foo" {
  name = "tftest"
}
```
## Argument Reference
The following arguments are supported:
* `ids` - (Optional) A list of cloud monitor contact group ids.
* `name` - (Optional) The keyword of the contact group names. Fuzzy match is supported.
* `output_file` - (Optional) File name where to save data source results.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `groups` - The collection of query.
    * `account_id` - The id of the account.
    * `contacts` - Contact information in the contact group.
        * `email` - The email of contact.
        * `id` - The id of the contact.
        * `name` - The name of contact.
        * `phone` - The phone of contact.
    * `created_at` - The create time.
    * `description` - The description of the contact group.
    * `id` - The id of the contact group.
    * `name` - The name of the contact group.
    * `updated_at` - The update time.
* `total_count` - The total count of query.


