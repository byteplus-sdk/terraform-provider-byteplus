package cr_endpoint_acl_policy

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
CrEndpointAclPolicy can be imported using the registry:entry, e.g.
```
$ terraform import byteplus_cr_endpoint_acl_policy.default resource_id
```

*/

func ResourceByteplusCrEndpointAclPolicy() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusCrEndpointAclPolicyCreate,
		Read:   resourceByteplusCrEndpointAclPolicyRead,
		Delete: resourceByteplusCrEndpointAclPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: crEndpointAclPolicyImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"registry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The registry name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of the acl policy. Valid values: `Public`.",
			},
			"entry": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ip list of the acl policy.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The description of the acl policy.",
			},
		},
	}
	return resource
}

func resourceByteplusCrEndpointAclPolicyCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointAclPolicyService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusCrEndpointAclPolicy())
	if err != nil {
		return fmt.Errorf("error on creating cr_endpoint_acl_policy %q, %s", d.Id(), err)
	}
	return resourceByteplusCrEndpointAclPolicyRead(d, meta)
}

func resourceByteplusCrEndpointAclPolicyRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointAclPolicyService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusCrEndpointAclPolicy())
	if err != nil {
		return fmt.Errorf("error on reading cr_endpoint_acl_policy %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusCrEndpointAclPolicyUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointAclPolicyService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusCrEndpointAclPolicy())
	if err != nil {
		return fmt.Errorf("error on updating cr_endpoint_acl_policy %q, %s", d.Id(), err)
	}
	return resourceByteplusCrEndpointAclPolicyRead(d, meta)
}

func resourceByteplusCrEndpointAclPolicyDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewCrEndpointAclPolicyService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusCrEndpointAclPolicy())
	if err != nil {
		return fmt.Errorf("error on deleting cr_endpoint_acl_policy %q, %s", d.Id(), err)
	}
	return err
}

func crEndpointAclPolicyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'registry:entry'")
	}
	if err := d.Set("registry", items[0]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	if err := d.Set("entry", items[1]); err != nil {
		return []*schema.ResourceData{d}, err
	}
	return []*schema.ResourceData{d}, nil
}
