package waf_instance_ctl

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusWafInstanceCtlService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewWafInstanceCtlService(c *bp.SdkClient) *ByteplusWafInstanceCtlService {
	return &ByteplusWafInstanceCtlService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusWafInstanceCtlService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusWafInstanceCtlService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusWafInstanceCtlService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Region":      id,
		"ProjectName": resourceData.Get("project_name"),
	}
	client := s.Client.UniversalClient
	action := "GetInstanceCtl"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err = client.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}

	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = results.(map[string]interface{}); !ok {
		return data, errors.New("Value is not map ")
	}

	return data, err
}

func (s *ByteplusWafInstanceCtlService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusWafInstanceCtlService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateInstanceCtl",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Region"] = client.Region
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(client.Region)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusWafInstanceCtlService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusWafInstanceCtlService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateInstanceCtl",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"project_name": {
					TargetField: "ProjectName",
					ForceGet:    true,
				},
				"allow_enable": {
					TargetField: "AllowEnable",
					ForceGet:    true,
				},
				"block_enable": {
					TargetField: "BlockEnable",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Region"] = client.Region
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

func (s *ByteplusWafInstanceCtlService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	//callback := bp.Callback{
	//	Call: bp.SdkCall{
	//		Action:      "UpdateInstanceCtl",
	//		ConvertMode: bp.RequestConvertIgnore,
	//		ContentType: bp.ContentTypeJson,
	//		BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
	//			(*call.SdkParam)["Region"] = client.Region
	//			(*call.SdkParam)["AllowEnable"] = 0
	//			(*call.SdkParam)["BlockEnable"] = 0
	//			return true, nil
	//		},
	//		ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
	//			logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
	//			return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
	//		},
	//		AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
	//			return s.checkResourceUtilRemoved(d, 5*time.Minute)
	//		},
	//		CallError: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall, baseErr error) error {
	//			//出现错误后重试
	//			return resource.Retry(5*time.Minute, func() *resource.RetryError {
	//				_, callErr := s.ReadResource(d, "")
	//				if callErr != nil {
	//					if bp.ResourceNotFoundError(callErr) {
	//						return nil
	//					} else {
	//						return resource.NonRetryableError(fmt.Errorf("error on  reading waf domain on delete %q, %w", d.Id(), callErr))
	//					}
	//				}
	//				_, callErr = call.ExecuteCall(d, client, call)
	//				if callErr == nil {
	//					return nil
	//				}
	//				return resource.RetryableError(callErr)
	//			})
	//		},
	//	},
	//}
	logger.Debug(logger.ReqFormat, "RemoveResource", "Remove only from tf management")
	return []bp.Callback{}
}

func (s *ByteplusWafInstanceCtlService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusWafInstanceCtlService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "waf",
		Version:     "2023-12-25",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
		RegionType:  bp.Global,
	}
}

func (s *ByteplusWafInstanceCtlService) checkResourceUtilRemoved(d *schema.ResourceData, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		instanceCtl, _ := s.ReadResource(d, d.Id())
		logger.Debug(logger.RespFormat, "instanceCtl", instanceCtl)

		// 能查询成功代表还在删除中，重试
		allowEnableInt, ok := instanceCtl["AllowEnable"].(float64)
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("AllowEnable is not float64"))
		}
		blockEnable, ok := instanceCtl["BlockEnable"].(float64)
		if !ok {
			return resource.NonRetryableError(fmt.Errorf("BlockEnable is not float64"))
		}
		if int(allowEnableInt) == 1 || int(blockEnable) == 1 {
			return resource.RetryableError(fmt.Errorf("resource still in removing status "))
		} else {
			if int(allowEnableInt) == 0 && int(blockEnable) == 0 {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("instanceCtl status is not disable "))
			}
		}
	})
}
