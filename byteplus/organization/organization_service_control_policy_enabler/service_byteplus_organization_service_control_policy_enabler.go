package organization_service_control_policy_enabler

import (
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusOrganizationServiceControlPolicyEnablerService struct {
	Client *bp.SdkClient
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	universalClient := s.Client.UniversalClient
	action := "GetServiceControlPolicyEnablement"
	resp, err := universalClient.DoCall(getUniversalInfo(action), &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	logger.Debug(logger.ReqFormat, action, *resp)
	status, err := bp.ObtainSdkValue("Result.Status", *resp)
	if err != nil {
		return nil, err
	}

	if status != "Enabled" {
		return data, fmt.Errorf(" Organization Service Control Policy is not Enabled")
	}
	return map[string]interface{}{
		"Status": status,
	}, nil
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "EnableServiceControlPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				// 先检查是否开启
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("GetServiceControlPolicyEnablement"), &map[string]interface{}{})
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.ReqFormat, "GetServiceControlPolicyEnablement", *resp)
				status, err := bp.ObtainSdkValue("Result.Status", *resp)
				if err != nil {
					return nil, err
				}
				if status == "Enabled" {
					return nil, nil // 不需要再重复开启了
				}

				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId("organization:service_control_policy_enable")
				time.Sleep(3 * time.Second)
				return nil
			},
			// 必须顺序执行，否则并发失败
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisableServiceControlPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				time.Sleep(3 * time.Second)
				return nil
			},
			// 必须顺序执行，否则并发失败
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) ReadResourceId(id string) string {
	return id
}

func NewService(client *bp.SdkClient) *ByteplusOrganizationServiceControlPolicyEnablerService {
	return &ByteplusOrganizationServiceControlPolicyEnablerService{
		Client: client,
	}
}

func (s *ByteplusOrganizationServiceControlPolicyEnablerService) GetClient() *bp.SdkClient {
	return s.Client
}

func getUniversalInfo(action string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "organization",
		Action:      action,
		Version:     "2022-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
