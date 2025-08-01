package cr_vpc_endpoint

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCrVpcEndpoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCrVpcEndpointsRead,
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CR registry name.",
			},
			"statuses": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"Enabling",
						"Enabled",
						"Disabling",
						"Failed",
					}, false),
				},
				Set:         schema.HashString,
				Description: "VPC access entry state array, used to filter out VPC access entries in the specified state. Available values are Enabling, Enabled, Disabling, Failed.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of CR vpc endpoints query.",
			},
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of CR vpc endpoints.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"registry": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of CR registry.",
						},
						"vpcs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of vpc information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the vpc.",
									},
									"subnet_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The ID of the subnet.",
									},
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The region id.",
									},
									"account_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The id of the account.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address of the mirror repository in the VPC.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The status of the vpc endpoint.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The creation time.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCrVpcEndpointsRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCrVpcEndpointService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(service, d, DataSourceByteplusCrVpcEndpoints())
}
