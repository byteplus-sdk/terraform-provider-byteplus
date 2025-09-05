package waf_cdn_domain

import (
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
WafCdnDomain can be imported using the id, e.g.
```
$ terraform import byteplus_waf_cdn_domain.default resource_id
```

*/

func ResourceByteplusWafCdnDomain() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusWafCdnDomainCreate,
		Read:   resourceByteplusWafCdnDomainRead,
		Update: resourceByteplusWafCdnDomainUpdate,
		Delete: resourceByteplusWafCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of project.",
			},
			"project_follow": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: "Set whether to follow the project to which other resources belong, such as the CDN's project. The default value is set to 0." +
					"0: Disabled" +
					"1: Enabled " +
					"If ProjectFollow is set to 1, you don't need to enter a ProjectName.",
			},
			"tls_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Whether to log requests to the protected domain name. The default value is set to 0.",
			},
			"tls_fields_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "It can not be empty when you choose to log all headers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headers_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							ForceNew:    true,
							Description: "The configuration of Headers. Works only on modified scenes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable": {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "Whether to log all headers.",
									},
									"excluded_key_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Set:      schema.HashString,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "To avoid excessive log storage, you can configure header fields that do not need to be logged.",
									},
									"statistical_key_list": {
										Type:     schema.TypeSet,
										Optional: true,
										Computed: true,
										ForceNew: true,
										Set:      schema.HashString,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "Set the header fields that need to be calculated, analyzed, and alerted within the logged headers.",
									},
								},
							},
						},
					},
				},
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CDN domain name that need to be protected by WAF.",
			},
			// 防护网站开关相关参数
			"bot_repeat_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the bot frequency limit policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"bot_dytoken_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the bot dynamic token. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"auto_cc_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the intelligent CC protection strategy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"bot_sequence_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the bot behavior map. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"bot_sequence_default_action": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Set the default actions of the bot behavior map strategy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"bot_frequency_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the bot frequency limit policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"waf_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the vulnerability protection strategy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"cc_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the CC protection policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"white_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the access list policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"black_ip_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the access ban list policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"black_lct_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the geographical location access control policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"waf_white_req_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the whitening strategy for vulnerability protection requests. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"white_field_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the whitening strategy for vulnerability protection fields. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"custom_rsp_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the custom response interception policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"system_bot_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the managed Bot classification strategy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"custom_bot_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the custom Bot classification strategy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"api_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the API protection policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"tamper_proof_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable the page tamper-proof policy. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			"dlp_enable": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Whether to activate the strategy for preventing the leakage of sensitive information. Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
			},
			// 更新域名防护模式
			"extra_defence_mode_lb_instance": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "The protection mode of the exception instance. " +
					"It takes effect when the access mode is accessed through an application load balancing (ALB) instance (AccessMode=20)." +
					" Works only on modified scenes.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"defence_mode": {
							Type:     schema.TypeInt,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// 创建时不存在这个参数，修改时存在这个参数
								return d.Id() == ""
							},
							Description: "Set the protection mode for exceptional ALB instances. Works only on modified scenes.",
						},
						"instance_id": {
							Type:     schema.TypeString,
							Optional: true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								// 创建时不存在这个参数，修改时存在这个参数
								return d.Id() == ""
							},
							Description: "The Id of ALB instance. Works only on modified scenes.",
						},
					},
				},
			},
			"defence_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
				Description: "The protection mode of the instance. Works only on modified scenes.",
			},
			"defence_mode_computed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The protection mode of the instance.",
			},
			"advanced_defense_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of Advanced DDoS Protection instance. Displayed if the instance is provisioned via Advanced DDoS Protection, otherwise, it is null.",
			},
			"advanced_defense_ipv6": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Advanced Defense IPv6. Displayed if the instance is provisioned via Advanced DDoS Protection, otherwise, it is null.",
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CNAME value generated by the WAF instance.",
			},
			"certificate_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Certificate ID, displayed when the protocol type includes HTTPS.",
			},
			"certificate_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the certificate.",
			},
			"lb_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The types of load balancing algorithms.",
			},
			"access_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Access mode. If your domain is added from BytePlus CDN, this parameter is set to 6.",
			},
			"cloud_access_config": {
				Type:     schema.TypeList,
				Computed: true,
				Description: "Displayed when cloud WAF instance is provisioned through load balancing, otherwise, " +
					"it is null. If your domain is added from BytePlus CDN, the value is null.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of listener.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol type of the forwarding rule.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port number of the forwarding rule.",
						},
						"access_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol type.",
						},
					},
				},
			},
			"public_real_server": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Back-to-origin mode of CNAME provisioning. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "VPC ID, displayed when the back-to-origin method is set as Private IP address within VPC ( PublicRealServer=0)." +
					" If your domain is added from BytePlus CDN, the value is null.",
			},
			"protocol_ports": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Computed:    true,
				Description: "Back-to-origin port. If your domain is added from BytePlus CDN, the value is null.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Set:         schema.HashInt,
							Description: "Back-to-origin port numbers for the HTTP protocol.",
						},
						"https": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Set:         schema.HashInt,
							Description: "Back-to-origin port numbers for the HTTPS protocol.",
						},
					},
				},
			},
			"enable_http2": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether HTTP 2.0 is enabled. If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"protocol_follow": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether the protocol following is enabled. If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"enable_ipv6": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Whether IPv6 request is enabled. If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"backend_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Details of the origin servergroup. If your domain is added from BytePlus CDN, the value is null.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_port": {
							Type:     schema.TypeSet,
							Computed: true,
							Set:      schema.HashInt,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Description: "Port number.",
						},
						"backends": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details of the origin server group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Protocol of origin server.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP address of origin server.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port number of the origin server.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The weight of the origin server.",
									},
								},
							},
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the origin server group.",
						},
					},
				},
			},
			"proxy_config": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Whether the proxy configuration is enabled. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"client_ip_location": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Method to obtain client IP." +
					" If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"custom_header": {
				Type:     schema.TypeList,
				Computed: true,
				Description: "After setting the client IP acquisition method to a custom field, it will display. " +
					"If your domain is added from BytePlus CDN, the value is null.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ssl_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "TLS protocol version. " +
					"If your domain is added from BytePlus CDN, the value is null.",
			},
			"ssl_ciphers": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "TLS encryption cipher suite." +
					" If your domain is added from BytePlus CDN, the value is null.",
			},
			"keep_alive_time_out": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Long connection keep-alive time. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"keep_alive_request": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Long connection reuse count. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"client_max_body_size": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Maximum client request body size. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_connect_time_out": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Timeout for establishing connection between WAF and the backend server. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_read_time_out": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Timeout for WAF to read response from the backend server." +
					" If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_retry": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Retry count from WAF to the origin server. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_keep_alive_time_out": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Idle persistent connection timeout." +
					" If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_write_time_out": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The timeout for WAF to transfer requests to the backend server." +
					" If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"proxy_keep_alive": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "The number of reusable long connections for WAF to the origin server." +
					" If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"attack_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The status of the attack.",
			},
			"enable_sni": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Whether the SNI configuration is enabled. " +
					"If your domain is added from BytePlus CDN, the default value is set to 0.",
			},
			"custom_sni": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Custom SNI domain name. " +
					"If your domain is added from BytePlus CDN, the value is null.",
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Connection status.",
			},
			"server_ips": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF instance IP. If your domain is added from BytePlus CDN, the value is null.",
			},
			"protocols": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Protocols of provisioning. " +
					"If your domain is added from BytePlus CDN, the value is null.",
			},
			"src_ips": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Back-to-origin IP of WAF instance. If your domain is added from BytePlus CDN, the value is null.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time of updating.",
			},
		},
	}
	return resource
}

func resourceByteplusWafCdnDomainCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusWafCdnDomain())
	if err != nil {
		return fmt.Errorf("error on creating waf_cdn_domain %q, %s", d.Id(), err)
	}
	return resourceByteplusWafCdnDomainRead(d, meta)
}

func resourceByteplusWafCdnDomainRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusWafCdnDomain())
	if err != nil {
		return fmt.Errorf("error on reading waf_cdn_domain %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusWafCdnDomainUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusWafCdnDomain())
	if err != nil {
		return fmt.Errorf("error on updating waf_cdn_domain %q, %s", d.Id(), err)
	}
	return resourceByteplusWafCdnDomainRead(d, meta)
}

func resourceByteplusWafCdnDomainDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewWafCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusWafCdnDomain())
	if err != nil {
		return fmt.Errorf("error on deleting waf_cdn_domain %q, %s", d.Id(), err)
	}
	return err
}
