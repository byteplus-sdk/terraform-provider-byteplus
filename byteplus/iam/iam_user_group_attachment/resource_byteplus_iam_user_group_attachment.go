package iam_user_group_attachment

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
IamUserGroupAttachment can be imported using the id, e.g.
```
$ terraform import byteplus_iam_user_group_attachment.default user_group_id:user_id
```

*/

func ResourceByteplusIamUserGroupAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusIamUserGroupAttachmentCreate,
		Read:   resourceByteplusIamUserGroupAttachmentRead,
		Delete: resourceByteplusIamUserGroupAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importIamUserGroupAttachment,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"user_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user group.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the user.",
			},
		},
	}
	return resource
}

func resourceByteplusIamUserGroupAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on creating iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return resourceByteplusIamUserGroupAttachmentRead(d, meta)
}

func resourceByteplusIamUserGroupAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on reading iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusIamUserGroupAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewIamUserGroupAttachmentService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusIamUserGroupAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting iam_user_group_attachment %q, %s", d.Id(), err)
	}
	return err
}

func importIamUserGroupAttachment(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form user_group_id:user_id")
	}
	err = data.Set("user_group_name", items[0])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	err = data.Set("user_name", items[1])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
