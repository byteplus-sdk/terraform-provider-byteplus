package ipv6_address

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func DataSourceByteplusIpv6Addresses() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusIpv6AddressesRead,
		Schema: map[string]*schema.Schema{
			"associated_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the ECS instance that is assigned the IPv6 address.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of Ipv6Address query.",
			},
			"ipv6_addresses": {
				Description: "The collection of Ipv6Address query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ipv6_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPv6 address.",
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusIpv6AddressesRead(d *schema.ResourceData, meta interface{}) error {
	ipv6AddressService := NewIpv6AddressService(meta.(*bp.SdkClient))
	return ipv6AddressService.Dispatcher.Data(ipv6AddressService, d, DataSourceByteplusIpv6Addresses())
}
