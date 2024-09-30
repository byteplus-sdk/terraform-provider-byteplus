package iam_saml_provider

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

type ByteplusIamSamlProviderService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewIamSamlProviderService(c *bp.SdkClient) *ByteplusIamSamlProviderService {
	return &ByteplusIamSamlProviderService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusIamSamlProviderService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusIamSamlProviderService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Offset", 100, 0, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListSAMLProviders"

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
		results, err = bp.ObtainSdkValue("Result.SAMLProviders", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.SAMLProviders is not Slice")
		}
		for index, ele := range data {
			provider := ele.(map[string]interface{})
			query := map[string]interface{}{
				// 跟API文档不对应，文档写的返回字段为SAMLProviderName，实际为ProviderName
				"SAMLProviderName": provider["ProviderName"],
			}
			action = "GetSAMLProvider"
			logger.Debug(logger.ReqFormat, action, query)
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &query)
			if err != nil {
				return data, err
			}
			logger.Debug(logger.RespFormat, action, query, *resp)
			document, err := bp.ObtainSdkValue("Result.EncodedSAMLMetadataDocument", *resp)
			if err != nil {
				return data, err
			}
			data[index].(map[string]interface{})["EncodedSAMLMetadataDocument"] = document
		}
		return data, err
	})
}

func (s *ByteplusIamSamlProviderService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		temp    map[string]interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		if temp, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		} else if temp["ProviderName"].(string) == id {
			data = temp
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("iam_saml_provider %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusIamSamlProviderService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusIamSamlProviderService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateSAMLProvider",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"encoded_saml_metadata_document": {
					TargetField: "EncodedSAMLMetadataDocument",
				},
				"saml_provider_name": {
					TargetField: "SAMLProviderName",
				},
				"sso_type": {
					TargetField: "SSOType",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := d.Get("saml_provider_name")
				d.SetId(id.(string))
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusIamSamlProviderService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusIamSamlProviderService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UpdateSAMLProvider",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"encoded_saml_metadata_document": {
					TargetField: "NewEncodedSAMLMetadataDocument",
				},
				"description": {
					TargetField: "Description",
				},
				"status": {
					TargetField: "Status",
				},
				"sso_type": {
					TargetField: "SSOType",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["SAMLProviderName"] = d.Id()
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

func (s *ByteplusIamSamlProviderService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteSAMLProvider",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"SAMLProviderName": resourceData.Id(),
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

func (s *ByteplusIamSamlProviderService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		NameField:    "SAMLProviderName",
		IdField:      "Trn",
		CollectField: "providers",
		ResponseConverts: map[string]bp.ResponseConvert{
			"ProviderName": {
				TargetField: "saml_provider_name",
			},
			"EncodedSAMLMetadataDocument": {
				TargetField: "encoded_saml_metadata_document",
			},
			"SSOType": {
				TargetField: "sso_type",
			},
		},
	}
}

func (s *ByteplusIamSamlProviderService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "iam",
		Version:     "2018-01-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
