package iam_login_profile

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Login profile can be imported using the UserName, e.g.
```
$ terraform import byteplus_iam_login_profile.default user_name
```

*/

func ResourceByteplusIamLoginProfile() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIamLoginProfileCreate,
		Read:   resourceByteplusIamLoginProfileRead,
		Update: resourceByteplusIamLoginProfileUpdate,
		Delete: resourceByteplusIamLoginProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user name.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The password.",
			},
			"login_allowed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The flag of login allowed.",
			},
			"password_reset_required": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is required reset password when next time login in.",
			},
		},
	}
	return resource
}

func resourceByteplusIamLoginProfileCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on creating login profile %q, %s", d.Id(), err)
	}
	return resourceByteplusIamLoginProfileRead(d, meta)
}

func resourceByteplusIamLoginProfileRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on reading login profile %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusIamLoginProfileUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on updating login profile %q, %s", d.Id(), err)
	}
	return resourceByteplusIamLoginProfileRead(d, meta)
}

func resourceByteplusIamLoginProfileDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamLoginProfileService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusIamLoginProfile())
	if err != nil {
		return fmt.Errorf("error on deleting login profile %q, %s", d.Id(), err)
	}
	return err
}
