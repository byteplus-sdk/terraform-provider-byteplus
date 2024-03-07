package ecs_instance_state

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
State Instance can be imported using the id, e.g.
```
$ terraform import byteplus_ecs_instance_state.default state:i-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceByteplusEcsInstanceState() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusEcsInstanceStateDelete,
		Create: resourceByteplusEcsInstanceStateCreate,
		Read:   resourceByteplusEcsInstanceStateRead,
		Update: resourceByteplusEcsInstanceStateUpdate,
		Importer: &schema.ResourceImporter{
			State: ecsInstanceStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Update: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Start", "Stop", "ForceStop"}, false),
				Description:  "Start or Stop of Instance Action, the value can be `Start`, `Stop` or `ForceStop`.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of Instance.",
			},
			"stopped_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 如开机行为，该字段修改忽略
					return d.Get("action").(string) == "Start"
				},
				Description: "Stop Mode of Instance, the value can be `KeepCharging` or `StopCharging`.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Instance.",
			},
		},
	}
}

func resourceByteplusEcsInstanceStateCreate(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(instanceStateService, d, ResourceByteplusEcsInstanceState()); err != nil {
		return fmt.Errorf("error on creating instance state %q, %w", d.Id(), err)
	}
	return resourceByteplusEcsInstanceStateRead(d, meta)
}

func resourceByteplusEcsInstanceStateRead(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(instanceStateService, d, ResourceByteplusEcsInstanceState()); err != nil {
		return fmt.Errorf("error on reading instance state %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusEcsInstanceStateUpdate(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(instanceStateService, d, ResourceByteplusEcsInstanceState()); err != nil {
		return fmt.Errorf("error on updating instance state %q, %w", d.Id(), err)
	}
	return resourceByteplusEcsInstanceStateRead(d, meta)
}

func resourceByteplusEcsInstanceStateDelete(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(instanceStateService, d, ResourceByteplusEcsInstanceState()); err != nil {
		return fmt.Errorf("error on deleting instance state %q, %w", d.Id(), err)
	}
	return nil
}
