package cloud_monitor_contact

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCloudMonitorContactService struct {
	Client *bp.SdkClient
}

func NewService(c *bp.SdkClient) *ByteplusCloudMonitorContactService {
	return &ByteplusCloudMonitorContactService{
		Client: c,
	}
}

func (s *ByteplusCloudMonitorContactService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCloudMonitorContactService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListContacts"
		if condition != nil {
			if ids, exist := condition["Ids"]; exist && len(ids.([]interface{})) != 0 {
				action = "ListContactsByIds"
			}
		}

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
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Data", *resp)

		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Data is not Slice")
		}

		return data, err
	})
}

func (s *ByteplusCloudMonitorContactService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"Ids": []interface{}{id},
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("value is not map")
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Contact %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCloudMonitorContactService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (ByteplusCloudMonitorContactService) WithResourceResponseHandlers(data map[string]interface{}) []bp.ResourceResponseHandler {
	// 防止读取不一致
	email, exist := data["Email"]
	if exist && email == "-" {
		delete(data, "Email")
	}
	return []bp.ResourceResponseHandler{}

}

func (s *ByteplusCloudMonitorContactService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateContacts",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam, *resp)
				return resp, nil
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				data, err := bp.ObtainSdkValue("Result.Data", *resp)
				if err != nil {
					return err
				}
				if len(data.([]interface{})) == 0 {
					return fmt.Errorf("create contact failed")
				}
				d.SetId(data.([]interface{})[0].(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorContactService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var callbacks []bp.Callback

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateContacts",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["Id"] = d.Id()
				(*call.SdkParam)["Name"] = d.Get("name")
				(*call.SdkParam)["Email"] = d.Get("email")
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	callbacks = append(callbacks, callback)
	return callbacks
}

func (s *ByteplusCloudMonitorContactService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteContactsByIds",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Ids": []string{resourceData.Id()},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCloudMonitorContactService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "Ids",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		ContentType:  bp.ContentTypeJson,
		NameField:    "Name",
		IdField:      "Id",
		CollectField: "contacts",
	}
}

func (s *ByteplusCloudMonitorContactService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Volc_Observe",
		Version:     "2018-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
