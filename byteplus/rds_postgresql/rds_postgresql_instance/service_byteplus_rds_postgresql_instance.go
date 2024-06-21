package rds_postgresql_instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsPostgresqlInstanceService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsPostgresqlInstanceService(c *bp.SdkClient) *ByteplusRdsPostgresqlInstanceService {
	return &ByteplusRdsPostgresqlInstanceService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsPostgresqlInstanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsPostgresqlInstanceService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstances"
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
		// append details
		for _, v := range data {
			var (
				basicInfo    interface{}
				endpointInfo interface{}
				nodeInfo     interface{}
			)
			action = "DescribeDBInstanceDetail"
			instance := v.(map[string]interface{})

			// DescribeDBInstanceDetail
			req := map[string]interface{}{
				"InstanceId": instance["InstanceId"],
			}
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
			if err != nil {
				logger.Info("DescribeDBInstanceDetail error:", err)
				continue
			}
			respBytes, _ = json.Marshal(resp)
			logger.Debug(logger.RespFormat, action, req, string(respBytes))

			// append basic info
			basicInfo, err = bp.ObtainSdkValue("Result.BasicInfo", *resp)
			if err != nil {
				logger.Info("ObtainSdkValue Result.BasicInfo error:", err)
				continue
			}
			if basicInfoMap, ok := basicInfo.(map[string]interface{}); ok {
				instance["VCPU"] = basicInfoMap["VCPU"]
				instance["Memory"] = basicInfoMap["Memory"]
				instance["UpdateTime"] = basicInfoMap["UpdateTime"]
				instance["BackupUse"] = basicInfoMap["BackupUse"]
				instance["DataSyncMode"] = basicInfoMap["DataSyncMode"]
			}

			// append endpoint info
			endpointInfo, err = bp.ObtainSdkValue("Result.Endpoints", *resp)
			if err != nil {
				logger.Info("ObtainSdkValue Result.Endpoints error:", err)
				continue
			}
			if infos, ok := endpointInfo.([]interface{}); ok {
				instance["Endpoints"] = infos
			} else {
				// 接口返回nil
				instance["Endpoints"] = []interface{}{}
			}

			// append node info
			nodeInfo, err = bp.ObtainSdkValue("Result.Nodes", *resp)
			if err != nil {
				logger.Info("ObtainSdkValue Result.Nodes error:", err)
				continue
			}
			if infos, ok := nodeInfo.([]interface{}); ok {
				instance["Nodes"] = infos
			} else {
				// 接口返回nil
				instance["Nodes"] = []interface{}{}
			}
		}
		return data, err
	})
}

func (s *ByteplusRdsPostgresqlInstanceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"InstanceId": id,
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
		return data, fmt.Errorf("Rds PostgreSQL instance %s not exist ", id)
	}

	if nodeArr, ok := data["Nodes"].([]interface{}); ok {
		for _, node := range nodeArr {
			if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
				if nodeMap["NodeType"] == "Primary" {
					data["PrimaryZoneId"] = nodeMap["ZoneId"]
				} else if nodeMap["NodeType"] == "Secondary" {
					data["SecondaryZoneId"] = nodeMap["ZoneId"]
				}
			}
		}
	}

	// Set特殊处理
	if parameterSet, ok := resourceData.GetOk("parameters"); ok {
		data["Parameters"] = parameterSet.(*schema.Set).List()
	}

	data["ChargeInfo"] = data["ChargeDetail"]

	return data, err
}

func (s *ByteplusRdsPostgresqlInstanceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      10 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error")

			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				demo, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("InstanceStatus", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Rds PostgreSQL instance status error, status:%s ", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}
}

func (s *ByteplusRdsPostgresqlInstanceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback
	// instance callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBInstance",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"db_engine_version": {
					TargetField: "DBEngineVersion",
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"charge_info": {
					ConvertType: bp.ConvertJsonObject,
				},
				// node ignore
				"node_spec": {
					Ignore: true,
				},
				"primary_zone_id": {
					Ignore: true,
				},
				"secondary_zone_id": {
					Ignore: true,
				},
				"parameters": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				var (
					nodeInfos []interface{}
					subnets   []interface{}
					results   interface{}
					ok        bool
				)
				// add vpc id
				subnetId := d.Get("subnet_id")
				req := map[string]interface{}{
					"SubnetIds.1": subnetId,
				}
				action := "DescribeSubnets"
				resp, err := s.Client.UniversalClient.DoCall(getVPCUniversalInfo(action), &req)
				if err != nil {
					return false, err
				}
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
				vpcId := subnets[0].(map[string]interface{})["VpcId"]

				(*call.SdkParam)["VpcId"] = vpcId

				// add NodeInfo
				primaryNodeInfo := make(map[string]interface{})
				primaryNodeInfo["NodeType"] = "Primary"
				primaryNodeInfo["ZoneId"] = d.Get("primary_zone_id")
				primaryNodeInfo["NodeSpec"] = d.Get("node_spec")
				nodeInfos = append(nodeInfos, primaryNodeInfo)

				secondaryNodeInfo := make(map[string]interface{})
				secondaryNodeInfo["NodeType"] = "Secondary"
				secondaryNodeInfo["ZoneId"] = d.Get("secondary_zone_id")
				secondaryNodeInfo["NodeSpec"] = d.Get("node_spec")
				nodeInfos = append(nodeInfos, secondaryNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos

				// add StorageType
				(*call.SdkParam)["StorageType"] = "LocalSSD"

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
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
	callbacks = append(callbacks, callback)

	// parameters callback
	if _, ok := resourceData.GetOk("parameters"); ok {
		paramCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceParameters",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"parameters": {
						ConvertType: bp.ConvertJsonObjectArray,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["InstanceId"] = d.Id()
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, paramCallback)
	}
	return callbacks
}

func (ByteplusRdsPostgresqlInstanceService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"DBEngineVersion": {
				TargetField: "db_engine_version",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsPostgresqlInstanceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	// ModifyDBInstanceName
	if resourceData.HasChange("instance_name") {
		nameCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceName",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"instance_name": {
						TargetField: "InstanceNewName",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) > 0 {
						(*call.SdkParam)["InstanceId"] = d.Id()
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					common, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					if err != nil {
						return common, err
					}
					return common, nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, nameCallback)
	}

	// ModifyDBInstanceSpec
	if resourceData.HasChanges("node_spec", "storage_space") {
		instanceCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceSpec",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()

					if d.HasChange("storage_space") {
						(*call.SdkParam)["StorageType"] = "LocalSSD"
						(*call.SdkParam)["StorageSpace"] = d.Get("storage_space")
					}

					if d.HasChange("node_spec") {
						nodeInfos := make([]interface{}, 0)
						primaryNodeInfo := make(map[string]interface{})
						secondaryNodeInfo := make(map[string]interface{})

						instance, err := s.ReadResource(resourceData, d.Id())
						if err != nil {
							return false, err
						}
						if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
							for _, node := range nodeArr {
								if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
									if nodeMap["NodeType"] == "Primary" {
										primaryNodeInfo["NodeId"] = nodeMap["NodeId"]
									} else if nodeMap["NodeType"] == "Secondary" {
										secondaryNodeInfo["NodeId"] = nodeMap["NodeId"]
									} else if nodeMap["NodeType"] == "ReadOnly" {
										readonlyNodeInfo := make(map[string]interface{})
										readonlyNodeInfo["NodeId"] = nodeMap["NodeId"]
										readonlyNodeInfo["NodeType"] = nodeMap["NodeType"]
										readonlyNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
										readonlyNodeInfo["ZoneId"] = nodeMap["ZoneId"]
										nodeInfos = append(nodeInfos, readonlyNodeInfo)
									}
								}
							}
						}

						primaryNodeInfo["NodeType"] = "Primary"
						primaryNodeInfo["ZoneId"] = d.Get("primary_zone_id")
						primaryNodeInfo["NodeSpec"] = d.Get("node_spec")
						primaryNodeInfo["NodeOperateType"] = "Modify"
						nodeInfos = append(nodeInfos, primaryNodeInfo)

						secondaryNodeInfo["NodeType"] = "Secondary"
						secondaryNodeInfo["ZoneId"] = d.Get("secondary_zone_id")
						secondaryNodeInfo["NodeSpec"] = d.Get("node_spec")
						secondaryNodeInfo["NodeOperateType"] = "Modify"
						nodeInfos = append(nodeInfos, secondaryNodeInfo)

						(*call.SdkParam)["NodeInfo"] = nodeInfos
					}

					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					common, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					if err != nil {
						return common, err
					}
					return common, nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, instanceCallback)
	}

	// ModifyDBInstanceParameters
	if resourceData.HasChange("parameters") {
		modifiedParams, _, _, _ := bp.GetSetDifference("parameters", resourceData, parameterHash, false)

		parameterCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceParameters",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if modifiedParams != nil && len(modifiedParams.List()) > 0 {
						(*call.SdkParam)["InstanceId"] = d.Id()
						(*call.SdkParam)["Parameters"] = make([]map[string]interface{}, 0)
						for _, v := range modifiedParams.List() {
							paramMap, ok := v.(map[string]interface{})
							if !ok {
								return false, fmt.Errorf("Parameter is not map ")
							}
							(*call.SdkParam)["Parameters"] = append((*call.SdkParam)["Parameters"].([]map[string]interface{}), map[string]interface{}{
								"Name":  paramMap["name"],
								"Value": paramMap["value"],
							})
						}
						return true, nil
					}
					return false, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutUpdate),
				},
			},
		}
		callbacks = append(callbacks, parameterCallback)
	}

	// Tags
	callbacks = s.setResourceTags(resourceData, callbacks)

	return callbacks
}

func (s *ByteplusRdsPostgresqlInstanceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBInstance",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading rds postgre instance on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusRdsPostgresqlInstanceService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"db_engine_version": {
				TargetField: "DBEngineVersion",
			},
			"tags": {
				TargetField: "TagFilters",
				ConvertType: bp.ConvertJsonObjectArray,
			},
		},
		NameField:    "InstanceName",
		IdField:      "InstanceId",
		CollectField: "instances",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"InstanceId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"DBEngineVersion": {
				TargetField: "db_engine_version",
			},
			"IPAddress": {
				TargetField: "ip_address",
			},
			"DNSVisibility": {
				TargetField: "dns_visibility",
			},
			"VCPU": {
				TargetField: "v_cpu",
			},
		},
	}
}

func (s *ByteplusRdsPostgresqlInstanceService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusRdsPostgresqlInstanceService) setResourceTags(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveTagsFromResource",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["InstanceIds"] = []string{resourceData.Id()}
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
			Action:      "AddTagsToResource",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if addedTags != nil && len(addedTags.List()) > 0 {
					(*call.SdkParam)["InstanceIds"] = []string{resourceData.Id()}
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

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_postgresql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func (s *ByteplusRdsPostgresqlInstanceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "rds_postgresql",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
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

func (s *ByteplusRdsPostgresqlInstanceService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("charge_info.0.charge_type").(string) == "PrePaid" {
		info.NeedUnsubscribe = true
		info.Products = []string{"RDS for PostgreSQL"}
	}
	return &info, nil
}
