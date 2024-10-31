---
subcategory: "ORGANIZATION"
layout: "byteplus"
page_title: "Byteplus: byteplus_organization"
sidebar_current: "docs-byteplus-resource-organization"
description: |-
  Provides a resource to manage organization
---
# byteplus_organization
Provides a resource to manage organization
## Example Usage
```hcl
resource "byteplus_organization" "foo" {

}
```
## Argument Reference
The following arguments are supported:


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `account_id` - The account id of the organization owner.
* `account_name` - The account name of the organization owner.
* `created_time` - The created time of the organization.
* `delete_uk` - The delete uk of the organization.
* `deleted_time` - The deleted time of the organization.
* `description` - The description of the organization.
* `main_name` - The main name of the organization owner.
* `name` - The name of the organization.
* `owner` - The owner id of the organization.
* `status` - The status of the organization.
* `type` - The type of the organization.
* `updated_time` - The updated time of the organization.


## Import
Organization can be imported using the id, e.g.
```
$ terraform import byteplus_organization.default resource_id
```

