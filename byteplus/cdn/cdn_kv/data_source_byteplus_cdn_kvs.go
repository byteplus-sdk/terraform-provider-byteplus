package cdn_kv

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnKvs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnKvsRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the kv namespace.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the kv namespace.",
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

			"namespace_keys": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the kv namespace key.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the kv namespace key.",
						},
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The key of the kv namespace key.",
						},
						"key_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The status of the kv namespace key.",
						},
						"ddl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data expiration time. After the data expires, the Value in the Key will be inaccessible.\nDisplayed in UNIX timestamp format.\n0: Permanent storage.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the kv namespace key.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The creation time of the kv namespace key. Displayed in UNIX timestamp format.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The update time of the kv namespace key. Displayed in UNIX timestamp format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCdnKvsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnKvService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnKvs())
}
