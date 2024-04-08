package network_interface_attach

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

type ByteplusNetworkInterfaceAttachService struct {
	Client *bp.SdkClient
}

func NewNetworkInterfaceAttachService(c *bp.SdkClient) *ByteplusNetworkInterfaceAttachService {
	return &ByteplusNetworkInterfaceAttachService{
		Client: c,
	}
}

func (s *ByteplusNetworkInterfaceAttachService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusNetworkInterfaceAttachService) ReadResources(condition map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusNetworkInterfaceAttachService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		resp               *map[string]interface{}
		results            interface{}
		deviceId           interface{}
		eniType            interface{}
		ok                 bool
		networkInterfaceId string
		targetInstanceId   string
		ids                []string
	)

	if id == "" {
		id = resourceData.Id()
	}

	ids = strings.Split(id, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid network interface attach id: %v", id)
	}
	networkInterfaceId = ids[0]
	targetInstanceId = ids[1]

	req := map[string]interface{}{
		"NetworkInterfaceId": networkInterfaceId,
	}
	action := "DescribeNetworkInterfaceAttributes"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}

	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = results.(map[string]interface{}); !ok {
		return data, errors.New("value is not map")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("network interface attributes %s does not exist ", networkInterfaceId)
	}
	eniId, ok := data["NetworkInterfaceId"]
	if !ok || eniId == "" {
		return data, fmt.Errorf("network interface attributes %s does not exist ", networkInterfaceId)
	}

	eniType, ok = data["Type"]
	if !ok || eniType == "" {
		return data, errors.New("eni type does not exist")
	}
	if eniType.(string) != "secondary" {
		return data, errors.New("only secondary eni support attach/detach")
	}

	deviceId, ok = data["DeviceId"]
	if !ok {
		return data, errors.New("device id does not exist")
	}
	if len(deviceId.(string)) == 0 {
		return data, errors.New("not associate")
	}
	if deviceId.(string) != targetInstanceId {
		return data, fmt.Errorf("network interface %s does not bound target device. bound_instance_id %s, target_instance_id %s",
			networkInterfaceId, deviceId.(string), targetInstanceId)
	}
	return data, err
}

func (s *ByteplusNetworkInterfaceAttachService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			failStates = append(failStates, "Error")
			demo, err = s.ReadResource(resourceData, id)
			if err != nil && !strings.Contains(err.Error(), "not associate") {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("network interface attach status error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}
}

func (ByteplusNetworkInterfaceAttachService) WithResourceResponseHandlers(networkInterface map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return networkInterface, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}

}

func (s *ByteplusNetworkInterfaceAttachService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachNetworkInterface",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["NetworkInterfaceId"], ":", (*call.SdkParam)["InstanceId"]))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"InUse"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusNetworkInterfaceAttachService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusNetworkInterfaceAttachService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	ids := strings.Split(resourceData.Id(), ":")
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachNetworkInterface",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"NetworkInterfaceId": ids[0],
				"InstanceId":         ids[1],
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				//出现错误后重试
				return resource.Retry(15*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						return resource.NonRetryableError(fmt.Errorf("error on reading network interface on delete %q, %w", d.Id(), callErr))
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					return resource.RetryableError(callErr)
				})
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusNetworkInterfaceAttachService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusNetworkInterfaceAttachService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "vpc",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
