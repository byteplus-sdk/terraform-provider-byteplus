package dnat_entry

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Dnat entry can be imported using the id, e.g.
```
$ terraform import byteplus_dnat_entry.default dnat-3fvhk47kf56****
```

*/

func ResourceByteplusDnatEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceByteplusDnatEntryCreate,
		Update: resourceByteplusDnatEntryUpdate,
		Read:   resourceByteplusDnatEntryRead,
		Delete: resourceByteplusDnatEntryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"nat_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the nat gateway to which the entry belongs.",
			},
			"external_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provides the public IP address for public network access.",
			},
			"external_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port or port segment that receives requests from the public network. If InternalPort is passed into the port segment, ExternalPort must also be passed into the port segment.",
			},
			"internal_ip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Provides the internal IP address.",
			},
			"internal_port": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The port or port segment on which the cloud server instance provides services to the public network.",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
				Description:  "The network protocol.",
			},
			"dnat_entry_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the DNAT rule.",
			},
			"dnat_entry_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the DNAT rule.",
			},
		},
	}
}

func resourceByteplusDnatEntryCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusDnatEntry())
	if err != nil {
		return fmt.Errorf("error on creating dnat entry: %q, %w", d.Id(), err)
	}
	return resourceByteplusDnatEntryRead(d, meta)
}

func resourceByteplusDnatEntryRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusDnatEntry())
	if err != nil {
		return fmt.Errorf("error on reading dnat entry: %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusDnatEntryUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusDnatEntry())
	if err != nil {
		return fmt.Errorf("error on updating dnat entry: %q, %w", d.Id(), err)
	}
	return resourceByteplusDnatEntryRead(d, meta)
}

func resourceByteplusDnatEntryDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewDnatEntryService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusDnatEntry())
	if err != nil {
		return fmt.Errorf("error on deleting dnat entry: %q, %w", d.Id(), err)
	}
	return nil
}
