package iam_role_policy_attachment

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Iam role policy attachment can be imported using the id, e.g.
```
$ terraform import byteplus_iam_role_policy_attachment.default TerraformTestRole:TerraformTestPolicy:Custom
```

*/

func ResourceByteplusIamRolePolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusIamRolePolicyAttachmentCreate,
		Read:   resourceByteplusIamRolePolicyAttachmentRead,
		Delete: resourceByteplusIamRolePolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: iamRolePolicyAttachmentImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"role_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Role.",
			},
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the Policy.",
			},
			"policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"System", "Custom"}, false),
				Description:  "The type of the Policy.",
			},
		},
	}
}

func resourceByteplusIamRolePolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	iamRolePolicyAttachmentService := NewIamRolePolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamRolePolicyAttachmentService, d, ResourceByteplusIamRolePolicyAttachment()); err != nil {
		return fmt.Errorf("error on creating iam role policy attachment %q, %w", d.Id(), err)
	}
	return resourceByteplusIamRolePolicyAttachmentRead(d, meta)
}

func resourceByteplusIamRolePolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	iamRolePolicyAttachmentService := NewIamRolePolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamRolePolicyAttachmentService, d, ResourceByteplusIamRolePolicyAttachment()); err != nil {
		return fmt.Errorf("error on reading iam role policy attachment %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusIamRolePolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	iamRolePolicyAttachmentService := NewIamRolePolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamRolePolicyAttachmentService, d, ResourceByteplusIamRolePolicyAttachment()); err != nil {
		return fmt.Errorf("error on deleting iam role policy attachment %q, %w", d.Id(), err)
	}
	return nil
}
