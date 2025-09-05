package waf_system_bot

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
WafSystemBot can be imported using the id, e.g.
```
$ terraform import byteplus_waf_system_bot.default BotType:Host
```

*/

func ResourceByteplusWafSystemBot() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusWafSystemBotCreate,
		Read:   resourceByteplusWafSystemBotRead,
		Update: resourceByteplusWafSystemBotUpdate,
		Delete: resourceByteplusWafSystemBotDelete,
		Importer: &schema.ResourceImporter{
			State: wafSystemBotImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"bot_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of bot.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain name information.",
			},
			"enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable bot.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of the affiliated project resource.",
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The execution action of the Bot.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Bot.",
			},
			"rule_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Bot rule.",
			},
		},
	}
	return resource
}

func resourceByteplusWafSystemBotCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafSystemBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusWafSystemBot())
	if err != nil {
		return fmt.Errorf("error on creating waf_system_bot %q, %s", d.Id(), err)
	}
	return resourceByteplusWafSystemBotRead(d, meta)
}

func resourceByteplusWafSystemBotRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafSystemBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusWafSystemBot())
	if err != nil {
		return fmt.Errorf("error on reading waf_system_bot %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusWafSystemBotUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafSystemBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusWafSystemBot())
	if err != nil {
		return fmt.Errorf("error on updating waf_system_bot %q, %s", d.Id(), err)
	}
	return resourceByteplusWafSystemBotRead(d, meta)
}

func resourceByteplusWafSystemBotDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafSystemBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusWafSystemBot())
	if err != nil {
		return fmt.Errorf("error on deleting waf_system_bot %q, %s", d.Id(), err)
	}
	return err
}
