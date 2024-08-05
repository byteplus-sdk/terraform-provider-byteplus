package cdn_cipher_template

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnCipherTemplate can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_cipher_template.default resource_id
```

*/

func ResourceByteplusCdnCipherTemplate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnCipherTemplateCreate,
		Read:   resourceByteplusCdnCipherTemplateRead,
		Update: resourceByteplusCdnCipherTemplateUpdate,
		Delete: resourceByteplusCdnCipherTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the name of the encryption policy you want to create. The name must not exceed 100 characters.",
			},
			"message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates the description of the encryption policy, which must not exceed 120 characters.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates the project to which this encryption policy belongs. The default value of the parameter is default, indicating the Default project.",
			},
			"https": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Indicates the configuration module for the HTTPS encryption service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disable_http": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							Description: "Indicates whether the CDN accepts HTTP user requests. " +
								"This parameter can take the following values: " +
								"true: Indicates that it does not accept. If an HTTP request is received, the CDN will reject the request. " +
								"false: Indicates that it accepts. The default value for this parameter is false.",
						},
						"forced_redirect": {
							Type:          schema.TypeList,
							MaxItems:      1,
							Optional:      true,
							ConflictsWith: []string{"http_forced_redirect"},
							Description:   "Indicates the configuration for the mandatory redirection from HTTP to HTTPS. This feature is disabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_forced_redirect": {
										Type:     schema.TypeBool,
										Required: true,
										Description: "Indicates the switch for the Forced Redirect configuration. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable Forced Redirect. " +
											"false: Indicates to disable Forced Redirect.",
									},
									"status_code": {
										Type:     schema.TypeString,
										Required: true,
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
							Optional: true,
							Computed: true,
							Description: "Indicates the switch for HTTP/2 configuration. " +
								"This parameter can take the following values: " +
								"true: Indicates to enable HTTP/2. " +
								"false: Indicates to disable HTTP/2. " +
								"The default value for this parameter is true.",
						},
						"ocsp": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
							Description: "Indicates whether to enable OCSP Stapling. " +
								"This parameter can take the following values: " +
								"true: Indicates to enable OCSP Stapling. " +
								"false: Indicates to disable OCSP Stapling. " +
								"The default value for this parameter is false.",
						},
						"tls_version": {
							Type:     schema.TypeSet,
							Optional: true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Indicates the HSTS (HTTP Strict Transport Security) configuration module. This feature is disabled by default.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subdomain": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Indicates whether the HSTS configuration should also be applied to the subdomains of the domain name. " +
											"This parameter can take the following values: " +
											"include: Indicates that HSTS settings apply to subdomains. " +
											"exclude: Indicates that HSTS settings do not apply to subdomains. " +
											"The default value for this parameter is exclude.",
									},
									"switch": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
										Description: "Indicates whether to enable HSTS. " +
											"This parameter can take the following values: " +
											"true: Indicates to enable HSTS. " +
											"false: Indicates to disable HSTS. " +
											"The default value for this parameter is false.",
									},
									"ttl": {
										Type:     schema.TypeInt,
										Optional: true,
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
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"https.0.forced_redirect"},
				Description:   "Indicates the configuration module for the forced redirection from HTTPS to HTTP. This feature is disabled by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enable_forced_redirect": {
							Type:     schema.TypeBool,
							Required: true,
							Description: "Indicates whether to enable the forced redirection from HTTPS. " +
								"This parameter can take the following values: " +
								"true: Indicates to enable the forced redirection from HTTPS. " +
								"Once enabled, the content delivery network will respond with StatusCode to inform the browser to send an HTTPS request when it receives an HTTP request from a user. " +
								"false: Indicates to disable the forced redirection from HTTPS.",
						},
						"status_code": {
							Type:     schema.TypeString,
							Required: true,
							Description: "Indicates the status code returned by the content delivery network when forced redirection from HTTPS occurs. " +
								"The default value for this parameter is 301.",
						},
					},
				},
			},
			"quic": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Indicates the QUIC configuration module. This feature is disabled by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Required: true,
							Type:     schema.TypeBool,
							Description: "Indicates whether to enable QUIC. " +
								"This parameter can take the following values: " +
								"true: Indicates to enable QUIC. " +
								"false: Indicates to disable QUIC.",
						},
					},
				},
			},
			"lock_template": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "Whether to lock the template. " +
					"If you set this field to true, then the template will be locked. Please note that the template cannot be modified or unlocked after it is locked. " +
					"When you want to use this template to create a domain name, please lock the template first. The default value is false.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnCipherTemplateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCipherTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnCipherTemplate())
	if err != nil {
		return fmt.Errorf("error on creating cdn_cipher_template %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCipherTemplateRead(d, meta)
}

func resourceByteplusCdnCipherTemplateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCipherTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnCipherTemplate())
	if err != nil {
		return fmt.Errorf("error on reading cdn_cipher_template %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnCipherTemplateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCipherTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnCipherTemplate())
	if err != nil {
		return fmt.Errorf("error on updating cdn_cipher_template %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnCipherTemplateRead(d, meta)
}

func resourceByteplusCdnCipherTemplateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnCipherTemplateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnCipherTemplate())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_cipher_template %q, %s", d.Id(), err)
	}
	return err
}
