package iam_user

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Iam user can be imported using the UserName, e.g.
```
$ terraform import byteplus_iam_user.default user_name
```

*/

func ResourceByteplusIamUser() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIamUserCreate,
		Read:   resourceByteplusIamUserRead,
		Update: resourceByteplusIamUserUpdate,
		Delete: resourceByteplusIamUserDelete,
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
				Description: "The name of the user.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the user.",
			},
			"mobile_phone": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "The mobile phone of the user.",
				DiffSuppressFunc: phoneDiffSuppressFunc,
			},
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The email of the user.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the user.",
			},
		},
	}
	bp.MergeDateSourceToResource(DataSourceByteplusIamUsers().Schema["users"].Elem.(*schema.Resource).Schema, &resource.Schema)
	return resource
}

func resourceByteplusIamUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusIamUser())
	if err != nil {
		return fmt.Errorf("error on creating iam user  %q, %s", d.Id(), err)
	}
	return resourceByteplusIamUserRead(d, meta)
}

func resourceByteplusIamUserRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusIamUser())
	if err != nil {
		return fmt.Errorf("error on reading iam user %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusIamUserUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusIamUser())
	if err != nil {
		return fmt.Errorf("error on updating iam user %q, %s", d.Id(), err)
	}
	return resourceByteplusIamUserRead(d, meta)
}

func resourceByteplusIamUserDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusIamUser())
	if err != nil {
		return fmt.Errorf("error on deleting iam user %q, %s", d.Id(), err)
	}
	return err
}
