package cr_endpoint

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CR endpoints can be imported using the endpoint:registryName, e.g.
```
$ terraform import byteplus_cr_endpoint.default endpoint:cr-basic
```

*/

func crEndpointImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must start with 'endpoint:',eg: 'endpoint:[registry-1]'")
	}
	if err := d.Set("registry", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}

func ResourceByteplusCrEndpoint() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusCrEndpointRead,
		Create: resourceByteplusCrEndpointCreate,
		Update: resourceByteplusCrEndpointUpdate,
		Delete: resourceByteplusCrEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: crEndpointImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CrRegistry name.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether enable public endpoint.",
			},
		},
	}
	dataSource := DataSourceByteplusCrEndpoints().Schema["endpoints"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusCrEndpointCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusCrEndpoint())
	if err != nil {
		return fmt.Errorf("Error on creating CrEndpoint %q,%s", d.Id(), err)
	}
	return resourceByteplusCrEndpointRead(d, meta)
}

func resourceByteplusCrEndpointUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusCrEndpoint())
	if err != nil {
		return fmt.Errorf("error on updating CrEndpoint  %q, %s", d.Id(), err)
	}
	return resourceByteplusCrEndpointRead(d, meta)
}

func resourceByteplusCrEndpointDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusCrEndpoint())
	if err != nil {
		return fmt.Errorf("error on deleting CrEndpoint %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrEndpointRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusCrEndpoint())
	if err != nil {
		return fmt.Errorf("Error on reading CrEndpoint %q,%s", d.Id(), err)
	}
	return err
}
