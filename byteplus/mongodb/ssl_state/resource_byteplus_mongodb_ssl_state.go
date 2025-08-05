package ssl_state

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
mongodb ssl state can be imported using the ssl:instanceId, e.g.
```
$ terraform import byteplus_mongodb_ssl_state.default ssl:mongo-shard-d050db19xxx
```

*/

func ResourceByteplusMongoDBSSLState() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusMongoDBSSLStateCreate,
		Read:   resourceByteplusMongoDBSSLStateRead,
		Update: resourceByteplusMongoDBSSLStateUpdate,
		Delete: resourceByteplusMongoDBSSLStateDelete,
		Importer: &schema.ResourceImporter{
			State: mongoDBSSLStateImporter,
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
				Description: "The ID of mongodb instance.",
			},
			"ssl_action": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Id() == ""
				},
				ValidateFunc: validation.StringInSlice([]string{
					"Update",
				}, false),
				Description: "The action of ssl, valid value contains `Update`. Set `ssl_action` to `Update` will will trigger an SSL update operation when executing `terraform apply`." +
					"When the current time is less than 30 days from the `ssl_expired_time`, executing `terraform apply` will automatically renew the SSL.",
			},
		},
	}
	dataSource := DataSourceByteplusMongoDBSSLStates().Schema["ssl_state"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusMongoDBSSLStateCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBSSLStateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusMongoDBSSLState())
	if err != nil {
		return fmt.Errorf("Error on opening ssl %q, %s ", d.Id(), err)
	}
	return resourceByteplusMongoDBSSLStateRead(d, meta)
}

func resourceByteplusMongoDBSSLStateUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBSSLStateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(service, d, ResourceByteplusMongoDBSSLState())
	if err != nil {
		return fmt.Errorf("Error on updating ssl %q, %s ", d.Id(), err)
	}
	return resourceByteplusMongoDBSSLStateRead(d, meta)
}

func resourceByteplusMongoDBSSLStateRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBSSLStateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusMongoDBSSLState())
	if err != nil {
		return fmt.Errorf("Error on reading ssl state %q, %s ", d.Id(), err)
	}
	return err
}

func resourceByteplusMongoDBSSLStateDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewMongoDBSSLStateService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusMongoDBSSLState())
	if err != nil {
		return fmt.Errorf("Error on deleting ssl state %q, %s ", d.Id(), err)
	}
	return err
}
