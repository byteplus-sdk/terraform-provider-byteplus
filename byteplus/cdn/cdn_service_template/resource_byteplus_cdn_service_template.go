package cdn_service_template

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnServiceTemplate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_service_template.default resource_id
```

*/

func ResourceByteplusCdnServiceTemplate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnServiceTemplateCreate,
		Read:   resourceByteplusCdnServiceTemplateRead,
		Update: resourceByteplusCdnServiceTemplateUpdate,
		Delete: resourceByteplusCdnServiceTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the name of the encryption policy you want to create. The name must not exceed 100 characters.",
			},
			"message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the description of the encryption policy, which must not exceed 120 characters.",
			},
			"project": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: "Indicates the project to which this encryption policy belongs. The default value of the parameter is default, indicating the Default project.",
			},
			"service_template_config": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The service template configuration. " +
					"Please convert the configuration module structure into json and pass it into a string. " +
					"You must specify the Origin module. The OriginProtocol parameter, and other domain configuration modules are optional. " +
					"For detailed parameter introduction, please refer to `https://docs.byteplus.com/en/docs/byteplus-cdn/reference-updateservicetemplate`.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnServiceTemplateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnServiceTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnServiceTemplate())
	if err != nil {
		return fmt.Errorf("error on creating cdn_service_template %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnServiceTemplateRead(d, meta)
}

func resourceByteplusCdnServiceTemplateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnServiceTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnServiceTemplate())
	if err != nil {
		return fmt.Errorf("error on reading cdn_service_template %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnServiceTemplateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnServiceTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnServiceTemplate())
	if err != nil {
		return fmt.Errorf("error on updating cdn_service_template %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnServiceTemplateRead(d, meta)
}

func resourceByteplusCdnServiceTemplateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnServiceTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnServiceTemplate())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_service_template %q, %s", d.Id(), err)
	}
	return err
}
