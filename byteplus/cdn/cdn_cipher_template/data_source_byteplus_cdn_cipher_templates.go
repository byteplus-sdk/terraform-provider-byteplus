package cdn_cipher_template

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnCipherTemplates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnCipherTemplatesRead,
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

						"https": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicates the configuration module for the HTTPS encryption service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"disable_http": {
										Type:     schema.TypeBool,
										Computed: true,
										Description: "Indicates whether the CDN accepts HTTP user requests. " +
											"This parameter can take the following values: " +
											"true: Indicates that it does not accept. If an HTTP request is received, the CDN will reject the request. " +
											"false: Indicates that it accepts. The default value for this parameter is false.",
									},
									"forced_redirect": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Indicates the configuration for the mandatory redirection from HTTP to HTTPS. This feature is disabled by default.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_forced_redirect": {
													Type:     schema.TypeBool,
													Computed: true,
													Description: "Indicates the switch for the Forced Redirect configuration. " +
														"This parameter can take the following values: " +
														"true: Indicates to enable Forced Redirect. " +
														"false: Indicates to disable Forced Redirect.",
												},
												"status_code": {
													Type:     schema.TypeString,
													Computed: true,
													Description: "Indicates the status code returned to the client by the CDN when forced redirect occurs. " +
														"This parameter can take the following values: " +
														"301: Indicates that the returned status code is 301. " +
														"302: Indicates that the returned status code is 302. " +
														"The default value for this parameter is 301.",
												},
											},
										},
									},
									"http2": {
										Type:     schema.TypeBool,
										Computed: true,
										Description: "Indicates the switch for HTTP/2 configuration. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable HTTP/2. " +
											"false: Indicates to disable HTTP/2. " +
											"The default value for this parameter is true.",
									},
									"ocsp": {
										Type:     schema.TypeBool,
										Computed: true,
										Description: "Indicates whether to enable OCSP Stapling. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable OCSP Stapling. " +
											"false: Indicates to disable OCSP Stapling. " +
											"The default value for this parameter is false.",
									},
									"tls_version": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Indicates a list that specifies the TLS versions supported by the domain name. " +
											"This parameter can take the following values: " +
											"tlsv1.0: Indicates TLS 1.0. " +
											"tlsv1.1: Indicates TLS 1.1. " +
											"tlsv1.2: Indicates TLS 1.2. " +
											"tlsv1.3: Indicates TLS 1.3. " +
											"The default value for this parameter is [\"tlsv1.1\", \"tlsv1.2\", \"tlsv1.3\"].",
									},
									"hsts": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Indicates the HSTS (HTTP Strict Transport Security) configuration module. This feature is disabled by default.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"subdomain": {
													Type:     schema.TypeString,
													Computed: true,
													Description: "Indicates whether the HSTS configuration should also be applied to the subdomains of the domain name. " +
														"This parameter can take the following values: " +
														"include: Indicates that HSTS settings apply to subdomains. " +
														"exclude: Indicates that HSTS settings do not apply to subdomains. " +
														"The default value for this parameter is exclude.",
												},
												"switch": {
													Type:     schema.TypeBool,
													Computed: true,
													Description: "Indicates whether to enable HSTS. " +
														"This parameter can take the following values: " +
														"true: Indicates to enable HSTS. " +
														"false: Indicates to disable HSTS. " +
														"The default value for this parameter is false.",
												},
												"ttl": {
													Type:     schema.TypeInt,
													Computed: true,
													Description: "Indicates the expiration time for the Strict-Transport-Security response header in the browser cache, in seconds. " +
														"If Switch is true, this parameter is required. " +
														"The value range for this parameter is 0 - 31,536,000 seconds, where 31,536,000 seconds represents 365 days. " +
														"If the value of this parameter is 0, it is equivalent to disabling the HSTS settings.",
												},
											},
										},
									},
								},
							},
						},

						"http_forced_redirect": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicates the configuration module for the forced redirection from HTTPS to HTTP. This feature is disabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_forced_redirect": {
										Type:     schema.TypeBool,
										Computed: true,
										Description: "Indicates whether to enable the forced redirection from HTTPS. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable the forced redirection from HTTPS. " +
											"Once enabled, the content delivery network will respond with StatusCode to inform the browser to send an HTTPS request when it receives an HTTP request from a user. " +
											"false: Indicates to disable the forced redirection from HTTPS.",
									},
									"status_code": {
										Type:     schema.TypeString,
										Computed: true,
										Description: "Indicates the status code returned by the content delivery network when forced redirection from HTTPS occurs. " +
											"The default value for this parameter is 301.",
									},
								},
							},
						},
						"quic": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Indicates the QUIC configuration module. This feature is disabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"switch": {
										Computed: true,
										Type:     schema.TypeBool,
										Description: "Indicates whether to enable QUIC. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable QUIC. " +
											"false: Indicates to disable QUIC.",
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

func dataSourceByteplusCdnCipherTemplatesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnCipherTemplateService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnCipherTemplates())
}
