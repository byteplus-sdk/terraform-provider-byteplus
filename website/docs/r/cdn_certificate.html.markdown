---
subcategory: "CLASSIC_CDN"
layout: "byteplus"
page_title: "Byteplus: byteplus_cdn_certificate"
sidebar_current: "docs-byteplus-resource-cdn_certificate"
description: |-
  Provides a resource to manage cdn certificate
---
# byteplus_cdn_certificate
Provides a resource to manage cdn certificate
## Example Usage
```hcl
resource "byteplus_cdn_certificate" "foo" {
  certificate = "-----BEGIN CERTIFICATE----- *** -----END CERTIFICATE-----"
  private_key = "-----BEGIN PRIVATE KEY----- *** -----END PRIVATE KEY-----"
  desc        = "tf-test"
  repeatable  = true
}
```
## Argument Reference
The following arguments are supported:
* `certificate` - (Required, ForceNew) Indicates the content of the certificate file, which must include the complete certificate chain. The line breaks in the content should be replaced with \r\n. The certificate file must have an extension of either `.crt` or `.pem`.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `private_key` - (Required, ForceNew) Indicates the content of the certificate private key file. The line breaks in the content should be replaced with \r\n. The certificate private key file must have an extension of either `.key` or `.pem`.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.
* `desc` - (Optional, ForceNew) Indicates the remarks for the certificate.
* `repeatable` - (Optional, ForceNew) Indicates whether uploading the same certificate is allowed. If the fingerprints of two certificates are the same, these certificates are considered identical. This parameter can take the following values:

true: Allows the upload of the same certificate.
false: Does not allow the upload of the same certificate. When calling this API, the CDN will check for the existence of an identical certificate. If one exists, you will not be able to upload the certificate, and the Error structure in the response body will include the ID of the existing certificate.
The default value of this parameter is true.
When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.
* `cert_fingerprint` - Indicates the fingerprint information of the certificate.
    * `sha1` - Indicates a fingerprint based on the SHA-1 encryption algorithm, composed of 40 hexadecimal characters.
    * `sha256` - Indicates a fingerprint based on the SHA-256 encryption algorithm, composed of 64 hexadecimal characters.
* `cert_name` - Indicates the content of the Common Name (CN) field of the certificate.
* `configured_domain` - Indicates the list of domain names associated with the certificate. If the certificate has not been associated with any domain name, the parameter value is null.
* `dns_name` - Indicates the domain names in the SAN field of the certificate.
* `effective_time` - Indicates the issuance time of the certificate. The unit is Unix timestamp.
* `expire_time` - Indicates the expiration time of the certificate. The unit is Unix timestamp.
* `source` - The source of the certificate.
* `status` - Indicates the status of the certificate. The parameter can have the following values:
running: indicates the certificate has a remaining validity period of more than 30 days.
expired: indicates the certificate has expired.
expiring_soon: indicates the certificate has a remaining validity period of 30 days or less but has not yet expired.


## Import
CdnCertificate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_certificate.default resource_id
```
You can delete the certificate hosted on the content delivery network.
You can configure the HTTPS module to associate the certificate and domain name through the domain_config field of byteplus_cdn_domain.
If the certificate to be deleted is already associated with a domain name, the deletion will fail.
To remove the association between the domain name and the certificate, you can disable the HTTPS function for the domain name in the Content Delivery Network console.

