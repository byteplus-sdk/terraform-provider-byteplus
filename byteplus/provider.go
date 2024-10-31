package byteplus

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/byteplusutil"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/credentials"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/session"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_activity"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_configuration"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_configuration_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_group"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_group_enabler"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_instance_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_lifecycle_hook"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/autoscaling/scaling_policy"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/bandwidth_package/bandwidth_package"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/bandwidth_package/bandwidth_package_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_certificate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_cipher_template"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_cron_job"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_cron_job_state"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_domain"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_domain_enabler"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_edge_function"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_edge_function_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_edge_function_publish"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_kv"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_kv_namespace"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cdn/cdn_service_template"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_attach_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_bandwidth_package"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_bandwidth_package_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_grant_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_inter_region_bandwidth"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_route_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/cen/cen_service_route_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/classic_cdn/classic_cdn_certificate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/classic_cdn/classic_cdn_config"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/classic_cdn/classic_cdn_domain"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/classic_cdn/classic_cdn_shared_config"
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
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ebs/volume_attach"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_available_resource"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_deployment_set_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance_state"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_instance_type"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_key_pair"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_key_pair_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/ecs_launch_template"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/image"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/region"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/ecs/zone"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_associate"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_access_key"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_login_profile"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_policy"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_role"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_role_policy_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_saml_provider"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_group_policy_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_user_policy_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/dnat_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/nat_gateway"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/nat/snat_entry"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization_account"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization_service_control_policy"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization_service_control_policy_attachment"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization_service_control_policy_enabler"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/organization/organization_unit"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint_connection"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint_service"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint_service_permission"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/privatelink/vpc_endpoint_zone"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_account"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_database"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_instance"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_instance_readonly_node"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_schema"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/addon"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/cluster"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/default_node_pool"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/default_node_pool_batch_attach"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/kubeconfig"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/node"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/node_pool"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/support_addon"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vke/support_resource_types"
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
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/customer_gateway"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/ssl_vpn_client_cert"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/ssl_vpn_server"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_connection"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_gateway"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/vpn/vpn_gateway_route"
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
			"byteplus_zones":                   zone.DataSourceByteplusZones(),
			"byteplus_images":                  image.DataSourceByteplusImages(),
			"byteplus_regions":                 region.DataSourceByteplusRegions(),
			"byteplus_ecs_instances":           ecs_instance.DataSourceByteplusEcsInstances(),
			"byteplus_ecs_deployment_sets":     ecs_deployment_set.DataSourceByteplusEcsDeploymentSets(),
			"byteplus_ecs_key_pairs":           ecs_key_pair.DataSourceByteplusEcsKeyPairs(),
			"byteplus_ecs_instance_types":      ecs_instance_type.DataSourceByteplusEcsInstanceTypes(),
			"byteplus_ecs_available_resources": ecs_available_resource.DataSourceByteplusEcsAvailableResources(),
			"byteplus_ecs_launch_templates":    ecs_launch_template.DataSourceByteplusEcsLaunchTemplates(),

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

			// ============= Bandwidth Package =============
			"byteplus_bandwidth_packages": bandwidth_package.DataSourceByteplusBandwidthPackages(),

			// ================ VKE ================
			"byteplus_vke_clusters":               cluster.DataSourceByteplusVkeVkeClusters(),
			"byteplus_vke_node_pools":             node_pool.DataSourceByteplusNodePools(),
			"byteplus_vke_nodes":                  node.DataSourceByteplusVkeNodes(),
			"byteplus_vke_addons":                 addon.DataSourceByteplusVkeAddons(),
			"byteplus_vke_support_addons":         support_addon.DataSourceByteplusVkeVkeSupportedAddons(),
			"byteplus_vke_kubeconfigs":            kubeconfig.DataSourceByteplusVkeKubeconfigs(),
			"byteplus_vke_support_resource_types": support_resource_types.DataSourceByteplusVkeVkeSupportResourceTypes(),

			// ================ AutoScaling ================
			"byteplus_scaling_groups":          scaling_group.DataSourceByteplusScalingGroups(),
			"byteplus_scaling_configurations":  scaling_configuration.DataSourceByteplusScalingConfigurations(),
			"byteplus_scaling_policies":        scaling_policy.DataSourceByteplusScalingPolicies(),
			"byteplus_scaling_lifecycle_hooks": scaling_lifecycle_hook.DataSourceByteplusScalingLifecycleHooks(),
			"byteplus_scaling_activities":      scaling_activity.DataSourceByteplusScalingActivities(),
			"byteplus_scaling_instances":       scaling_instance.DataSourceByteplusScalingInstances(),

			// ================ Postgresql ================
			"byteplus_rds_postgresql_databases": rds_postgresql_database.DataSourceByteplusRdsPostgresqlDatabases(),
			"byteplus_rds_postgresql_accounts":  rds_postgresql_account.DataSourceByteplusRdsPostgresqlAccounts(),
			"byteplus_rds_postgresql_instances": rds_postgresql_instance.DataSourceByteplusRdsPostgresqlInstances(),
			"byteplus_rds_postgresql_schemas":   rds_postgresql_schema.DataSourceByteplusRdsPostgresqlSchemas(),

			// ================ CDN ================
			"byteplus_cdn_domains":                 cdn_domain.DataSourceByteplusCdnDomains(),
			"byteplus_cdn_cipher_templates":        cdn_cipher_template.DataSourceByteplusCdnCipherTemplates(),
			"byteplus_cdn_service_templates":       cdn_service_template.DataSourceByteplusCdnServiceTemplates(),
			"byteplus_cdn_certificates":            cdn_certificate.DataSourceByteplusCdnCertificates(),
			"byteplus_cdn_edge_functions":          cdn_edge_function.DataSourceByteplusCdnEdgeFunctions(),
			"byteplus_cdn_edge_function_publishes": cdn_edge_function_publish.DataSourceByteplusCdnEdgeFunctionPublishs(),
			"byteplus_cdn_cron_jobs":               cdn_cron_job.DataSourceByteplusCdnCronJobs(),
			"byteplus_cdn_kv_namespaces":           cdn_kv_namespace.DataSourceByteplusCdnKvNamespaces(),
			"byteplus_cdn_kvs":                     cdn_kv.DataSourceByteplusCdnKvs(),

			// ================ VPN ================
			"byteplus_vpn_gateways":         vpn_gateway.DataSourceByteplusVpnGateways(),
			"byteplus_customer_gateways":    customer_gateway.DataSourceByteplusCustomerGateways(),
			"byteplus_vpn_connections":      vpn_connection.DataSourceByteplusVpnConnections(),
			"byteplus_vpn_gateway_routes":   vpn_gateway_route.DataSourceByteplusVpnGatewayRoutes(),
			"byteplus_ssl_vpn_servers":      ssl_vpn_server.DataSourceByteplusSslVpnServers(),
			"byteplus_ssl_vpn_client_certs": ssl_vpn_client_cert.DataSourceByteplusSslVpnClientCerts(),

			// ================ Classic CDN ================
			"byteplus_classic_cdn_certificates":   classic_cdn_certificate.DataSourceByteplusCdnCertificates(),
			"byteplus_classic_cdn_domains":        classic_cdn_domain.DataSourceByteplusCdnDomains(),
			"byteplus_classic_cdn_configs":        classic_cdn_config.DataSourceByteplusCdnConfigs(),
			"byteplus_classic_cdn_shared_configs": classic_cdn_shared_config.DataSourceByteplusCdnSharedConfigs(),

			// ================ Cen ================
			"byteplus_cens":                        cen.DataSourceByteplusCens(),
			"byteplus_cen_attach_instances":        cen_attach_instance.DataSourceByteplusCenAttachInstances(),
			"byteplus_cen_route_entries":           cen_route_entry.DataSourceByteplusCenRouteEntries(),
			"byteplus_cen_bandwidth_packages":      cen_bandwidth_package.DataSourceByteplusCenBandwidthPackages(),
			"byteplus_cen_inter_region_bandwidths": cen_inter_region_bandwidth.DataSourceByteplusCenInterRegionBandwidths(),
			"byteplus_cen_service_route_entries":   cen_service_route_entry.DataSourceByteplusCenServiceRouteEntries(),

			// ================ IAM ================
			"byteplus_iam_policies":                      iam_policy.DataSourceByteplusIamPolicies(),
			"byteplus_iam_roles":                         iam_role.DataSourceByteplusIamRoles(),
			"byteplus_iam_users":                         iam_user.DataSourceByteplusIamUsers(),
			"byteplus_iam_user_groups":                   iam_user_group.DataSourceByteplusIamUserGroups(),
			"byteplus_iam_user_group_policy_attachments": iam_user_group_policy_attachment.DataSourceByteplusIamUserGroupPolicyAttachments(),
			"byteplus_iam_saml_providers":                iam_saml_provider.DataSourceByteplusIamSamlProviders(),
			"byteplus_iam_access_keys":                   iam_access_key.DataSourceByteplusIamAccessKeys(),

			// ================ PrivateLink ==================
			"byteplus_privatelink_vpc_endpoints":                    vpc_endpoint.DataSourceByteplusPrivatelinkVpcEndpoints(),
			"byteplus_privatelink_vpc_endpoint_services":            vpc_endpoint_service.DataSourceByteplusPrivatelinkVpcEndpointServices(),
			"byteplus_privatelink_vpc_endpoint_service_permissions": vpc_endpoint_service_permission.DataSourceByteplusPrivatelinkVpcEndpointServicePermissions(),
			"byteplus_privatelink_vpc_endpoint_connections":         vpc_endpoint_connection.DataSourceByteplusPrivatelinkVpcEndpointConnections(),
			"byteplus_privatelink_vpc_endpoint_zones":               vpc_endpoint_zone.DataSourceByteplusPrivatelinkVpcEndpointZones(),

			// ================ Organization ================
			"byteplus_organization_units":                    organization_unit.DataSourceByteplusOrganizationUnits(),
			"byteplus_organization_service_control_policies": organization_service_control_policy.DataSourceByteplusServiceControlPolicies(),
			"byteplus_organization_accounts":                 organization_account.DataSourceByteplusOrganizationAccounts(),
			"byteplus_organizations":                         organization.DataSourceByteplusOrganizations(),
		},
		ResourcesMap: map[string]*schema.Resource{
			// ================ ECS ================
			"byteplus_ecs_instance":                 ecs_instance.ResourceByteplusEcsInstance(),
			"byteplus_ecs_instance_state":           ecs_instance_state.ResourceByteplusEcsInstanceState(),
			"byteplus_ecs_deployment_set":           ecs_deployment_set.ResourceByteplusEcsDeploymentSet(),
			"byteplus_ecs_deployment_set_associate": ecs_deployment_set_associate.ResourceByteplusEcsDeploymentSetAssociate(),
			"byteplus_ecs_key_pair":                 ecs_key_pair.ResourceByteplusEcsKeyPair(),
			"byteplus_ecs_key_pair_associate":       ecs_key_pair_associate.ResourceByteplusEcsKeyPairAssociate(),
			"byteplus_ecs_launch_template":          ecs_launch_template.ResourceByteplusEcsLaunchTemplate(),

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
			"byteplus_volume":        volume.ResourceByteplusVolume(),
			"byteplus_volume_attach": volume_attach.ResourceByteplusVolumeAttach(),

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

			// ============= Bandwidth Package =============
			"byteplus_bandwidth_package":            bandwidth_package.ResourceByteplusBandwidthPackage(),
			"byteplus_bandwidth_package_attachment": bandwidth_package_attachment.ResourceByteplusBandwidthPackageAttachment(),

			// ================ VKE ================
			"byteplus_vke_cluster":                        cluster.ResourceByteplusVkeCluster(),
			"byteplus_vke_node_pool":                      node_pool.ResourceByteplusNodePool(),
			"byteplus_vke_node":                           node.ResourceByteplusVkeNode(),
			"byteplus_vke_addon":                          addon.ResourceByteplusVkeAddon(),
			"byteplus_vke_default_node_pool":              default_node_pool.ResourceByteplusDefaultNodePool(),
			"byteplus_vke_default_node_pool_batch_attach": default_node_pool_batch_attach.ResourceByteplusDefaultNodePoolBatchAttach(),
			"byteplus_vke_kubeconfig":                     kubeconfig.ResourceByteplusVkeKubeconfig(),

			// ================ AutoScaling ================
			"byteplus_scaling_group":                    scaling_group.ResourceByteplusScalingGroup(),
			"byteplus_scaling_configuration":            scaling_configuration.ResourceByteplusScalingConfiguration(),
			"byteplus_scaling_configuration_attachment": scaling_configuration_attachment.ResourceByteplusScalingConfigurationAttachment(),
			"byteplus_scaling_policy":                   scaling_policy.ResourceByteplusScalingPolicy(),
			"byteplus_scaling_lifecycle_hook":           scaling_lifecycle_hook.ResourceByteplusScalingLifecycleHook(),
			"byteplus_scaling_group_enabler":            scaling_group_enabler.ResourceByteplusScalingGroupEnabler(),
			"byteplus_scaling_instance_attachment":      scaling_instance_attachment.ResourceByteplusScalingInstanceAttachment(),

			// ================ Postgresql ================
			"byteplus_rds_postgresql_database":               rds_postgresql_database.ResourceByteplusRdsPostgresqlDatabase(),
			"byteplus_rds_postgresql_account":                rds_postgresql_account.ResourceByteplusRdsPostgresqlAccount(),
			"byteplus_rds_postgresql_instance":               rds_postgresql_instance.ResourceByteplusRdsPostgresqlInstance(),
			"byteplus_rds_postgresql_instance_readonly_node": rds_postgresql_instance_readonly_node.ResourceByteplusRdsPostgresqlInstanceReadonlyNode(),
			"byteplus_rds_postgresql_schema":                 rds_postgresql_schema.ResourceByteplusRdsPostgresqlSchema(),

			// ================ CDN ================
			"byteplus_cdn_domain":                  cdn_domain.ResourceByteplusCdnDomain(),
			"byteplus_cdn_cipher_template":         cdn_cipher_template.ResourceByteplusCdnCipherTemplate(),
			"byteplus_cdn_service_template":        cdn_service_template.ResourceByteplusCdnServiceTemplate(),
			"byteplus_cdn_domain_enabler":          cdn_domain_enabler.ResourceByteplusCdnDomainEnabler(),
			"byteplus_cdn_certificate":             cdn_certificate.ResourceByteplusCdnCertificate(),
			"byteplus_cdn_edge_function":           cdn_edge_function.ResourceByteplusCdnEdgeFunction(),
			"byteplus_cdn_edge_function_publish":   cdn_edge_function_publish.ResourceByteplusCdnEdgeFunctionPublish(),
			"byteplus_cdn_edge_function_associate": cdn_edge_function_associate.ResourceByteplusCdnEdgeFunctionAssociate(),
			"byteplus_cdn_cron_job":                cdn_cron_job.ResourceByteplusCdnCronJob(),
			"byteplus_cdn_cron_job_state":          cdn_cron_job_state.ResourceByteplusCdnCronJobState(),
			"byteplus_cdn_kv_namespace":            cdn_kv_namespace.ResourceByteplusCdnKvNamespace(),
			"byteplus_cdn_kv":                      cdn_kv.ResourceByteplusCdnKv(),

			// ================ VPN ================
			"byteplus_vpn_gateway":         vpn_gateway.ResourceByteplusVpnGateway(),
			"byteplus_customer_gateway":    customer_gateway.ResourceByteplusCustomerGateway(),
			"byteplus_vpn_connection":      vpn_connection.ResourceByteplusVpnConnection(),
			"byteplus_vpn_gateway_route":   vpn_gateway_route.ResourceByteplusVpnGatewayRoute(),
			"byteplus_ssl_vpn_server":      ssl_vpn_server.ResourceByteplusSslVpnServer(),
			"byteplus_ssl_vpn_client_cert": ssl_vpn_client_cert.ResourceByteplusSslClientCertServer(),

			// ================ Classic CDN ================
			"byteplus_classic_cdn_certificate":   classic_cdn_certificate.ResourceByteplusCdnCertificate(),
			"byteplus_classic_cdn_domain":        classic_cdn_domain.ResourceByteplusCdnDomain(),
			"byteplus_classic_cdn_shared_config": classic_cdn_shared_config.ResourceByteplusCdnSharedConfig(),

			// ================ Cen ================
			"byteplus_cen":                             cen.ResourceByteplusCen(),
			"byteplus_cen_attach_instance":             cen_attach_instance.ResourceByteplusCenAttachInstance(),
			"byteplus_cen_grant_instance":              cen_grant_instance.ResourceByteplusCenGrantInstance(),
			"byteplus_cen_route_entry":                 cen_route_entry.ResourceByteplusCenRouteEntry(),
			"byteplus_cen_bandwidth_package":           cen_bandwidth_package.ResourceByteplusCenBandwidthPackage(),
			"byteplus_cen_bandwidth_package_associate": cen_bandwidth_package_associate.ResourceByteplusCenBandwidthPackageAssociate(),
			"byteplus_cen_inter_region_bandwidth":      cen_inter_region_bandwidth.ResourceByteplusCenInterRegionBandwidth(),
			"byteplus_cen_service_route_entry":         cen_service_route_entry.ResourceByteplusCenServiceRouteEntry(),

			// ================ IAM ================
			"byteplus_iam_policy":                       iam_policy.ResourceByteplusIamPolicy(),
			"byteplus_iam_role":                         iam_role.ResourceByteplusIamRole(),
			"byteplus_iam_role_policy_attachment":       iam_role_policy_attachment.ResourceByteplusIamRolePolicyAttachment(),
			"byteplus_iam_access_key":                   iam_access_key.ResourceByteplusIamAccessKey(),
			"byteplus_iam_user":                         iam_user.ResourceByteplusIamUser(),
			"byteplus_iam_login_profile":                iam_login_profile.ResourceByteplusIamLoginProfile(),
			"byteplus_iam_user_policy_attachment":       iam_user_policy_attachment.ResourceByteplusIamUserPolicyAttachment(),
			"byteplus_iam_user_group":                   iam_user_group.ResourceByteplusIamUserGroup(),
			"byteplus_iam_user_group_attachment":        iam_user_group_attachment.ResourceByteplusIamUserGroupAttachment(),
			"byteplus_iam_user_group_policy_attachment": iam_user_group_policy_attachment.ResourceByteplusIamUserGroupPolicyAttachment(),
			"byteplus_iam_saml_provider":                iam_saml_provider.ResourceByteplusIamSamlProvider(),

			// ================ PrivateLink ==================
			"byteplus_privatelink_vpc_endpoint":                    vpc_endpoint.ResourceByteplusPrivatelinkVpcEndpoint(),
			"byteplus_privatelink_vpc_endpoint_service":            vpc_endpoint_service.ResourceByteplusPrivatelinkVpcEndpointService(),
			"byteplus_privatelink_vpc_endpoint_service_permission": vpc_endpoint_service_permission.ResourceByteplusPrivatelinkVpcEndpointServicePermission(),
			"byteplus_privatelink_vpc_endpoint_connection":         vpc_endpoint_connection.ResourceByteplusPrivatelinkVpcEndpointConnectionService(),
			"byteplus_privatelink_vpc_endpoint_zone":               vpc_endpoint_zone.ResourceByteplusPrivatelinkVpcEndpointZone(),

			// ================ Organization ================
			"byteplus_organization_unit":                              organization_unit.ResourceByteplusOrganizationUnit(),
			"byteplus_organization_service_control_policy_enabler":    organization_service_control_policy_enabler.ResourceByteplusOrganizationServiceControlPolicyEnabler(),
			"byteplus_organization_service_control_policy":            organization_service_control_policy.ResourceByteplusServiceControlPolicy(),
			"byteplus_organization_service_control_policy_attachment": organization_service_control_policy_attachment.ResourceByteplusServiceControlPolicyAttachment(),
			"byteplus_organization_account":                           organization_account.ResourceByteplusOrganizationAccount(),
			"byteplus_organization":                                   organization.ResourceByteplusOrganization(),
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
		config.AccessKey = cred["AccessKeyId"].(string)
		config.SecretKey = cred["SecretAccessKey"].(string)
		config.SessionToken = cred["SessionToken"].(string)
	}

	client, err := config.Client()
	return client, err
}

func assumeRole(c bp.Config, arTrn, arSessionName, arPolicy string, arDurationSeconds int) (map[string]interface{}, error) {
	version := fmt.Sprintf("%s/%s", bp.TerraformProviderName, bp.TerraformProviderVersion)
	conf := byteplus.NewConfig().
		WithRegion(c.Region).
		WithExtraUserAgent(byteplus.String(version)).
		WithCredentials(credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.SessionToken)).
		WithDisableSSL(c.DisableSSL).
		WithExtendHttpRequest(func(ctx context.Context, request *http.Request) {
			if len(c.CustomerHeaders) > 0 {
				for k, v := range c.CustomerHeaders {
					request.Header.Add(k, v)
				}
			}
		}).
		WithEndpoint(byteplusutil.NewEndpoint().WithCustomerEndpoint(c.Endpoint).GetEndpoint())

	if c.ProxyUrl != "" {
		u, _ := url.Parse(c.ProxyUrl)
		t := &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		httpClient := http.DefaultClient
		httpClient.Transport = t
		httpClient.Timeout = time.Duration(30000) * time.Millisecond
	}

	sess, err := session.NewSession(conf)
	if err != nil {
		return nil, err
	}

	universalClient := bp.NewUniversalClient(sess, c.CustomerEndpoints)

	action := "AssumeRole"
	req := map[string]interface{}{
		"RoleTrn":         arTrn,
		"RoleSessionName": arSessionName,
		"DurationSeconds": arDurationSeconds,
		"Policy":          arPolicy,
	}
	resp, err := universalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return nil, fmt.Errorf("AssumeRole failed, error: %s", err.Error())
	}
	results, err := bp.ObtainSdkValue("Result.Credentials", *resp)
	if err != nil {
		return nil, err
	}
	cred, ok := results.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("AssumeRole Result.Credentials is not Map")
	}
	return cred, nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "sts",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
