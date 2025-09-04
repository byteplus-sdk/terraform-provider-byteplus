package instance_parameter

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
mongodb parameter can be imported using the param:instanceId:parameterName:parameterRole, e.g.
```
$ terraform import byteplus_mongodb_instance_parameter.default param:mongo-replica-e405f8e2****:connPoolMaxConnsPerHost
```

*/

func ResourceByteplusMongoDBInstanceParameter() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusMongoDBInstanceParameterCreate,
		Read:   resourceByteplusMongoDBInstanceParameterRead,
		Update: resourceByteplusMongoDBInstanceParameterUpdate,
		Delete: resourceByteplusMongoDBInstanceParameterDelete,
		Importer: &schema.ResourceImporter{
			State: mongoDBParameterImporter,
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
				Description: "The instance ID.",
			},
			"parameter_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of parameter. The parameter resource can only be added or modified, deleting this resource will not actually execute any operation.",
			},
			"parameter_role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The node type to which the parameter belongs. The value range is as follows: Node, Shard, ConfigServer, Mongos.",
			},
			"parameter_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The value of parameter.",
			},
		},
	}

	return resource
}

func resourceByteplusMongoDBInstanceParameterCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBInstanceParameterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusMongoDBInstanceParameter())
	if err != nil {
		return fmt.Errorf("Error on creating instance parameters %q, %s ", d.Id(), err)
	}
	return nil
}

func resourceByteplusMongoDBInstanceParameterUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBInstanceParameterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusMongoDBInstanceParameter())
	if err != nil {
		return fmt.Errorf("Error on updating instance parameters %q, %s ", d.Id(), err)
	}
	return resourceByteplusMongoDBInstanceParameterRead(d, meta)
}

func resourceByteplusMongoDBInstanceParameterDelete(d *schema.ResourceData, meta interface{}) (err error) {
	return nil
}

func resourceByteplusMongoDBInstanceParameterRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBInstanceParameterService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusMongoDBInstanceParameter())
	if err != nil {
		return fmt.Errorf("Error on reading instance parameters %q, %s ", d.Id(), err)
	}
	return err
}

func mongoDBParameterImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 4 || items[0] != "param" {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'param:instanceId:parameterName'")
	}
	_ = d.Set("instance_id", items[1])
	_ = d.Set("parameter_name", items[2])
	_ = d.Set("parameter_role", items[3])
	return []*schema.ResourceData{d}, nil
}
