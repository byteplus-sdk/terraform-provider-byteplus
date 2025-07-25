package kafka_public_address

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
KafkaPublicAddress can be imported using the instance_id:eip_id, e.g.
```
$ terraform import byteplus_kafka_public_address.default instance_id:eip_id
```

*/

func ResourceByteplusKafkaPublicAddress() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusKafkaPublicAddressCreate,
		Read:   resourceByteplusKafkaPublicAddressRead,
		Delete: resourceByteplusKafkaPublicAddressDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("eip_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("instance_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of kafka instance.",
				ForceNew:    true,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of eip.",
				ForceNew:    true,
			},
			"endpoint_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The endpoint type of instance.",
			},
			"network_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network type of instance.",
			},
			"public_endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public endpoint of instance.",
			},
		},
	}
	return resource
}

func resourceByteplusKafkaPublicAddressCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaInternetEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusKafkaPublicAddress())
	if err != nil {
		return fmt.Errorf("error on creating kafka public address %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaPublicAddressRead(d, meta)
}

func resourceByteplusKafkaPublicAddressRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaInternetEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusKafkaPublicAddress())
	if err != nil {
		return fmt.Errorf("error on reading kafka public address %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusKafkaPublicAddressDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaInternetEnablerService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusKafkaPublicAddress())
	if err != nil {
		return fmt.Errorf("error on deleting kafka public address %q, %s", d.Id(), err)
	}
	return err
}
