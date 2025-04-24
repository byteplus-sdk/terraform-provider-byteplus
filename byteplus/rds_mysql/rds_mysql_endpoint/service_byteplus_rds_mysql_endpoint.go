package rds_mysql_endpoint

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_instance"
	"strconv"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlEndpointService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsMysqlEndpointService(c *bp.SdkClient) *ByteplusRdsMysqlEndpointService {
	return &ByteplusRdsMysqlEndpointService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsMysqlEndpointService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlEndpointService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(m, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeDBInstanceDetail"

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
		results, err = bp.ObtainSdkValue("Result.Endpoints", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Endpoints is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusRdsMysqlEndpointService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
		temp    map[string]interface{}
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"InstanceId": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else {
			if endpointId, ok := temp["EndpointId"]; ok {
				if ids[1] == endpointId.(string) {
					data = temp
				}
			}
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("rds_mysql_endpoint %s not exist ", id)
	}
	logger.Debug(logger.ReqFormat, "Before Trans Data", data)
	// nodes 不读
	transEnableToBool("AutoAddNewNodes", data)
	transEnableToBool("EnableReadWriteSplitting", data)
	addresses, ok := data["Addresses"]
	if ok {
		addressesList := addresses.([]interface{})
		for _, address := range addressesList {
			addressMap := address.(map[string]interface{})
			if addressMap["NetworkType"].(string) == "Private" {
				port, err := strconv.Atoi(addressMap["Port"].(string))
				if err != nil {
					return data, fmt.Errorf("Port is not a number ")
				}
				data["Port"] = port
				data["DnsVisibility"] = addressMap["DNSVisibility"]
				data["Domain"] = addressMap["Domain"]
				break
			}
		}
	}
	data["ReadWriteSpliting"] = data["EnableReadWriteSplitting"]
	// 删除不对应的Nodes
	delete(data, "Nodes")
	if nodesSet, ok := resourceData.GetOk("nodes"); ok {
		data["Nodes"] = nodesSet.(*schema.Set).List()
	}
	// 防止自增节点
	delete(data, "ReadOnlyNodeWeight")
	if w, ok := resourceData.GetOk("read_only_node_weight"); ok {
		weights := make([]interface{}, 0)
		for _, v := range w.(*schema.Set).List() {
			weight := make(map[string]interface{})
			vMap := v.(map[string]interface{})
			if nodeId, ok := vMap["node_id"]; ok {
				weight["NodeId"] = nodeId
			}
			if nodeType, ok := vMap["node_type"]; ok {
				weight["NodeType"] = nodeType
			}
			if we, ok := vMap["weight"]; ok {
				weight["Weight"] = we
			}
			weights = append(weights, weight)
		}
		data["ReadOnlyNodeWeight"] = weights
	}
	logger.Debug(logger.ReqFormat, "After Trans Data", data)
	return data, err
}

func transEnableToBool(field string, data map[string]interface{}) {
	var (
		v   interface{}
		str string
		ok  bool
	)
	if v, ok = data[field]; ok {
		if str, ok = v.(string); ok {
			if str == "Enable" {
				data[field] = true
			} else if str == "Disable" {
				data[field] = false
			}
		}
	}
}

func (s *ByteplusRdsMysqlEndpointService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusRdsMysqlEndpointService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	if id, ok := resourceData.GetOk("endpoint_id"); ok {
		return []bp.Callback{{
			Call: bp.SdkCall{
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					return nil, nil
				},
				AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
					time.Sleep(30 * time.Second)
					instanceId := d.Get("instance_id").(string)
					d.SetId(fmt.Sprintf("%s:%s", instanceId, id))
					return nil
				},
			},
		}}
	}
	var callbacks []bp.Callback
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDBEndpoint",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"read_write_spliting": {
					Ignore: true,
				},
				"read_only_node_max_delay_time": {
					Ignore: true,
				},
				"read_only_node_distribution_type": {
					Ignore: true,
				},
				"read_only_node_weight": {
					Ignore: true,
				},
				"nodes": {
					Ignore: true,
				},
				"dns_visibility": {
					Ignore: true,
				},
				"domain": {
					Ignore: true,
				},
				"port": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				if nodes, ok := d.GetOk("nodes"); ok {
					nodesList := nodes.(*schema.Set).List()
					nodesArr := make([]string, 0)
					for _, node := range nodesList {
						nodesArr = append(nodesArr, node.(string))
					}
					(*call.SdkParam)["Nodes"] = strings.Join(nodesArr, ",")
				}
				(*call.SdkParam)["EndpointType"] = "Custom"
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(30 * time.Second)
				endpointId, _ := bp.ObtainSdkValue("Result.EndpointId", *resp)
				instanceId := d.Get("instance_id").(string)
				d.SetId(fmt.Sprintf("%s:%s", instanceId, endpointId))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					ResourceId: resourceData.Get("instance_id").(string),
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		},
	}
	callbacks = append(callbacks, callback)
	// 调用 ModifyDBEndpointDNS 接口修改私网地址的解析方式。
	//if dns, ok := resourceData.GetOk("dns_visibility"); ok {
	//	dnsCallback := bp.Callback{
	//		Call: bp.SdkCall{
	//			Action:      "ModifyDBEndpointDNS",
	//			ConvertMode: bp.RequestConvertIgnore,
	//			ContentType: bp.ContentTypeJson,
	//			SdkParam: &map[string]interface{}{
	//				"NetworkType":   "Private",
	//				"DNSVisibility": dns,
	//			},
	//			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
	//				ids := strings.Split(d.Id(), ":")
	//				(*call.SdkParam)["InstanceId"] = ids[0]
	//				(*call.SdkParam)["EndpointId"] = ids[1]
	//				return true, nil
	//			},
	//			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
	//				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
	//				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
	//				logger.Debug(logger.RespFormat, call.Action, resp, err)
	//				return resp, err
	//			},
	//			LockId: func(d *schema.ResourceData) string {
	//				return d.Get("instance_id").(string)
	//			},
	//			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
	//				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
	//					ResourceId: resourceData.Get("instance_id").(string),
	//					Target:     []string{"Running"},
	//					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
	//				},
	//			},
	//		},
	//	}
	//	callbacks = append(callbacks, dnsCallback)
	//}
	// 调用 ModifyDBEndpointAddress 接口修改连接地址的前缀或端口。 仅针对私网。
	port, ok := resourceData.GetOk("port")
	domain, ok1 := resourceData.GetOk("domain")
	if ok || ok1 {
		addressCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBEndpointAddress",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					ids := strings.Split(d.Id(), ":")
					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["EndpointId"] = ids[1]
					(*call.SdkParam)["NetworkType"] = "Private"
					// 默认3306
					if ok && port.(int) != 3306 {
						(*call.SdkParam)["Port"] = port
					}
					if ok1 {
						arr := strings.Split(domain.(string), ".")
						if len(arr) < 2 {
							return false, fmt.Errorf("domain is not valid")
						}
						(*call.SdkParam)["DomainPrefix"] = arr[0]
					}
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
						ResourceId: resourceData.Get("instance_id").(string),
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					},
				},
			},
		}
		callbacks = append(callbacks, addressCallback)
	}
	// 调用 ModifyDBEndpoint 接口修改 Endpoint。
	/*
		read_write_splitting
		read_only_node_max_delay_time
		read_only_node_distribution_type
		read_only_node_weight
	*/
	_, spExist := resourceData.GetOk("read_write_spliting")
	_, timeExist := resourceData.GetOk("read_only_node_max_delay_time")
	_, typeExist := resourceData.GetOk("read_only_node_distribution_type")
	_, weightExist := resourceData.GetOk("read_only_node_weight")
	if spExist || timeExist || typeExist || weightExist {
		modifyCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBEndpoint",
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertAll,
				Convert: map[string]bp.RequestConvert{
					"nodes": {
						Ignore: true,
					},
					"read_only_node_weight": {
						ConvertType: bp.ConvertJsonObjectArray,
					},
					"dns_visibility": {
						Ignore: true,
					},
					"domain": {
						Ignore: true,
					},
					"port": {
						Ignore: true,
					},
				},
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					ids := strings.Split(d.Id(), ":")
					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["EndpointId"] = ids[1]
					if nodes, ok := d.GetOk("nodes"); ok {
						nodesList := nodes.(*schema.Set).List()
						nodesArr := make([]string, 0)
						for _, node := range nodesList {
							nodesArr = append(nodesArr, node.(string))
						}
						(*call.SdkParam)["Nodes"] = strings.Join(nodesArr, ",")
					}
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
						ResourceId: resourceData.Get("instance_id").(string),
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					},
				},
			},
		}
		callbacks = append(callbacks, modifyCallback)
	}
	return callbacks
}

func (ByteplusRdsMysqlEndpointService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlEndpointService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callbacks := make([]bp.Callback, 0)
	ids := strings.Split(resourceData.Id(), ":")
	// 调用 ModifyDBEndpointAddress 接口修改连接地址的前缀或端口。 仅针对私网。
	if resourceData.HasChanges("port", "domain") {
		addressCallback := bp.Callback{
			Call: bp.SdkCall{
				Action:      "ModifyDBEndpointAddress",
				ConvertMode: bp.RequestConvertIgnore,
				ContentType: bp.ContentTypeJson,
				BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
					(*call.SdkParam)["InstanceId"] = ids[0]
					(*call.SdkParam)["EndpointId"] = ids[1]
					(*call.SdkParam)["NetworkType"] = "Private"
					port, ok := resourceData.GetOk("port")
					domain, ok1 := resourceData.GetOk("domain")
					if ok && d.HasChange("port") {
						(*call.SdkParam)["Port"] = port
					}
					if ok1 && d.HasChange("domain") {
						arr := strings.Split(domain.(string), ".")
						if len(arr) < 2 {
							return false, fmt.Errorf("domain is not valid")
						}
						(*call.SdkParam)["DomainPrefix"] = arr[0]
					}
					return true, nil
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				},
				LockId: func(d *schema.ResourceData) string {
					return d.Get("instance_id").(string)
				},
				ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
					rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
						ResourceId: resourceData.Get("instance_id").(string),
						Target:     []string{"Running"},
						Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					},
				},
			},
		}
		callbacks = append(callbacks, addressCallback)
	}
	// 调用 ModifyDBEndpointDNS 接口修改私网地址的解析方式。
	//if resourceData.HasChange("dns_visibility") {
	//	dnsCallback := bp.Callback{
	//		Call: bp.SdkCall{
	//			Action:      "ModifyDBEndpointDNS",
	//			ConvertMode: bp.RequestConvertIgnore,
	//			ContentType: bp.ContentTypeJson,
	//			SdkParam: &map[string]interface{}{
	//				"InstanceId":    ids[0],
	//				"EndpointId":    ids[1],
	//				"NetworkType":   "Private",
	//				"DNSVisibility": resourceData.Get("dns_visibility"),
	//			},
	//			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
	//				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
	//				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
	//				logger.Debug(logger.RespFormat, call.Action, resp, err)
	//				return resp, err
	//			},
	//			LockId: func(d *schema.ResourceData) string {
	//				return d.Get("instance_id").(string)
	//			},
	//			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
	//				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
	//					ResourceId: resourceData.Get("instance_id").(string),
	//					Target:     []string{"Running"},
	//					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
	//				},
	//			},
	//		},
	//	}
	//	callbacks = append(callbacks, dnsCallback)
	//}
	// 调用 ModifyDBEndpoint 接口修改 Endpoint。
	modifyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyDBEndpoint",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"nodes": {
					Ignore: true,
				},
				"read_write_mode": {
					TargetField: "ReadWriteMode",
				},
				"endpoint_name": {
					TargetField: "EndpointName",
				},
				"description": {
					TargetField: "Description",
				},
				"auto_add_new_nodes": {
					TargetField: "AutoAddNewNodes",
					ForceGet:    true,
				},
				"read_write_spliting": {
					TargetField: "ReadWriteSpliting",
					ForceGet:    true,
				},
				"read_only_node_max_delay_time": {
					TargetField: "ReadOnlyNodeMaxDelayTime",
				},
				"read_only_node_distribution_type": {
					TargetField: "ReadOnlyNodeDistributionType",
					// 不传会变default
					ForceGet: true,
				},
				"read_only_node_weight": {
					Ignore: true,
				},
				"dns_visibility": {
					Ignore: true,
				},
				"domain": {
					Ignore: true,
				},
				"port": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["EndpointId"] = ids[1]
				if nodes, ok := d.GetOk("nodes"); ok {
					nodesList := nodes.(*schema.Set).List()
					nodesArr := make([]string, 0)
					for _, node := range nodesList {
						nodesArr = append(nodesArr, node.(string))
					}
					(*call.SdkParam)["Nodes"] = strings.Join(nodesArr, ",")
				}
				if d.HasChange("read_only_node_weight") ||
					d.Get("read_only_node_distribution_type").(string) == "Custom" {
					weights := make([]interface{}, 0)
					w := d.Get("read_only_node_weight")
					for _, v := range w.(*schema.Set).List() {
						weight := make(map[string]interface{})
						vMap := v.(map[string]interface{})
						if nodeId, ok := vMap["node_id"]; ok {
							weight["NodeId"] = nodeId
						}
						if nodeType, ok := vMap["node_type"]; ok {
							weight["NodeType"] = nodeType
						}
						if we, ok := vMap["weight"]; ok {
							weight["Weight"] = we
						}
						weights = append(weights, weight)
					}
					(*call.SdkParam)["ReadOnlyNodeWeight"] = weights
				}
				return true, nil
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					ResourceId: resourceData.Get("instance_id").(string),
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				},
			},
		},
	}
	callbacks = append(callbacks, modifyCallback)
	return callbacks
}

func (s *ByteplusRdsMysqlEndpointService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDBEndpoint",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceId": ids[0],
				"EndpointId": ids[1],
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					ResourceId: resourceData.Get("instance_id").(string),
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				return bp.CheckResourceUtilRemoved(d, s.ReadResource, 5*time.Minute)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlEndpointService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "EndpointName",
		IdField:      "EndpointId",
		CollectField: "endpoints",
		ResponseConverts: map[string]bp.ResponseConvert{
			"EndpointId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"IPAddress": {
				TargetField: "ip_address",
			},
			"DNSVisibility": {
				TargetField: "dns_visibility",
			},
		},
	}
}

func (s *ByteplusRdsMysqlEndpointService) ReadResourceId(id string) string {
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
