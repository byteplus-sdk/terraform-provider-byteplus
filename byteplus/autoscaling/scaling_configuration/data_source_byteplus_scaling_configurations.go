package scaling_configuration

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusScalingConfigurations() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusScalingConfigurationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of scaling configuration ids.",
			},
			"scaling_configuration_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of scaling configuration names.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An id of scaling group.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of scaling configuration.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of scaling configuration query.",
			},
			"scaling_configurations": {
				Description: "The collection of scaling configuration query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the scaling configuration.",
						},
						"scaling_configuration_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the scaling configuration.",
						},
						"scaling_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the scaling group to which the scaling configuration belongs.",
						},
						"scaling_configuration_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the scaling configuration.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the scaling configuration.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the scaling configuration.",
						},
						"volumes": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of volume of the scaling configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"volume_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of volume.",
									},
									"size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The size of volume.",
									},
									"delete_with_instance": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The delete with instance flag of volume.",
									},
								},
							},
						},
						"host_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS hostname which the scaling configuration set.",
						},
						"image_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS image id which the scaling configuration set.",
						},
						"instance_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS instance description which the scaling configuration set.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS instance name which the scaling configuration set.",
						},
						"instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of the ECS instance type which the scaling configuration set.",
						},
						"key_pair_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS key pair name which the scaling configuration set.",
						},
						"lifecycle_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the scaling configuration.",
						},
						"security_enhancement_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Ecs security enhancement strategy which the scaling configuration set.",
						},
						"security_group_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "The list of the security group id of the networkInterface which the scaling configuration set.",
						},
						"eip_bandwidth": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The EIP bandwidth which the scaling configuration set.",
						},
						"eip_isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP ISP which the scaling configuration set.",
						},
						"eip_billing_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The EIP ISP which the scaling configuration set.",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ECS user data which the scaling configuration set.",
						},
						"tags": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "The label of the instance created by the scaling configuration.",
							Set:         bp.TagsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Key of Tags.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The Value of Tags.",
									},
								},
							},
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project to which the instance created by the scaling configuration belongs.",
						},
						"hpc_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the HPC cluster to which the instance belongs. Valid only when InstanceTypes.N specifies High Performance Computing GPU Type.",
						},
						"spot_strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The preemption policy of the instance. Valid Value: NoSpot (default), SpotAsPriceGo.",
						},
						"ipv6_address_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Assign IPv6 address to instance network card. Possible values:\n0: Do not assign IPv6 address.\n1: Assign IPv6 address and the system will automatically assign an IPv6 subnet for you.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusScalingConfigurationsRead(d *schema.ResourceData, meta interface{}) error {
	configurationService := NewScalingConfigurationService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(configurationService, d, DataSourceByteplusScalingConfigurations())
}
