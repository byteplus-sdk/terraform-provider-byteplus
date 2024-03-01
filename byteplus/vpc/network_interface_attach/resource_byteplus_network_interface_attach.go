package network_interface_attach

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Network interface attach can be imported using the network_interface_id:instance_id.
```
$ terraform import byteplus_network_interface_attach.default eni-bp1fg655nh68xyz9***:i-wijfn35c****
```

*/

func ResourceByteplusNetworkInterfaceAttach() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusNetworkInterfaceAttachDelete,
		Create: resourceByteplusNetworkInterfaceAttachCreate,
		Read:   resourceByteplusNetworkInterfaceAttachRead,
		Importer: &schema.ResourceImporter{
			State: networkInterfaceAttachImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"network_interface_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the ENI.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the instance to which the ENI is bound.",
			},
		},
	}
}

func resourceByteplusNetworkInterfaceAttachCreate(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceAttachService := NewNetworkInterfaceAttachService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(networkInterfaceAttachService, d, ResourceByteplusNetworkInterfaceAttach()); err != nil {
		return fmt.Errorf("error on creating network interface attach %q, %w", d.Id(), err)
	}
	return resourceByteplusNetworkInterfaceAttachRead(d, meta)
}

func resourceByteplusNetworkInterfaceAttachRead(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceAttachService := NewNetworkInterfaceAttachService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(networkInterfaceAttachService, d, ResourceByteplusNetworkInterfaceAttach()); err != nil {
		return fmt.Errorf("error on reading network interface attach %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusNetworkInterfaceAttachDelete(d *schema.ResourceData, meta interface{}) error {
	networkInterfaceAttachService := NewNetworkInterfaceAttachService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(networkInterfaceAttachService, d, ResourceByteplusNetworkInterfaceAttach()); err != nil {
		return fmt.Errorf("error on deleting network interface attach %q, %w", d.Id(), err)
	}
	return nil
}
