package organization_service_control_policy

import (
	"encoding/json"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Service Control Policy can be imported using the id, e.g.
```
$ terraform import byteplus_organization_service_control_policy.default 1000001
```

*/

func ResourceByteplusServiceControlPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusServiceControlPolicyCreate,
		Read:   resourceByteplusServiceControlPolicyRead,
		Update: resourceByteplusServiceControlPolicyUpdate,
		Delete: resourceByteplusServiceControlPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Policy.",
			},
			"statement": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The statement of the Policy.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					oldMap := make(map[string]interface{})
					newMap := make(map[string]interface{})

					_ = json.Unmarshal([]byte(old), &oldMap)
					_ = json.Unmarshal([]byte(new), &newMap)

					oldStr, _ := json.MarshalIndent(oldMap, "", "\t")
					newStr, _ := json.MarshalIndent(newMap, "", "\t")
					return string(oldStr) == string(newStr)
				},
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Policy.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the Policy.",
			},
			"create_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of the Policy.",
			},
			"update_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the Policy.",
			},
		},
	}
}

func resourceByteplusServiceControlPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamPolicyService, d, ResourceByteplusServiceControlPolicy()); err != nil {
		return fmt.Errorf("error on creating policy %q, %w", d.Id(), err)
	}
	return resourceByteplusServiceControlPolicyRead(d, meta)
}

func resourceByteplusServiceControlPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamPolicyService, d, ResourceByteplusServiceControlPolicy()); err != nil {
		return fmt.Errorf("error on reading policy %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusServiceControlPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(iamPolicyService, d, ResourceByteplusServiceControlPolicy()); err != nil {
		return fmt.Errorf("error on updating policy %q, %w", d.Id(), err)
	}
	return resourceByteplusServiceControlPolicyRead(d, meta)
}

func resourceByteplusServiceControlPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamPolicyService, d, ResourceByteplusServiceControlPolicy()); err != nil {
		return fmt.Errorf("error on deleting policy %q, %w", d.Id(), err)
	}
	return nil
}
