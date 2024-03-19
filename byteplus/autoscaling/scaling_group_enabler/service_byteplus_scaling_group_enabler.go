package scaling_group_enabler

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

type ByteplusScalingGroupEnablerService struct {
	Client *bp.SdkClient
}

func (s *ByteplusScalingGroupEnablerService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		universalClient := s.Client.UniversalClient
		action := "DescribeScalingGroups"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = universalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = universalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, err
			}
		}
		logger.Debug(logger.RespFormat, action, action, *resp)
		results, err = bp.ObtainSdkValue("Result.ScalingGroups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.ScalingGroups is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusScalingGroupEnablerService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("Invalid ScalingGroupEnable Id ")
	}
	req := map[string]interface{}{
		"ScalingGroupIds.1": ids[1],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("ScalingGroup %s not exist ", ids[1])
	}
	state := data["LifecycleState"].(string)
	if state != "Active" && state != "Locked" {
		return data, fmt.Errorf("ScalingGroup %s is not active", ids[1])
	}
	return data, err
}

func (s *ByteplusScalingGroupEnablerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusScalingGroupEnablerService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusScalingGroupEnablerService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	param := &map[string]interface{}{
		"ScalingGroupId": data.Get("scaling_group_id").(string),
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "EnableScalingGroup",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam:    param,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.ScalingGroupId", *resp)
				d.SetId(fmt.Sprintf("enable:%s", id.(string)))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusScalingGroupEnablerService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusScalingGroupEnablerService) RemoveResource(data *schema.ResourceData, r *schema.Resource) []bp.Callback {
	param := &map[string]interface{}{
		"ScalingGroupId": data.Get("scaling_group_id").(string),
	}
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisableScalingGroup",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam:    param,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
				return resource.Retry(5*time.Minute, func() *resource.RetryError {
					_, callErr := s.ReadResource(d, "")
					if callErr != nil {
						if bp.ResourceNotFoundError(callErr) {
							return nil
						} else {
							return resource.NonRetryableError(fmt.Errorf("error on  reading scaling group enabler on delete %q, %w", d.Id(), callErr))
						}
					}
					_, callErr = call.ExecuteCall(d, client, call)
					if callErr == nil {
						return nil
					}
					// 伸缩组lock，重试
					if strings.Contains(callErr.Error(), "ErrInvalidGroupStatus") {
						return resource.RetryableError(callErr)
					}
					return resource.NonRetryableError(callErr)
				})
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusScalingGroupEnablerService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusScalingGroupEnablerService) ReadResourceId(id string) string {
	return id
}

func NewScalingGroupEnablerService(client *bp.SdkClient) *ByteplusScalingGroupEnablerService {
	return &ByteplusScalingGroupEnablerService{
		Client: client,
	}
}

func (s *ByteplusScalingGroupEnablerService) GetClient() *bp.SdkClient {
	return s.Client
}

func getUniversalInfo(action string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "auto_scaling",
		Action:      action,
		Version:     "2020-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
	}
}
