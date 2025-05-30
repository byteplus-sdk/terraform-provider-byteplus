package classic_cdn_certificate

import (
	"fmt"
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
You can delete the certificate hosted on the content delivery network.
You can configure the HTTPS module to associate the certificate and domain name through the domain_config field of byteplus_cdn_domain.
If the certificate to be deleted is already associated with a domain name, the deletion will fail.
To remove the association between the domain name and the certificate, you can disable the HTTPS function for the domain name in the Content Delivery Network console.
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
				Description: "Content of the specified certificate public key file. " +
					"Line breaks in the content should be replaced with `\\r\\n`. " +
					"The file extension for the certificate public key is `.crt` or `.pem`. " +
					"The public key must include the complete certificate chain. " +
					"When importing resources, this attribute will not be imported. " +
					"If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"private_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The content of the specified certificate private key file. " +
					"Replace line breaks in the content with `\\r\\n`. " +
					"The file extension for the certificate private key is `.key` or `.pem`. " +
					"The private key must be unencrypted. " +
					"When importing resources, this attribute will not be imported. " +
					"If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Specify the location for storing the certificate. " +
					"The parameter can take the following values: " +
					"`cert_center`: indicates that the certificate will be stored in the certificate center." +
					"`cdn_cert_hosting`: indicates that the certificate will be hosted on the content delivery network.",
			},
			"desc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Note on the certificate.",
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

func resourceByteplusCdnCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnCertificate())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_certificate %q, %s", d.Id(), err)
	}
	return err
}
