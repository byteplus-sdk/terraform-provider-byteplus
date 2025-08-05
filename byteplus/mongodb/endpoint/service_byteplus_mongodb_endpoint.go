package endpoint

import (
	"errors"
	"fmt"
	"strings"
	"time"

	mongodbInstance "github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/mongodb/instance"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusMongoDBEndpointService struct {
	Client *bp.SdkClient
}

func NewMongoDBEndpointService(c *bp.SdkClient) *ByteplusMongoDBEndpointService {
	return &ByteplusMongoDBEndpointService{
		Client: c,
	}
}

func (s *ByteplusMongoDBEndpointService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongoDBEndpointService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action := "DescribeDBEndpoint"

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

	logger.Debug(logger.RespFormat, action, resp)
	results, err = bp.ObtainSdkValue("Result.DBEndpoints", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New("Result.DBEndpoints is not Slice")
	}
	return data, err
}

func (s *ByteplusMongoDBEndpointService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		instanceId   string
		endpointId   string
		objectId     string
		tempObjectId string
		networkType  string
	)

	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return data, fmt.Errorf("format of mongodb endpoint resource id is invalid,%s", id)
	}
	instanceId = parts[0]
	endpointId = parts[1]

	req := map[string]interface{}{
		"InstanceId": instanceId,
	}
	results, err := s.ReadResources(req)
	if err != nil {
		return nil, err
	}

	var targetEndpoint map[string]interface{}
	if a, ok := resourceData.GetOkExists("network_type"); ok {
		networkType = a.(string)
	}
	if a, ok := resourceData.GetOkExists("object_id"); ok {
		objectId = a.(string)
	}
	for _, v := range results {
		dbEndpoint, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("dbEndpoint value is not map")
		}
		eId, err := bp.ObtainSdkValue("EndpointId", dbEndpoint)
		if err != nil {
			return nil, err
		}
		logger.DebugInfo("---- endpointId:%s,eid:%s", endpointId, eId)
		if endpointId != "" { // check by EndpointId
			if endpointId == eId.(string) {
				logger.DebugInfo("get endpoint of endpointId:%s", endpointId)
				targetEndpoint = dbEndpoint
				break
			}
		} else { // check by NetworkType and ObjectId
			nType, err := bp.ObtainSdkValue("NetworkType", dbEndpoint)
			if err != nil {
				return data, err
			}
			oId, err := bp.ObtainSdkValue("ObjectId", dbEndpoint)
			if err != nil {
				return data, err
			}
			if oId != nil {
				tempObjectId = oId.(string)
			}
			if networkType == nType.(string) { // Private or Public
				if objectId == "" || objectId == tempObjectId { //endpointType is ReplicaSet if objectId is empty
					logger.DebugInfo("get mongodb endpoint by  network type and object id %s", networkType, objectId)
					targetEndpoint = dbEndpoint
					break
				}
			}
		}
	}

	if targetEndpoint == nil {
		return data, fmt.Errorf("mongodb endpoint not found")
	}

	nodeIds := make([]string, 0)
	eipIds := make([]string, 0)
	addresses, err := bp.ObtainSdkValue("DBAddresses", targetEndpoint)
	if err != nil {
		return data, err
	}
	endpointType, err := bp.ObtainSdkValue("EndpointType", targetEndpoint)
	if err != nil {
		return data, err
	}
	nType, err := bp.ObtainSdkValue("NetworkType", targetEndpoint)
	if err != nil {
		return data, err
	}
	for _, address := range addresses.([]interface{}) {
		logger.DebugInfo("address %v :", address)
		if nodeId, ok := address.(map[string]interface{})["NodeId"]; ok && nodeId.(string) != "" &&
			endpointType == "Mongos" && nType == "Public" {
			nodeIds = append(nodeIds, nodeId.(string))
		}
		if eipId, ok := address.(map[string]interface{})["EipId"]; ok && eipId.(string) != "" {
			eipIds = append(eipIds, eipId.(string))
		}
	}
	targetEndpoint["MongosNodeIds"] = nodeIds
	targetEndpoint["EipIds"] = eipIds

	return targetEndpoint, err
}

func (s *ByteplusMongoDBEndpointService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusMongoDBEndpointService) WithResourceResponseHandlers(instance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return instance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongoDBEndpointService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBEndpoint",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"eip_ids": {
					ConvertType: bp.ConvertJsonArray,
				},
				"mongos_node_ids": {
					ConvertType: bp.ConvertJsonArray,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				networkType := d.Get("network_type")
				eipIds := d.Get("eip_ids")
				if networkType != nil && networkType.(string) == "Public" {
					if eipIds == nil {
						return false, fmt.Errorf("eip_ids is required when network_type is 'Public'")
					}
				}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// 在 LockId 执行后再进行已有 Endpoint 信息的查询
				endpoint, err := s.ReadResource(d, fmt.Sprintf("%s:", instanceId))
				if err != nil && !strings.Contains(err.Error(), "mongodb endpoint not found") {
					return nil, err
				} else if len(endpoint) != 0 {
					return nil, fmt.Errorf("the instance already contains this endpoint, and duplicate creation is not allowed")
				}

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				logger.Debug("lock instance id:%s", instanceId, "")
				return instanceId
			},
		},
	}
	obtainEndpointIdCallback := bp.Callback{
		Call: bp.SdkCall{
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				endpoint, err := s.ReadResource(d, fmt.Sprintf("%s:", instanceId))
				if err != nil {
					return nil, err
				}
				endpointId := endpoint["EndpointId"].(string)
				_ = d.Set("endpoint_id", endpointId)
				d.SetId(fmt.Sprintf("%s:%s", instanceId, endpointId))
				return nil, nil
			},
		},
	}

	return []bp.Callback{callback, obtainEndpointIdCallback}
}

func (s *ByteplusMongoDBEndpointService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongoDBEndpointService) RemoveResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := resourceData.Get("instance_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBEndpoint",
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"mongos_node_ids": {
					ConvertType: bp.ConvertJsonArray,
					ForceGet:    true,
				},
			},
			SdkParam: &map[string]interface{}{
				"InstanceId": instanceId,
				"EndpointId": resourceData.Get("endpoint_id").(string),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				mongodbInstance.NewMongoDBInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongoDBEndpointService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		IdField:      "EndpointId",
		CollectField: "endpoints",
		ContentType:  bp.ContentTypeJson,
		ResponseConverts: map[string]bp.ResponseConvert{
			"DBAddresses": {
				TargetField: "db_addresses",
			},
			"AddressIP": {
				TargetField: "address_ip",
			},
		},
	}
}

func (s *ByteplusMongoDBEndpointService) ReadResourceId(id string) string {
	return id
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
