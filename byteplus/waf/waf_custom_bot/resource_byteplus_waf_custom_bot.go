package waf_custom_bot

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
WafCustomBot can be imported using the id, e.g.
```
$ terraform import byteplus_waf_custom_bot.default resource_id
```

*/

func ResourceByteplusWafCustomBot() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusWafCustomBotCreate,
		Read:   resourceByteplusWafCustomBotRead,
		Update: resourceByteplusWafCustomBotUpdate,
		Delete: resourceByteplusWafCustomBotDelete,
		Importer: &schema.ResourceImporter{
			State: wafCustomBotImporter,
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
				Description: "bot name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of bot.",
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The execution action of the Bot.",
			},
			"enable": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Whether to enable bot.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of the affiliated project resource.",
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Domain name information.",
			},
			"accurate": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Advanced conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accurate_rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Details of advanced conditions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http_obj": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The HTTP object to be added to the advanced conditions.",
									},
									"obj_type": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The matching field for HTTP objects.",
									},
									"opretar": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "The logical operator for the condition.",
									},
									"property": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Operate the properties of the http object.",
									},
									"value_string": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The value to be matched.",
									},
								},
							},
						},
						"logic": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "The logical relationship of advanced conditions.",
						},
					},
				},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time.",
			},
			"advanced": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether to set advanced conditions.",
			},
			"rule_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Rule unique identifier.",
			},
		},
	}
	return resource
}

func resourceByteplusWafCustomBotCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCustomBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusWafCustomBot())
	if err != nil {
		return fmt.Errorf("error on creating waf_custom_bot %q, %s", d.Id(), err)
	}
	return resourceByteplusWafCustomBotRead(d, meta)
}

func resourceByteplusWafCustomBotRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCustomBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusWafCustomBot())
	if err != nil {
		return fmt.Errorf("error on reading waf_custom_bot %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusWafCustomBotUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCustomBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusWafCustomBot())
	if err != nil {
		return fmt.Errorf("error on updating waf_custom_bot %q, %s", d.Id(), err)
	}
	return resourceByteplusWafCustomBotRead(d, meta)
}

func resourceByteplusWafCustomBotDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCustomBotService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusWafCustomBot())
	if err != nil {
		return fmt.Errorf("error on deleting waf_custom_bot %q, %s", d.Id(), err)
	}
	return err
}
