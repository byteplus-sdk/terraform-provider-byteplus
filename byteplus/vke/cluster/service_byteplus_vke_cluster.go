package cluster

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusVkeClusterService struct {
	Client *bp.SdkClient
}

func NewVkeClusterService(c *bp.SdkClient) *ByteplusVkeClusterService {
	return &ByteplusVkeClusterService{
		Client: c,
	}
}

func (s *ByteplusVkeClusterService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusVkeClusterService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	if filter, filterExist := condition["Filter"]; filterExist {
		if podsConfig, exist := filter.(map[string]interface{})["PodsConfig"]; exist {
			if podNetworkMode, ex := podsConfig.(map[string]interface{})["PodNetworkMode"]; ex {
				condition["Filter"].(map[string]interface{})["PodsConfig.PodNetworkMode"] = podNetworkMode
				delete(condition["Filter"].(map[string]interface{}), "PodsConfig")
			}
		}
	}

	// 适配 Conditions.Type 字段
	if filter, filterExist := condition["Filter"]; filterExist {
		if statuses, exist := filter.(map[string]interface{})["Statuses"]; exist {
			for index, status := range statuses.([]interface{}) {
				if ty, ex := status.(map[string]interface{})["ConditionsType"]; ex {
					condition["Filter"].(map[string]interface{})["Statuses"].([]interface{})[index].(map[string]interface{})["Conditions.Type"] = ty
					delete(condition["Filter"].(map[string]interface{})["Statuses"].([]interface{})[index].(map[string]interface{}), "ConditionsType")
				}
			}
		}
	}

	data, err = bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "ListClusters"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.Items", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Items is not Slice")
		}
		data, err = removeSystemTags(data)
		return data, err
	})
	if err != nil {
		return data, err
	}

	// get extra data
	for _, d := range data {
		if cluster, ok := d.(map[string]interface{}); ok {
			// 1. cluster status
			status, err := bp.ObtainSdkValue("Status.Phase", cluster)
			if err != nil {
				logger.Info("Get cluster status failed, cluster: %+v, err: %s", cluster, err.Error())
				return data, err
			}
			if !clusterReadyStatuses[status.(string)] {
				logger.Info("Cluster not ready, cluster: %+v", cluster)
				continue
			}

			// 2. get kubeconfig and eip allocation id
			clusterId := cluster["Id"].(string)
			publicAccess, err := bp.ObtainSdkValue("ClusterConfig.ApiServerPublicAccessEnabled", cluster)
			if err != nil {
				logger.Info("Get cluster public access failed, cluster: %+v, err: %s", cluster, err.Error())
				return data, err
			}
			publicIp, err := bp.ObtainSdkValue("ClusterConfig.ApiServerEndpoints.PublicIp.Ipv4", cluster)
			if err != nil || publicIp == "" {
				logger.Info("Get cluster public ip error or public ip is empty, cluster: %+v, err: %v", cluster, err)
			} else {
				if publicAccess, ok := publicAccess.(bool); ok && publicAccess {
					// a. get public kubeconfig
					publicKubeconfigResp, err := s.getKubeconfig(clusterId, "Public")
					if err != nil {
						logger.Info("Get public kubeconfig error, cluster: %+v, err: %s", cluster, err.Error())
						return data, err
					}

					kubeconfigs, err := bp.ObtainSdkValue("Result.Items", *publicKubeconfigResp)
					if err != nil {
						return data, err
					}
					if len(kubeconfigs.([]interface{})) > 0 {
						cluster["KubeconfigPublic"] = kubeconfigs.([]interface{})[0].(map[string]interface{})["Kubeconfig"]
					}

					// b. get eip data
					action := "DescribeEipAddresses"
					req := map[string]interface{}{
						"EipAddresses.1": publicIp,
					}
					eipAddressResp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
					if err != nil {
						return data, err
					}
					eipAddresses, err := bp.ObtainSdkValue("Result.EipAddresses", *eipAddressResp)
					if err != nil {
						return data, err
					}

					if eipAddresses, ok := eipAddresses.([]interface{}); !ok {
						return data, errors.New("Result.EipAddresses is not Slice")
					} else if len(eipAddresses) == 0 {
						return data, fmt.Errorf("Eip %s not found ", publicIp)
					} else {
						// get eip allocation id
						cluster["EipAllocationId"] = eipAddresses[0].(map[string]interface{})["AllocationId"].(string)

						// get eip bandwidth, billing_type, isp
						if clusterConfig, exist := cluster["ClusterConfig"]; exist {
							if apiServerPublicAccessConfig, exist := clusterConfig.(map[string]interface{})["ApiServerPublicAccessConfig"]; exist {
								if publicAccessNetworkConfig, exist := apiServerPublicAccessConfig.(map[string]interface{})["PublicAccessNetworkConfig"]; exist {
									if eipConfig, ok := publicAccessNetworkConfig.(map[string]interface{}); ok {
										eipConfig["BillingType"] = eipAddresses[0].(map[string]interface{})["BillingType"]
										eipConfig["Bandwidth"] = eipAddresses[0].(map[string]interface{})["Bandwidth"]
										eipConfig["Isp"] = eipAddresses[0].(map[string]interface{})["ISP"]
									}
								}
							}
						}
					}
				}
			}

			privateKubeconfigResp, err := s.getKubeconfig(clusterId, "Private")
			if err != nil {
				logger.Info("Get private kubeconfig error, cluster: %+v, err: %s", cluster, err.Error())
				return data, err
			}

			kubeconfigs, err := bp.ObtainSdkValue("Result.Items", *privateKubeconfigResp)
			if err != nil {
				return data, err
			}
			if len(kubeconfigs.([]interface{})) > 0 {
				cluster["KubeconfigPrivate"] = kubeconfigs.([]interface{})[0].(map[string]interface{})["Kubeconfig"]
			}
		}
	}

	return data, err
}

func (s *ByteplusVkeClusterService) getKubeconfig(clusterId, accessType string) (*map[string]interface{}, error) {
	kubeconfigReq := &map[string]interface{}{
		"Filter": map[string]interface{}{
			"ClusterIds": []string{clusterId},
			"Types":      []string{accessType},
		},
	}
	logger.Debug(logger.ReqFormat, "ListKubeconfigs", kubeconfigReq)
	return s.Client.UniversalClient.DoCall(getUniversalInfo("ListKubeconfigs"), kubeconfigReq)
}

func (s *ByteplusVkeClusterService) ReadResource(resourceData *schema.ResourceData, clusterId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if clusterId == "" {
		clusterId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Filter": map[string]interface{}{
			"Ids": []string{clusterId},
		},
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Vke Cluster %s not exist ", clusterId)
	}

	return data, err
}

func (s *ByteplusVkeClusterService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status.Phase", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf(" Vke Cluster status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}

}

func (ByteplusVkeClusterService) WithResourceResponseHandlers(cluster map[string]interface{}) []bp.ResourceResponseHandler {
	if clusterConfig, ok := cluster["ClusterConfig"]; ok {
		if apiServerPublicAccessConfig, ok := clusterConfig.(map[string]interface{})["ApiServerPublicAccessConfig"]; ok {
			if publicAccessNetworkConfig, ok := apiServerPublicAccessConfig.(map[string]interface{})["PublicAccessNetworkConfig"]; ok {
				apiServerPublicAccessConfig.(map[string]interface{})["PublicAccessNetworkConfig"] = []interface{}{publicAccessNetworkConfig}
			}
			clusterConfig.(map[string]interface{})["ApiServerPublicAccessConfig"] = []interface{}{apiServerPublicAccessConfig}
		}
	}

	if podsConfig, ok := cluster["PodsConfig"]; ok {
		if flannelConfig, ok := podsConfig.(map[string]interface{})["FlannelConfig"]; ok {
			podsConfig.(map[string]interface{})["FlannelConfig"] = []interface{}{flannelConfig}
		}
		if vpcCniConfig, ok := podsConfig.(map[string]interface{})["VpcCniConfig"]; ok {
			podsConfig.(map[string]interface{})["VpcCniConfig"] = []interface{}{vpcCniConfig}
		}
	}

	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return cluster, map[string]bp.ResponseConvert{
			"ClusterConfig": {
				TargetField: "cluster_config",
				Convert: func(i interface{}) interface{} {
					clusterConfig := i.(map[string]interface{})

					if apiServerPublicAccessConfig, ok := clusterConfig["ApiServerPublicAccessConfig"].([]interface{}); !ok {
						return i
					} else if publicAccessNetworkConfig, ok := apiServerPublicAccessConfig[0].(map[string]interface{})["PublicAccessNetworkConfig"].([]interface{}); !ok {
						return i
					} else {
						billingType := publicAccessNetworkConfig[0].(map[string]interface{})["BillingType"]
						if billingType == nil {
							return i
						}
						publicAccessNetworkConfig[0].(map[string]interface{})["BillingType"] = billingTypeResponseConvert(billingType)
					}

					return i
				},
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusVkeClusterService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateCluster",
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"cluster_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"subnet_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"api_server_public_access_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"public_access_network_config": {
									ConvertType: bp.ConvertJsonObject,
								},
							},
						},
					},
				},
				"pods_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"flannel_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"pod_cidrs": {
									ConvertType: bp.ConvertJsonArray,
								},
							},
						},
						"vpc_cni_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"subnet_ids": {
									ConvertType: bp.ConvertJsonArray,
								},
								"vpc_id": {
									Ignore: true,
								},
							},
						},
					},
				},
				"services_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"service_cidrsv4": {
							ConvertType: bp.ConvertJsonArray,
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"logging_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"log_setups": {
							ConvertType: bp.ConvertJsonObjectArray,
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if billingType, ok := (*call.SdkParam)["ClusterConfig.ApiServerPublicAccessConfig.PublicAccessNetworkConfig.BillingType"]; ok {
					realBillingType := billingTypeRequestConvert(d, billingType)
					(*call.SdkParam)["ClusterConfig.ApiServerPublicAccessConfig.PublicAccessNetworkConfig.BillingType"] = realBillingType
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//创建cluster
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				//注意 获取内容 这个地方不能是指针 需要转一次
				id, _ := bp.ObtainSdkValue("Result.Id", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}

}

func (s *ByteplusVkeClusterService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateClusterConfig",
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"cluster_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"subnet_ids": {
							ConvertType: bp.ConvertJsonArray,
						},
						"api_server_public_access_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"public_access_network_config": {
									ConvertType: bp.ConvertJsonObject,
								},
							},
						},
					},
				},
				"pods_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"flannel_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"pod_cidrs": {
									ConvertType: bp.ConvertJsonArray,
								},
							},
						},
						"vpc_cni_config": {
							ConvertType: bp.ConvertJsonObject,
							NextLevelConvert: map[string]bp.RequestConvert{
								"subnet_ids": {
									ConvertType: bp.ConvertJsonArray,
								},
								"vpc_id": {
									Ignore: true,
								},
							},
						},
					},
				},
				"services_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"service_cidrsv4": {
							ConvertType: bp.ConvertJsonArray,
						},
					},
				},
				"logging_config": {
					ConvertType: bp.ConvertJsonObject,
					NextLevelConvert: map[string]bp.RequestConvert{
						"log_setups": {
							ConvertType: bp.ConvertJsonObjectArray,
							NextLevelConvert: map[string]bp.RequestConvert{
								"log_type": {
									ConvertType: bp.ConvertDefault,
									ForceGet:    true,
								},
								"log_ttl": {
									ConvertType: bp.ConvertDefault,
									ForceGet:    true,
								},
								"enabled": {
									ConvertType: bp.ConvertDefault,
									ForceGet:    true,
								},
							},
						},
					},
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if billingType, ok := (*call.SdkParam)["ClusterConfig.ApiServerPublicAccessConfig.PublicAccessNetworkConfig.BillingType"]; ok {
					realBillingType := billingTypeRequestConvert(d, billingType)
					(*call.SdkParam)["ClusterConfig.ApiServerPublicAccessConfig.PublicAccessNetworkConfig.BillingType"] = realBillingType
				}
				(*call.SdkParam)["Id"] = d.Id()

				delete(*call.SdkParam, "Tags")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				err := validateLogSetups(d)
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//修改cluster属性
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Running"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	callbacks = append(callbacks, callback)

	if resourceData.HasChange("cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth") &&
		!resourceData.HasChange("cluster_config.0.api_server_public_access_enabled") {
		// enable public access, vke will create eip automatic
		eipAllocationId := resourceData.Get("eip_allocation_id").(string)
		modifyEipCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyEipAddressAttributes",
				ContentType: bp.ContentTypeDefault,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["AllocationId"] = eipAllocationId
					(*call.SdkParam)["Bandwidth"] = d.Get("cluster_config.0.api_server_public_access_config.0.public_access_network_config.0.bandwidth")
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					//修改eip属性
					return s.Client.UniversalClient.DoCall(getVpcUniversalInfo(call.Action), call.SdkParam)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					eip_address.NewEipAddressService(s.Client): {
						Target:     []string{"Available", "Attached", "Attaching", "Detaching"},
						Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
						ResourceId: eipAllocationId,
					},
				},
			},
		}
		callbacks = append(callbacks, modifyEipCallback)
	}

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, "Cluster", callbacks)

	return callbacks
}

func (s *ByteplusVkeClusterService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCluster",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Id"] = resourceData.Id()
				(*call.SdkParam)["RetainResources"] = []string{}
				(*call.SdkParam)["CascadingDeleteResources"] = []string{"All"}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				//删除Cluster
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				if protection, ok := d.Get("delete_protection_enabled").(bool); ok && protection {
					// 开启集群保护，直接返回失败
					return baseErr
				}

				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading cluster on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusVkeClusterService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Filter.Ids",
				ConvertType: bp.ConvertJsonArray,
			},
			"delete_protection_enabled": {
				TargetField: "Filter.DeleteProtectionEnabled",
			},
			"name": {
				TargetField: "Filter.Name",
			},
			"pods_config_pod_network_mode": {
				TargetField: "Filter.PodsConfig.PodNetworkMode",
			},
			"statuses": {
				TargetField: "Filter.Statuses",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"phase": {
						TargetField: "Phase",
					},
					"conditions_type": {
						TargetField: "ConditionsType",
					},
				},
			},
			"create_client_token": {
				TargetField: "Filter.CreateClientToken",
			},
			"update_client_token": {
				TargetField: "Filter.UpdateClientToken",
			},
			"project_name": {
				TargetField: "ProjectName",
			},
			"tags": {
				TargetField: "Tags",
				ConvertType: bp.ConvertJsonObjectArray,
			},
		},
		ContentType:  bp.ContentTypeJson,
		NameField:    "Name",
		IdField:      "Id",
		CollectField: "clusters",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ClusterConfig": {
				TargetField: "cluster_config",
				Convert: func(i interface{}) interface{} {
					realBillingType, err := bp.ObtainSdkValue("ApiServerPublicAccessConfig.PublicAccessNetworkConfig.BillingType", i)
					if err != nil {
						return i
					}
					billingType := billingTypeResponseConvert(realBillingType)

					if clusterConfig, ok := i.(map[string]interface{}); !ok {
						return i
					} else if apiServerPublicAccessConfig, ok := clusterConfig["ApiServerPublicAccessConfig"].(map[string]interface{}); !ok {
						return i
					} else if publicAccessNetworkConfig, ok := apiServerPublicAccessConfig["PublicAccessNetworkConfig"].(map[string]interface{}); !ok {
						return i
					} else {
						publicAccessNetworkConfig["BillingType"] = billingType
					}

					return i
				},
			},
		},
	}
}

func (s *ByteplusVkeClusterService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusVkeClusterService) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UntagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["TagKeys"] = make([]string, 0)
					for _, tag := range removedTags.List() {
						(*call.SdkParam)["TagKeys"] = append((*call.SdkParam)["TagKeys"].([]string), tag.(map[string]interface{})["key"].(string))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, removeCallback)

	addCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "TagResources",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["ResourceIds"] = []string{resourceData.Id()}
					(*call.SdkParam)["ResourceType"] = resourceType
					(*call.SdkParam)["Tags"] = make([]map[string]interface{}, 0)
					for _, tag := range addedTags.List() {
						(*call.SdkParam)["Tags"] = append((*call.SdkParam)["Tags"].([]map[string]interface{}), tag.(map[string]interface{}))
					}
					return true, nil
				}
				return false, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, addCallback)

	return callbacks
}

func (s *ByteplusVkeClusterService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "vke",
		ResourceType:         "cluster",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vke",
		Version:     "2022-05-12",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func getVpcUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func validateLogSetups(d *schema.ResourceData) error {
	if d.HasChange("logging_config.0.log_setups") {
		oldSet, newSet := d.GetChange("logging_config.0.log_setups")
		logger.DebugInfo("set get change", oldSet, newSet)
		oldTypeArr := make([]string, 0)
		newTypeArr := make([]string, 0)
		// 取到old和new的去重log type数组
		for _, o := range oldSet.(*schema.Set).List() {
			if oMap, ok := o.(map[string]interface{}); ok {
				if !ContainsInSlice(oldTypeArr, oMap["log_type"].(string)) {
					oldTypeArr = append(oldTypeArr, oMap["log_type"].(string))
				}
			}
		}
		for _, n := range newSet.(*schema.Set).List() {
			if nMap, ok := n.(map[string]interface{}); ok {
				if !ContainsInSlice(newTypeArr, nMap["log_type"].(string)) {
					newTypeArr = append(newTypeArr, nMap["log_type"].(string))
				}
			}
		}
		/*
			1. old数组长度大，必出现了减少 报错
			2. old数组长度小，需判断old所有type是否都在new中，如有缺失，报错
			3. old和new长度相等，需判断old和new完全相等
		*/
		if len(oldTypeArr) > len(newTypeArr) {
			return fmt.Errorf("logging setups can only be modified and added, and cannot be deleted")
		}
		if len(oldTypeArr) < len(newTypeArr) {
			for _, o := range oldTypeArr {
				if !ContainsInSlice(newTypeArr, o) {
					return fmt.Errorf("logging setups can only be modified and added, and cannot be deleted")
				}
			}
		} else {
			sort.Strings(newTypeArr)
			sort.Strings(oldTypeArr)
			if !reflect.DeepEqual(oldTypeArr, newTypeArr) {
				return fmt.Errorf("logging setups can only be modified and added, and cannot be deleted")
			}
		}
	}
	return nil
}

func ContainsInSlice(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func removeSystemTags(data []interface{}) ([]interface{}, error) {
	var (
		ok      bool
		result  map[string]interface{}
		results []interface{}
		tags    []interface{}
	)
	for _, d := range data {
		if result, ok = d.(map[string]interface{}); !ok {
			return results, errors.New("The elements in data are not map ")
		}
		tags, ok = result["Tags"].([]interface{})
		if ok {
			tags = bp.FilterSystemTags(tags)
			result["Tags"] = tags
		}
		results = append(results, result)
	}
	return results, nil
}
