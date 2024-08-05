package cdn_domain

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnDomain can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_domain.default resource_id
```

*/

func ResourceByteplusCdnDomain() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnDomainCreate,
		Read:   resourceByteplusCdnDomainRead,
		Update: resourceByteplusCdnDomainUpdate,
		Delete: resourceByteplusCdnDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "Indicates a domain name you want to add. " +
					"The domain name you add must meet all of the following conditions: " +
					"Length does not exceed 100 characters. " +
					"Cannot contain uppercase letters. " +
					"Does not include any of these suffixes: zjgslb.com, yangyi19.com, volcgslb.com, veew-alb-cn1.com, bplgslb.com, bplslb.com, ttgslb.com. " +
					"When you bind your domain name with a delivery policy, " +
					"the origin address specified in the policy must not be the same as your domain name.",
			},
			"service_template_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Indicates a delivery policy to be bound with this domain name. " +
					"You can use DescribeTemplates to obtain the ID of the delivery policy you want to bind.",
			},
			"cert_id": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					httpsSwitch, ok := d.GetOk("https_switch")
					if ok {
						if httpsSwitch.(string) == "off" {
							return true
						} else if httpsSwitch.(string) == "on" {
							return false
						}
					}
					return true
				},
				Description: "Indicates the ID of a certificate. " +
					"This certificate is stored in the BytePlus Certificate Center and will be associated with the domain name. " +
					"If HTTPSSwitch is on, this parameter is required. " +
					"Before using this API, you need to grant CDN access to the Certificate Center, then upload your certificate to the BytePlus Certificate Center to obtain the ID of the certificate. " +
					"It is recommended to authorize CDN access to the Certificate Center using the primary account. " +
					"You can use ListCertInfo to obtain the ID of the certificate you want to associate. " +
					"If HTTPSSwitch is off, this parameter does not take effect.",
			},
			"cipher_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Indicates an encryption policy to be bound with this domain name. " +
					"You can use DescribeTemplates to obtain the ID of the encryption policy you want to bind. " +
					"If this parameter is not specified, it means that the domain name will not be bound to any encryption policy at present.",
			},
			"https_switch": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Indicates whether to enable \"HTTPS Encryption Service\" for this domain name. " +
					"This parameter can take the following values: " +
					"on: Indicates to enable this service. " +
					"off: Indicates not to enable this service. " +
					"The default value of this parameter is off.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates the project to which this domain name belongs, with the default value being default.",
			},
			"service_region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Indicates the service region enabled for this domain name. " +
					"This parameter can take the following values: " +
					"outside_chinese_mainland：Indicates \"Global (excluding Chinese Mainland)\"." +
					" chinese_mainland：Indicates \"Chinese Mainland\". " +
					"global: Indicates \"Global\". " +
					"The default value of this parameter is outside_chinese_mainland. " +
					"Note that chinese_mainland or global are not available by default. " +
					"To make the two service regions available, please submit a ticket. " +
					"Also, since both regions include Chinese Mainland, " +
					"you must complete the following additional actions: " +
					"Perform real-name authentication for your BytePlus account. " +
					"Perform ICP filing for your domain name.",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Indicates the status of the domain name. " +
					"This parameter can be: " +
					"online: Indicates the status is Enabled. " +
					"offline: Indicates the status is Disabled. " +
					"configuring: Indicates the status is Configuring.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnDomainCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnDomain())
	if err != nil {
		return fmt.Errorf("error on creating cdn_domain %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnDomainRead(d, meta)
}

func resourceByteplusCdnDomainRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnDomain())
	if err != nil {
		return fmt.Errorf("error on reading cdn_domain %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnDomainUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnDomain())
	if err != nil {
		return fmt.Errorf("error on updating cdn_domain %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnDomainRead(d, meta)
}

func resourceByteplusCdnDomainDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnDomain())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_domain %q, %s", d.Id(), err)
	}
	return err
}
