package endpoint

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/eip/eip_address"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/redis/instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisEndpointService struct {
	Client *bp.SdkClient
}

const (
	ActionCreateDBEndpointPublicAddress = "CreateDBEndpointPublicAddress"
	ActionDeleteDBEndpointPublicAddress = "DeleteDBEndpointPublicAddress"
	ActionDescribeDBInstanceDetail      = "DescribeDBInstanceDetail"
)

func NewRedisEndpointService(c *bp.SdkClient) *ByteplusRedisEndpointService {
	return &ByteplusRedisEndpointService{
		Client: c,
	}
}

func (s *ByteplusRedisEndpointService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisEndpointService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusRedisEndpointService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		ids        []string
		instanceId string
		req        map[string]interface{}
		output     *map[string]interface{}
		results    interface{}
		ok         bool
	)
	if tmpId == "" {
		tmpId = s.ReadResourceId(resourceData.Id())
	}
	ids = strings.Split(tmpId, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid redis endpoint id: %v", tmpId)
	}
	instanceId = ids[0]
	req = map[string]interface{}{
		"InstanceId": instanceId,
	}

	logger.Debug(logger.ReqFormat, ActionDescribeDBInstanceDetail, req)
	output, err = s.Client.UniversalClient.DoCall(getUniversalInfo(ActionDescribeDBInstanceDetail), &req)
	logger.Debug(logger.RespFormat, ActionDescribeDBInstanceDetail, req, *output)

	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue("Result", *output)
	if err != nil {
		return data, err
	}
	if data, ok = results.(map[string]interface{}); !ok {
		return data, errors.New("value is not map")
	}

	if _, exist := data["VisitAddrs"]; !exist {
		return nil, fmt.Errorf("not associated instance and eip. %s", tmpId)
	}

	attached := false
	for _, address := range data["VisitAddrs"].([]interface{}) {
		addr := address.(map[string]interface{})
		if addr["AddrType"].(string) == "Public" && addr["EipId"].(string) == ids[1] {
			attached = true
			break
		}
	}
	if !attached {
		return nil, fmt.Errorf("not associated instance and eip. %s", tmpId)
	}

	return map[string]interface{}{
		"InstanceId": ids[0],
		"EipId":      ids[1],
	}, nil
}

func (s *ByteplusRedisEndpointService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      time.Second,
		Pending:    []string{},
		Target:     target,
		Timeout:    timeout,
		MinTimeout: time.Second,

		Refresh: nil,
	}
}

func (s *ByteplusRedisEndpointService) WithResourceResponseHandlers(endpoint map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusRedisEndpointService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionCreateDBEndpointPublicAddress,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["InstanceId"], ":", (*call.SdkParam)["EipId"]))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				eip_address.NewEipAddressService(s.Client): {
					Target:     []string{"Attached"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("eip_id").(string),
				},
				instance.NewRedisDbInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisEndpointService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRedisEndpointService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionDeleteDBEndpointPublicAddress,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := s.ReadResourceId(d.Id())
				ids := strings.Split(id, ":")
				instanceId := ids[0]
				eipId := ids[1]
				(*call.SdkParam)["InstanceId"] = instanceId
				(*call.SdkParam)["EipId"] = eipId
				(*call.SdkParam)["ClientToken"] = uuid.New().String()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				eip_address.NewEipAddressService(s.Client): {
					Target:     []string{"Available"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: resourceData.Get("eip_id").(string),
				},
				instance.NewRedisDbInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutDelete),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisEndpointService) DatasourceResources(data *schema.ResourceData, resource2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType: bp.ContentTypeJson,
	}
}

func (s *ByteplusRedisEndpointService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
