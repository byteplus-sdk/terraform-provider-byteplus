package alb_ca_certificate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
AlbCaCertificate can be imported using the id, e.g.
```
$ terraform import byteplus_alb_ca_certificate.default cert-*****
```

*/

func ResourceByteplusAlbCaCertificate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusAlbCaCertificateCreate,
		Read:   resourceByteplusAlbCaCertificateRead,
		Update: resourceByteplusAlbCaCertificateUpdate,
		Delete: resourceByteplusAlbCaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"ca_certificate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the CA certificate.",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The content of the CA certificate.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the CA certificate.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the CA certificate.",
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
				Description: "The create time of the CA Certificate.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expire time of the CA Certificate.",
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain name of the CA Certificate.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the CA Certificate.",
			},
			// 文档与接口实际返回不同
			"certificate_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the CA Certificate.",
			},
		},
	}
	return resource
}

func resourceByteplusAlbCaCertificateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCaCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusAlbCaCertificate())
	if err != nil {
		return fmt.Errorf("error on creating alb_ca_certificate %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbCaCertificateRead(d, meta)
}

func resourceByteplusAlbCaCertificateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCaCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusAlbCaCertificate())
	if err != nil {
		return fmt.Errorf("error on reading alb_ca_certificate %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusAlbCaCertificateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCaCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusAlbCaCertificate())
	if err != nil {
		return fmt.Errorf("error on updating alb_ca_certificate %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbCaCertificateRead(d, meta)
}

func resourceByteplusAlbCaCertificateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCaCertificateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusAlbCaCertificate())
	if err != nil {
		return fmt.Errorf("error on deleting alb_ca_certificate %q, %s", d.Id(), err)
	}
	return err
}
