package cdn_edge_function_associate

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusCdnEdgeFunctionAssociateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnEdgeFunctionAssociateService(c *bp.SdkClient) *ByteplusCdnEdgeFunctionAssociateService {
	return &ByteplusCdnEdgeFunctionAssociateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnEdgeFunctionAssociateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Page", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListSparrowDomains"

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

		results, err = bp.ObtainSdkValue("Result.Domains", *resp)
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

func (s *ByteplusCdnEdgeFunctionAssociateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf(" Invalid CdnFunctionAssociate Id %s ", id)
	}

	req := map[string]interface{}{
		"FunctionId": ids[0],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		var result map[string]interface{}
		if result, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if result["Domain"].(string) == ids[1] {
			data = result
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("cdn_edge_function_associate %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnEdgeFunctionAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				d      map[string]interface{}
				status interface{}
			)
			d, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("BindStatus", d)
			if err != nil {
				return nil, "", err
			}

			return d, strconv.Itoa(int(status.(float64))), err
		},
	}
}

func (ByteplusCdnEdgeFunctionAssociateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "SparrowBindDomains",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"domain": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				domain := d.Get("domain").(string)
				(*call.SdkParam)["Domains"] = []string{domain}

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				functionId := d.Get("function_id").(string)
				domain := d.Get("domain").(string)
				d.SetId(fmt.Sprintf(functionId + ":" + domain))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"1"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "SparrowUnBindDomains",
			ConvertMode: bp.RequestConvertIgnore,
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"Id": resourceData.Id(),
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 2 {
					return nil, fmt.Errorf(" Invalid CdnFunctionAssociate Id %s ", d.Id())
				}
				(*call.SdkParam)["FunctionId"] = ids[0]
				(*call.SdkParam)["Domains"] = []string{ids[1]}

				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"2"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusCdnEdgeFunctionAssociateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}

func getPostUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "CDN",
		Version:     "2021-03-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
