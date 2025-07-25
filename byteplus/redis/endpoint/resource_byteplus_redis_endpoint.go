package endpoint

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
Redis Endpoint can be imported using the instanceId:eipId, e.g.
```
$ terraform import byteplus_redis_endpoint.default redis-asdljioeixxxx:eip-2fef2qcfbfw8w5oxruw3w****
```
*/

func ResourceByteplusRedisEndpoint() *schema.Resource {
	resource := &schema.Resource{
		Read:   resourceByteplusRedisEndpointRead,
		Create: resourceByteplusRedisEndpointCreate,
		Delete: resourceByteplusRedisEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: redisEndpointAssociateImporter,
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
				Description: "Id of instance.",
			},
			"eip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of eip.",
			},
		},
	}

	return resource
}

func resourceByteplusRedisEndpointRead(d *schema.ResourceData, meta interface{}) (err error) {
	redisEndpointService := NewRedisEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(redisEndpointService, d, ResourceByteplusRedisEndpoint())
	if err != nil {
		return fmt.Errorf("error on reading redis endpoint %v, %v", d.Id(), err)
	}
	return nil
}

func resourceByteplusRedisEndpointCreate(d *schema.ResourceData, meta interface{}) (err error) {
	redisEndpointService := NewRedisEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(redisEndpointService, d, ResourceByteplusRedisEndpoint())
	if err != nil {
		return fmt.Errorf("error on creating redis endpoint %v, %v", d.Id(), err)
	}
	return resourceByteplusRedisEndpointRead(d, meta)
}

func resourceByteplusRedisEndpointDelete(d *schema.ResourceData, meta interface{}) (err error) {
	redisEndpointService := NewRedisEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(redisEndpointService, d, ResourceByteplusRedisEndpoint())
	if err != nil {
		return fmt.Errorf("error on deleting redis endpoint %q, %s", d.Id(), err)
	}
	return err
}
