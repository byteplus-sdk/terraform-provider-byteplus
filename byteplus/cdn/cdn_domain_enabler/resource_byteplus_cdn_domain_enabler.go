package cdn_domain_enabler

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnDomainEnabler can be imported using the id, e.g.
```
$ terraform import byteplus_cdn_domain_enabler.default enabler:resource_id
```

*/

func ResourceByteplusCdnDomainEnabler() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnDomainEnablerCreate,
		Read:   resourceByteplusCdnDomainEnablerRead,
		Delete: resourceByteplusCdnDomainEnablerDelete,
		Importer: &schema.ResourceImporter{
			State: cdnDomainEnablerState,
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
				Description: "Indicate the domain name you want to enable.",
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

func resourceByteplusCdnDomainEnablerCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnDomainEnabler())
	if err != nil {
		return fmt.Errorf("error on creating cdn_domain_enabler %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnDomainEnablerRead(d, meta)
}

func resourceByteplusCdnDomainEnablerRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnDomainEnabler())
	if err != nil {
		return fmt.Errorf("error on reading cdn_domain_enabler %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnDomainEnablerDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnDomainEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnDomainEnabler())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_domain_enabler %q, %s", d.Id(), err)
	}
	return err
}

var cdnDomainEnablerState = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("domain", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
