package organization_service_control_policy_attachment

import (
	"errors"
	"fmt"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
	"time"
)

type ByteplusServiceControlPolicyAttachmentService struct {
	Client *bp.SdkClient
}

func NewServiceControlPolicyAttachmentService(c *bp.SdkClient) *ByteplusServiceControlPolicyAttachmentService {
	return &ByteplusServiceControlPolicyAttachmentService{
		Client: c,
	}
}

func (s *ByteplusServiceControlPolicyAttachmentService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusServiceControlPolicyAttachmentService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	action := "ListTargetAttachmentsForServiceControlPolicy"
	logger.Debug(logger.ReqFormat, action, m)
	if m == nil {
		resp, err = s.Client.UniversalClient.DoCall(postUniversalInfo(action), nil)
		if err != nil {
			return data, err
		}
	} else {
		resp, err = s.Client.UniversalClient.DoCall(postUniversalInfo(action), &m)
		if err != nil {
			return data, err
		}
	}

	logger.Debug(logger.RespFormat, action, m, *resp)

	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if results == nil {
		results = []interface{}{}
	}
	if data, ok = results.([]interface{}); !ok {
		return data, errors.New(" Result is not Slice")
	}
	return data, err
}

func (s *ByteplusServiceControlPolicyAttachmentService) ReadResource(resourceData *schema.ResourceData, roleId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if roleId == "" {
		roleId = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(roleId, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("import id is invalid")
	}
	req := map[string]interface{}{
		"PolicyID": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		} else if ids[1] == data["TargetID"].(string) {
			return data, err
		}
	}
	return data, fmt.Errorf("service control policy attachment %s not exist ", roleId)
}

func (s *ByteplusServiceControlPolicyAttachmentService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusServiceControlPolicyAttachmentService) WithResourceResponseHandlers(rolePolicyAttachment map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return rolePolicyAttachment, map[string]bp.ResponseConvert{
			"TargetID": {
				TargetField: "target_id",
			},
			"PolicyID": {
				TargetField: "policy_id",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusServiceControlPolicyAttachmentService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	createPolicyAttachmentCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AttachServiceControlPolicy",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"policy_id": {
					TargetField: "PolicyID",
				},
				"target_id": {
					TargetField: "TargetID",
				},
				"target_type": {
					Convert: func(data *schema.ResourceData, old interface{}) interface{} {
						ty := 0
						switch old.(string) {
						case "OU":
							ty = 1
						case "Account":
							ty = 2
						}
						return ty
					},
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(postUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("%s:%s", d.Get("policy_id"), d.Get("target_id")))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{createPolicyAttachmentCallback}
}

func (s *ByteplusServiceControlPolicyAttachmentService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusServiceControlPolicyAttachmentService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	deleteRoleCallback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DetachServiceControlPolicy",
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["PolicyID"] = ids[0]
				(*call.SdkParam)["TargetID"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(postUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return "lock-Organization"
			},
		},
	}
	return []bp.Callback{deleteRoleCallback}
}

func (s *ByteplusServiceControlPolicyAttachmentService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusServiceControlPolicyAttachmentService) ReadResourceId(id string) string {
	return id
}

func postUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "organization",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
