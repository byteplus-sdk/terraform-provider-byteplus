package alb_health_check_template

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
AlbHealthCheckTemplate can be imported using the id, e.g.
```
$ terraform import byteplus_alb_health_check_template.default hctpl-123*****432
```

*/

func ResourceByteplusAlbHealthCheckTemplate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusAlbHealthCheckTemplateCreate,
		Read:   resourceByteplusAlbHealthCheckTemplateRead,
		Update: resourceByteplusAlbHealthCheckTemplateUpdate,
		Delete: resourceByteplusAlbHealthCheckTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"health_check_template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The health check template name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of health check template.",
			},
			"health_check_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The interval for performing health checks, the default value is 2, and the value is 1-300.",
			},
			"health_check_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The timeout of health check response,the default value is 2, and the value is 1-60.",
			},
			"healthy_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The healthy threshold of the health check, the default is 3, the value is 2-10.",
			},
			"unhealthy_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The unhealthy threshold of the health check, the default is 3, the value is 2-10.",
			},
			"health_check_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The health check method,default is `GET`, support `GET` and `HEAD`.",
			},
			"health_check_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The domain name to health check.",
			},
			"health_check_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The uri to health check,default is `/`.",
			},
			"health_check_http_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The normal HTTP status code for health check, the default is http_2xx, http_3xx, separated by commas.",
			},
			"health_check_protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "THe protocol of health check,only support HTTP.",
			},
			"health_check_http_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The HTTP version of health check.",
			},
		},
	}
	return resource
}

func resourceByteplusAlbHealthCheckTemplateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbHealthCheckTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusAlbHealthCheckTemplate())
	if err != nil {
		return fmt.Errorf("error on creating alb_health_check_template %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbHealthCheckTemplateRead(d, meta)
}

func resourceByteplusAlbHealthCheckTemplateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbHealthCheckTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusAlbHealthCheckTemplate())
	if err != nil {
		return fmt.Errorf("error on reading alb_health_check_template %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusAlbHealthCheckTemplateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbHealthCheckTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusAlbHealthCheckTemplate())
	if err != nil {
		return fmt.Errorf("error on updating alb_health_check_template %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbHealthCheckTemplateRead(d, meta)
}

func resourceByteplusAlbHealthCheckTemplateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbHealthCheckTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusAlbHealthCheckTemplate())
	if err != nil {
		return fmt.Errorf("error on deleting alb_health_check_template %q, %s", d.Id(), err)
	}
	return err
}
