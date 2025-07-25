package kafka_allow_list

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
KafkaAllowList can be imported using the id, e.g.
```
$ terraform import byteplus_kafka_allow_list.default resource_id
```

*/

func ResourceByteplusKafkaAllowList() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusKafkaAllowListCreate,
		Read:   resourceByteplusKafkaAllowListRead,
		Update: resourceByteplusKafkaAllowListUpdate,
		Delete: resourceByteplusKafkaAllowListDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"allow_list_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the allow list.",
			},
			"allow_list_desc": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the allow list.",
			},
			"allow_list": {
				Type:     schema.TypeSet,
				Required: true,
				Description: "Whitelist rule list. " +
					"Supports specifying as IP addresses or IP network segments. " +
					"Each whitelist can be configured with a maximum of 300 IP addresses or network segments.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	return resource
}

func resourceByteplusKafkaAllowListCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusKafkaAllowList())
	if err != nil {
		return fmt.Errorf("error on creating kafka_allow_list %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaAllowListRead(d, meta)
}

func resourceByteplusKafkaAllowListRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusKafkaAllowList())
	if err != nil {
		return fmt.Errorf("error on reading kafka_allow_list %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusKafkaAllowListUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusKafkaAllowList())
	if err != nil {
		return fmt.Errorf("error on updating kafka_allow_list %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaAllowListRead(d, meta)
}

func resourceByteplusKafkaAllowListDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusKafkaAllowList())
	if err != nil {
		return fmt.Errorf("error on deleting kafka_allow_list %q, %s", d.Id(), err)
	}
	return err
}
