package eip_address

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

/*

Import
Eip address can be imported using the id, e.g.
```
$ terraform import byteplus_eip_address.default eip-274oj9a8rs9a87fap8sf9515b
```

*/

func ResourceByteplusEipAddress() *schema.Resource {
	return &schema.Resource{
		Delete: resourceByteplusEipAddressDelete,
		Create: resourceByteplusEipAddressCreate,
		Read:   resourceByteplusEipAddressRead,
		Update: resourceByteplusEipAddressUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"billing_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaidByBandwidth", "PostPaidByTraffic"}, false),
				Description:  "The billing type of the EIP Address. And optional choice contains `PostPaidByBandwidth` or `PostPaidByTraffic` or `PrePaid`.",
			},
			//"period_unit": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Default:  "Month",
			//	ValidateFunc: validation.StringInSlice([]string{
			//		"Month", "Year",
			//	}, false),
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		// 创建时，只有付费类型为 PrePaid 时生效
			//		if d.Id() == "" {
			//			if d.Get("billing_type").(string) == "PrePaid" {
			//				return false
			//			}
			//		} else { // 修改时，只有付费类型由按量付费转为 PrePaid 时生效
			//			if d.HasChange("billing_type") && d.Get("billing_type").(string) == "PrePaid" {
			//				return false
			//			}
			//		}
			//		return true
			//	},
			//	Description: "The period unit of the EIP Address. Optional choice contains `Month` or `Year`. Default is `Month`." +
			//		"This field is only effective when creating a PrePaid Eip or changing the billing_type from PostPaid to PrePaid.",
			//},
			"period": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  12,
				ValidateFunc: validation.Any(
					validation.IntBetween(1, 9),
					validation.IntInSlice([]int{12, 36})),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// 创建时，只有付费类型为 PrePaid 时生效
					if d.Id() == "" {
						if d.Get("billing_type").(string) == "PrePaid" {
							return false
						}
					} else { // 修改时，只有付费类型由按量付费转为 PrePaid 时生效
						if d.HasChange("billing_type") && d.Get("billing_type").(string) == "PrePaid" {
							return false
						}
					}
					return true
				},
				Description: "The period of the EIP Address, the valid value range in 1~9 or 12 or 36. Default value is 12. The period unit defaults to `Month`." +
					"This field is only effective when creating a PrePaid Eip or changing the billing_type from PostPaid to PrePaid.",
			},
			"bandwidth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				//ValidateFunc: validation.IntBetween(1, 500),
				//Description:  "The peek bandwidth of the EIP, the value range in 1~500 for PostPaidByBandwidth, and 1~200 for PostPaidByTraffic.",
				Description: "The peek bandwidth of the EIP.",
			},
			"isp": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ISP of the EIP, the value can be `BGP` or `ChinaMobile` or `ChinaUnicom` or `ChinaTelecom` or `SingleLine_BGP` or `Static_BGP` or `Fusion_BGP`.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the EIP Address.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the EIP.",
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ProjectName of the EIP.",
			},
			"security_protection_types": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Security protection types for public IP addresses. " +
					"Parameter - N: Indicates the number of security protection types, currently only supports taking 1. Value: `AntiDDoS_Enhanced` or left blank." +
					"If the value is `AntiDDoS_Enhanced`, then will create an eip with enhanced protection," +
					"(can be added to DDoS native protection (enterprise version) instance). " +
					"If left blank, it indicates an eip with basic protection.",
			},
			"tags": bp.TagsSchema(),
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the EIP.",
			},
			"eip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the EIP.",
			},
			"overdue_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The overdue time of the EIP.",
			},
			"deleted_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deleted time of the EIP.",
			},
			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expired time of the EIP.",
			},
		},
	}
}

func resourceByteplusEipAddressCreate(d *schema.ResourceData, meta interface{}) error {
	eipAddressService := NewEipAddressService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Create(eipAddressService, d, ResourceByteplusEipAddress()); err != nil {
		return fmt.Errorf("error on creating eip address  %q, %w", d.Id(), err)
	}
	return resourceByteplusEipAddressRead(d, meta)
}

func resourceByteplusEipAddressRead(d *schema.ResourceData, meta interface{}) error {
	eipAddressService := NewEipAddressService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Read(eipAddressService, d, ResourceByteplusEipAddress()); err != nil {
		return fmt.Errorf("error on reading  eip address %q, %w", d.Id(), err)
	}
	return nil
}

func resourceByteplusEipAddressUpdate(d *schema.ResourceData, meta interface{}) error {
	eipAddressService := NewEipAddressService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Update(eipAddressService, d, ResourceByteplusEipAddress()); err != nil {
		return fmt.Errorf("error on updating  eip address %q, %w", d.Id(), err)
	}
	return resourceByteplusEipAddressRead(d, meta)
}

func resourceByteplusEipAddressDelete(d *schema.ResourceData, meta interface{}) error {
	eipAddressService := NewEipAddressService(meta.(*bp.SdkClient))
	if err := bp.DefaultDispatcher().Delete(eipAddressService, d, ResourceByteplusEipAddress()); err != nil {
		return fmt.Errorf("error on deleting  eip address %q, %w", d.Id(), err)
	}
	return nil
}
