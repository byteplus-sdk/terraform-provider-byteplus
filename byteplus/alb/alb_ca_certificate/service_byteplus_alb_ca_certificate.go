package alb_ca_certificate

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

type ByteplusAlbCaCertificateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewAlbCaCertificateService(c *bp.SdkClient) *ByteplusAlbCaCertificateService {
	return &ByteplusAlbCaCertificateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusAlbCaCertificateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusAlbCaCertificateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeCACertificates"
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
		results, err = bp.ObtainSdkValue("Result.CACertificates", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.CACertificates is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusAlbCaCertificateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"CACertificateIds.1": id,
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
		return data, fmt.Errorf("alb_ca_certificate %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusAlbCaCertificateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d          map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Failed")

			if err = resource.Retry(5*time.Minute, func() *resource.RetryError {
				d, err = s.ReadResource(resourceData, id)
				if err != nil {
					if bp.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("alb_ca_certificate status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusAlbCaCertificateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "UploadCACertificate",
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"ca_certificate_name": {
					TargetField: "CACertificateName",
				},
				"ca_certificate": {
					TargetField: "CACertificate",
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id, _ := bp.ObtainSdkValue("Result.CACertificateId", *resp)
				d.SetId(id.(string))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusAlbCaCertificateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusAlbCaCertificateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyCACertificateAttributes",
			ConvertMode: bp.RequestConvertInConvert,
			Convert: map[string]bp.RequestConvert{
				"ca_certificate_name": {
					TargetField: "CACertificateName",
					ForceGet:    true,
				},
				"description": {
					TargetField: "Description",
					ForceGet:    true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["CACertificateId"] = d.Id()
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Active"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusAlbCaCertificateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteCACertificate",
			ConvertMode: bp.RequestConvertIgnore,
			SdkParam: &map[string]interface{}{
				"CACertificateId": resourceData.Id(),
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

func (s *ByteplusAlbCaCertificateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "CACertificateIds",
				ConvertType: bp.ConvertWithN,
			},
			"ca_certificate_name": {
				TargetField: "CACertificateName",
			},
		},
		NameField:    "CACertificateName",
		IdField:      "CACertificateId",
		CollectField: "certificates",
		ResponseConverts: map[string]bp.ResponseConvert{
			"CACertificateId": {
				TargetField: "ca_certificate_id",
			},
			"CACertificateName": {
				TargetField: "ca_certificate_name",
			},
		},
	}
}

func (s *ByteplusAlbCaCertificateService) ReadResourceId(id string) string {
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

func (s *ByteplusAlbCaCertificateService) ProjectTrn() *bp.ProjectTrn {
	return &bp.ProjectTrn{
		ServiceName:          "alb",
		ResourceType:         "certificate",
		ProjectResponseField: "ProjectName",
		ProjectSchemaField:   "project_name",
	}
}
