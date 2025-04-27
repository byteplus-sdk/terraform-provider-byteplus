package escloud_instance_v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusEscloudInstanceV2Service struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewEscloudInstanceV2Service(c *bp.SdkClient) *ByteplusEscloudInstanceV2Service {
	return &ByteplusEscloudInstanceV2Service{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusEscloudInstanceV2Service) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusEscloudInstanceV2Service) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeInstances"

		// 重新组织 Filter 的格式
		if filter, filterExist := condition["Filters"]; filterExist {
			newFilter := make([]interface{}, 0)
			for k, v := range filter.(map[string]interface{}) {
				newFilter = append(newFilter, map[string]interface{}{
					"Name":   k,
					"Values": v,
				})
			}
			condition["Filters"] = newFilter
		}
		if tags, tagsExist := condition["Tags"]; tagsExist {
			tagFilter := make(map[string]interface{})
			tagFilter["Tags"] = tags
			condition["TagFilter"] = tagFilter
			delete(condition, "Tags")
		}

		bytes, _ := json.Marshal(condition)
		logger.Debug(logger.ReqFormat, action, string(bytes))
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
		respBytes, _ := json.Marshal(resp)
		logger.Debug(logger.RespFormat, action, condition, string(respBytes))

		results, err = bp.ObtainSdkValue("Result.Instances", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Instances is not Slice")
		}

		// get instance node and plugin info
		for index, ele := range data {
			ins := ele.(map[string]interface{})
			// 只在 Running 状态下才查询
			if ins["Status"] != "Running" {
				continue
			}

			con := &map[string]interface{}{
				"InstanceId": ins["InstanceId"],
			}
			bytes, _ = json.Marshal(con)
			logger.Debug(logger.ReqFormat, "DescribeInstanceNodes", string(bytes))
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeInstanceNodes"), con)
			if err != nil {
				return data, err
			}
			respBytes, _ = json.Marshal(resp)
			logger.Debug(logger.RespFormat, "DescribeInstanceNodes", con, string(respBytes))
			results, err = bp.ObtainSdkValue("Result.Nodes", *resp)
			if err != nil {
				return data, err
			}
			if results == nil {
				results = []interface{}{}
			}
			data[index].(map[string]interface{})["Nodes"] = results

			bytes, _ = json.Marshal(con)
			logger.Debug(logger.ReqFormat, "DescribeInstancePlugins", string(bytes))
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo("DescribeInstancePlugins"), con)
			if err != nil {
				return data, err
			}
			respBytes, _ = json.Marshal(resp)
			logger.Debug(logger.RespFormat, "DescribeInstancePlugins", con, string(respBytes))
			results, err = bp.ObtainSdkValue("Result.InstancePlugins", *resp)
			if err != nil {
				return data, err
			}
			if results == nil {
				results = []interface{}{}
			}
			data[index].(map[string]interface{})["Plugins"] = results
		}

		return data, err
	})
}

func (s *ByteplusEscloudInstanceV2Service) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	req := map[string]interface{}{
		"Filters": map[string]interface{}{
			"InstanceId": []string{id},
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
		return data, fmt.Errorf("escloud_instance_v2 %s not exist ", id)
	}

	configMap, ok := data["InstanceConfiguration"].(map[string]interface{})
	if !ok {
		return data, fmt.Errorf("InstanceConfiguration is not map")
	}
	for k, v := range configMap {
		data[k] = v
	}
	if subnet, ok := configMap["Subnet"]; ok {
		data["SubnetId"] = subnet.(map[string]interface{})["SubnetId"]
	}
	if zoneId, ok := configMap["ZoneId"]; ok {
		data["ZoneIds"] = strings.Split(zoneId.(string), ",")
	}

	// 查询 configuration_code
	action := "DescribeNodeAvailableSpecs"
	con := &map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.ReqFormat, action, con)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), con)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	configurationCode, err := bp.ObtainSdkValue("Result.ConfigurationCode", *resp)
	if err != nil {
		return data, err
	}
	if configurationCode == nil {
		configurationCode = ""
	}
	data["ConfigurationCode"] = configurationCode

	// 回填 NodeSpecsAssigns & NetworkSpecs
	assigns := resourceData.Get("node_specs_assigns")
	if assigns != nil && assigns.(*schema.Set).Len() > 0 {
		data["NodeSpecsAssigns"] = assigns.(*schema.Set).List()
	}

	network := resourceData.Get("network_specs")
	if network != nil && network.(*schema.Set).Len() > 0 {
		data["NetworkSpecs"] = network.(*schema.Set).List()
	}

	return data, err
}

func (s *ByteplusEscloudInstanceV2Service) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("escloud_instance_v2 status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (ByteplusEscloudInstanceV2Service) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"EnableESPublicNetwork": {
				TargetField: "enable_es_public_network",
			},
			"EnableESPrivateNetwork": {
				TargetField: "enable_es_private_network",
			},
			"ESPublicDomain": {
				TargetField: "es_public_domain",
			},
			"ESPrivateDomain": {
				TargetField: "es_private_domain",
			},
			"ESPrivateEndpoint": {
				TargetField: "es_private_endpoint",
			},
			"ESPublicEndpoint": {
				TargetField: "es_public_endpoint",
			},
			"ESInnerEndpoint": {
				TargetField: "es_inner_endpoint",
			},
			"ESPublicIpWhitelist": {
				TargetField: "es_public_ip_whitelist",
			},
			"ESPrivateIpWhitelist": {
				TargetField: "es_private_ip_whitelist",
			},
			"EnableESPrivateDomainPublic": {
				TargetField: "enable_es_private_domain_public",
			},
			"CPU": {
				TargetField: "cpu",
			},
			"VPC": {
				TargetField: "vpc",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusEscloudInstanceV2Service) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateInstanceInOneStep",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"zone_ids": {
					Ignore: true,
				},
				"enable_https": {
					TargetField: "EnableHttps",
					ForceGet:    true,
				},
				"network_specs": {
					TargetField: "NetworkSpecs",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"node_specs_assigns": {
					TargetField: "NodeSpecsAssigns",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"storage_spec_name": {
							ForceGet: true,
						},
						"storage_size": {
							ForceGet: true,
						},
						"extra_performance": {
							ForceGet:    true,
							ConvertType: bp.ConvertJsonObject,
						},
					},
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var (
					results interface{}
					subnets []interface{}
					vpcs    []interface{}
					ok      bool
				)

				// zone id
				zoneIdsArr := d.Get("zone_ids").(*schema.Set).List()
				zoneIds := make([]string, 0)
				for _, id := range zoneIdsArr {
					zoneIds = append(zoneIds, id.(string))
				}
				zoneId := strings.Join(zoneIds, ",")
				(*call.SdkParam)["ZoneId"] = zoneId

				// region & vpc & subnet
				subnetId := (*call.SdkParam)["SubnetId"]
				req := map[string]interface{}{
					"SubnetIds.1": subnetId,
				}
				action := "DescribeSubnets"
				resp, err := s.Client.UniversalClient.DoCall(getVPCUniversalInfo(action), &req)
				if err != nil {
					return false, err
				}
				logger.Debug(logger.RespFormat, action, req, *resp)
				results, err = bp.ObtainSdkValue("Result.Subnets", *resp)
				if err != nil {
					return false, err
				}
				if results == nil {
					results = []interface{}{}
				}
				if subnets, ok = results.([]interface{}); !ok {
					return false, errors.New("Result.Subnets is not Slice")
				}
				if len(subnets) == 0 {
					return false, fmt.Errorf("subnet %s not exist", subnetId.(string))
				}
				subnetName := subnets[0].(map[string]interface{})["SubnetName"]
				vpcId := subnets[0].(map[string]interface{})["VpcId"]

				req = map[string]interface{}{
					"VpcIds.1": vpcId,
				}
				action = "DescribeVpcs"
				resp, err = s.Client.UniversalClient.DoCall(getVPCUniversalInfo(action), &req)
				if err != nil {
					return false, err
				}
				logger.Debug(logger.RespFormat, action, req, *resp)
				results, err = bp.ObtainSdkValue("Result.Vpcs", *resp)
				if err != nil {
					return false, err
				}
				if results == nil {
					results = []interface{}{}
				}
				if vpcs, ok = results.([]interface{}); !ok {
					return false, errors.New("Result.Vpcs is not Slice")
				}
				if len(vpcs) == 0 {
					return false, fmt.Errorf("vpc %s not exist", subnetId.(string))
				}
				vpcName := vpcs[0].(map[string]interface{})["VpcName"]

				(*call.SdkParam)["VPC"] = map[string]interface{}{
					"VpcId":   vpcId,
					"VpcName": vpcName,
				}
				(*call.SdkParam)["Subnet"] = map[string]interface{}{
					"SubnetId":   subnetId,
					"SubnetName": subnetName,
				}
				(*call.SdkParam)["RegionId"] = s.Client.Region

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				instanceConfig := *call.SdkParam
				param := make(map[string]interface{})
				param["InstanceConfiguration"] = instanceConfig
				*call.SdkParam = param
				(*call.SdkParam)["Tags"] = instanceConfig["Tags"]
				(*call.SdkParam)["ClientToken"] = uuid.New().String()

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.InstanceId", *resp)
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

func (s *ByteplusEscloudInstanceV2Service) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if resourceData.HasChange("instance_name") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "RenameInstance",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"instance_name": {
						TargetField: "NewName",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChanges("maintenance_time", "maintenance_day") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyMaintenanceSetting",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"maintenance_time": {
						TargetField: "MaintenanceTime",
						ForceGet:    true,
					},
					"maintenance_day": {
						TargetField: "MaintenanceDay",
						ConvertType: bp.ConvertJsonArray,
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("admin_password") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ResetAdminPassword",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"admin_password": {
						TargetField: "NewPassword",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["UserName"] = "admin"
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("node_specs_assigns") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyNodeSpecInOneStep",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"node_specs_assigns": {
						ConvertType: bp.ConvertJsonObjectArray,
						ForceGet:    true,
						NextLevelConvert: map[string]bp.RequestConvert{
							"storage_spec_name": {
								ForceGet: true,
							},
							"storage_size": {
								ForceGet: true,
							},
							"extra_performance": {
								ForceGet:    true,
								ConvertType: bp.ConvertJsonObject,
							},
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)

					// 异步任务，等待 5s
					time.Sleep(5 * time.Second)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("charge_type") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyChargeCode",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"charge_type": {
						TargetField: "ToChargeType",
						ForceGet:    true,
					},
					"auto_renew": {
						TargetField: "AutoRenew",
						ForceGet:    true,
					},
					"period": {
						TargetField: "IncludeMonths",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					old, _ := d.GetChange("charge_type")
					if old == "PrePaid" {
						return false, fmt.Errorf("The operation is not permitted due to the instance charge type is prepaid. ")
					}

					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("deletion_protection") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDeletionProtection",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"deletion_protection": {
						TargetField: "DeletionProtection",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, "instance", callbacks)

	return callbacks
}

func (s *ByteplusEscloudInstanceV2Service) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ReleaseInstance",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				// 开启删除保护时，跳过 CallError
				if d.Get("deletion_protection").(bool) {
					return baseErr
				}

				//出现错误后重试
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading ESCloud instance v2 on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusEscloudInstanceV2Service) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Filters.InstanceId",
				ConvertType: bp.ConvertJsonArray,
			},
			"statuses": {
				TargetField: "Filters.Status",
				ConvertType: bp.ConvertJsonArray,
			},
			"charge_types": {
				TargetField: "Filters.ChargeType",
				ConvertType: bp.ConvertJsonArray,
			},
			"instance_names": {
				TargetField: "Filters.InstanceName",
				ConvertType: bp.ConvertJsonArray,
			},
			"versions": {
				TargetField: "Filters.Version",
				ConvertType: bp.ConvertJsonArray,
			},
			"zone_ids": {
				TargetField: "Filters.ZoneId",
				ConvertType: bp.ConvertJsonArray,
			},
			"tags": {
				TargetField: "Tags",
				ConvertType: bp.ConvertJsonObjectArray,
				NextLevelConvert: map[string]bp.RequestConvert{
					"key": {
						TargetField: "Key",
					},
					"values": {
						TargetField: "Values",
						ConvertType: bp.ConvertJsonArray,
					},
				},
			},
		},
		ContentType:  bp.ContentTypeJson,
		IdField:      "InstanceId",
		CollectField: "instances",
		ResponseConverts: map[string]bp.ResponseConvert{
			"InstanceId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"EnableESPublicNetwork": {
				TargetField: "enable_es_public_network",
			},
			"EnableESPrivateNetwork": {
				TargetField: "enable_es_private_network",
			},
			"ESPublicDomain": {
				TargetField: "es_public_domain",
			},
			"ESPrivateDomain": {
				TargetField: "es_private_domain",
			},
			"ESPrivateEndpoint": {
				TargetField: "es_private_endpoint",
			},
			"ESPublicEndpoint": {
				TargetField: "es_public_endpoint",
			},
			"ESInnerEndpoint": {
				TargetField: "es_inner_endpoint",
			},
			"ESPublicIpWhitelist": {
				TargetField: "es_public_ip_whitelist",
			},
			"ESPrivateIpWhitelist": {
				TargetField: "es_private_ip_whitelist",
			},
			"EnableESPrivateDomainPublic": {
				TargetField: "enable_es_private_domain_public",
			},
			"CPU": {
				TargetField: "cpu",
			},
			"VPC": {
				TargetField: "vpc",
			},
		},
	}
}

func (s *ByteplusEscloudInstanceV2Service) setResourceTags(resourceData *schema.ResourceData, resourceType string, callbacks []bp.Callback) []bp.Callback {
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
					for _, v := range addedTags.List() {
						tagMap, ok := v.(map[string]interface{})
						if !ok {
							return false, fmt.Errorf("Tags is not map ")
						}
						tag := make(map[string]interface{})
						tag["Key"] = tagMap["key"]
						tag["Value"] = tagMap["value"]
						(*call.SdkParam)["Tags"] = append((*call.SdkParam)["Tags"].([]map[string]interface{}), tag)
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

func (s *ByteplusEscloudInstanceV2Service) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusEscloudInstanceV2Service) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "ESCloud",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusEscloudInstanceV2Service) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("charge_type").(string) == "PrePaid" {
		info.NeedUnsubscribe = true
		info.Products = []string{"ESCloud"}
	}
	return &info, nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ESCloud",
		Version:     "2023-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func getVPCUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
