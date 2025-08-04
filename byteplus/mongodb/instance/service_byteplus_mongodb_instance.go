package instance

import (
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusMongoDBInstanceService struct {
	Client *bp.SdkClient
}

func NewMongoDBInstanceService(c *bp.SdkClient) *ByteplusMongoDBInstanceService {
	return &ByteplusMongoDBInstanceService{
		Client: c,
	}
}

func (s *ByteplusMongoDBInstanceService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongoDBInstanceService) readInstanceDetails(id string) (instance interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	action := "DescribeDBInstanceDetail"
	cond := map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	if err != nil {
		return instance, err
	}
	logger.Debug(logger.RespFormat, action, resp)

	instance, err = bp.ObtainSdkValue("Result.DBInstance", *resp)
	if err != nil {
		return instance, err
	}

	return instance, err
}

func (s *ByteplusMongoDBInstanceService) readSSLDetails(id string) (ssl interface{}, err error) {
	var (
		resp *map[string]interface{}
	)
	action := "DescribeDBInstanceSSL"
	cond := map[string]interface{}{
		"InstanceId": id,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	if err != nil {
		return ssl, err
	}
	logger.Debug(logger.RespFormat, action, resp)

	return bp.ObtainSdkValue("Result", *resp)
}

func (s *ByteplusMongoDBInstanceService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)
	withoutDetail, ok := condition["WithoutDetail"]
	if !ok {
		withoutDetail = false
	}
	return bp.WithPageNumberQuery(condition, "PageSize", "PageNumber", 20, 1, func(m map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstances"
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
		logger.Debug(logger.RespFormat, action, condition, *resp)
		results, err = bp.ObtainSdkValue("Result.DBInstances", *resp)
		if err != nil {
			logger.DebugInfo("bp.ObtainSdkValue return :%v", err)
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		instances, ok := results.([]interface{})
		if !ok {
			return data, fmt.Errorf("DescribeDBInstances response instances is not a slice")
		}

		for _, ele := range instances {
			ins := ele.(map[string]interface{})
			instanceId, err := bp.ObtainSdkValue("InstanceId", ele)
			if err != nil {
				return data, err
			}
			// do not get detail when refresh status
			if withoutDetail.(bool) {
				data = append(data, ins)
				continue
			}

			detail, err := s.readInstanceDetails(instanceId.(string))
			if err != nil {
				logger.DebugInfo("read instance %s detail failed,err:%v.", instanceId, err)
				data = append(data, ele)
				continue
			}
			ssl, err := s.readSSLDetails(instanceId.(string))
			if err != nil {
				logger.DebugInfo("read instance ssl information of %s failed,err:%v.", instanceId, err)
				data = append(data, ele)
				continue
			}
			ConfigServers, err := bp.ObtainSdkValue("ConfigServers", detail)
			if err != nil {
				return data, err
			}
			Nodes, err := bp.ObtainSdkValue("Nodes", detail)
			if err != nil {
				return data, err
			}
			Mongos, err := bp.ObtainSdkValue("Mongos", detail)
			if err != nil {
				return data, err
			}
			Shards, err := bp.ObtainSdkValue("Shards", detail)
			if err != nil {
				return data, err
			}
			SSLEnable, err := bp.ObtainSdkValue("SSLEnable", ssl)
			if err != nil {
				return data, err
			}
			SSLIsValid, err := bp.ObtainSdkValue("SSLIsValid", ssl)
			if err != nil {
				return data, err
			}
			SSLExpiredTime, err := bp.ObtainSdkValue("SSLExpiredTime", ssl)
			if err != nil {
				return data, err
			}

			ins["ConfigServers"] = ConfigServers
			ins["Nodes"] = Nodes
			ins["Mongos"] = Mongos
			ins["Shards"] = Shards
			ins["SSLEnable"] = SSLEnable
			ins["SSLIsValid"] = SSLIsValid
			ins["SSLExpiredTime"] = SSLExpiredTime
			data = append(data, ins)
		}
		return data, nil
	})
}

func (s *ByteplusMongoDBInstanceService) readResource(resourceData *schema.ResourceData, id string, withoutDetail bool) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	req := map[string]interface{}{
		"InstanceId":    id,
		"WithoutDetail": withoutDetail,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("instance %s is not exist", id)
	}

	if zoneId, ok := data["ZoneId"]; ok {
		zoneIds := strings.Split(zoneId.(string), ",")
		data["ZoneIds"] = zoneIds
	}

	if nodeZoneSet, ok := resourceData.GetOk("node_availability_zone"); ok {
		data["NodeAvailabilityZone"] = nodeZoneSet.(*schema.Set).List()
	}

	if withoutDetail {
		return data, nil
	}
	instanceType, _ := bp.ObtainSdkValue("InstanceType", data)
	if instanceType.(string) == "ReplicaSet" {
		n, err := bp.ObtainSdkValue("Nodes", data)
		if err != nil || n == nil {
			data["NodeNumber"] = 0
		} else {
			nodes := n.([]interface{})
			data["NodeNumber"] = len(nodes)
			data["NodeSpec"] = nodes[0].(map[string]interface{})["NodeSpec"]
			data["StorageSpaceGb"] = nodes[0].(map[string]interface{})["TotalStorageGB"]
		}
	} else if instanceType.(string) == "ShardedCluster" {
		m, err := bp.ObtainSdkValue("Mongos", data)
		if err != nil || m == nil {
			data["MongosNodeNumber"] = 0
		} else {
			mongos := m.([]interface{})
			data["MongosNodeNumber"] = len(mongos)
			data["MongosNodeSpec"] = mongos[0].(map[string]interface{})["NodeSpec"]
		}
		s, err := bp.ObtainSdkValue("Shards", data)
		if err != nil || s == nil {
			data["ShardNumber"] = 0
			data["StorageSpaceGb"] = 0
		} else {
			shards := s.([]interface{})
			data["ShardNumber"] = len(shards)
			if tmp, ok := shards[0].(map[string]interface{})["Nodes"]; ok {
				nodes := tmp.([]interface{})
				data["StorageSpaceGb"] = nodes[0].(map[string]interface{})["TotalStorageGB"]
				data["NodeSpec"] = nodes[0].(map[string]interface{})["NodeSpec"]
				data["NodeNumber"] = len(nodes)
			}
		}
	}
	return data, err
}

func (s *ByteplusMongoDBInstanceService) readResourceWithoutDetail(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return s.readResource(resourceData, id, true)
}

func (s *ByteplusMongoDBInstanceService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	return s.readResource(resourceData, id, false)
}

func (s *ByteplusMongoDBInstanceService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      1 * time.Second,
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
			failStates = append(failStates, "CreateFailed", "Failed")

			logger.DebugInfo("start refresh :%s", id)
			instance, err = s.readResourceWithoutDetail(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			logger.DebugInfo("Refresh instance status resp: %v", instance)

			status, err = bp.ObtainSdkValue("InstanceStatus", instance)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("instance status error,status %s", status.(string))
				}
			}

			// 判断下实例的计费类型
			if chargeType, ok := resourceData.GetOk("charge_type"); ok && chargeType == "Prepaid" {
				dataChargeType, err := bp.ObtainSdkValue("ChargeType", instance)
				if err != nil || dataChargeType != "Prepaid" {
					return nil, "", err
				}
			}

			logger.DebugInfo("refresh status:%v", status)
			return instance, status.(string), err
		},
	}
}

func (s *ByteplusMongoDBInstanceService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, map[string]bp.ResponseConvert{
			"DBEngine": {
				TargetField: "db_engine",
			},
			"DBEngineVersion": {
				TargetField: "db_engine_version",
			},
			"DBEngineVersionStr": {
				TargetField: "db_engine_version_str",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBInstanceService) getVpcIdAndZoneIdBySubnet(subnetId string) (vpcId, zoneId string, err error) {
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

func (s *ByteplusMongoDBInstanceService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBInstance",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				// "db_engine": {
				// 	TargetField: "DBEngine",
				// },
				"db_engine_version": {
					TargetField: "DBEngineVersion",
				},
				"storage_space_gb": {
					TargetField: "StorageSpaceGB",
				},
				"config_server_node_spec": {
					TargetField: "ConfigServerNodeSpec",
				},
				"config_server_storage_space_gb": {
					TargetField: "ConfigServerStorageSpaceGB",
				},
				"tags": {
					TargetField: "Tags",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"node_availability_zone": {
					TargetField: "NodeAvailabilityZone",
					ConvertType: bp.ConvertJsonObjectArray,
				},
				"auto_renew": {
					TargetField: "AutoRenew",
					ForceGet:    true,
				},
				"zone_id": {
					Ignore: true,
				},
				"zone_ids": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				// describe subnet
				subnetId := d.Get("subnet_id")
				vpcId, zoneId, err := s.getVpcIdAndZoneIdBySubnet(subnetId.(string))
				if err != nil {
					return false, fmt.Errorf("get vpc ID by subnet ID %s failed", subnetId)
				}
				// check custom
				if vpcIdCustom, ok := (*call.SdkParam)["VpcId"]; ok {
					if vpcIdCustom.(string) != vpcId {
						return false, fmt.Errorf("vpc ID does not match")
					}
				}
				if zoneIdCustom, ok := (*call.SdkParam)["ZoneId"]; ok {
					if zoneIdCustom.(string) != zoneId {
						return false, fmt.Errorf("zone ID does not match")
					}
				}

				var zoneIdsStr string
				zoneIdsArr, ok := d.Get("zone_ids").([]interface{})
				if !ok {
					return false, fmt.Errorf("zone_ids is not slice")
				}
				if len(zoneIdsArr) > 0 {
					zoneIds := make([]string, 0)
					for _, id := range zoneIdsArr {
						zoneIds = append(zoneIds, id.(string))
					}
					zoneIdsStr = strings.Join(zoneIds, ",")
				} else {
					zoneIdsStr = zoneId
				}

				(*call.SdkParam)["VpcId"] = vpcId
				(*call.SdkParam)["ZoneId"] = zoneIdsStr
				// (*call.SdkParam)["DBEngine"] = "MongoDB"
				// (*call.SdkParam)["DBEngineVersion"] = "MongoDB_4_2"
				// (*call.SdkParam)["NodeNumber"] = 3
				// (*call.SdkParam)["SuperAccountName"] = "root"

				if (*call.SdkParam)["InstanceType"] == "ShardedCluster" {
					if _, ok := (*call.SdkParam)["MongosNodeSpec"]; !ok {
						return false, fmt.Errorf("mongos_node_spec must exist for ShardedCluster")
					}
					if _, ok := (*call.SdkParam)["MongosNodeNumber"]; !ok {
						return false, fmt.Errorf("mongos_node_number must exist for ShardedCluster")
					}
					if _, ok := (*call.SdkParam)["ShardNumber"]; !ok {
						return false, fmt.Errorf("shard_number must exist for ShardedCluster")
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, resp)
				id, _ := bp.ObtainSdkValue("Result.InstanceId", *resp)
				d.SetId(id.(string))
				time.Sleep(time.Second * 10) //如果创建之后立即refresh，DescribeDBInstances会查找不到这个实例直接返回错误..
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

func (s *ByteplusMongoDBInstanceService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)

	if resourceData.HasChange("instance_name") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceName",
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["InstanceNewName"] = d.Get("instance_name")
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

	if resourceData.HasChange("instance_type") || resourceData.HasChange("node_spec") ||
		resourceData.HasChange("mongos_node_spec") || resourceData.HasChange("shard_number") ||
		resourceData.HasChange("mongos_node_number") || resourceData.HasChange("storage_space_gb") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceSpec",
				ConvertMode: bp.RequestConvertIgnore,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					(*call.SdkParam)["InstanceType"] = d.Get("instance_type")
					if resourceData.HasChange("node_spec") {
						(*call.SdkParam)["NodeSpec"] = d.Get("node_spec")
					}
					if resourceData.HasChange("mongos_node_spec") {
						(*call.SdkParam)["MongosNodeSpec"] = d.Get("mongos_node_spec")
					}
					if resourceData.HasChange("shard_number") {
						(*call.SdkParam)["ShardNumber"] = d.Get("shard_number")
					}
					if resourceData.HasChange("mongos_node_number") {
						(*call.SdkParam)["MongosNodeNumber"] = d.Get("mongos_node_number")
					}
					if resourceData.HasChange("storage_space_gb") {
						(*call.SdkParam)["StorageSpaceGB"] = d.Get("storage_space_gb")
					}
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					time.Sleep(time.Second * 10) //变更之后立即refresh，实例状态还是Running将立即返回..
					return nil
				},
				Refresh: &bp.StateRefresh{
					Target:  []string{"Running"},
					Timeout: resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		}
		callbacks = append(callbacks, callback)
	}

	if resourceData.HasChange("charge_type") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceChargeType",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceIds"] = []interface{}{d.Id()}
					chargeType := d.Get("charge_type")
					if chargeType.(string) != "Prepaid" {
						return false, fmt.Errorf("only supports PostPaid to PrePaid currently")
					}
					(*call.SdkParam)["ChargeType"] = chargeType
					(*call.SdkParam)["PeriodUnit"] = d.Get("period_unit")
					(*call.SdkParam)["Period"] = d.Get("period")
					(*call.SdkParam)["AutoRenew"] = d.Get("auto_renew")
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
		callbacks = append(callbacks, callback)
	}
	if resourceData.HasChange("super_account_password") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ResetDBAccount",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = d.Id()
					//暂时写死 当前不支持这个字段 只能是root
					(*call.SdkParam)["AccountName"] = "root"
					(*call.SdkParam)["AccountPassword"] = d.Get("super_account_password")
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
		callbacks = append(callbacks, callback)
	}

	// 更新Tags
	callbacks = s.setResourceTags(resourceData, callbacks)

	return callbacks
}

func (s *ByteplusMongoDBInstanceService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBInstance",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"InstanceId": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 15*time.Minute)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading mongodb on delete %q, %w", d.Id(), callErr))
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

func (s *ByteplusMongoDBInstanceService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"db_engine": {
				TargetField: "DBEngine",
			},
			"db_engine_version": {
				TargetField: "DBEngineVersion",
			},
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
			"DBEngine": {
				TargetField: "db_engine",
			},
			"DBEngineVersion": {
				TargetField: "db_engine_version",
			},
			"DBEngineVersionStr": {
				TargetField: "db_engine_version_str",
			},
			"TotalMemoryGB": {
				TargetField: "total_memory_gb",
			},
			"TotalvCPU": {
				TargetField: "total_vcpu",
			},
			"UsedMemoryGB": {
				TargetField: "used_memory_gb",
			},
			"UsedvCPU": {
				TargetField: "used_vcpu",
			},
			"TotalStorageGB": {
				TargetField: "total_storage_gb",
			},
			"UsedStorageGB": {
				TargetField: "used_storage_gb",
			},
			"SSLEnable": {
				TargetField: "ssl_enable",
			},
			"SSLIsValid": {
				TargetField: "ssl_is_valid",
			},
			"SSLExpireTime": {
				TargetField: "ssl_expire_time",
			},
		},
	}
}

func (s *ByteplusMongoDBInstanceService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusMongoDBInstanceService) setResourceTags(resourceData *schema.ResourceData, callbacks []bp.Callback) []bp.Callback {
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

func (s *ByteplusMongoDBInstanceService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "mongodb",
		ResourceType:         "instance",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}

func (s *ByteplusMongoDBInstanceService) UnsubscribeInfo(resourceData *schema.ResourceData, resource *schema.Resource) (*bp.UnsubscribeInfo, error) {
	info := bp.UnsubscribeInfo{
		InstanceId: s.ReadResourceId(resourceData.Id()),
	}
	if resourceData.Get("charge_type").(string) == "Prepaid" {
		info.NeedUnsubscribe = true
		info.Products = []string{"veDB for DocumentDB"}
	}
	return &info, nil
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "mongodb",
		Action:      actionName,
		Version:     "2022-01-01",
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
