---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_service_templates"
sidebar_current: "docs-byteplus-datasource-cdn_service_templates"
description: |-
  Use this data source to query detailed information of cdn service templates
---
# byteplus_cdn_service_templates
Use this data source to query detailed information of cdn service templates
## Example Usage
```hcl
data "byteplus_cdn_service_templates" "foo" {}
```
## Argument Reference
The following arguments are supported:
* `filters` - (Optional) Indicates a set of filtering criteria used to obtain a list of policies that meet these criteria. If you do not specify any filtering criteria, this API returns all policies under your account. Multiple filtering criteria are related by AND, meaning only policies that meet all filtering criteria will be included in the list returned by this API. In the API response, the actual policies returned are affected by PageNum and PageSize.
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

The `filters` object supports the following:

* `fuzzy` - (Optional) Indicates the matching method. This parameter can take the following values: true: Indicates fuzzy matching. A policy is considered to meet the filtering criteria if the corresponding value of Name contains any value in the Value array. false: Indicates exact matching. A policy is considered to meet the filtering criteria if the corresponding value of Name matches any value in the Value array. Moreover, the Fuzzy value you can specify is affected by the Name value. See the description of Name. The default value of this parameter is false. Note that the matching process is case-sensitive.
* `name` - (Optional) Represents the filtering type. This parameter can take the following values: Title: Filters policies by name. Id: Filters policies by ID. For this parameter value, the value of Fuzzy can only be false. Domain: Filters policies by the bound domain name. Type: Filters policies by type. For this parameter value, the value of Fuzzy can only be false. Status: Filters policies by status. For this parameter value, the value of Fuzzy can only be false. You can specify multiple filtering criteria simultaneously, but the Name in different filtering criteria cannot be the same.
* `value` - (Optional) Represents the values corresponding to Name, which is an array. When Name is Title, Id, or Domain, each value in the Value array should not exceed 100 characters in length. When Name is Type, the Value array can include one or more of the following values: cipher: Indicates a encryption policy. service: Indicates a delivery policy. When Name is Status, the Value array can include one or more of the following values: locked: Indicates the status is "published". editing: Indicates the status is "draft". When Fuzzy is false, you can specify multiple values in the array. When Fuzzy is true, you can only specify one value in the array.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `templates` - The collection of query.
    * `bound_domains` - Represents a list of domain names bound to the policy specified by TemplateId. If the policy is not bound to any domain names, the value of this parameter is null.
        * `bound_time` - Indicates the time when the policy was bound to the domain name specified by Domain, in Unix timestamp format.
        * `domain` - Represents one of the domain names bound to the policy.
    * `create_time` - Indicates the creation time of the policy, in Unix timestamp format.
    * `exception` - Indicates whether the policy includes special configurations. Special configurations refer to those not operated by users but by BytePlus engineers. This parameter can take the following values: true: Indicates it includes special configurations. false: Indicates it does not include special configurations.
    * `message` - Indicates the description of the policy.
    * `project` - Indicates the project to which the policy belongs.
    * `status` - Indicates the status of the policy. This parameter can take the following values: locked: Indicates the status is "published". editing: Indicates the status is "draft".
    * `template_id` - Indicates the ID of a policy in the list of policies returned by the API.
    * `title` - Indicates the name of the policy.
    * `type` - Indicates the type of the policy. This parameter can take the following values: cipher: Indicates an encryption policy. service: Indicates a distribution policy.
    * `update_time` - Indicates the last modification time of the policy, in Unix timestamp format. If the policy has not been updated since its creation, the value of this parameter is the same as CreateTime.
* `total_count` - The total count of query.


