package cdn_certificate

import (
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func DataSourceByteplusCdnCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceByteplusCdnCertificatesRead,
		Schema: map[string]*schema.Schema{
			"cert_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates a certificate ID to retrieve the certificate with that ID.",
			},
			"configured_domain": {
				Type:     schema.TypeSet,
				Optional: true,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Indicates a list of domain names for acceleration, to obtain certificates that have been bound to any domain name on the list.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates a domain name used to obtain certificates that include that domain name in the SAN field. The domain name can be a wildcard domain. For example, *.example.com can match certificates containing img.example.com or www.example.com, etc., in the SAN field.",
			},
			"fuzzy_match": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "When Name is specified, FuzzyMatch indicates the matching method used by the CDN when filtering certificates by Name. The parameter can have the following values:\ntrue: indicates fuzzy matching.\nfalse: indicates exact matching.\nIf you don not specify Name, FuzzyMatch is not effective.\nThe default value of FuzzyMatch is false.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Indicates a list of states to retrieve certificates that are in any of the states on the list. The parameter can have the following values:\nrunning: indicates certificates with a remaining validity period of more than 30 days.\nexpired: indicates certificates that have expired.\nexpiring_soon: indicates certificates with a remaining validity period of 30 days or less but have not yet expired.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A Name Regex of Resource.",
			},
			"output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File name where to save data source results.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of query.",
			},

			"certificates": {
				Description: "The collection of query.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the ID of the certificate.",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the ID of the certificate.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source of the certificate.",
						},
						"cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the content of the Common Name (CN) field of the certificate.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the status of the certificate. The parameter can have the following values:\nrunning: indicates the certificate has a remaining validity period of more than 30 days.\nexpired: indicates the certificate has expired.\nexpiring_soon: indicates the certificate has a remaining validity period of 30 days or less but has not yet expired.",
						},
						"dns_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the domain names in the SAN field of the certificate.",
						},
						"desc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the remark of the certificate.",
						},
						"configured_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the list of domain names associated with the certificate. If the certificate has not been associated with any domain name, the parameter value is null.",
						},
						"effective_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the issuance time of the certificate. The unit is Unix timestamp.",
						},
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the expiration time of the certificate. The unit is Unix timestamp.",
						},
						"cert_fingerprint": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sha1": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates a fingerprint based on the SHA-1 encryption algorithm, composed of 40 hexadecimal characters.",
									},
									"sha256": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates a fingerprint based on the SHA-256 encryption algorithm, composed of 64 hexadecimal characters.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceByteplusCdnCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	service := NewCdnCertificateService(meta.(*bp.SdkClient))
	return service.Dispatcher.Data(service, d, DataSourceByteplusCdnCertificates())
}
