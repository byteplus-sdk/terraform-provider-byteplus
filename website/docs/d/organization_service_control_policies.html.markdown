---
subcategory: "ORGANIZATION"
layout: "byteplus"
page_title: "Byteplus: byteplus_organization_service_control_policies"
sidebar_current: "docs-byteplus-datasource-organization_service_control_policies"
description: |-
  Use this data source to query detailed information of organization service control policies
---
# byteplus_organization_service_control_policies
Use this data source to query detailed information of organization service control policies
## Example Usage
```hcl
data "byteplus_organization_service_control_policies" "foo" {
  policy_type = "Custom"
  query       = "test"
}
```
## Argument Reference
The following arguments are supported:
* `output_file` - (Optional) File name where to save data source results.
* `policy_type` - (Optional) The type of policy. The value can be System or Custom.
* `query` - (Optional) Query policies, support policy name or description.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `policies` - The collection of Policy query.
    * `create_date` - The create time of the Policy.
    * `description` - The description of the Policy.
    * `id` - The ID of the Policy.
    * `policy_name` - The name of the Policy.
    * `policy_type` - The type of the Policy.
    * `statement` - The statement of the Policy.
    * `update_date` - The update time of the Policy.
* `total_count` - The total count of Policy query.


