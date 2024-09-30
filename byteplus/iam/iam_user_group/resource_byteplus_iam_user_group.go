package iam_user_group

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
IamUserGroup can be imported using the id, e.g.
```
$ terraform import byteplus_iam_user_group.default user_group_name
```

*/

func ResourceByteplusIamUserGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIamUserGroupCreate,
		Read:   resourceByteplusIamUserGroupRead,
		Update: resourceByteplusIamUserGroupUpdate,
		Delete: resourceByteplusIamUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the user group.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The display name of the user group.",
			},
		},
	}
	return resource
}

func resourceByteplusIamUserGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on creating iam_user_group %q, %s", d.Id(), err)
	}
	return resourceByteplusIamUserGroupRead(d, meta)
}

func resourceByteplusIamUserGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on reading iam_user_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusIamUserGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on updating iam_user_group %q, %s", d.Id(), err)
	}
	return resourceByteplusIamUserGroupRead(d, meta)
}

func resourceByteplusIamUserGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusIamUserGroup())
	if err != nil {
		return fmt.Errorf("error on deleting iam_user_group %q, %s", d.Id(), err)
	}
	return err
}
