package cdn_edge_function

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnEdgeFunction can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_edge_function.default resource_id
```

*/

func ResourceByteplusCdnEdgeFunction() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnEdgeFunctionCreate,
		Read:   resourceByteplusCdnEdgeFunctionRead,
		Update: resourceByteplusCdnEdgeFunctionUpdate,
		Delete: resourceByteplusCdnEdgeFunctionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the edge function.",
			},
			"remark": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The remark of the edge function.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the project to which the edge function belongs, defaulting to `default`.",
			},
			"source_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Code content. The input requirements are as follows：\nNot empty.\nValue after base64 encoding.",
			},
			"envs": {
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         envsHash,
				Description: "The environment variables of the edge function.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the environment variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the environment variable.",
						},
					},
				},
			},
			"canary_countries": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The array of countries where the canary cluster is located.",
			},

			// computed fields
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the edge function.",
			},
			"account_identity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The account id of the edge function.",
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creator of the edge function.",
			},
			"user_identity": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The user id of the edge function.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The create time of the edge function. Displayed in UNIX timestamp format.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time of the edge function. Displayed in UNIX timestamp format.",
			},
			"domain": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The domain name bound to the edge function.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnEdgeFunctionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnEdgeFunction())
	if err != nil {
		return fmt.Errorf("error on creating cdn_edge_function %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnEdgeFunction())
	if err != nil {
		return fmt.Errorf("error on reading cdn_edge_function %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnEdgeFunctionUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnEdgeFunction())
	if err != nil {
		return fmt.Errorf("error on updating cdn_edge_function %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnEdgeFunction())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_edge_function %q, %s", d.Id(), err)
	}
	return err
}

var envsHash = func(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	var (
		buf bytes.Buffer
	)
	buf.WriteString(fmt.Sprintf("%v#%v", m["key"], m["value"]))
	return hashcode.String(buf.String())
}
