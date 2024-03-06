package certificate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Certificate can be imported using the id, e.g.
```
$ terraform import byteplus_certificate.default cert-2fe5k****c16o5oxruvtk3qf5
```

*/

func ResourceByteplusCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusCertificateCreate,
		Read:   resourceByteplusCertificateRead,
		Update: resourceByteplusCertificateUpdate,
		Delete: resourceByteplusCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the Certificate.",
			},
			"public_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public key of the Certificate. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The private key of the Certificate. When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
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
				ForceNew:    true,
				Description: "The ProjectName of the Certificate.",
			},
			"tags": bp.TagsSchema(),
		},
	}
}

func resourceByteplusCertificateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(certificateService, d, ResourceByteplusCertificate())
	if err != nil {
		return fmt.Errorf("error on creating certificate  %q, %w", d.Id(), err)
	}
	return resourceByteplusCertificateRead(d, meta)
}

func resourceByteplusCertificateRead(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(certificateService, d, ResourceByteplusCertificate())
	if err != nil {
		return fmt.Errorf("error on reading certificate %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusCertificateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(certificateService, d, ResourceByteplusCertificate())
	if err != nil {
		return fmt.Errorf("error on updating certificate  %q, %w", d.Id(), err)
	}
	return resourceByteplusCertificateRead(d, meta)
}

func resourceByteplusCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	certificateService := NewCertificateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(certificateService, d, ResourceByteplusCertificate())
	if err != nil {
		return fmt.Errorf("error on deleting certificate %q, %w", d.Id(), err)
	}
	return err
}
