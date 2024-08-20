package cdn_edge_function_associate

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CdnEdgeFunctionAssociate can be imported using the function_id:domain, e.g.
```
$ terraform import byteplus_cdn_edge_function_associate.default function_id:domain
```

*/

func ResourceByteplusCdnEdgeFunctionAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCdnEdgeFunctionAssociateCreate,
		Read:   resourceByteplusCdnEdgeFunctionAssociateRead,
		Delete: resourceByteplusCdnEdgeFunctionAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: functionAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"function_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the function for which you want to bind to domain.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The domain name which you wish to bind with the function.",
			},
		},
	}
	return resource
}

func resourceByteplusCdnEdgeFunctionAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCdnEdgeFunctionAssociate())
	if err != nil {
		return fmt.Errorf("error on creating cdn_edge_function_associate %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionAssociateRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCdnEdgeFunctionAssociate())
	if err != nil {
		return fmt.Errorf("error on reading cdn_edge_function_associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCdnEdgeFunctionAssociateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCdnEdgeFunctionAssociate())
	if err != nil {
		return fmt.Errorf("error on updating cdn_edge_function_associate %q, %s", d.Id(), err)
	}
	return resourceByteplusCdnEdgeFunctionAssociateRead(d, meta)
}

func resourceByteplusCdnEdgeFunctionAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCdnEdgeFunctionAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCdnEdgeFunctionAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting cdn_edge_function_associate %q, %s", d.Id(), err)
	}
	return err
}

var functionAssociateImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("function_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("domain", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
