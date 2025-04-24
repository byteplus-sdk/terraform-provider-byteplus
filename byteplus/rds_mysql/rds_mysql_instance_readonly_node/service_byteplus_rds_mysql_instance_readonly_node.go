package rds_mysql_instance_readonly_node

import (
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlInstanceReadonlyNodeService struct {
	Client             *bp.SdkClient
	Dispatcher         *bp.Dispatcher
	RdsInstanceService *rds_mysql_instance.ByteplusRdsMysqlInstanceService
}

func NewRdsMysqlInstanceReadonlyNodeService(c *bp.SdkClient) *ByteplusRdsMysqlInstanceReadonlyNodeService {
	return &ByteplusRdsMysqlInstanceReadonlyNodeService{
		Client:             c,
		Dispatcher:         &bp.Dispatcher{},
		RdsInstanceService: rds_mysql_instance.NewRdsMysqlInstanceService(c),
	}
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, err
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) ReadResource(resourceData *schema.ResourceData, rdsInstanceNodeId string) (data map[string]interface{}, err error) {
	if rdsInstanceNodeId == "" {
		rdsInstanceNodeId = resourceData.Id()
	}

	ids := strings.Split(rdsInstanceNodeId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid rdsInstanceNodeId: %s", rdsInstanceNodeId)
	}

	instanceId := ids[0]
	nodeId := ids[1]

	result, err := s.RdsInstanceService.ReadResource(resourceData, instanceId)
	if err != nil {
		return result, err
	}
	if len(result) == 0 {
		return result, fmt.Errorf("Rds instance %s not exist ", instanceId)
	}

	if nodeArr, ok := result["Nodes"].([]interface{}); ok {
		for _, node := range nodeArr {
			if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
				if nodeMap["NodeId"] == nodeId {
					data = nodeMap
				}
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Rds instance readonly node %s is not exist ", nodeId)
	}
	data["NodeId"] = nodeId

	return data, err
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (*ByteplusRdsMysqlInstanceReadonlyNodeService) WithResourceResponseHandlers(rdsInstance map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rdsInstance, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) CreateResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var (
		callbacks               []bp.Callback
		existingReadOnlyNodeIds = make(map[string]bool)
	)

	nodeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSpec",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (common *map[string]interface{}, err error) {
				// 在LockId执行后再进行已有Node信息的查询
				(*call.SdkParam)["InstanceId"] = d.Get("instance_id").(string)

				nodeInfos := make([]interface{}, 0)
				// 1. 获取当前RdsInstance已有的Node信息
				instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
				if err != nil {
					return common, err
				}
				if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
					for _, node := range nodeArr {
						if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
							if nodeMap["NodeType"] == "Primary" {
								primaryNodeInfo := make(map[string]interface{})
								primaryNodeInfo["NodeId"] = nodeMap["NodeId"]
								primaryNodeInfo["NodeType"] = nodeMap["NodeType"]
								primaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
								primaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
								nodeInfos = append(nodeInfos, primaryNodeInfo)
							} else if nodeMap["NodeType"] == "Secondary" {
								secondaryNodeInfo := make(map[string]interface{})
								secondaryNodeInfo["NodeId"] = nodeMap["NodeId"]
								secondaryNodeInfo["NodeType"] = nodeMap["NodeType"]
								secondaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
								secondaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
								nodeInfos = append(nodeInfos, secondaryNodeInfo)
							} else if nodeMap["NodeType"] == "ReadOnly" {
								readonlyNodeInfo := make(map[string]interface{})
								readonlyNodeInfo["NodeId"] = nodeMap["NodeId"]
								readonlyNodeInfo["NodeType"] = nodeMap["NodeType"]
								readonlyNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
								readonlyNodeInfo["ZoneId"] = nodeMap["ZoneId"]
								nodeInfos = append(nodeInfos, readonlyNodeInfo)

								existingReadOnlyNodeIds[readonlyNodeInfo["NodeId"].(string)] = true
							}
						}
					}
				}

				// 2. 新增 readonly node
				newReadonlyNodeInfo := make(map[string]interface{})
				newReadonlyNodeInfo["NodeType"] = "ReadOnly"
				newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
				newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
				newReadonlyNodeInfo["NodeOperateType"] = "Create"
				nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				common, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				if err != nil {
					return common, err
				}
				return common, nil
			},
			AfterRefresh: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
				var (
					instance map[string]interface{}
					err      error
				)
				instance, err = s.RdsInstanceService.ReadResource(d, d.Get("instance_id").(string))
				if err != nil {
					return err
				}
				var newReadonlyNodeId string
				if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
					for _, node := range nodeArr {
						if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
							if nodeMap["NodeType"] == "ReadOnly" {
								if _, ok := existingReadOnlyNodeIds[nodeMap["NodeId"].(string)]; !ok {
									newReadonlyNodeId = nodeMap["NodeId"].(string)
								}
							}
						}
					}
				}
				// ResourceData中，rds_mysql_instance_readonly_node的Id形式为'instance_id:node_id'
				logger.Debug(logger.ReqFormat, "newReadonlyNodeId", newReadonlyNodeId)
				if newReadonlyNodeId == "" {
					return fmt.Errorf(" Failed to create readonly node ")
				}
				id := fmt.Sprintf("%s:%s", d.Get("instance_id"), newReadonlyNodeId)
				d.SetId(id)
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				s.RdsInstanceService: {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, nodeCallback)

	return callbacks
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	if resourceData.HasChange("node_spec") {
		nodeCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBInstanceSpec",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertIgnore,
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (common *map[string]interface{}, err error) {
					// 在 LockId 后再进行已有 Node 信息的查询
					ids := strings.Split(d.Id(), ":")
					if len(ids) != 2 {
						return common, fmt.Errorf("invalid rdsInstanceNodeId: %s", d.Id())
					}
					instanceId := ids[0]
					nodeId := ids[1]
					(*call.SdkParam)["InstanceId"] = instanceId

					nodeInfos := make([]interface{}, 0)
					// 1. 获取当前RdsInstance已有的Node信息
					instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
					if err != nil {
						return common, err
					}
					if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
						for _, node := range nodeArr {
							if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
								if nodeMap["NodeType"] == "Primary" {
									primaryNodeInfo := make(map[string]interface{})
									primaryNodeInfo["NodeId"] = nodeMap["NodeId"]
									primaryNodeInfo["NodeType"] = nodeMap["NodeType"]
									primaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
									primaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
									nodeInfos = append(nodeInfos, primaryNodeInfo)
								} else if nodeMap["NodeType"] == "Secondary" {
									secondaryNodeInfo := make(map[string]interface{})
									secondaryNodeInfo["NodeId"] = nodeMap["NodeId"]
									secondaryNodeInfo["NodeType"] = nodeMap["NodeType"]
									secondaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
									secondaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
									nodeInfos = append(nodeInfos, secondaryNodeInfo)
								} else if nodeMap["NodeType"] == "ReadOnly" && nodeMap["NodeId"] != nodeId {
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

					// 2. 修改当前 readonly node
					newReadonlyNodeInfo := make(map[string]interface{})
					newReadonlyNodeInfo["NodeId"] = nodeId
					newReadonlyNodeInfo["NodeType"] = "ReadOnly"
					newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
					newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
					newReadonlyNodeInfo["NodeOperateType"] = "Modify"
					nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

					(*call.SdkParam)["NodeInfo"] = nodeInfos

					logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
					common, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					if err != nil {
						return common, err
					}
					return common, nil
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					s.RdsInstanceService: {
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutUpdate),
						ResourceId: resourceData.Get("instance_id").(string),
					},
				},
			},
		}
		callbacks = append(callbacks, nodeCallback)
	}

	return callbacks
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	nodeCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSpec",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (common *map[string]interface{}, err error) {
				// 在 LockId 后再进行已有 Node 信息的查询
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return common, fmt.Errorf("invalid rdsInstanceNodeId: %s", d.Id())
				}
				instanceId := ids[0]
				nodeId := ids[1]
				(*call.SdkParam)["InstanceId"] = instanceId

				nodeInfos := make([]interface{}, 0)
				// 1. 获取当前RdsInstance已有的Node信息
				instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
				if err != nil {
					return common, err
				}
				if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
					for _, node := range nodeArr {
						if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
							if nodeMap["NodeType"] == "Primary" {
								primaryNodeInfo := make(map[string]interface{})
								primaryNodeInfo["NodeId"] = nodeMap["NodeId"]
								primaryNodeInfo["NodeType"] = nodeMap["NodeType"]
								primaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
								primaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
								nodeInfos = append(nodeInfos, primaryNodeInfo)
							} else if nodeMap["NodeType"] == "Secondary" {
								secondaryNodeInfo := make(map[string]interface{})
								secondaryNodeInfo["NodeId"] = nodeMap["NodeId"]
								secondaryNodeInfo["NodeType"] = nodeMap["NodeType"]
								secondaryNodeInfo["NodeSpec"] = nodeMap["NodeSpec"]
								secondaryNodeInfo["ZoneId"] = nodeMap["ZoneId"]
								nodeInfos = append(nodeInfos, secondaryNodeInfo)
							} else if nodeMap["NodeType"] == "ReadOnly" && nodeMap["NodeId"] != nodeId {
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

				// 2. 删除 readonly node
				newReadonlyNodeInfo := make(map[string]interface{})
				newReadonlyNodeInfo["NodeId"] = nodeId
				newReadonlyNodeInfo["NodeType"] = "ReadOnly"
				newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
				newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
				newReadonlyNodeInfo["NodeOperateType"] = "Delete"
				nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				common, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				if err != nil {
					return common, err
				}
				return common, nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				s.RdsInstanceService: {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	callbacks = append(callbacks, nodeCallback)

	return callbacks
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusRdsMysqlInstanceReadonlyNodeService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_mysql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
