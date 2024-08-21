---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_domains"
sidebar_current: "docs-byteplus-datasource-cdn_domains"
description: |-
  Use this data source to query detailed information of cdn domains
---
# byteplus_cdn_domains
Use this data source to query detailed information of cdn domains
## Example Usage
```hcl
data "byteplus_cdn_domains" "foo" {}
```
## Argument Reference
The following arguments are supported:
* `filters` - (Optional) Indicates a set of filtering criteria used to obtain a list of policies that meet these criteria. If you do not specify any filtering criteria, this API returns all policies under your account. Multiple filtering criteria are related by AND, meaning only policies that meet all filtering criteria will be included in the list returned by this API. In the API response, the actual policies returned are affected by PageNum and PageSize.
* `name_regex` - (Optional) A Name Regex of Resource.
* `output_file` - (Optional) File name where to save data source results.

The `filters` object supports the following:

* `fuzzy` - (Optional) Indicates the matching method. This parameter can take the following values: true: Indicates fuzzy matching. A policy is considered to meet the filtering criteria if the corresponding value of Name contains any value in the Value array. false: Indicates exact matching. A policy is considered to meet the filtering criteria if the corresponding value of Name matches any value in the Value array. Moreover, the Fuzzy value you can specify is affected by the Name value. See the description of Name. The default value of this parameter is false. Note that the matching process is case-sensitive.
* `name` - (Optional) Indicates the filtering type. This parameter can take the following values: TemplateTitle: Filters domain names by the name of the bound policy. TemplateId: Filters domain names by the ID of the bound policy. For this parameter value, the Fuzzy value can only be false. TemplateType: Filters domain names by the type of the bound policy. For this parameter value, the Fuzzy value can only be false. Domain: Filters domain names by name. Status: Filters domain names by status. For this parameter value, the Fuzzy value can only be false. HTTPSSwitch: Filters domain names by the status of the HTTPS encryption service. For this parameter value, the Fuzzy value can only be false. WAFStatus: Filters domain names by the status of WAF protection. For this parameter value, the Fuzzy value can only be false. Multiple filtering conditions can be specified at the same time, but the Name in different filtering conditions cannot be the same.
* `value` - (Optional) Indicates the values corresponding to Name, which is an array. The values in the array are used to match against the object value. If the object value matches any value in the array, it is considered a match. Values are case-sensitive when matching. When Name is TemplateTitle or Domain, each value in Value cannot exceed 100 characters. Furthermore, When Fuzzy is false, Value can contain multiple values. When Fuzzy is true, Value can only contain one value. When Name is TemplateId, Value can only contain one value. When Name is TemplateType, Value can contain one or more of the following values: cipher: Indicates an encryption policy. service: Indicates a delivery policy. When Name is Status, Value can contain one of the following values: online: Indicates that the status of the domain name is Enabled. offline: Indicates that the status of the domain name is Disabled. configuring: Indicates that the status of the domain name is Configuring. When Name is HTTPSSwitch, Value can contain one of the following values: on: Indicates that the HTTPS encryption service is enabled. off: Indicates that the HTTPS encryption service is disabled. When Name is WAFStatus, Value can contain one of the following values: on: Indicates that WAF protection is enabled. off: Indicates that WAF protection is disabled.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `domains` - The collection of query.
    * `cert_info` - Indicates information about the certificate associated with the domain name specified by Domain. If HTTPSSwitch is off, the value of this parameter is null.
        * `cert_id` - Indicates the ID of the certificate associated with the domain name. This certificate is hosted in the BytePlus Certificate Center.
    * `cname` - Indicates the CNAME assigned to the domain name by CDN.
    * `domain` - Represents one of the domain names in the Domains list.
    * `https_switch` - Indicates whether the domain name has enabled HTTPS Encryption Service. This parameter can be: on: Indicates that the service is enabled. off: Indicates that the service is disabled.
    * `lock_status` - Indicates whether the domain name is locked. This parameter can be: on: Indicates that the domain name is locked. In this case, you cannot use UpdateTemplateDomain to change the configurations of this domain name. off: Indicates that the domain name is not locked.
    * `project` - Indicates the project to which the domain name belongs.
    * `remark` - Indicates the reason why the domain name is locked. If LockStatus is on, this parameter indicates the reason why the domain name is locked. If LockStatus is off, the value of this parameter is empty ("").
    * `service_region` - Indicates the service region enabled for the domain name. This parameter can be: outside_chinese_mainland: Indicates Global (Excluding Chinese Mainland). chinese_mainland: Indicates Chinese Mainland. global: Indicates Global.
    * `status` - Indicates the status of the domain name. This parameter can be: online: Indicates the status is Enabled. offline: Indicates the status is Disabled. configuring: Indicates the status is Configuring.
    * `templates` - Indicates the list of policies bound to the domain name. A domain name must and can only be bound to one delivery policy, and optionally to one encryption policy.
        * `exception` - Indicates whether the policy contains special configurations. Special configurations refer to those configurations that are operated by BytePlus engineers instead of users. This parameter can be: true: Indicates it contains special configurations. false: Indicates it does not contain special configurations.
        * `template_id` - Indicates the ID of a policy.
        * `title` - Indicates the name of the policy.
        * `type` - Indicates the type of the policy. This parameter can be: cipher: Indicates an encryption policy. service: Indicates a delivery policy.
    * `waf_status` - Indicates whether the domain name has enabled WAF Protection. This parameter can be: on: Indicates that WAF Protection is enabled. off: Indicates that WAF Protection is disabled.
* `total_count` - The total count of query.


