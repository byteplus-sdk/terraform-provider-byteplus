package tos_bucket_cors

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
TosBucketCors can be imported using the id, e.g.
```
$ terraform import byteplus_tos_bucket_cors.default resource_id
```

*/

func ResourceByteplusTosBucketCors() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusTosBucketCorsCreate,
		Read:   resourceByteplusTosBucketCorsRead,
		Update: resourceByteplusTosBucketCorsUpdate,
		Delete: resourceByteplusTosBucketCorsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bucket.",
			},
			"cors_rules": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The CORS rules of the bucket.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of origins that are allowed to make requests to the bucket.",
						},
						"allowed_methods": {
							Type:     schema.TypeSet,
							Required: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of HTTP methods that are allowed in a preflight request. Valid values: `PUT`, `POST`, `DELETE`, `GET`, `HEAD`.",
						},
						"allowed_headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of headers that are allowed in a preflight request.",
						},
						"expose_headers": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      schema.HashString,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of headers that are exposed in the response to a preflight request. It is recommended to add two expose headers, X-Tos-Request-Id and ETag.",
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The maximum amount of time that a preflight request can be cached. Unit: second. Default value: 3600.",
						},
						"response_vary": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether the bucket returns the 'Vary: Origin' header in the response to preflight requests. Default value: false.",
						},
					},
				},
			},
		},
	}
	return resource
}

func resourceByteplusTosBucketCorsCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCorsService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusTosBucketCors())
	if err != nil {
		return fmt.Errorf("error on creating tos_bucket_cors %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketCorsRead(d, meta)
}

func resourceByteplusTosBucketCorsRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCorsService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusTosBucketCors())
	if err != nil {
		return fmt.Errorf("error on reading tos_bucket_cors %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusTosBucketCorsUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCorsService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusTosBucketCors())
	if err != nil {
		return fmt.Errorf("error on updating tos_bucket_cors %q, %s", d.Id(), err)
	}
	return resourceByteplusTosBucketCorsRead(d, meta)
}

func resourceByteplusTosBucketCorsDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewTosBucketCorsService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusTosBucketCors())
	if err != nil {
		return fmt.Errorf("error on deleting tos_bucket_cors %q, %s", d.Id(), err)
	}
	return err
}
