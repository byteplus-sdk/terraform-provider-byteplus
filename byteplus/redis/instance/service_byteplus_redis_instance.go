package instance

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisDbInstanceService struct {
	Client *bp.SdkClient
}

func NewRedisDbInstanceService(c *bp.SdkClient) *ByteplusRedisDbInstanceService {
	return &ByteplusRedisDbInstanceService{
		Client: c,
	}
}

func (s *ByteplusRedisDbInstanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisDbInstanceService) readInstanceDetails(id string) (instance interface{}, err error) {
	action := "DescribeDBInstanceDetail"
	cond := map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	logger.Debug(logger.RespFormat, action, *resp)
	if err != nil {
		return instance, err
	}

	instance, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return instance, err
	}
	if instance == nil {
		return instance, fmt.Errorf("instance %s is not exist", id)
	}
	return instance, err
}

func (s *ByteplusRedisDbInstanceService) readInstanceBandwidthPerShard(id string) (instance interface{}, err error) {
	action := "DescribeDBInstanceBandwidthPerShard"
	cond := map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	logger.Debug(logger.RespFormat, action, *resp)
	if err != nil {
		return instance, err
	}

	instance, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return instance, err
	}
	if instance == nil {
		return instance, fmt.Errorf("instance %s is not exist", id)
	}
	return instance, err
}

//
//func (s *ByteplusRedisDbInstanceService) readInstanceAclCategories(id string) (instance interface{}, err error) {
//	action := "DescribeDBInstanceAclCategories"
//	cond := map[string]interface{}{
//		"InstanceId": id,
//	}
//	logger.Debug(logger.RespFormat, action, cond)
//	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
//	logger.Debug(logger.RespFormat, action, *resp)
//	if err != nil {
//		return instance, err
//	}
//
//	instance, err = bp.ObtainSdkValue("Result", *resp)
//	if err != nil {
//		return instance, err
//	}
//	if instance == nil {
//		return instance, fmt.Errorf("instance %s is not exist", id)
//	}
//	return instance, err
//}

func (s *ByteplusRedisDbInstanceService) readInstanceShards(id string) (instance interface{}, err error) {
	var (
		resp *map[string]interface{}
		ok   bool
	)
	cond := map[string]interface{}{
		"InstanceId": id,
	}

	action := "DescribeDBInstanceShards"
	pageCall := func(condition map[string]interface{}) (data []interface{}, err error) {
		logger.Debug(logger.RespFormat, action, condition)
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
		logger.Debug(logger.RespFormat, action, condition, *resp)

		paramsResults, err := bp.ObtainSdkValue("Result.InstanceShards", *resp)
		if err != nil {
			return data, err
		}
		if paramsResults == nil {
			paramsResults = []interface{}{}
		}
		if data, ok = paramsResults.([]interface{}); !ok {
			return data, errors.New("Results.Params is not slice")
		}
		return data, nil
	}
	params, err := bp.WithPageNumberQuery(cond, "PageSize", "PageNumber", 100, 1, pageCall)
	if err != nil {
		return params, err
	}
	return params, nil
}

func (s *ByteplusRedisDbInstanceService) readInstanceParams(id string) (params interface{}, err error) {
	var (
		resp *map[string]interface{}
		ok   bool
	)
	cond := map[string]interface{}{
		"InstanceId": id,
	}

	action := "DescribeDBInstanceParams"
	pageCall := func(condition map[string]interface{}) (data []interface{}, err error) {
		logger.Debug(logger.RespFormat, action, condition)
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
		logger.Debug(logger.RespFormat, action, condition, *resp)

		paramsResults, err := bp.ObtainSdkValue("Result.Params", *resp)
		if err != nil {
			return data, err
		}
		if paramsResults == nil {
			paramsResults = []interface{}{}
		}
		if data, ok = paramsResults.([]interface{}); !ok {
			return data, errors.New("Results.Params is not slice")
		}
		return data, nil
	}
	params, err = bp.WithPageNumberQuery(cond, "PageSize", "PageNumber", 100, 1, pageCall)
	if err != nil {
		return params, err
	}
	return params, nil
}

func (s *ByteplusRedisDbInstanceService) readInstanceBackupPlan(id string) (backupPlan interface{}, err error) {
	cond := map[string]interface{}{
		"InstanceId": id,
	}

	action := "DescribeBackupPlan"
	logger.Debug(logger.RespFormat, action, cond)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	if err != nil {
		return backupPlan, err
	}
	logger.Debug(logger.RespFormat, action, *resp)

	backupPlan, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return backupPlan, err
	}
	logger.DebugInfo("backup plan:%v", backupPlan)

	return backupPlan, err
}

func (s *ByteplusRedisDbInstanceService) readInstanceNodeIds(id string) (nodeIds interface{}, err error) {
	cond := map[string]interface{}{
		"InstanceId": id,
	}

	action := "DescribeNodeIds"
	logger.Debug(logger.RespFormat, action, cond)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	if err != nil {
		return nodeIds, err
	}
	logger.Debug(logger.RespFormat, action, cond, *resp)

	nodeIds, err = bp.ObtainSdkValue("Result.NodeIds", *resp)
	if err != nil {
		return nodeIds, err
	}
	logger.DebugInfo("node ids:%v", nodeIds)

	return nodeIds, err
}

//func idInTargets(id string, ids interface{}) bool {
//	for _, targetId := range ids.([]interface{}) {
//		if id == targetId.(string) {
//			return true
//		}
//	}
//	return false
//}

func (s *ByteplusRedisDbInstanceService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	regionId := s.Client.Region
	condition["RegionId"] = regionId

	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat, action, condition)
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
		if err != nil {
			return data, err
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.Instances", *resp)
		if err != nil {
			logger.DebugInfo("bp.ObtainSdkValue return :%v", err)
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		instances, ok := results.([]interface{})
		if !ok {
			return data, fmt.Errorf("DescribeDBInstances responsed instances is not a slice")
		}

		for _, ele := range instances {
			ins := ele.(map[string]interface{})
			instanceId := ins["InstanceId"].(string)

			ins["ShardCapacity"] = ins["ShardCapacity"].(float64) * 1024

			params, err := s.readInstanceParams(instanceId)
			if err != nil {
				return data, err
			}
			ins["Params"] = params

			// 单节点实例不支持查询 Backup plan
			if nodeNumber, exist := ins["NodeNumber"]; exist && nodeNumber.(float64) > 1 {
				backupPlan, err := s.readInstanceBackupPlan(instanceId)
				if err != nil {
					return data, err
				}
				ins["BackupPlan"] = backupPlan
			}

			nodeIds, err := s.readInstanceNodeIds(instanceId)
			if err != nil {
				return data, err
			}
			ins["NodeIds"] = nodeIds

			detail, err := s.readInstanceDetails(instanceId)
			if err != nil {
				return data, err
			}
			ins["DeletionProtection"] = detail.(map[string]interface{})["DeletionProtection"]
			ins["MaintenanceTime"] = detail.(map[string]interface{})["MaintenanceTime"]
			ins["SubnetId"] = detail.(map[string]interface{})["SubnetId"]
			ins["VisitAddrs"] = detail.(map[string]interface{})["VisitAddrs"]
			ins["VpcAuthMode"] = detail.(map[string]interface{})["VpcAuthMode"]
			if nodes, ok := detail.(map[string]interface{})["ConfigureNodes"]; ok {
				ins["ConfigureNodes"] = nodes
			}

			instanceClass, ok := detail.(map[string]interface{})["InstanceClass"]
			if ok {
				ins["InstanceClass"] = instanceClass
			}

			maxConnections, ok := detail.(map[string]interface{})["MaxConnections"]
			if ok {
				ins["MaxConnections"] = maxConnections
			}

			dataLayout, ok := detail.(map[string]interface{})["DataLayout"]
			if ok {
				ins["DataLayout"] = dataLayout
			}

			shardCapacityV2, ok := detail.(map[string]interface{})["ShardCapacityV2"]
			if ok {
				ins["ShardCapacityV2"] = shardCapacityV2
			}

			maintenanceTime, ok := detail.(map[string]interface{})["MaintenanceTime"]
			if ok {
				ins["TimeScope"] = maintenanceTime
			}
			instanceShards, err := s.readInstanceShards(instanceId)
			if err != nil {
				return data, err
			}
			ins["InstanceShards"] = instanceShards

			bandwidthPerShard, err := s.readInstanceBandwidthPerShard(instanceId)
			if err != nil {
				return data, err
			}
			if defaultBandwidthPerShard, ok := bandwidthPerShard.(map[string]interface{})["DefaultBandwidthPerShard"]; ok {
				ins["DefaultBandwidthPerShard"] = defaultBandwidthPerShard
			}
			if additionalBandwidthPerShard, ok := bandwidthPerShard.(map[string]interface{})["AdditionalBandwidthPerShard"]; ok {
				ins["AdditionalBandwidthPerShard"] = additionalBandwidthPerShard
			}

			//instanceAclCategories, err := s.readInstanceAclCategories(instanceId)
			//if err != nil {
			//	return data, err
			//}
			//if categories, ok := instanceAclCategories.(map[string]interface{})["Categories"]; ok {
			//	ins["Categories"] = categories
			//}

			//data = append(data, ins)
		}
		data = instances
		return data, nil
	})
}

func (s *ByteplusRedisDbInstanceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
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
		var instanceMap map[string]interface{}
		if instanceMap, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
		if id == instanceMap["InstanceId"].(string) {
			data = instanceMap
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("instance %s is not exist", id)
	}

	if addrs, ok := data["VisitAddrs"].([]interface{}); ok {
		for _, address := range addrs {
			addr := address.(map[string]interface{})
			if addr["AddrType"].(string) == "Private" {
				data["Port"], _ = strconv.Atoi(addr["Port"].(string))
				break
			}
		}
	}

	if backupPlan, exist := data["BackupPlan"]; exist {
		if backupMap, ok := backupPlan.(map[string]interface{}); ok {
			data["BackupHour"] = backupMap["BackupHour"]
			data["BackupActive"] = backupMap["Active"]
			data["BackupPeriod"] = backupMap["Period"]
		}
	}

	if parameterSet, ok := resourceData.GetOk("param_values"); ok {
		data["ParamValues"] = parameterSet.(*schema.Set).List()
	}

	// 接口返回会乱序，所以这里只能兼容处理一下
	if nodes, ok := data["ConfigureNodes"]; ok {
		configNodes, ok := resourceData.GetOk("configure_nodes")
		if ok {
			if len(nodes.([]interface{})) == len(configNodes.([]interface{})) {
				// 数量相等则用本地的
				data["ConfigureNodes"] = configNodes
			}
		}
	}

	if visitAddrs, exist := data["VisitAddrs"]; exist {
		visitAddr, ok := visitAddrs.([]interface{})
		if ok {
			for _, addr := range visitAddr {
				addrMap := addr.(map[string]interface{})
				addrType, ok := resourceData.Get("addr_type").(string)
				if ok {
					if addrMap["AddrType"].(string) == addrType {
						port := addrMap["Port"]
						portString, ok := port.(string)
						if ok {
							newport, err := strconv.Atoi(portString)
							if err != nil {
								data["NewPort"] = newport
							}
						}
						break
					}
				}

			}
		}
	}

	if visitAddrs, exist := data["VisitAddrs"]; exist {
		visitAddr, ok := visitAddrs.([]interface{})
		if ok {
			for _, addr := range visitAddr {
				addrMap := addr.(map[string]interface{})
				addrType, ok := resourceData.Get("addr_type").(string)
				if ok {
					if addrMap["AddrType"].(string) == addrType {
						address := addrMap["Address"]
						addressString, ok := address.(string)
						if ok {
							regionEndPointVolces := ".redis." + s.Client.Region + ".volces.com"
							regionEndPoinvIVolces := ".redis." + s.Client.Region + ".ivolces.com"
							targets := []string{".redis.ivolces.com", ".redis.volces.com", regionEndPointVolces, regionEndPoinvIVolces}
							for _, target := range targets {
								addressPrefix := getPrefixBeforeEndPoint(addressString, target)
								if addressPrefix != "" {
									data["NewAddressPrefix"] = addressPrefix
								}
							}
						}
						break
					}
				}

			}
		}
	}

	return data, err
}

func (s *ByteplusRedisDbInstanceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      10 * time.Second,
		Pending:    []string{},
		Target:     target,
		Timeout:    timeout,
		MinTimeout: 1 * time.Second,

		Refresh: func() (result interface{}, state string, err error) {
			var (
				instance   map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "CreateFailed", "TaskFailed")

			logger.DebugInfo("start refresh :%s", id)
			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				instance, err = s.ReadResource(resourceData, id)
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
			logger.DebugInfo("Refresh instance status resp: %v", instance)

			status, err = bp.ObtainSdkValue("Status", instance)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("instance %s status error, status %s", id, status.(string))
				}
			}
			logger.DebugInfo("refresh status:%v", status)
			return instance, status.(string), err
		},
	}
}

func (s *ByteplusRedisDbInstanceService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, map[string]bp.ResponseConvert{
			"MultiAZ": {
				TargetField: "multi_az",
			},
			"AZ": {
				TargetField: "az",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRedisDbInstanceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	instanceCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBInstance",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"multi_az": {
					TargetField: "MultiAZ",
				},
				"configure_nodes": {
					TargetField: "ConfigureNodes",
					ConvertType: bp.ConvertJsonObjectArray,
					NextLevelConvert: map[string]bp.RequestConvert{
						"az": {
							TargetField: "AZ",
						},
					},
				},
				"param_values": {
					Ignore: true,
				},
				"vpc_auth_mode": {
					Ignore: true,
				},
				"backup_period": {
					Ignore: true,
				},
				"backup_hour": {
					Ignore: true,
				},
				"backup_active": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if nodeNumber, ok := (*call.SdkParam)["NodeNumber"]; ok && nodeNumber.(int) == 1 {
					// 单节点实例不支持设置 backup plan
					_, exist1 := d.GetOkExists("backup_hour")
					_, exist2 := d.GetOkExists("backup_active")
					period := d.Get("backup_period")
					periodSet, ok := period.(*schema.Set)
					if !ok {
						return false, fmt.Errorf("backup_period is not set ")
					}
					if exist1 || exist2 || periodSet.Len() > 0 {
						return false, fmt.Errorf("The single node instance cannot specify any fields related to backup plan, including `backup_period`, `backup_hour` and `backup_active`. ")
					}
				}

				if _, ok := (*call.SdkParam)["ShardedCluster"]; !ok {
					(*call.SdkParam)["ShardedCluster"] = 0
				}
				// describe subnet
				subnetId := d.Get("subnet_id")
				vpcId, _, err := s.getVpcIdAndZoneIdBySubnet(subnetId.(string))
				if err != nil {
					return false, fmt.Errorf("get vpc ID by subnet ID %s failed", subnetId)
				}
				(*call.SdkParam)["VpcId"] = vpcId
				(*call.SdkParam)["RegionId"] = s.Client.Region
				(*call.SdkParam)["ClientToken"] = uuid.New().String()

				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, *resp)
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
	callbacks = append(callbacks, instanceCallback)

	// parameters
	parameterCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceParams",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"param_values": {
					TargetField: "ParamValues",
					ConvertType: bp.ConvertJsonObjectArray,
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) == 0 {
					return false, nil
				}
				(*call.SdkParam)["InstanceId"] = d.Id()
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
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
	callbacks = append(callbacks, parameterCallback)

	// vpc_auth_mode
	vpcAuthModeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceVpcAuthMode",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"vpc_auth_mode": {
					TargetField: "VpcAuthMode",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if len(*call.SdkParam) == 0 {
					return false, nil
				}
				(*call.SdkParam)["InstanceId"] = d.Id()
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
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
	callbacks = append(callbacks, vpcAuthModeCallback)

	// backup plan
	if nodeNumber := resourceData.Get("node_number"); nodeNumber.(int) > 1 {
		backupPlanCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyBackupPlan",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					active, exist := d.GetOkExists("backup_active")
					if !exist {
						active = true
					}
					(*call.SdkParam)["Active"] = active

					period := d.Get("backup_period")
					periodSet, ok := period.(*schema.Set)
					if !ok {
						return false, fmt.Errorf("backup_period is not set ")
					}
					if periodSet.Len() > 0 {
						(*call.SdkParam)["Period"] = periodSet.List()
					} else {
						(*call.SdkParam)["Period"] = []interface{}{1, 2, 3, 4, 5, 6, 7}
					}

					(*call.SdkParam)["BackupHour"] = d.Get("backup_hour")
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, backupPlanCallback)
	}

	return callbacks
}

func (s *ByteplusRedisDbInstanceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	// 调用 EnableShardedCluster 接口将目标 Redis 实例变更为启用分片集群实例。
	if resourceData.HasChange("sharded_cluster") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "EnableShardedCluster",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"create_backup": {
						TargetField: "CreateBackup",
						ForceGet:    true,
					},
					"backup_point_name": {
						TargetField: "BackupPointName",
						ForceGet:    true,
					},
					"sharded_cluster": {
						TargetField: "ShardedCluster",
					},
					"apply_immediately": {
						TargetField: "ApplyImmediately",
						ForceGet:    true,
					},
					"shard_number": {
						TargetField: "ShardNumber",
						ForceGet:    true,
					},
					"shard_capacity": {
						TargetField: "ShardCapacity",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					nodeNumber := d.Get("node_number").(int)
					if nodeNumber > 1 {
						(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
					}
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChanges("multi_az", "configure_nodes") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceAZConfigure",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"multi_az": {
						TargetField: "MultiAZ",
						ForceGet:    true,
					},
					"create_backup": {
						TargetField: "CreateBackup",
						ForceGet:    true,
					},
					"backup_point_name": {
						TargetField: "BackupPointName",
						ForceGet:    true,
					},
					"configure_nodes": {
						TargetField: "ConfigureNodes",
						ForceGet:    true,
						ConvertType: bp.ConvertJsonObjectArray,
						NextLevelConvert: map[string]bp.RequestConvert{
							"az": {
								TargetField: "AZ",
								ForceGet:    true,
							},
						},
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if d.HasChange("node_number") && d.HasChange("configure_nodes") &&
						!d.HasChange("multi_az") && d.Get("multi_az").(string) != "disabled" {
						// 特殊情况，node number 已经修改了configure nodes，并且multi_az没有发生变更，那么就不执行这个callback，避免报错
						// 单可用区情况要允许修改configure nodes
						return false, nil
					}
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					_, ok1 := d.GetOkExists("multi_az")
					_, ok2 := d.GetOkExists("configure_nodes")
					// 这俩字段是必填字段，即使关闭多可用区，configure nodes也需要传入node number对应的相同az
					if !ok1 || !ok2 {
						return false, fmt.Errorf("MultiAZ and ConfigureNodes are required parameters ")
					}
					apply := d.Get("apply_immediately").(bool)
					nodeNumber := d.Get("node_number").(int)
					if nodeNumber > 1 {
						(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
					}
					(*call.SdkParam)["ApplyImmediately"] = apply
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("instance_name") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceName",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"instance_name": {
						TargetField: "InstanceName",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("deletion_protection") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceDeletionProtectionPolicy",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"deletion_protection": {
						TargetField: "DeletionProtection",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("vpc_auth_mode") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceVpcAuthMode",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"vpc_auth_mode": {
						TargetField: "VpcAuthMode",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("subnet_id") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceSubnet",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"subnet_id": {
						TargetField: "SubnetId",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}

					// describe subnet
					subnetId := (*call.SdkParam)["SubnetId"]
					vpcId, _, err := s.getVpcIdAndZoneIdBySubnet(subnetId.(string))
					if err != nil {
						return false, fmt.Errorf("get vpc ID by subnet ID %s failed", subnetId)
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["VpcId"] = vpcId
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("charge_type") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceChargeType",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"charge_type": {
						TargetField: "ChargeType",
					},
					"purchase_months": {
						TargetField: "PurchaseMonths",
					},
					"auto_renew": {
						TargetField: "AutoRenew",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					if chargeType, ok := (*call.SdkParam)["ChargeType"]; ok {
						(*call.SdkParam)["ChargeType"] = chargeType

					}
					(*call.SdkParam)["InstanceIds"] = []interface{}{d.Id()}
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("param_values") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceParams",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"param_values": {
						TargetField: "ParamValues",
						ConvertType: bp.ConvertJsonObjectArray,
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChanges("backup_period", "backup_hour", "backup_active") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyBackupPlan",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"backup_period": {
						TargetField: "Period",
						ConvertType: bp.ConvertJsonArray,
						ForceGet:    true,
					},
					"backup_hour": {
						TargetField: "BackupHour",
						ForceGet:    true,
					},
					"backup_active": {
						TargetField: "Active",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					logger.DebugInfo("call.sdk param:%v", call.SdkParam)
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("node_number") {
		// ModifyDBInstanceNodeNumber 接口废弃，使用IncreaseDBInstanceNodeNumber和DecreaseDBInstanceNodeNumber代替
		//callback := bp.Callback{
		//	Call: bp.SdkCall{
		//		Action:      "ModifyDBInstanceNodeNumber",
		//		ConvertMode: bp.RequestConvertInConvert,
		//		ContentType: bp.ContentTypeJson,
		//		Convert: map[string]bp.RequestConvert{
		//			"node_number": {
		//				TargetField: "NodeNumber",
		//			},
		//			"apply_immediately": {
		//				TargetField: "ApplyImmediately",
		//				ForceGet:    true,
		//			},
		//			"create_backup": {
		//				TargetField: "CreateBackup",
		//				ForceGet:    true,
		//			},
		//		},
		//		BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
		//			if len(*call.SdkParam) == 0 {
		//				return false, nil
		//			}
		//			(*call.SdkParam)["InstanceId"] = d.Id()
		//			(*call.SdkParam)["ClientToken"] = uuid.New().String()
		//			return true, nil
		//		},
		//		ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
		//			logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
		//			return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
		//		},
		//		Refresh: &bp.StateRefresh{
		//			Target:  []string{"Running"},
		//			Timeout: resourceData.Timeout(schema.TimeoutUpdate),
		//		},
		//	},
		//}
		var action string
		oldNum, newNum := resourceData.GetChange("node_number")
		if oldNum.(int) > newNum.(int) {
			action = "DecreaseDBInstanceNodeNumber"
		} else {
			action = "IncreaseDBInstanceNodeNumber"
		}
		changeNum := abs(oldNum.(int) - newNum.(int))
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      action,
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"apply_immediately": {
						TargetField: "ApplyImmediately",
						ForceGet:    true,
					},
					"create_backup": {
						TargetField: "CreateBackup",
						ForceGet:    true,
					},
					"backup_point_name": {
						TargetField: "BackupPointName",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					nodeNumber := d.Get("node_number").(int)
					if nodeNumber > 1 {
						(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
					}
					multiAZ := d.Get("multi_az").(string)
					oldNodeList, newNodeList := d.GetChange("configure_nodes")
					addNodes, removeNodes := compareMaps(oldNodeList.([]interface{}), newNodeList.([]interface{}))
					if multiAZ == "enabled" && len(addNodes) != 0 && len(removeNodes) != 0 {
						return false, fmt.Errorf("A single operation can only add or reduce nodes, and cannot add and reduce nodes simultaneously")
					}
					if action == "IncreaseDBInstanceNodeNumber" {
						(*call.SdkParam)["NodesNumberToIncrease"] = changeNum
						if multiAZ == "enabled" && len(addNodes) > 0 {
							nodes := make([]map[string]interface{}, 0)
							for _, n := range addNodes {
								nodes = append(nodes, map[string]interface{}{
									"AZ": n["az"],
								})
							}
							(*call.SdkParam)["ConfigureNewNodes"] = nodes
						}
					} else {
						(*call.SdkParam)["NodesNumberToDecrease"] = changeNum
						if multiAZ == "enabled" && len(removeNodes) > 0 {
							nodes := make([]map[string]interface{}, 0)
							for _, n := range removeNodes {
								nodes = append(nodes, map[string]interface{}{
									"AZ": n["az"],
								})
							}
							(*call.SdkParam)["NodesToRemove"] = nodes
						}
					}
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("shard_number") && !resourceData.HasChange("sharded_cluster") {
		// 如果触发了EnableShardedCluster，就不用调用ModifyDBInstanceShardNumber了，因为EnableShardedCluster会同步修改shard_number
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceShardNumber",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"shard_number": {
						TargetField: "ShardNumber",
					},
					"apply_immediately": {
						TargetField: "ApplyImmediately",
						ForceGet:    true,
					},
					"create_backup": {
						TargetField: "CreateBackup",
						ForceGet:    true,
					},
					"backup_point_name": {
						TargetField: "BackupPointName",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					nodeNumber := d.Get("node_number").(int)
					if nodeNumber > 1 {
						(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
					}
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("additional_bandwidth") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceAdditionalBandwidthPerShard",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"additional_bandwidth": {
						TargetField: "AdditionalBandwidth",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("shard_capacity") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceShardCapacity",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"shard_capacity": {
						TargetField: "ShardCapacity",
					},
					"apply_immediately": {
						TargetField: "ApplyImmediately",
						ForceGet:    true,
					},
					"create_backup": {
						TargetField: "CreateBackup",
						ForceGet:    true,
					},
					"backup_point_name": {
						TargetField: "BackupPointName",
						ForceGet:    true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					nodeNumber := d.Get("node_number").(int)
					if nodeNumber > 1 {
						(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
					}
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("password") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBAccount",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"password": {
						TargetField: "Password",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["AccountName"] = "default"
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("time_scope") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyMaintenanceTime",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"time_scope": {
						TargetField: "TimeScope",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["TimeScope"] = d.Get("time_scope")
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("max_connections") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceMaxConn",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"max_connections": {
						TargetField: "MaxConnections",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["MaxConnections"] = d.Get("max_connections")
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChanges("addr_type", "new_address_prefix", "new_port", "upgrade_region_domain") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceVisitAddress",
				ConvertMode: bp.RequestConvertInConvert,
				ContentType: bp.ContentTypeJson,
				Convert: map[string]bp.RequestConvert{
					"addr_type": {
						TargetField: "AddrType",
					},
					"new_address_prefix": {
						TargetField: "NewAddressPrefix",
					},
					"new_port": {
						TargetField: "NewPort",
					},
					"upgrade_region_domain": {
						TargetField: "UpgradeRegionDomain",
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					if len(*call.SdkParam) == 0 {
						return false, nil
					}
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["AddrType"] = d.Get("addr_type")
					(*call.SdkParam)["NewAddressPrefix"] = d.Get("new_address_prefix")
					(*call.SdkParam)["NewPort"] = d.Get("new_port")
					(*call.SdkParam)["UpgradeRegionDomain"] = d.Get("upgrade_region_domain")
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
					return true, nil
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
		callbacks = append(callbacks, callback)
	}

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, callbacks)

	return callbacks
}

func (s *ByteplusRedisDbInstanceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBInstance",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if deletionProjection := d.Get("deletion_protection"); deletionProjection.(string) == "enabled" {
					return false, fmt.Errorf("can not delete protected redis instance")
				}
				(*call.SdkParam)["InstanceId"] = d.Id()
				nodeNumber := d.Get("node_number").(int)
				if nodeNumber > 1 {
					(*call.SdkParam)["CreateBackup"] = d.Get("create_backup")
				}
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				backupPointName, ok := d.Get("backup_point_name").(string)
				if ok {
					if backupPointName != "" {
						(*call.SdkParam)["BackupPointName"] = backupPointName
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				// 开启删除保护时，跳过 CallError
				if strings.Contains(baseErr.Error(), "can not delete protected redis instance") {
					return baseErr
				}
				// 出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on reading redis instance on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisDbInstanceService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"tags": {
				TargetField: "Tags",
				ConvertType: bp.ConvertJsonObjectArray,
			},
		},
		IdField:      "InstanceId",
		NameField:    "InstanceName",
		CollectField: "instances",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"InstanceId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"TTL": {
				TargetField: "ttl",
			},
			"MultiAZ": {
				TargetField: "multi_az",
			},
			"AZ": {
				TargetField: "az",
			},
			"VIP": {
				TargetField: "vip",
			},
			"VIPv6": {
				TargetField: "vip_v6",
			},
			"ShardCapacityV2": {
				TargetField: "shard_capacity_v2",
			},
		},
	}
}

func (s *ByteplusRedisDbInstanceService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusRedisDbInstanceService) getVpcIdAndZoneIdBySubnet(subnetId string) (vpcId, zoneId string, err error) {
	// describe subnet
	req := map[string]interface{}{
		"SubnetIds.1": subnetId,
	}
	action := "DescribeSubnets"
	resp, err := s.Client.UniversalClient.DoCall(getVpcUniversalInfo(action), &req)
	if err != nil {
		return "", "", err
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	results, err := bp.ObtainSdkValue("Result.Subnets", *resp)
	if err != nil {
		return "", "", err
	}
	if results == nil {
		results = []interface{}{}
	}
	subnets, ok := results.([]interface{})
	if !ok {
		return "", "", errors.New("Result.Subnets is not Slice")
	}
	if len(subnets) == 0 {
		return "", "", fmt.Errorf("subnet %s not exist", subnetId)
	}
	vpcId = subnets[0].(map[string]interface{})["VpcId"].(string)
	zoneId = subnets[0].(map[string]interface{})["ZoneId"].(string)
	return vpcId, zoneId, nil
}

func (s *ByteplusRedisDbInstanceService) setResourceTags(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
	addedTags, removedTags, _, _ := bp.GetSetDifference("tags", resourceData, bp.TagsHash, false)

	removeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "RemoveTagsFromResource",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if removedTags != nil && len(removedTags.List()) > 0 {
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
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
					(*call.SdkParam)["ClientToken"] = uuid.New().String()
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

func (s *ByteplusRedisDbInstanceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "Redis",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Action:      actionName,
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
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

func getPrefixBeforeEndPoint(domain, target string) string {
	index := strings.Index(domain, target)
	if index == -1 {
		return ""
	}
	return domain[:index]
}
