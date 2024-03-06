package rule

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Rule can be imported using the id, e.g.
Notice: resourceId is ruleId, due to the lack of describeRuleAttributes in openapi, for import resources, please use ruleId:listenerId to import.
we will fix this problem later.
```
$ terraform import byteplus_clb_rule.foo rule-273zb9hzi1gqo7fap8u1k3utb:lsn-273ywvnmiu70g7fap8u2xzg9d
```

*/

func resourceParseId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of Id (%s), expected ruleId:listenerId", id)
	}
	return parts[0], parts[1], nil
}

func ResourceByteplusRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusRuleCreate,
		Read:   resourceByteplusRuleRead,
		Update: resourceByteplusRuleUpdate,
		Delete: resourceByteplusRuleDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				ruleId, listenerId, err := resourceParseId(d.Id())
				if err != nil {
					return nil, err
				}
				_ = d.Set("listener_id", listenerId)
				d.SetId(ruleId)
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of listener.",
			},
			"domain": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"domain", "url"},
				Description:  "The domain of Rule.",
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"domain", "url"},
				// 若指定Domain，则Url不传入数值时，默认为“/”
				Default:     "/",
				Description: "The Url of Rule.",
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Server Group Id.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Rule.",
			},
		},
	}
}

func resourceByteplusRuleCreate(d *schema.ResourceData, meta interface{}) (err error) {
	ruleService := NewRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(ruleService, d, ResourceByteplusRule())
	if err != nil {
		return fmt.Errorf("error on creating rule %q, %w", d.Id(), err)
	}
	return resourceByteplusRuleRead(d, meta)
}

func resourceByteplusRuleRead(d *schema.ResourceData, meta interface{}) (err error) {
	ruleService := NewRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(ruleService, d, ResourceByteplusRule())
	if err != nil {
		return fmt.Errorf("error on reading rule %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusRuleUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	ruleService := NewRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(ruleService, d, ResourceByteplusRule())
	if err != nil {
		return fmt.Errorf("error on updating rule %q, %w", d.Id(), err)
	}
	return resourceByteplusRuleRead(d, meta)
}

func resourceByteplusRuleDelete(d *schema.ResourceData, meta interface{}) (err error) {
	ruleService := NewRuleService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(ruleService, d, ResourceByteplusRule())
	if err != nil {
		return fmt.Errorf("error on deleting rule %q, %w", d.Id(), err)
	}
	return err
}
