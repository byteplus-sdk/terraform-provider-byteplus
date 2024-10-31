package organization_service_control_policy

import (
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusPolicyService struct {
	Client *bp.SdkClient
}

func NewService(c *bp.SdkClient) *ByteplusPolicyService {
	return &ByteplusPolicyService{
		Client: c,
	}
}

func (s *ByteplusPolicyService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusPolicyService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) (data []interface{}, err error) {
		action := "ListServiceControlPolicies"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return nil, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return nil, err
			}
		}

		results, err = bp.ObtainSdkValue("Result.ServiceControlPolicies", *resp)
		if err != nil {
			return nil, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ServiceControlPolicies is not Slice")
		}

		// 获取每一个策略内容
		for _, ele := range data {
			temp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("GetServiceControlPolicy"), &map[string]interface{}{
				"PolicyID": ele.(map[string]interface{})["PolicyID"],
			})
			if err != nil {
				return nil, err
			}
			statement, err := bp.ObtainSdkValue("Result.Statement", *temp)
			if err != nil {
				return nil, err
			}
			ele.(map[string]interface{})["Statement"] = statement
		}
		return data, err
	})
}

func (s *ByteplusPolicyService) ReadResource(resourceData *schema.ResourceData, policyId string) (data map[string]interface{}, err error) {
	if policyId == "" {
		policyId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"PolicyID": policyId,
	}
	temp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("GetServiceControlPolicy"), &req)
	if err != nil {
		return nil, err
	}
	res, err := bp.ObtainSdkValue("Result", *temp)
	if err != nil {
		return nil, err
	}
	return res.(map[string]interface{}), nil
}

func (s *ByteplusPolicyService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusPolicyService) WithResourceResponseHandlers(policy map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusPolicyService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createIamPolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateServiceControlPolicy",
			ConvertMode: bp.RequestConvertAll,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(postUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				policyId, err := bp.ObtainSdkValue("Result.PolicyId", *resp)
				if err != nil {
					return err
				}
				d.SetId(policyId.(string))

				// 单独处理
				time.Sleep(2 * time.Second)
				return nil
			},
			// 必须顺序执行，否则并发失败
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{createIamPolicyCallback}
}

func (s *ByteplusPolicyService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	updatePolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateServiceControlPolicy",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"policy_name": {
					ForceGet: true,
				},
				"description": {
					ForceGet: true,
				},
				"statement": {
					ForceGet: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["PolicyID"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(postUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				// 单独处理
				time.Sleep(2 * time.Second)
				return nil
			},
			// 必须顺序执行，否则并发失败
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{updatePolicyCallback}
}

func (s *ByteplusPolicyService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deletePolicyCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteServiceControlPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["PolicyID"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(postUniversalInfo(call.Action), call.SdkParam)
			},
			// 必须顺序执行，否则并发失败
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{deletePolicyCallback}
}

func (s *ByteplusPolicyService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ResponseConverts: map[string]bp.ResponseConvert{
			"PolicyID": {
				TargetField: "id",
			},
		},
		NameField:    "PolicyName",
		IdField:      "PolicyID",
		CollectField: "policies",
	}
}

func (s *ByteplusPolicyService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "organization",
		Version:     "2022-01-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}

func postUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "organization",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		Action:      actionName,
		ContentType: bp.ApplicationJSON,
	}
}
