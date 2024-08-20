package cdn_cron_job_state

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnCronJobState can be imported using the state:function_id:job_name, e.g.
```
$ terraform import byteplus_cdn_cron_job_state.default state:function_id:job_name
```

*/

func ResourceByteplusCdnCronJobState() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnCronJobStateCreate,
		Read:   resourceByteplusCdnCronJobStateRead,
		Update: resourceByteplusCdnCronJobStateUpdate,
		Delete: resourceByteplusCdnCronJobStateDelete,
		Importer: &schema.ResourceImporter{
			State: cronJobStateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the function.",
			},
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the cron job.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Start", "Stop"}, false),
				Description: "Start or Stop of corn job, the value can be `Start` or `Stop`. \n" +
					"If the target status of the action is consistent with the current status of the corn job, the action will not actually be executed.\n" +
					"When importing resources, this attribute will not be imported. If this attribute is set, please use lifecycle and ignore_changes ignore changes in fields.",
			},

			// computed_fields
			"job_state": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the cron job.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnCronJobStateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobStateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnCronJobState())
	if err != nil {
		return fmt.Errorf("error on creating cdn_cron_job_state %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCronJobStateRead(d, meta)
}

func resourceByteplusCdnCronJobStateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobStateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnCronJobState())
	if err != nil {
		return fmt.Errorf("error on reading cdn_cron_job_state %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnCronJobStateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobStateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnCronJobState())
	if err != nil {
		return fmt.Errorf("error on updating cdn_cron_job_state %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCronJobStateRead(d, meta)
}

func resourceByteplusCdnCronJobStateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobStateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnCronJobState())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_cron_job_state %q, %s", d.Id(), err)
	}
	return err
}

var cronJobStateImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 3 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("function_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("job_name", items[2]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
