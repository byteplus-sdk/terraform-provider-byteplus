package iam_policy

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Iam policy can be imported using the id, e.g.
```
$ terraform import byteplus_iam_policy.default TerraformTestPolicy
```

*/

func ResourceByteplusIamPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusIamPolicyCreate,
		Read:   resourceByteplusIamPolicyRead,
		Update: resourceByteplusIamPolicyUpdate,
		Delete: resourceByteplusIamPolicyDelete,
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
			"policy_document": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The document of the Policy.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Policy.",
			},
			"policy_trn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource name of the Policy.",
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

func resourceByteplusIamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamPolicyService, d, ResourceByteplusIamPolicy()); err != nil {
		return fmt.Errorf("error on creating iam policy %q, %w", d.Id(), err)
	}
	return resourceByteplusIamPolicyRead(d, meta)
}

func resourceByteplusIamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamPolicyService, d, ResourceByteplusIamPolicy()); err != nil {
		return fmt.Errorf("error on reading iam policy %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusIamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(iamPolicyService, d, ResourceByteplusIamPolicy()); err != nil {
		return fmt.Errorf("error on updating iam policy %q, %w", d.Id(), err)
	}
	return resourceByteplusIamPolicyRead(d, meta)
}

func resourceByteplusIamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	iamPolicyService := NewIamPolicyService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamPolicyService, d, ResourceByteplusIamPolicy()); err != nil {
		return fmt.Errorf("error on deleting iam policy %q, %w", d.Id(), err)
	}
	return nil
}
