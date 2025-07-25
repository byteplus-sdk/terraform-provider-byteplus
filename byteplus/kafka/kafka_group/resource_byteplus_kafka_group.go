package kafka_group

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
KafkaGroup can be imported using the instance_id:group_id, e.g.
```
$ terraform import byteplus_kafka_group.default kafka-****x:groupId
```

*/

func ResourceByteplusKafkaGroup() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusKafkaGroupCreate,
		Read:   resourceByteplusKafkaGroupRead,
		Update: resourceByteplusKafkaGroupUpdate,
		Delete: resourceByteplusKafkaGroupDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
				items := strings.Split(data.Id(), ":")
				if len(items) != 2 {
					return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
				}
				if err := data.Set("instance_id", items[0]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				if err := data.Set("group_id", items[1]); err != nil {
					return []*schema.ResourceData{data}, err
				}
				return []*schema.ResourceData{data}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance id of kafka group.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of kafka group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of kafka group.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of kafka group.",
			},
		},
	}
	return resource
}

func resourceByteplusKafkaGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusKafkaGroup())
	if err != nil {
		return fmt.Errorf("error on creating kafka_group %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaGroupRead(d, meta)
}

func resourceByteplusKafkaGroupRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusKafkaGroup())
	if err != nil {
		return fmt.Errorf("error on reading kafka_group %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusKafkaGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusKafkaGroup())
	if err != nil {
		return fmt.Errorf("error on updating kafka_group %q, %s", d.Id(), err)
	}
	return resourceByteplusKafkaGroupRead(d, meta)
}

func resourceByteplusKafkaGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewKafkaGroupService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusKafkaGroup())
	if err != nil {
		return fmt.Errorf("error on deleting kafka_group %q, %s", d.Id(), err)
	}
	return err
}
