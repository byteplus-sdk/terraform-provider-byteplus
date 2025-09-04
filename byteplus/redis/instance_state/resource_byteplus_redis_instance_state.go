package instance_state

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Redis State Instance can be imported using the id, e.g.
```
$ terraform import byteplus_redis_instance_state.default state:redis-mizl7m1kqccg5smt1bdpijuj
```

*/

func ResourceByteplusRedisInstanceState() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusRedisInstanceStateDelete,
		Create: resourceByteplusRedisInstanceStateCreate,
		Read:   resourceByteplusRedisInstanceStateRead,
		Importer: &schema.ResourceImporter{
			State: instanceStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Hour),
			Delete: schema.DefaultTimeout(1 * time.Hour),
		},
		Schema: map[string]*schema.Schema{
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Restart"}, false),
				Description:  "Instance Action, the value can be `Restart`.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of Instance.",
			},
		},
	}
}

func resourceByteplusRedisInstanceStateCreate(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewRedisInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(instanceStateService, d, ResourceByteplusRedisInstanceState()); err != nil {
		return fmt.Errorf("error on creating instance state %q, %w", d.Id(), err)
	}
	return resourceByteplusRedisInstanceStateRead(d, meta)
}

func resourceByteplusRedisInstanceStateRead(d *schema.ResourceData, meta interface{}) error {
	instanceStateService := NewRedisInstanceStateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(instanceStateService, d, ResourceByteplusRedisInstanceState()); err != nil {
		return fmt.Errorf("error on reading instance state %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusRedisInstanceStateDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

var instanceStateImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("instance_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
