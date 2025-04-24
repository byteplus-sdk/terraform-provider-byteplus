package rds_mysql_endpoint_public_address

import (
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
RdsMysqlEndpointPublicAddress can be imported using the instance id, endpoint id and eip id, e.g.
```
$ terraform import byteplus_rds_mysql_endpoint_public_address.default instanceId:endpointId:eipId
```

*/

func ResourceByteplusRdsMysqlEndpointPublicAddress() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusRdsMysqlEndpointPublicAddressCreate,
		Read:   resourceByteplusRdsMysqlEndpointPublicAddressRead,
		Update: resourceByteplusRdsMysqlEndpointPublicAddressUpdate,
		Delete: resourceByteplusRdsMysqlEndpointPublicAddressDelete,
		Importer: &schema.ResourceImporter{
			State: addressImporter,
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
				Description: "The id of mysql instance.",
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the endpoint.",
			},
			"eip_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the eip.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The domain.",
			},
		},
	}
	return resource
}

func resourceByteplusRdsMysqlEndpointPublicAddressCreate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointPublicAddressService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Create(service, d, ResourceByteplusRdsMysqlEndpointPublicAddress())
	if err != nil {
		return fmt.Errorf("error on creating rds_mysql_endpoint_public_address %q, %s", d.Id(), err)
	}
	return resourceByteplusRdsMysqlEndpointPublicAddressRead(d, meta)
}

func resourceByteplusRdsMysqlEndpointPublicAddressRead(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointPublicAddressService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Read(service, d, ResourceByteplusRdsMysqlEndpointPublicAddress())
	if err != nil {
		return fmt.Errorf("error on reading rds_mysql_endpoint_public_address %q, %s", d.Id(), err)
	}
	return err
}

func resourceByteplusRdsMysqlEndpointPublicAddressUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointPublicAddressService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Update(service, d, ResourceByteplusRdsMysqlEndpointPublicAddress())
	if err != nil {
		return fmt.Errorf("error on updating rds_mysql_endpoint_public_address %q, %s", d.Id(), err)
	}
	return resourceByteplusRdsMysqlEndpointPublicAddressRead(d, meta)
}

func resourceByteplusRdsMysqlEndpointPublicAddressDelete(d *schema.ResourceData, meta interface{}) (err error) {
	service := NewRdsMysqlEndpointPublicAddressService(meta.(*bp.SdkClient))
	err = service.Dispatcher.Delete(service, d, ResourceByteplusRdsMysqlEndpointPublicAddress())
	if err != nil {
		return fmt.Errorf("error on deleting rds_mysql_endpoint_public_address %q, %s", d.Id(), err)
	}
	return err
}

func addressImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	items := strings.Split(d.Id(), ":")
	if len(items) != 3 {
		return []*schema.ResourceData{d}, fmt.Errorf("the format of import id must be 'instanceId:endpointId:eipId'")
	}
	instanceId := items[0]
	endpointId := items[1]
	eipId := items[2]
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("endpoint_id", endpointId)
	_ = d.Set("eip_id", eipId)
	return []*schema.ResourceData{d}, nil
}
