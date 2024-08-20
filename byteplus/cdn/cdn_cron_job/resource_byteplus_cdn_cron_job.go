package cdn_cron_job

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnCronJob can be imported using the function_id:job_name, e.g.
```
$ terraform import byteplus_cdn_cron_job.default function_id:job_name
```

*/

func ResourceByteplusCdnCronJob() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnCronJobCreate,
		Read:   resourceByteplusCdnCronJobRead,
		Update: resourceByteplusCdnCronJobUpdate,
		Delete: resourceByteplusCdnCronJobDelete,
		Importer: &schema.ResourceImporter{
			State: cronJobImporter,
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
				Description: "The name of the cron job. The name must meet the following requirements:\nEach cron job name for a function must be unique\nLength should not exceed 128 characters.",
			},
			"cron_expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The execution expression. The expression must meet the following requirements:\nSupports cron expressions (does not support second-level triggers).",
			},
			"cron_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The schedule type of the cron job. Possible values:\n1: Global schedule.\n2: Single point schedule.",
			},
			"parameter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The execution parameter of the cron job.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the cron job.",
			},

			// computed fields
			"job_state": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the cron job.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The creation time of the cron job. Displayed in UNIX timestamp format.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time of the cron job. Displayed in UNIX timestamp format.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnCronJobCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnCronJob())
	if err != nil {
		return fmt.Errorf("error on creating cdn_cron_job %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCronJobRead(d, meta)
}

func resourceByteplusCdnCronJobRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnCronJob())
	if err != nil {
		return fmt.Errorf("error on reading cdn_cron_job %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnCronJobUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnCronJob())
	if err != nil {
		return fmt.Errorf("error on updating cdn_cron_job %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCronJobRead(d, meta)
}

func resourceByteplusCdnCronJobDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCronJobService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnCronJob())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_cron_job %q, %s", d.Id(), err)
	}
	return err
}

var cronJobImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("function_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("job_name", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
