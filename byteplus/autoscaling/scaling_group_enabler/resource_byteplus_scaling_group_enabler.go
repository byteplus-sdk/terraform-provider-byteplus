package scaling_group_enabler

import (
	"fmt"
	"strings"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Scaling Group enabler can be imported using the scaling_group_id, e.g.
```
$ terraform import byteplus_scaling_group_enabler.default enable:scg-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceByteplusScalingGroupEnabler() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusScalingGroupEnablerCreate,
		Read:   resourceByteplusScalingGroupEnablerRead,
		Delete: resourceByteplusScalingGroupEnablerDelete,
		Importer: &schema.ResourceImporter{
			State: scalingGroupEnablerImporter,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the scaling group.",
			},
		},
	}
	return resource
}

func resourceByteplusScalingGroupEnablerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingGroupEnablerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusScalingGroupEnabler())
	if err != nil {
		return fmt.Errorf("error on creating ScalingGroupEnable: %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingGroupEnablerRead(d, meta)
}

func resourceByteplusScalingGroupEnablerRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingGroupEnablerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusScalingGroupEnabler())
	if err != nil {
		return fmt.Errorf("error on reading ScalingGroupEnable: %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusScalingGroupEnablerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewScalingGroupEnablerService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusScalingGroupEnabler())
	if err != nil {
		return fmt.Errorf("erron on deleting ScalingGroupEnabler: %q, %s", d.Id(), err)
	}
	return err
}

func scalingGroupEnablerImporter(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("scaling_group_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
