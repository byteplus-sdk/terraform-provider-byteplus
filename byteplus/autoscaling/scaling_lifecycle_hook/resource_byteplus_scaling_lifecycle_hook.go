package scaling_lifecycle_hook

import (
	"fmt"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
ScalingLifecycleHook can be imported using the ScalingGroupId:LifecycleHookId, e.g.
```
$ terraform import byteplus_scaling_lifecycle_hook.default scg-yblfbfhy7agh9zn72iaz:sgh-ybqholahe4gso0ee88sd
```

*/

func ResourceByteplusScalingLifecycleHook() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusScalingLifecycleHookCreate,
		Read:   resourceByteplusScalingLifecycleHookRead,
		Update: resourceVetackScalingLifecycleHookUpdate,
		Delete: resourceVetackScalingLifecycleHookDelete,
		Importer: &schema.ResourceImporter{
			State: lifecycleHookImporter,
		},
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the scaling group.",
			},
			"lifecycle_hook_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the lifecycle hook.",
			},
			"lifecycle_hook_timeout": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(30, 21600),
				Description:  "The timeout of the lifecycle hook.",
			},
			"lifecycle_hook_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SCALE_IN", "SCALE_OUT"}, false),
				Description:  "The type of the lifecycle hook. Valid values: SCALE_IN, SCALE_OUT.",
			},
			"lifecycle_hook_policy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy of the lifecycle hook. Valid values: CONTINUE, REJECT, ROLLBACK.",
			},
			"lifecycle_command": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Batch job command.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Batch job command ID, which indicates the batch job command to be executed after triggering the lifecycle hook and installed in the instance.",
						},
						"parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Parameters and parameter values in batch job commands.\nThe number of parameters ranges from 0 to 60.",
						},
					},
				},
			},
		},
	}
	dataSource := DataSourceByteplusScalingLifecycleHooks().Schema["lifecycle_hooks"].Elem.(*schema.Resource).Schema
	delete(dataSource, "id")
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusScalingLifecycleHookCreate(d *schema.ResourceData, meta interface{}) (err error) {
	lifecycleHookService := NewScalingLifecycleHookService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(lifecycleHookService, d, ResourceByteplusScalingLifecycleHook())
	if err != nil {
		return fmt.Errorf("error on creating ScalingLifecycleHook %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingLifecycleHookRead(d, meta)
}

func resourceByteplusScalingLifecycleHookRead(d *schema.ResourceData, meta interface{}) (err error) {
	lifecycleHookService := NewScalingLifecycleHookService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(lifecycleHookService, d, ResourceByteplusScalingLifecycleHook())
	if err != nil {
		return fmt.Errorf("error on reading ScalingLifecycleHook %q, %s", d.Id(), err)
	}
	return err
}

func resourceVetackScalingLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	lifecycleHookService := NewScalingLifecycleHookService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(lifecycleHookService, d, ResourceByteplusScalingLifecycleHook())
	if err != nil {
		return fmt.Errorf("error on updating ScalingLifecycleHook %q, %s", d.Id(), err)
	}
	return resourceByteplusScalingLifecycleHookRead(d, meta)
}

func resourceVetackScalingLifecycleHookDelete(d *schema.ResourceData, meta interface{}) (err error) {
	lifecycleHookService := NewScalingLifecycleHookService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(lifecycleHookService, d, ResourceByteplusScalingLifecycleHook())
	if err != nil {
		return fmt.Errorf("error on deleting ScalingLifecycleHook %q, %s", d.Id(), err)
	}
	return err
}
