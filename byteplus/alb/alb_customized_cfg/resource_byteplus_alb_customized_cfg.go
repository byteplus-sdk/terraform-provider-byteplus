package alb_customized_cfg

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
AlbCustomizedCfg can be imported using the id, e.g.
```
$ terraform import byteplus_alb_customized_cfg.default ccfg-3cj44nv0jhhxc6c6rrtet****
```

*/

func ResourceByteplusAlbCustomizedCfg() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusAlbCustomizedCfgCreate,
		Read:   resourceByteplusAlbCustomizedCfgRead,
		Update: resourceByteplusAlbCustomizedCfgUpdate,
		Delete: resourceByteplusAlbCustomizedCfgDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"customized_cfg_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of CustomizedCfg.",
			},
			"customized_cfg_content": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The content of CustomizedCfg. The length cannot exceed 4096 characters. " +
					"Spaces and semicolons need to be escaped. " +
					"Currently supported configuration items are `ssl_protocols`, `ssl_ciphers`, `client_max_body_size`, `keepalive_timeout`, `proxy_request_buffering` and `proxy_connect_timeout`.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of CustomizedCfg.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The project name of the CustomizedCfg.",
			},
		},
	}
	return resource
}

func resourceByteplusAlbCustomizedCfgCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCustomizedCfgService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusAlbCustomizedCfg())
	if err != nil {
		return fmt.Errorf("error on creating alb_customized_cfg %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbCustomizedCfgRead(d, meta)
}

func resourceByteplusAlbCustomizedCfgRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCustomizedCfgService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusAlbCustomizedCfg())
	if err != nil {
		return fmt.Errorf("error on reading alb_customized_cfg %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusAlbCustomizedCfgUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCustomizedCfgService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusAlbCustomizedCfg())
	if err != nil {
		return fmt.Errorf("error on updating alb_customized_cfg %q, %s", d.Id(), err)
	}
	return resourceByteplusAlbCustomizedCfgRead(d, meta)
}

func resourceByteplusAlbCustomizedCfgDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewAlbCustomizedCfgService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusAlbCustomizedCfg())
	if err != nil {
		return fmt.Errorf("error on deleting alb_customized_cfg %q, %s", d.Id(), err)
	}
	return err
}
