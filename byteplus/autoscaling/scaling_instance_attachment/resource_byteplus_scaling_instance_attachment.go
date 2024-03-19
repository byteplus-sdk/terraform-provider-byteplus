package scaling_instance_attachment

import (
	"fmt"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Scaling instance attachment can be imported using the scaling_group_id and instance_id, e.g.
You can choose to remove or detach the instance according to the `delete_type` field.
```
$ terraform import byteplus_scaling_instance_attachment.default scg-mizl7m1kqccg5smt1bdpijuj:i-l8u2ai4j0fauo6mrpgk8
```

*/

func ResourceByteplusScalingInstanceAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusScalingInstanceAttachmentCreate,
		Read:   resourceByteplusScalingInstanceAttachmentRead,
		Update: resourceByteplusScalingInstanceAttachmentUpdate,
		Delete: resourceByteplusScalingInstanceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importScalingInstanceAttachment,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the scaling group.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the instance.",
			},
			"entrusted": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to host the instance to a scaling group. Default value is false.",
			},
			"delete_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Remove",
					"Detach",
				}, false),
				Description: "The type of delete activity. Valid values: Remove, Detach. Default value is Remove.",
			},
			"detach_option": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"both",
					"none",
				}, false),
				Description: "Whether to cancel the association of the instance with the load balancing and public network IP. Valid values: both, none. Default value is both.",
			},
		},
	}
	return resource
}

func resourceByteplusScalingInstanceAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	scalingInstanceAttachService := NewScalingInstanceAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(scalingInstanceAttachService, d, ResourceByteplusScalingInstanceAttachment())
	if err != nil {
		return fmt.Errorf("error on creating ScalingInstanceAttach %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingInstanceAttachmentRead(d, meta)
}

func resourceByteplusScalingInstanceAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	scalingInstanceAttachService := NewScalingInstanceAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(scalingInstanceAttachService, d, ResourceByteplusScalingInstanceAttachment())
	if err != nil {
		return fmt.Errorf("error on reading ScalingInstanceAttach %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusScalingInstanceAttachmentUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	scalingInstanceAttachService := NewScalingInstanceAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(scalingInstanceAttachService, d, ResourceByteplusScalingInstanceAttachment())
	if err != nil {
		return fmt.Errorf("error on updating ScalingInstanceAttach %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingInstanceAttachmentRead(d, meta)
}

func resourceByteplusScalingInstanceAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	scalingInstanceAttachService := NewScalingInstanceAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(scalingInstanceAttachService, d, ResourceByteplusScalingInstanceAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting ScalingInstanceAttach %q, %s", d.Id(), err)
	}
	return err
}
