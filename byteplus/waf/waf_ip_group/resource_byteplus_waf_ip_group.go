package waf_ip_group

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
WafIpGroup can be imported using the id, e.g.
```
$ terraform import byteplus_waf_ip_group.default resource_id
```

*/

func ResourceByteplusWafIpGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusWafIpGroupCreate,
		Read:   resourceByteplusWafIpGroupRead,
		Update: resourceByteplusWafIpGroupUpdate,
		Delete: resourceByteplusWafIpGroupDelete,
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
				Description: "The name of ip group.",
			},
			"add_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The way of addition.",
			},
			"ip_list": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The IP address to be added.",
			},
			"ip_group_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the ip group.",
			},
			"ip_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of IP addresses within the address group.",
			},
			"related_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of associated rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the rule.",
						},
						"rule_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the rule.",
						},
						"rule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the rule.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The information of the protected domain names associated with the rules.",
						},
					},
				},
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ip group update time.",
			},
		},
	}
	return resource
}

func resourceByteplusWafIpGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafIpGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusWafIpGroup())
	if err != nil {
		return fmt.Errorf("error on creating waf_ip_group %q, %s", d.Id(), err)
	}
	return resourceByteplusWafIpGroupRead(d, meta)
}

func resourceByteplusWafIpGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafIpGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusWafIpGroup())
	if err != nil {
		return fmt.Errorf("error on reading waf_ip_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusWafIpGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafIpGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusWafIpGroup())
	if err != nil {
		return fmt.Errorf("error on updating waf_ip_group %q, %s", d.Id(), err)
	}
	return resourceByteplusWafIpGroupRead(d, meta)
}

func resourceByteplusWafIpGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafIpGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusWafIpGroup())
	if err != nil {
		return fmt.Errorf("error on deleting waf_ip_group %q, %s", d.Id(), err)
	}
	return err
}
