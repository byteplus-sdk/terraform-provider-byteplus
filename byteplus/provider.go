package byteplus

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/byteplus-sdk/byteplus-sdk-golang/base"
	"github.com/byteplus-sdk/byteplus-sdk-golang/service/sts"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/acl_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/certificate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/clb"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/listener"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/rule"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/server_group"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/server_group_server"
	clbZone "github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/clb/zone"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ebs/volume"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/image"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/zone"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/dnat_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/nat_gateway"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/snat_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_address"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_address_bandwidth"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/ipv6_gateway"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_acl"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_acl_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_interface"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/network_interface_attach"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_table"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/route_table_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/security_group"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/security_group_rule"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/subnet"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpc/vpc"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_ACCESS_KEY", nil),
				Description: "The Access Key for BytePlus Provider",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_SECRET_KEY", nil),
				Description: "The Secret Key for BytePlus Provider",
			},
			"session_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_SESSION_TOKEN", nil),
				Description: "The Session Token for BytePlus Provider",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_REGION", nil),
				Description: "The Region for BytePlus Provider",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_ENDPOINT", nil),
				Description: "The Customer Endpoint for BytePlus Provider",
			},
			"disable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_DISABLE_SSL", nil),
				Description: "Disable SSL for BytePlus Provider",
			},
			"customer_headers": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_CUSTOMER_HEADERS", nil),
				Description: "CUSTOMER HEADERS for BytePlus Provider",
			},
			"customer_endpoints": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_CUSTOMER_ENDPOINTS", nil),
				Description: "CUSTOMER ENDPOINTS for BytePlus Provider",
			},
			"proxy_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_PROXY_URL", nil),
				Description: "PROXY URL for BytePlus Provider",
			},
			"assume_role": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "The ASSUME ROLE block for BytePlus Provider. If provided, terraform will attempt to assume this role using the supplied credentials.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"assume_role_trn": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_ASSUME_ROLE_TRN", nil),
							Description: "The TRN of the role to assume.",
						},
						"assume_role_session_name": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("BYTEPLUS_ASSUME_ROLE_SESSION_NAME", nil),
							Description: "The session name to use when making the AssumeRole call.",
						},
						"duration_seconds": {
							Type:     schema.TypeInt,
							Required: true,
							DefaultFunc: func() (interface{}, error) {
								if v := os.Getenv("BYTEPLUS_ASSUME_ROLE_DURATION_SECONDS"); v != "" {
									return strconv.Atoi(v)
								}
								return 3600, nil
							},
							ValidateFunc: validation.IntBetween(900, 43200),
							Description:  "The duration of the session when making the AssumeRole call. Its value ranges from 900 to 43200(seconds), and default is 3600 seconds.",
						},
						"policy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A more restrictive policy when making the AssumeRole call.",
						},
					},
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			// ================ ECS ================
			"byteplus_zones":               zone.DataSourceByteplusZones(),
			"byteplus_images":              image.DataSourceByteplusImages(),
			"byteplus_ecs_instances":       ecs_instance.DataSourceByteplusEcsInstances(),
			"byteplus_ecs_deployment_sets": ecs_deployment_set.DataSourceByteplusEcsDeploymentSets(),

			// ================ VPC ================
			"byteplus_vpcs":                        vpc.DataSourceByteplusVpcs(),
			"byteplus_subnets":                     subnet.DataSourceByteplusSubnets(),
			"byteplus_security_groups":             security_group.DataSourceByteplusSecurityGroups(),
			"byteplus_security_group_rules":        security_group_rule.DataSourceByteplusSecurityGroupRules(),
			"byteplus_network_interfaces":          network_interface.DataSourceByteplusNetworkInterfaces(),
			"byteplus_route_tables":                route_table.DataSourceByteplusRouteTables(),
			"byteplus_route_entries":               route_entry.DataSourceByteplusRouteEntries(),
			"byteplus_vpc_ipv6_gateways":           ipv6_gateway.DataSourceByteplusIpv6Gateways(),
			"byteplus_vpc_ipv6_addresses":          ipv6_address.DataSourceByteplusIpv6Addresses(),
			"byteplus_vpc_ipv6_address_bandwidths": ipv6_address_bandwidth.DataSourceByteplusIpv6AddressBandwidths(),
			"byteplus_network_acls":                network_acl.DataSourceByteplusNetworkAcls(),

			// ================ EBS ================
			"byteplus_volumes": volume.DataSourceByteplusVolumes(),

			// ================ NAT ================
			"byteplus_nat_gateways": nat_gateway.DataSourceByteplusNatGateways(),
			"byteplus_dnat_entries": dnat_entry.DataSourceByteplusDnatEntries(),
			"byteplus_snat_entries": snat_entry.DataSourceByteplusSnatEntries(),

			// ================ EIP ================
			"byteplus_eip_addresses": eip_address.DataSourceByteplusEipAddresses(),

			// ================ CLB ================
			"byteplus_clbs":                 clb.DataSourceByteplusClbs(),
			"byteplus_acls":                 acl.DataSourceByteplusAcls(),
			"byteplus_certificates":         certificate.DataSourceByteplusCertificates(),
			"byteplus_listeners":            listener.DataSourceByteplusListeners(),
			"byteplus_server_groups":        server_group.DataSourceByteplusServerGroups(),
			"byteplus_clb_rules":            rule.DataSourceByteplusRules(),
			"byteplus_server_group_servers": server_group_server.DataSourceByteplusServerGroupServers(),
			"byteplus_clb_zones":            clbZone.DataSourceByteplusClbZones(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// ================ ECS ================
			"byteplus_ecs_instance":                 ecs_instance.ResourceByteplusEcsInstance(),
			"byteplus_ecs_deployment_set":           ecs_deployment_set.ResourceByteplusEcsDeploymentSet(),
			"byteplus_ecs_deployment_set_associate": ecs_deployment_set_associate.ResourceByteplusEcsDeploymentSetAssociate(),

			// ================ VPC ================
			"byteplus_vpc":                        vpc.ResourceByteplusVpc(),
			"byteplus_subnet":                     subnet.ResourceByteplusSubnet(),
			"byteplus_security_group":             security_group.ResourceByteplusSecurityGroup(),
			"byteplus_security_group_rule":        security_group_rule.ResourceByteplusSecurityGroupRule(),
			"byteplus_network_interface":          network_interface.ResourceByteplusNetworkInterface(),
			"byteplus_network_interface_attach":   network_interface_attach.ResourceByteplusNetworkInterfaceAttach(),
			"byteplus_route_table":                route_table.ResourceByteplusRouteTable(),
			"byteplus_route_table_associate":      route_table_associate.ResourceByteplusRouteTableAssociate(),
			"byteplus_route_entry":                route_entry.ResourceByteplusRouteEntry(),
			"byteplus_vpc_ipv6_gateway":           ipv6_gateway.ResourceByteplusIpv6Gateway(),
			"byteplus_vpc_ipv6_address_bandwidth": ipv6_address_bandwidth.ResourceByteplusIpv6AddressBandwidth(),
			"byteplus_network_acl":                network_acl.ResourceByteplusNetworkAcl(),
			"byteplus_network_acl_associate":      network_acl_associate.ResourceByteplusNetworkAclAssociate(),

			// ================ EBS ================
			"byteplus_volume": volume.ResourceByteplusVolume(),

			// ================ NAT ================
			"byteplus_nat_gateway": nat_gateway.ResourceByteplusNatGateway(),
			"byteplus_dnat_entry":  dnat_entry.ResourceByteplusDnatEntry(),
			"byteplus_snat_entry":  snat_entry.ResourceByteplusSnatEntry(),

			// ================ EIP ================
			"byteplus_eip_address":   eip_address.ResourceByteplusEipAddress(),
			"byteplus_eip_associate": eip_associate.ResourceByteplusEipAssociate(),

			// ================ CLB ================
			"byteplus_clb":                 clb.ResourceByteplusClb(),
			"byteplus_acl":                 acl.ResourceByteplusAcl(),
			"byteplus_acl_entry":           acl_entry.ResourceByteplusAclEntry(),
			"byteplus_certificate":         certificate.ResourceByteplusCertificate(),
			"byteplus_listener":            listener.ResourceByteplusListener(),
			"byteplus_server_group":        server_group.ResourceByteplusServerGroup(),
			"byteplus_clb_rule":            rule.ResourceByteplusRule(),
			"byteplus_server_group_server": server_group_server.ResourceByteplusServerGroupServer(),
		},
		ConfigureFunc: ProviderConfigure,
	}
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	config := bp.Config{
		AccessKey:         d.Get("access_key").(string),
		SecretKey:         d.Get("secret_key").(string),
		SessionToken:      d.Get("session_token").(string),
		Region:            d.Get("region").(string),
		Endpoint:          d.Get("endpoint").(string),
		DisableSSL:        d.Get("disable_ssl").(bool),
		CustomerHeaders:   map[string]string{},
		CustomerEndpoints: map[string]string{},
		ProxyUrl:          d.Get("proxy_url").(string),
	}

	headers := d.Get("customer_headers").(string)
	if headers != "" {
		hs1 := strings.Split(headers, ",")
		for _, hh := range hs1 {
			hs2 := strings.Split(hh, ":")
			if len(hs2) == 2 {
				config.CustomerHeaders[hs2[0]] = hs2[1]
			}
		}
	}

	endpoints := d.Get("customer_endpoints").(string)
	if endpoints != "" {
		ends := strings.Split(endpoints, ",")
		for _, end := range ends {
			point := strings.Split(end, ":")
			if len(point) == 2 {
				config.CustomerEndpoints[point[0]] = point[1]
			}
		}
	}

	// get assume role
	var (
		arTrn             string
		arSessionName     string
		arPolicy          string
		arDurationSeconds int
	)

	if v, ok := d.GetOk("assume_role"); ok {
		assumeRoleList, ok := v.([]interface{})
		if !ok {
			return nil, fmt.Errorf("the assume_role is not slice ")
		}
		if len(assumeRoleList) == 1 {
			assumeRoleMap, ok := assumeRoleList[0].(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("the value of the assume_role is not map ")
			}
			arTrn = assumeRoleMap["assume_role_trn"].(string)
			arSessionName = assumeRoleMap["assume_role_session_name"].(string)
			arDurationSeconds = assumeRoleMap["duration_seconds"].(int)
			arPolicy = assumeRoleMap["policy"].(string)
		}
	} else {
		arTrn = os.Getenv("BYTEPLUS_ASSUME_ROLE_TRN")
		arSessionName = os.Getenv("BYTEPLUS_ASSUME_ROLE_SESSION_NAME")
		duration := os.Getenv("BYTEPLUS_ASSUME_ROLE_DURATION_SECONDS")
		if duration != "" {
			durationSeconds, err := strconv.Atoi(duration)
			if err != nil {
				return nil, err
			}
			arDurationSeconds = durationSeconds
		} else {
			arDurationSeconds = 3600
		}
	}

	if arTrn != "" && arSessionName != "" {
		cred, err := assumeRole(config, arTrn, arSessionName, arPolicy, arDurationSeconds)
		if err != nil {
			return nil, err
		}
		config.AccessKey = cred.AccessKeyId
		config.SecretKey = cred.SecretAccessKey
		config.SessionToken = cred.SessionToken
	}

	client, err := config.Client()
	return client, err
}

func assumeRole(c bp.Config, arTrn, arSessionName, arPolicy string, arDurationSeconds int) (*sts.Credentials, error) {
	ins := sts.NewInstance()
	if c.Region != "" {
		ins.SetRegion(c.Region)
	}
	if c.Endpoint != "" {
		ins.SetHost(c.Endpoint)
	}

	ins.Client.SetAccessKey(c.AccessKey)
	ins.Client.SetSecretKey(c.SecretKey)
	input := &sts.AssumeRoleRequest{
		RoleTrn:         arTrn,
		RoleSessionName: arSessionName,
		DurationSeconds: arDurationSeconds,
		Policy:          arPolicy,
	}
	output, statusCode, err := ins.AssumeRole(input)
	var (
		reqId  string
		errObj *base.ErrorObj
	)
	if output != nil {
		reqId = output.ResponseMetadata.RequestId
		errObj = output.ResponseMetadata.Error
	}
	if err != nil {
		return nil, fmt.Errorf("AssumeRole error, httpcode is %v and reqId is %s error is %s", statusCode, reqId, err.Error())
	}
	if errObj != nil {
		return nil, fmt.Errorf("AssumeRole error, code is %v and reqId is %s error is %s", errObj.Code, reqId, errObj.Message)
	}

	if output.Result == nil || output.Result.Credentials == nil {
		return nil, fmt.Errorf("assume role failed, result is nil")
	}

	return output.Result.Credentials, nil
}
