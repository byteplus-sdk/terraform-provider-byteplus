package iam_user_policy_attachment

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Iam user policy attachment can be imported using the UserName:PolicyName:PolicyType, e.g.
```
$ terraform import byteplus_iam_user_policy_attachment.default TerraformTestUser:TerraformTestPolicy:Custom
```

*/

func ResourceByteplusIamUserPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusIamUserPolicyAttachmentCreate,
		Read:   resourceByteplusIamUserPolicyAttachmentRead,
		Delete: resourceByteplusIamUserPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: iamUserPolicyAttachmentImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user.",
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

func resourceByteplusIamUserPolicyAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(iamUserPolicyAttachmentService, d, ResourceByteplusIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on creating iam user policy attachment %q, %w", d.Id(), err)
	}
	return resourceByteplusIamUserPolicyAttachmentRead(d, meta)
}

func resourceByteplusIamUserPolicyAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(iamUserPolicyAttachmentService, d, ResourceByteplusIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on reading iam user policy attachment %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusIamUserPolicyAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	iamUserPolicyAttachmentService := NewIamUserPolicyAttachmentService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(iamUserPolicyAttachmentService, d, ResourceByteplusIamUserPolicyAttachment()); err != nil {
		return fmt.Errorf("error on deleting iam user policy attachment %q, %w", d.Id(), err)
	}
	return nil
}
