package cdn_domain_enabler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnDomainEnablerService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnDomainEnablerService(c *bp.SdkClient) *ByteplusCdnDomainEnablerService {
	return &ByteplusCdnDomainEnablerService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnDomainEnablerService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnDomainEnablerService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageNumberQuery(m, "PageSize", "PageNum", 100, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListCdnDomains"

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

		results, err = bp.ObtainSdkValue("Result.Data", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Domains is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusCdnDomainEnablerService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return nil, fmt.Errorf("err cdn domain enabler id")
	}

	req := map[string]interface{}{
		"Domain": ids[1],
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
		return data, fmt.Errorf("cdn_domain_enabler %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnDomainEnablerService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("cdn_domain_enabler status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusCdnDomainEnablerService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "StartCdnDomain",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				domain := d.Get("domain").(string)
				req := map[string]interface{}{
					"Domain": domain,
				}
				logger.Debug(logger.ReqFormat, "ListCdnDomains", req)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("ListCdnDomains"), &req)
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.RespFormat, "ListCdnDomains", req, *resp)
				status, err := bp.ObtainSdkValue("Status", *resp)
				if err != nil {
					return nil, err
				}
				// 已经online就什么也不做，直接跳过
				if status.(string) == "online" {
					return nil, nil
				} else {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				}
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				domain := d.Get("domain").(string)
				d.SetId(fmt.Sprintf("enabler:%s", domain))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"online"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusCdnDomainEnablerService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnDomainEnablerService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCdnDomainEnablerService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "StopCdnDomain",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				(*call.SdkParam)["Domain"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				domain := d.Get("domain").(string)
				req := map[string]interface{}{
					"Domain": domain,
				}
				logger.Debug(logger.ReqFormat, "ListCdnDomains", req)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo("ListCdnDomains"), &req)
				if err != nil {
					return nil, err
				}
				logger.Debug(logger.RespFormat, "ListCdnDomains", req, *resp)
				status, err := bp.ObtainSdkValue("Status", *resp)
				if err != nil {
					return nil, err
				}
				if status.(string) == "offline" {
					return nil, nil
				} else {
					logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
					resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, resp, err)
					return resp, err
				}
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"offline"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnDomainEnablerService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusCdnDomainEnablerService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
