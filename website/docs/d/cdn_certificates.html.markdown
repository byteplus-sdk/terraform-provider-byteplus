---
subcategory: "CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_certificates"
sidebar_current: "docs-byteplus-datasource-cdn_certificates"
description: |-
  Use this data source to query detailed information of cdn certificates
---
# byteplus_cdn_certificates
Use this data source to query detailed information of cdn certificates
## Example Usage
```hcl
data "byteplus_cdn_certificates" "foo" {
  configured_domain = ["byteplus-demo.byte-test.com"]
  name              = "*.byte-test.com"
  fuzzy_match       = true
  status            = "running"
}
```
## Argument Reference
The following arguments are supported:
* `cert_id` - (Optional) Indicates a certificate ID to retrieve the certificate with that ID.
* `configured_domain` - (Optional) Indicates a list of domain names for acceleration, to obtain certificates that have been bound to any domain name on the list.
* `fuzzy_match` - (Optional) When Name is specified, FuzzyMatch indicates the matching method used by the CDN when filtering certificates by Name. The parameter can have the following values:
true: indicates fuzzy matching.
false: indicates exact matching.
If you don not specify Name, FuzzyMatch is not effective.
The default value of FuzzyMatch is false.
* `name_regex` - (Optional) A Name Regex of Resource.
* `name` - (Optional) Indicates a domain name used to obtain certificates that include that domain name in the SAN field. The domain name can be a wildcard domain. For example, *.example.com can match certificates containing img.example.com or www.example.com, etc., in the SAN field.
* `output_file` - (Optional) File name where to save data source results.
* `status` - (Optional) Indicates a list of states to retrieve certificates that are in any of the states on the list. The parameter can have the following values:
running: indicates certificates with a remaining validity period of more than 30 days.
expired: indicates certificates that have expired.
expiring_soon: indicates certificates with a remaining validity period of 30 days or less but have not yet expired.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `certificates` - The collection of query.
    * `cert_fingerprint` - Indicates the fingerprint information of the certificate.
        * `sha1` - Indicates a fingerprint based on the SHA-1 encryption algorithm, composed of 40 hexadecimal characters.
        * `sha256` - Indicates a fingerprint based on the SHA-256 encryption algorithm, composed of 64 hexadecimal characters.
    * `cert_id` - Indicates the ID of the certificate.
    * `cert_name` - Indicates the content of the Common Name (CN) field of the certificate.
    * `configured_domain` - Indicates the list of domain names associated with the certificate. If the certificate has not been associated with any domain name, the parameter value is null.
    * `desc` - Indicates the remark of the certificate.
    * `dns_name` - Indicates the domain names in the SAN field of the certificate.
    * `effective_time` - Indicates the issuance time of the certificate. The unit is Unix timestamp.
    * `expire_time` - Indicates the expiration time of the certificate. The unit is Unix timestamp.
    * `id` - Indicates the ID of the certificate.
    * `source` - The source of the certificate.
    * `status` - Indicates the status of the certificate. The parameter can have the following values:
running: indicates the certificate has a remaining validity period of more than 30 days.
expired: indicates the certificate has expired.
expiring_soon: indicates the certificate has a remaining validity period of 30 days or less but has not yet expired.
* `total_count` - The total count of query.


