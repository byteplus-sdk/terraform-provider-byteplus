package cdn_certificate

import (
	"fmt"
	"log"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnCertificate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_certificate.default resource_id
```

*/

func ResourceByteplusCdnCertificate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnCertificateCreate,
		Read:   resourceByteplusCdnCertificateRead,
		Delete: resourceByteplusCdnCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Indicates the content of the certificate file, which must include the complete certificate chain. The line breaks in the content should be replaced with \\r\\n. The certificate file must have an extension of either .crt or .pem." +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Indicates the content of the certificate private key file. The line breaks in the content should be replaced with \\r\\n. The certificate private key file must have an extension of either .key or .pem." +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"desc": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Indicates the remarks for the certificate.",
			},
			"repeatable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
				Description: "Indicates whether uploading the same certificate is allowed. If the fingerprints of two certificates are the same, these certificates are considered identical. This parameter can take the following values:\n\ntrue: Allows the upload of the same certificate.\nfalse: Does not allow the upload of the same certificate. When calling this API, the CDN will check for the existence of an identical certificate. If one exists, you will not be able to upload the certificate, and the Error structure in the response body will include the ID of the existing certificate.\nThe default value of this parameter is true." +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},

			// computed fields
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source of the certificate.",
			},
			"cert_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the content of the Common Name (CN) field of the certificate.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the status of the certificate. The parameter can have the following values:\nrunning: indicates the certificate has a remaining validity period of more than 30 days.\nexpired: indicates the certificate has expired.\nexpiring_soon: indicates the certificate has a remaining validity period of 30 days or less but has not yet expired.",
			},
			"dns_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the domain names in the SAN field of the certificate.",
			},
			"configured_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the list of domain names associated with the certificate. If the certificate has not been associated with any domain name, the parameter value is null.",
			},
			"effective_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the issuance time of the certificate. The unit is Unix timestamp.",
			},
			"expire_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Indicates the expiration time of the certificate. The unit is Unix timestamp.",
			},
			"cert_fingerprint": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sha1": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates a fingerprint based on the SHA-1 encryption algorithm, composed of 40 hexadecimal characters.",
						},
						"sha256": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates a fingerprint based on the SHA-256 encryption algorithm, composed of 64 hexadecimal characters.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusCdnCertificateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnCertificate())
	if err != nil {
		return fmt.Errorf("error on creating cdn_certificate %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCertificateRead(d, meta)
}

func resourceByteplusCdnCertificateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnCertificate())
	if err != nil {
		return fmt.Errorf("error on reading cdn_certificate %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnCertificateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnCertificate())
	if err != nil {
		return fmt.Errorf("error on updating cdn_certificate %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCertificateRead(d, meta)
}

func resourceByteplusCdnCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf("[DEBUG] deleting a byteplus_cdn_certificate resource will only remove the cdn certificate from terraform state.")
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnCertificate())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_certificate %q, %s", d.Id(), err)
	}
	return err
}
