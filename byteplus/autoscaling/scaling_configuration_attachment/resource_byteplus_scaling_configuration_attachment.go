package scaling_configuration_attachment

import (
	"fmt"
	"strings"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Scaling Configuration attachment can be imported using the scaling_configuration_id e.g.
The launch template and scaling configuration cannot take effect at the same time.
```
$ terraform import byteplus_scaling_configuration_attachment.default enable:scc-ybrurj4uw6gh9zecj327
```

*/

func ResourceByteplusScalingConfigurationAttachment() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusScalingConfigurationAttachmentCreate,
		Read:   resourceByteplusScalingConfigurationAttachmentRead,
		Delete: resourceByteplusScalingConfigurationAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: importScalingConfigurationAttachment,
		},
		Schema: map[string]*schema.Schema{
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the scaling configuration.",
			},
		},
	}
	return resource
}

func resourceByteplusScalingConfigurationAttachmentCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingConfigurationAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusScalingConfigurationAttachment())
	if err != nil {
		return fmt.Errorf("error on creating ScalingConfigurationEnable: %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingConfigurationAttachmentRead(d, meta)
}

func resourceByteplusScalingConfigurationAttachmentRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingConfigurationAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusScalingConfigurationAttachment())
	if err != nil {
		return fmt.Errorf("error on reading ScalingConfigurationEnable: %q, %s", d.Id(), err)
	}
	return nil
}

func resourceByteplusScalingConfigurationAttachmentDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingConfigurationAttachmentService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusScalingConfigurationAttachment())
	if err != nil {
		return fmt.Errorf("error on deleting ScalingConfigurationEnable: %q, %s", d.Id(), err)
	}
	return nil
}

func importScalingConfigurationAttachment(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form enable:ScalingConfigurationId")
	}
	err = data.Set("scaling_configuration_id", items[1])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
