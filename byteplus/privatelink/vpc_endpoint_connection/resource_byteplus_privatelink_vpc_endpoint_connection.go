package vpc_endpoint_connection

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
PrivateLink Vpc Endpoint Connection Service can be imported using the endpoint id and service id, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint_connection.default ep-3rel74u229dz45zsk2i6l69qa:epsvc-2byz5mykk9y4g2dx0efs4aqz3
```

*/

func ResourceByteplusPrivatelinkVpcEndpointConnectionService() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusPrivatelinkVpcEndpointConnectionCreate,
		Read:   resourceByteplusPrivatelinkVpcEndpointConnectionRead,
		Delete: resourceByteplusPrivatelinkVpcEndpointConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: vpcConnectionImporter,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the endpoint.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the security group.",
			},
		},
	}
	dataSource := DataSourceByteplusPrivatelinkVpcEndpointConnections().Schema["connections"].Elem.(*schema.Resource).Schema
	bp.MergeDateSourceToResource(dataSource, &resource.Schema)
	return resource
}

func resourceByteplusPrivatelinkVpcEndpointConnectionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcEndpointConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(service, d, ResourceByteplusPrivatelinkVpcEndpointConnectionService())
	if err != nil {
		return fmt.Errorf("error on creating private link VpcEndpointConnection service %q, %w", d.Id(), err)
	}
	return resourceByteplusPrivatelinkVpcEndpointConnectionRead(d, meta)
}

func resourceByteplusPrivatelinkVpcEndpointConnectionRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcEndpointConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(service, d, ResourceByteplusPrivatelinkVpcEndpointConnectionService())
	if err != nil {
		return fmt.Errorf("error on reading private link VpcEndpointConnection service %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusPrivatelinkVpcEndpointConnectionDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewVpcEndpointConnectionService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(service, d, ResourceByteplusPrivatelinkVpcEndpointConnectionService())
	if err != nil {
		return fmt.Errorf("error on deleting private link VpcEndpointConnection service %q, %w", d.Id(), err)
	}
	return nil
}

var vpcConnectionImporter = func(data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(data.Id(), ":")
	if len(items) != 2 {
		return []*schema.ResourceData{data}, fmt.Errorf("import id must split with ':'")
	}
	if err := data.Set("endpoint_id", items[0]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	if err := data.Set("service_id", items[1]); err != nil {
		return []*schema.ResourceData{data}, err
	}
	return []*schema.ResourceData{data}, nil
}
