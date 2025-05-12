package alb_certificate

import (
	"fmt"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Certificate can be imported using the id, e.g.
```
$ terraform import byteplus_alb_certificate.default cert-2fe5k****c16o5oxruvtk3qf5
```

*/

func ResourceByteplusAlbCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusCertificateCreate,
		Read:   resourceByteplusCertificateRead,
		Delete: resourceByteplusCertificateDelete,
		Update: resourceByteplusCertificateUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Certificate.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public key of the Certificate.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private key of the Certificate.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Certificate.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the Certificate.",
			},
			"listeners": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "The ID list of the Listener.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the Certificate.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expire time of the Certificate.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name of the Certificate.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the Certificate.",
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the Certificate.",
			},
		},
	}
}

func resourceByteplusCertificateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = certificateService.Dispatcher.Create(certificateService, d, ResourceByteplusAlbCertificate())
	if err != nil {
		return fmt.Errorf("error on creating certificate  %q, %w", d.Id(), err)
	}
	return resourceByteplusCertificateRead(d, meta)
}

func resourceByteplusCertificateRead(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = certificateService.Dispatcher.Read(certificateService, d, ResourceByteplusAlbCertificate())
	if err != nil {
		return fmt.Errorf("error on reading certificate %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusCertificateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = certificateService.Dispatcher.Update(certificateService, d, ResourceByteplusAlbCertificate())
	if err != nil {
		return fmt.Errorf("error on updating certificate  %q, %w", d.Id(), err)
	}
	return resourceByteplusCertificateRead(d, meta)
}

func resourceByteplusCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = certificateService.Dispatcher.Delete(certificateService, d, ResourceByteplusAlbCertificate())
	if err != nil {
		return fmt.Errorf("error on deleting certificate %q, %w", d.Id(), err)
	}
	return err
}
