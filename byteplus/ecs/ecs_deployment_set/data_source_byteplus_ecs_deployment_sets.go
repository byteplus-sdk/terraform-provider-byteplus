package ecs_deployment_set

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusEcsDeploymentSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusEcsDeploymentSetsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set:         schema.HashString,
				Description: "A list of ECS DeploymentSet IDs.",
			},
			"granularity": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"switch",
					"host",
					"rack",
				}, false),
				Description: "The granularity of ECS DeploymentSet.Valid values: switch, host, rack.",
			},

			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of ECS DeploymentSet.",
			},

			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of ECS DeploymentSet query.",
			},
			"deployment_sets": {
				Description: "The collection of ECS DeploymentSet query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_set_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of ECS DeploymentSet.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of ECS DeploymentSet.",
						},
						"granularity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The granularity of ECS DeploymentSet.",
						},
						"strategy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The strategy of ECS DeploymentSet.",
						},
						"deployment_set_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of ECS DeploymentSet.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusEcsDeploymentSetsRead(d *schema.ResourceData, meta interface{}) error {
	deploymentSetService := NewEcsDeploymentSetService(meta.(*bp.SdkClient))
	return bp.DefaultDispatcher().Data(deploymentSetService, d, DataSourceByteplusEcsDeploymentSets())
}
