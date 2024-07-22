package cdn_service_template

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnServiceTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnServiceTemplatesRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fuzzy": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							Description: "Indicates the matching method. This parameter can take the following values: " +
								"true: Indicates fuzzy matching. A policy is considered to meet the filtering criteria if the corresponding value of Name contains any value in the Value array. " +
								"false: Indicates exact matching. A policy is considered to meet the filtering criteria if the corresponding value of Name matches any value in the Value array. " +
								"Moreover, the Fuzzy value you can specify is affected by the Name value. See the description of Name. " +
								"The default value of this parameter is false. " +
								"Note that the matching process is case-sensitive.",
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Represents the filtering type. This parameter can take the following values: " +
								"Title: Filters policies by name. " +
								"Id: Filters policies by ID. For this parameter value, the value of Fuzzy can only be false. " +
								"Domain: Filters policies by the bound domain name. " +
								"Type: Filters policies by type. For this parameter value, the value of Fuzzy can only be false. " +
								"Status: Filters policies by status. For this parameter value, the value of Fuzzy can only be false. " +
								"You can specify multiple filtering criteria simultaneously, but the Name in different filtering criteria cannot be the same.",
						},
						"value": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "Represents the values corresponding to Name, which is an array. " +
								"When Name is Title, Id, or Domain, each value in the Value array should not exceed 100 characters in length. " +
								"When Name is Type, the Value array can include one or more of the following values: " +
								"cipher: Indicates a encryption policy. " +
								"service: Indicates a delivery policy. " +
								"When Name is Status, the Value array can include one or more of the following values: " +
								"locked: Indicates the status is \"published\". " +
								"editing: Indicates the status is \"draft\". " +
								"When Fuzzy is false, you can specify multiple values in the array. " +
								"When Fuzzy is true, you can only specify one value in the array.",
						},
					},
				},
				Description: "Indicates a set of filtering criteria used to obtain a list of policies that meet these criteria. " +
					"If you do not specify any filtering criteria, this API returns all policies under your account. " +
					"Multiple filtering criteria are related by AND, meaning only policies that meet all filtering criteria will be included in the list returned by this API. " +
					"In the API response, the actual policies returned are affected by PageNum and PageSize.",
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
			"templates": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bound_domains": {
							Type:     schema.TypeList,
							Computed: true,
							Description: "Represents a list of domain names bound to the policy specified by TemplateId. " +
								"If the policy is not bound to any domain names, the value of this parameter is null.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bound_time": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Indicates the time when the policy was bound to the domain name specified by Domain, in Unix timestamp format.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Represents one of the domain names bound to the policy.",
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the creation time of the policy, in Unix timestamp format.",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the description of the policy.",
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "Indicates the status of the policy. This parameter can take the following values: " +
								"locked: Indicates the status is \"published\". " +
								"editing: Indicates the status is \"draft\".",
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the ID of a policy in the list of policies returned by the API.",
						},
						"title": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the name of the policy.",
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "Indicates the type of the policy. This parameter can take the following values: " +
								"cipher: Indicates an encryption policy. " +
								"service: Indicates a distribution policy.",
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
							Description: "Indicates the last modification time of the policy, in Unix timestamp format. " +
								"If the policy has not been updated since its creation, the value of this parameter is the same as CreateTime.",
						},
						"exception": {
							Type:     schema.TypeBool,
							Computed: true,
							Description: "Indicates whether the policy includes special configurations. " +
								"Special configurations refer to those not operated by users but by BytePlus engineers. " +
								"This parameter can take the following values:" +
								" true: Indicates it includes special configurations. " +
								"false: Indicates it does not include special configurations.",
						},
						"project": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the project to which the policy belongs.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCdnServiceTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnServiceTemplateService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnServiceTemplates())
}
