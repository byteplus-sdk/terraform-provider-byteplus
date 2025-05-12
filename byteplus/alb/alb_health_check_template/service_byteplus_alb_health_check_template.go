package alb_health_check_template

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusAlbHealthCheckTemplateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbHealthCheckTemplateService(c *bp.SdkClient) *ByteplusAlbHealthCheckTemplateService {
	return &ByteplusAlbHealthCheckTemplateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbHealthCheckTemplateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbHealthCheckTemplateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeHealthCheckTemplates"

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
		results, err = bp.ObtainSdkValue("Result.HealthCheckTemplates", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.HealthCheckTemplates is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusAlbHealthCheckTemplateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"HealthCheckTemplateIds.1": id,
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
		return data, fmt.Errorf("alb_health_check_template %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbHealthCheckTemplateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusAlbHealthCheckTemplateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateHealthCheckTemplates",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"health_check_template_name": {
					TargetField: "HealthCheckTemplates.1.HealthCheckTemplateName",
				},
				"description": {
					TargetField: "HealthCheckTemplates.1.Description",
				},
				"health_check_interval": {
					TargetField: "HealthCheckTemplates.1.HealthCheckInterval",
				},
				"health_check_timeout": {
					TargetField: "HealthCheckTemplates.1.HealthCheckTimeout",
				},
				"healthy_threshold": {
					TargetField: "HealthCheckTemplates.1.HealthyThreshold",
				},
				"unhealthy_threshold": {
					TargetField: "HealthCheckTemplates.1.UnhealthyThreshold",
				},
				"health_check_method": {
					TargetField: "HealthCheckTemplates.1.HealthCheckMethod",
				},
				"health_check_domain": {
					TargetField: "HealthCheckTemplates.1.HealthCheckDomain",
				},
				"health_check_uri": {
					TargetField: "HealthCheckTemplates.1.HealthCheckURI",
				},
				"health_check_http_code": {
					TargetField: "HealthCheckTemplates.1.HealthCheckHttpCode",
				},
				"health_check_protocol": {
					TargetField: "HealthCheckTemplates.1.HealthCheckProtocol",
				},
				"health_check_http_version": {
					TargetField: "HealthCheckTemplates.1.HealthCheckHttpVersion",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				ids, err := bp.ObtainSdkValue("Result.HealthCheckTemplateIDs", *resp)
				if err != nil {
					return err
				}
				idArr, ok := ids.([]interface{})
				if !ok || len(idArr) == 0 {
					return fmt.Errorf("ids is invalid")
				}
				d.SetId(idArr[0].(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusAlbHealthCheckTemplateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbHealthCheckTemplateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyHealthCheckTemplatesAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"health_check_template_name": {
					TargetField: "HealthCheckTemplates.1.HealthCheckTemplateName",
				},
				"description": {
					TargetField: "HealthCheckTemplates.1.Description",
				},
				"health_check_interval": {
					TargetField: "HealthCheckTemplates.1.HealthCheckInterval",
				},
				"health_check_timeout": {
					TargetField: "HealthCheckTemplates.1.HealthCheckTimeout",
				},
				"healthy_threshold": {
					TargetField: "HealthCheckTemplates.1.HealthyThreshold",
				},
				"unhealthy_threshold": {
					TargetField: "HealthCheckTemplates.1.UnhealthyThreshold",
				},
				"health_check_method": {
					TargetField: "HealthCheckTemplates.1.HealthCheckMethod",
				},
				"health_check_domain": {
					TargetField: "HealthCheckTemplates.1.HealthCheckDomain",
				},
				"health_check_uri": {
					TargetField: "HealthCheckTemplates.1.HealthCheckURI",
				},
				"health_check_http_code": {
					TargetField: "HealthCheckTemplates.1.HealthCheckHttpCode",
				},
				"health_check_protocol": {
					TargetField: "HealthCheckTemplates.1.HealthCheckProtocol",
				},
				"health_check_http_version": {
					TargetField: "HealthCheckTemplates.1.HealthCheckHttpVersion",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["HealthCheckTemplates.1.HealthCheckTemplateId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbHealthCheckTemplateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteHealthCheckTemplates",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"HealthCheckTemplateIds.1": resourceData.Id(),
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

func (s *ByteplusAlbHealthCheckTemplateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "HealthCheckTemplateIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "HealthCheckTemplateName",
		IdField:      "HealthCheckTemplateId",
		CollectField: "health_check_templates",
		ResponseConverts: map[string]bp.ResponseConvert{
			"HealthCheckTemplateId": {
				TargetField: "id",
				KeepDefault: true,
			},
			"HealthCheckURI": {
				TargetField: "health_check_uri",
			},
		},
	}
}

func (s *ByteplusAlbHealthCheckTemplateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "alb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
