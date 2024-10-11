package vpc_endpoint

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

/*

Import
VpcEndpoint can be imported using the id, e.g.
```
$ terraform import byteplus_privatelink_vpc_endpoint.default ep-3rel74u229dz45zsk2i6l****
```

*/

func ResourceByteplusPrivatelinkVpcEndpoint() *schema.Resource {
	resource := &schema.Resource{
		Create: resourceByteplusPrivateLinkVpcEndpointCreate,
		Read:   resourceByteplusPrivateLinkVpcEndpointRead,
		Update: resourceByteplusPrivateLinkVpcEndpointUpdate,
		Delete: resourceByteplusPrivateLinkVpcEndpointDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "The security group ids of vpc endpoint. " +
					"It is recommended to bind security groups using the 'security_group_ids' field in this resource instead of using `byteplus_privatelink_security_group`.\n" +
					"For operations that jointly use this resource and `byteplus_privatelink_security_group`, use lifecycle ignore_changes to suppress changes to the 'security_group_ids' field.",
			},
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of vpc endpoint service.",
			},
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of vpc endpoint service.",
			},
			"endpoint_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of vpc endpoint.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The description of vpc endpoint.",
			},

			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The vpc id of vpc endpoint.",
			},
			"endpoint_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of vpc endpoint.",
			},
			"endpoint_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain of vpc endpoint.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of vpc endpoint.",
			},
			"business_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whether the vpc endpoint is locked.",
			},
			"connection_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The connection  status of vpc endpoint.",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The create time of vpc endpoint.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of vpc endpoint.",
			},
			"deleted_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The delete time of vpc endpoint.",
			},
		},
	}
	return resource
}

func resourceByteplusPrivateLinkVpcEndpointCreate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointService := NewVpcEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Create(vpcEndpointService, d, ResourceByteplusPrivatelinkVpcEndpoint())
	if err != nil {
		return fmt.Errorf("error on creating vpc endpoint %q, %w", d.Id(), err)
	}
	return resourceByteplusPrivateLinkVpcEndpointRead(d, meta)
}

func resourceByteplusPrivateLinkVpcEndpointRead(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointService := NewVpcEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Read(vpcEndpointService, d, ResourceByteplusPrivatelinkVpcEndpoint())
	if err != nil {
		return fmt.Errorf("error on reading vpc endpoint %q, %w", d.Id(), err)
	}
	return err
}

func resourceByteplusPrivateLinkVpcEndpointUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointService := NewVpcEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Update(vpcEndpointService, d, ResourceByteplusPrivatelinkVpcEndpoint())
	if err != nil {
		return fmt.Errorf("error on updating vpc endoint %q, %w", d.Id(), err)
	}
	return resourceByteplusPrivateLinkVpcEndpointRead(d, meta)
}

func resourceByteplusPrivateLinkVpcEndpointDelete(d *schema.ResourceData, meta interface{}) (err error) {
	vpcEndpointService := NewVpcEndpointService(meta.(*bp.SdkClient))
	err = bp.DefaultDispatcher().Delete(vpcEndpointService, d, ResourceByteplusPrivatelinkVpcEndpoint())
	if err != nil {
		return fmt.Errorf("error on deleting vpc endpoint %q, %w", d.Id(), err)
	}
	return err
}
