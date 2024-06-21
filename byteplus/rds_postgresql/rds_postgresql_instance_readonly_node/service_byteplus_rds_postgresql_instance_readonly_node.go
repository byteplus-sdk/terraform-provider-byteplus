package rds_postgresql_instance_readonly_node

import (
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_postgresql/rds_postgresql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsPostgresqlInstanceReadonlyNodeService struct {
	Client             *bp.SdkClient
	Dispatcher         *bp.Dispatcher
	RdsInstanceService *rds_postgresql_instance.ByteplusRdsPostgresqlInstanceService
}

func NewRdsPostgresqlInstanceReadonlyNodeService(c *bp.SdkClient) *ByteplusRdsPostgresqlInstanceReadonlyNodeService {
	return &ByteplusRdsPostgresqlInstanceReadonlyNodeService{
		Client:             c,
		Dispatcher:         &bp.Dispatcher{},
		RdsInstanceService: rds_postgresql_instance.NewRdsPostgresqlInstanceService(c),
	}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return data, nil
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		result map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	result, err = s.RdsInstanceService.ReadResource(resourceData, ids[0])
	if err != nil {
		return data, err
	}
	if len(result) == 0 {
		return result, fmt.Errorf("Rds instance %s not exist ", ids[0])
	}
	if nodeArr, ok := result["Nodes"].([]interface{}); ok {
		for _, node := range nodeArr {
			if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
				if nodeMap["NodeId"] == ids[1] {
					data = nodeMap
				}
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("rds_postgresql_instance_readonly_node %s not exist ", id)
	}
	// 接口返回InstanceId为节点ID
	data["InstanceId"] = ids[0]
	data["NodeId"] = ids[1]
	return data, err
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	existNodeIds := make(map[string]bool)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSpec",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// BeforeCall -> LockId -> ExecuteCall
				// 因此不在BeforeCall中写参数
				// ---------------------------- param trans start ------------------------
				(*call.SdkParam)["InstanceId"] = d.Get("instance_id").(string)

				nodeInfos := make([]interface{}, 0)

				// read node
				instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.ReqFormat, "Read Create ReadOnly Node", instance)
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

								existNodeIds[readonlyNodeInfo["NodeId"].(string)] = true
							}
						}
					}
				}
				logger.Debug(logger.ReqFormat, "New Create ReadOnly Node", nodeInfos)

				// 2. add node
				newReadonlyNodeInfo := make(map[string]interface{})
				newReadonlyNodeInfo["NodeType"] = "ReadOnly"
				newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
				newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
				newReadonlyNodeInfo["NodeOperateType"] = "Create"
				nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos
				// ---------------------------- param trans end------------------------

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterRefresh: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) error {
				// AfterCall -> ExtraRefresh -> AfterRefresh
				// 因此使用AfterRefresh做set id，避免没读到
				var (
					instance          map[string]interface{}
					err               error
					newReadonlyNodeId string
				)
				instance, err = s.RdsInstanceService.ReadResource(d, d.Get("instance_id").(string))
				if err != nil {
					return err
				}
				logger.Debug(logger.ReqFormat, "AfterCall ReadOnly Node", instance["Nodes"], existNodeIds)
				if nodeArr, ok := instance["Nodes"].([]interface{}); ok {
					for _, node := range nodeArr {
						if nodeMap, ok1 := node.(map[string]interface{}); ok1 {
							if nodeMap["NodeType"] == "ReadOnly" {
								if _, ok = existNodeIds[nodeMap["NodeId"].(string)]; !ok {
									newReadonlyNodeId = nodeMap["NodeId"].(string)
								}
							}
						}
					}
				}
				logger.Debug(logger.ReqFormat, "newReadonlyNodeId", newReadonlyNodeId, existNodeIds)
				if newReadonlyNodeId == "" {
					return fmt.Errorf(" Failed to create readonly node ")
				}
				id := fmt.Sprintf("%s:%s", d.Get("instance_id"), newReadonlyNodeId)
				d.SetId(id)
				return nil
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
	return []bp.Callback{callback}
}

func (ByteplusRdsPostgresqlInstanceReadonlyNodeService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSpec",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// BeforeCall -> LockId -> ExecuteCall
				// 因此不在BeforeCall中写参数
				// ---------------------------- param trans start ------------------------
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["InstanceId"] = ids[0]

				nodeInfos := make([]interface{}, 0)

				// read node
				instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.ReqFormat, "Read Modify ReadOnly Node", instance)
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
							} else if nodeMap["NodeType"] == "ReadOnly" && nodeMap["NodeId"] != ids[1] {
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

				// 2. modify node
				newReadonlyNodeInfo := make(map[string]interface{})
				newReadonlyNodeInfo["NodeId"] = ids[1]
				newReadonlyNodeInfo["NodeType"] = "ReadOnly"
				newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
				newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
				newReadonlyNodeInfo["NodeOperateType"] = "Modify"
				nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos
				// ---------------------------- param trans end------------------------

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
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
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBInstanceSpec",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// BeforeCall -> LockId -> ExecuteCall
				// 因此不在BeforeCall中写参数
				// ---------------------------- param trans start ------------------------
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["InstanceId"] = ids[0]

				nodeInfos := make([]interface{}, 0)

				// read node
				instance, err := s.RdsInstanceService.ReadResource(resourceData, d.Get("instance_id").(string))
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.ReqFormat, "Read Delete ReadOnly Node", instance)
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
							} else if nodeMap["NodeType"] == "ReadOnly" && nodeMap["NodeId"] != ids[1] {
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

				// 2. delete node
				newReadonlyNodeInfo := make(map[string]interface{})
				newReadonlyNodeInfo["NodeId"] = ids[1]
				newReadonlyNodeInfo["NodeType"] = "ReadOnly"
				newReadonlyNodeInfo["NodeSpec"] = d.Get("node_spec")
				newReadonlyNodeInfo["ZoneId"] = d.Get("zone_id")
				newReadonlyNodeInfo["NodeOperateType"] = "Delete"
				nodeInfos = append(nodeInfos, newReadonlyNodeInfo)

				(*call.SdkParam)["NodeInfo"] = nodeInfos
				// ---------------------------- param trans end------------------------

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
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
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusRdsPostgresqlInstanceReadonlyNodeService) ReadResourceId(id string) string {
	return id
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
