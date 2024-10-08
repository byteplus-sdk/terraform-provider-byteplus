package classic_cdn_domain

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnDomain can be imported using the domain, e.g.
```
$ terraform import byteplus_cdn_domain.default www.byteplus.com
```
Please note that when you execute destroy, we will first take the domain name offline and then delete it.
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "You need to add a domain. The main account can add up to 200 accelerated domains.",
			},
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The business type of the domain name is indicated by this parameter. " +
					"The possible values are: `download`: for file downloads. `web`: for web pages. " +
					"`video`: for audio and video on demand.",
			},
			"service_region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: "Indicates the acceleration area. The parameter can take the following values: " +
					"`chinese_mainland`: Indicates mainland China. `global`: Indicates global." +
					" `outside_chinese_mainland`: Indicates global (excluding mainland China).",
			},
			"domain_config": {
				Type:     schema.TypeString,
				Required: true,
				Description: "Accelerate domain configuration. " +
					"Please convert the configuration module structure into json and pass it into a string. " +
					"You must specify the Origin module. The OriginProtocol parameter, OriginHost parameter, " +
					"and other domain configuration modules are optional.",
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "default",
				Description: "The project to which this domain name belongs. Default is `default`.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Indicate the tags you have set for this domain name. You can set up to 10 tags.",
				Set:         TagsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The key of the tag.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The value of the tag.",
						},
					},
				},
			},
			"shared_cname": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Configuration for sharing CNAME.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: "Specify whether to enable shared CNAME.",
						},
						"cname": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Assign a CNAME to the accelerated domain.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the domain.",
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

var TagsHash = func(v interface{}) int {
	if v == nil {
		return hashcode.String("")
	}
	m := v.(map[string]interface{})
	var (
		buf bytes.Buffer
	)
	buf.WriteString(fmt.Sprintf("%v#%v", m["key"], m["value"]))
	return hashcode.String(buf.String())
}
