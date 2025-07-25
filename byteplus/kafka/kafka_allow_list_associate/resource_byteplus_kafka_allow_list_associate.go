package kafka_allow_list_associate

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
KafkaAllowListAssociate can be imported using the id, e.g.
```
$ terraform import byteplus_kafka_allow_list_associate.default kafka-cnitzqgn**:acl-d1fd76693bd54e658912e7337d5b****
```

*/

func ResourceByteplusKafkaAllowListAssociate() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusKafkaAllowListAssociateCreate,
		Read:   resourceByteplusKafkaAllowListAssociateRead,
		Delete: resourceByteplusKafkaAllowListAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: importAllowListAssociate,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the kafka instance.",
			},
			"allow_list_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the allow list.",
			},
		},
	}
	return resource
}

func resourceByteplusKafkaAllowListAssociateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusKafkaAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on creating kafka_allow_list_associate %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaAllowListAssociateRead(d, meta)
}

func resourceByteplusKafkaAllowListAssociateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusKafkaAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on reading kafka_allow_list_associate %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusKafkaAllowListAssociateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaAllowListAssociateService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusKafkaAllowListAssociate())
	if err != nil {
		return fmt.Errorf("error on deleting kafka_allow_list_associate %q, %s", d.Id(), err)
	}
	return err
}

func importAllowListAssociate(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	var err error
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must be of the form InstanceId:AllowListId")
	}
	err = data.Set("instance_id", items[0])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	err = data.Set("allow_list_id", items[1])
	if err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
