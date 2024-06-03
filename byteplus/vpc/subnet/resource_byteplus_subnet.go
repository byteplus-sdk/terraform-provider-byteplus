package subnet

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Subnet can be imported using the id, e.g.
```
$ terraform import byteplus_subnet.default subnet-274oj9a8rs9a87fap8sf9515b
```

*/

func ResourceByteplusSubnet() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusSubnetDelete,
		Create: resourceByteplusSubnetCreate,
		Read:   resourceByteplusSubnetRead,
		Update: resourceByteplusSubnetUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"cidr_block": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsCIDR,
				Description:  "A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).",
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Subnet.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the Subnet.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of Subnet.",
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the VPC.",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of the Zone.",
			},
			"enable_ipv6": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时不存在这个参数，修改时存在这个参数
					return d.Id() == ""
				},
				Description: "Specifies whether to enable the IPv6 CIDR block of the Subnet. This field is only valid when modifying the Subnet.",
			},
			"ipv6_cidr_block": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() == "" {
						return false
					} else {
						if d.HasChange("enable_ipv6") && d.Get("enable_ipv6").(bool) {
							return false
						}
						return true
					}
				},
				Description: "The last eight bits of the IPv6 CIDR block of the Subnet. Valid values: 0 - 255.",
			},
			"tags": bp.TagsSchema(),
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of Subnet.",
			},
		},
	}
}

func resourceByteplusSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(subnetService, d, ResourceByteplusSubnet()); err != nil {
		return fmt.Errorf("error on creating subnet  %q, %w", d.Id(), err)
	}
	return resourceByteplusSubnetRead(d, meta)
}

func resourceByteplusSubnetRead(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(subnetService, d, ResourceByteplusSubnet()); err != nil {
		return fmt.Errorf("error on reading subnet %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(subnetService, d, ResourceByteplusSubnet()); err != nil {
		return fmt.Errorf("error on updating subnet %q, %w", d.Id(), err)
	}
	return resourceByteplusSubnetRead(d, meta)
}

func resourceByteplusSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	subnetService := NewSubnetService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(subnetService, d, ResourceByteplusSubnet()); err != nil {
		return fmt.Errorf("error on deleting subnet %q, %w", d.Id(), err)
	}
	return nil
}
