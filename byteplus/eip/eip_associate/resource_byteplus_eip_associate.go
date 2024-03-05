package eip_associate

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Eip associate can be imported using the eip allocation_id:instance_id, e.g.
```
$ terraform import byteplus_eip_associate.default eip-274oj9a8rs9a87fap8sf9515b:i-cm9t9ug9lggu79yr5tcw
```

*/

func ResourceByteplusEipAssociate() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusEipAssociateDelete,
		Create: resourceByteplusEipAssociateCreate,
		Read:   resourceByteplusEipAssociateRead,
		Importer: &schema.ResourceImporter{
			State: eipAssociateImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allocation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The allocation id of the EIP.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance id which be associated to the EIP.",
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "Nat" && new == "NAT" {
						return true
					}
					return false
				},
				Description: "The type of the associated instance,the value is `Nat` or `NetworkInterface` or `ClbInstance` or `EcsInstance` or `HaVip`.",
			},
			"private_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The private IP address of the instance will be associated to the EIP.",
			},
		},
	}
}

func resourceByteplusEipAssociateCreate(d *schema.ResourceData, meta interface{}) error {
	eipAssociateService := NewEipAssociateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(eipAssociateService, d, ResourceByteplusEipAssociate()); err != nil {
		return fmt.Errorf("error on creating eip associate %q, %w", d.Id(), err)
	}
	return resourceByteplusEipAssociateRead(d, meta)
}

func resourceByteplusEipAssociateRead(d *schema.ResourceData, meta interface{}) error {
	eipAssociateService := NewEipAssociateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(eipAssociateService, d, ResourceByteplusEipAssociate()); err != nil {
		return fmt.Errorf("error on reading  eip associate %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusEipAssociateDelete(d *schema.ResourceData, meta interface{}) error {
	eipAssociateService := NewEipAssociateService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(eipAssociateService, d, ResourceByteplusEipAssociate()); err != nil {
		return fmt.Errorf("error on deleting  eip associate %q, %w", d.Id(), err)
	}
	return nil
}
