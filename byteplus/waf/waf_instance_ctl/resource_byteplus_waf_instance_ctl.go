package waf_instance_ctl

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
WafInstanceCtl can be imported using the id, e.g.
```
$ terraform import byteplus_waf_instance_ctl.default resource_id
```

*/

func ResourceByteplusWafInstanceCtl() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusWafInstanceCtlCreate,
		Read:   resourceByteplusWafInstanceCtlRead,
		Update: resourceByteplusWafInstanceCtlUpdate,
		Delete: resourceByteplusWafInstanceCtlDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the allowed access list policy for the instance corresponding to the current region.",
			},
			"block_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the prohibited access list policy for the instance corresponding to the current region.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the project associated with the current resource.",
			},
		},
	}
	return resource
}

func resourceByteplusWafInstanceCtlCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafInstanceCtlService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusWafInstanceCtl())
	if err != nil {
		return fmt.Errorf("error on creating waf_instance_ctl %q, %s", d.Id(), err)
	}
	return resourceByteplusWafInstanceCtlRead(d, meta)
}

func resourceByteplusWafInstanceCtlRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafInstanceCtlService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusWafInstanceCtl())
	if err != nil {
		return fmt.Errorf("error on reading waf_instance_ctl %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusWafInstanceCtlUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafInstanceCtlService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusWafInstanceCtl())
	if err != nil {
		return fmt.Errorf("error on updating waf_instance_ctl %q, %s", d.Id(), err)
	}
	return resourceByteplusWafInstanceCtlRead(d, meta)
}

func resourceByteplusWafInstanceCtlDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafInstanceCtlService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusWafInstanceCtl())
	if err != nil {
		return fmt.Errorf("error on deleting waf_instance_ctl %q, %s", d.Id(), err)
	}
	return err
}
