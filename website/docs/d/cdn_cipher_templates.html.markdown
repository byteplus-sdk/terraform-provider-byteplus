---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_cipher_templates"
sidebar_current: "docs-byteplus-datasource-cdn_cipher_templates"
description: |-
  Use this data source to query detailed information of cdn cipher templates
---
# byteplus_cdn_cipher_templates
Use this data source to query detailed information of cdn cipher templates
## Example Usage
```hcl
data "byteplus_cdn_cipher_templates" "foo" {}
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
    * `http_forced_redirect` - Indicates the configuration module for the forced redirection from HTTPS to HTTP. This feature is disabled by default.
        * `enable_forced_redirect` - Indicates whether to enable the forced redirection from HTTPS. This parameter can take the following values: true: Indicates to enable the forced redirection from HTTPS. Once enabled, the content delivery network will respond with StatusCode to inform the browser to send an HTTPS request when it receives an HTTP request from a user. false: Indicates to disable the forced redirection from HTTPS.
        * `status_code` - Indicates the status code returned by the content delivery network when forced redirection from HTTPS occurs. The default value for this parameter is 301.
    * `https` - Indicates the configuration module for the HTTPS encryption service.
        * `disable_http` - Indicates whether the CDN accepts HTTP user requests. This parameter can take the following values: true: Indicates that it does not accept. If an HTTP request is received, the CDN will reject the request. false: Indicates that it accepts. The default value for this parameter is false.
        * `forced_redirect` - Indicates the configuration for the mandatory redirection from HTTP to HTTPS. This feature is disabled by default.
            * `enable_forced_redirect` - Indicates the switch for the Forced Redirect configuration. This parameter can take the following values: true: Indicates to enable Forced Redirect. false: Indicates to disable Forced Redirect.
            * `status_code` - Indicates the status code returned to the client by the CDN when forced redirect occurs. This parameter can take the following values: 301: Indicates that the returned status code is 301. 302: Indicates that the returned status code is 302. The default value for this parameter is 301.
        * `hsts` - Indicates the HSTS (HTTP Strict Transport Security) configuration module. This feature is disabled by default.
            * `subdomain` - Indicates whether the HSTS configuration should also be applied to the subdomains of the domain name. This parameter can take the following values: include: Indicates that HSTS settings apply to subdomains. exclude: Indicates that HSTS settings do not apply to subdomains. The default value for this parameter is exclude.
            * `switch` - Indicates whether to enable HSTS. This parameter can take the following values: true: Indicates to enable HSTS. false: Indicates to disable HSTS. The default value for this parameter is false.
            * `ttl` - Indicates the expiration time for the Strict-Transport-Security response header in the browser cache, in seconds. If Switch is true, this parameter is required. The value range for this parameter is 0 - 31,536,000 seconds, where 31,536,000 seconds represents 365 days. If the value of this parameter is 0, it is equivalent to disabling the HSTS settings.
        * `http2` - Indicates the switch for HTTP/2 configuration. This parameter can take the following values: true: Indicates to enable HTTP/2. false: Indicates to disable HTTP/2. The default value for this parameter is true.
        * `ocsp` - Indicates whether to enable OCSP Stapling. This parameter can take the following values: true: Indicates to enable OCSP Stapling. false: Indicates to disable OCSP Stapling. The default value for this parameter is false.
        * `tls_version` - Indicates a list that specifies the TLS versions supported by the domain name. This parameter can take the following values: tlsv1.0: Indicates TLS 1.0. tlsv1.1: Indicates TLS 1.1. tlsv1.2: Indicates TLS 1.2. tlsv1.3: Indicates TLS 1.3. The default value for this parameter is ["tlsv1.1", "tlsv1.2", "tlsv1.3"].
    * `message` - Indicates the description of the policy.
    * `project` - Indicates the project to which the policy belongs.
    * `quic` - Indicates the QUIC configuration module. This feature is disabled by default.
        * `switch` - Indicates whether to enable QUIC. This parameter can take the following values: true: Indicates to enable QUIC. false: Indicates to disable QUIC.
    * `status` - Indicates the status of the policy. This parameter can take the following values: locked: Indicates the status is "published". editing: Indicates the status is "draft".
    * `template_id` - Indicates the ID of a policy in the list of policies returned by the API.
    * `title` - Indicates the name of the policy.
    * `type` - Indicates the type of the policy. This parameter can take the following values: cipher: Indicates an encryption policy. service: Indicates a distribution policy.
    * `update_time` - Indicates the last modification time of the policy, in Unix timestamp format. If the policy has not been updated since its creation, the value of this parameter is the same as CreateTime.
* `total_count` - The total count of query.


