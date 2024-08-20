package cdn_cron_job

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnCronJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnCronJobsRead,
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the function.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Resource.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of query.",
			},

			"cron_jobs": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the cron job.",
						},
						"cron_expression": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron expression of the cron job.",
						},
						"parameter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameter of the cron job.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the cron job.",
						},
						"job_state": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the cron job.",
						},
						"cron_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The type of the cron job.",
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
				},
			},
		},
	}
}

func dataSourceByteplusCdnCronJobsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnCronJobService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnCronJobs())
}
