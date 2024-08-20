package cdn_edge_function_publish

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

type ByteplusCdnEdgeFunctionPublishService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnEdgeFunctionPublishService(c *bp.SdkClient) *ByteplusCdnEdgeFunctionPublishService {
	return &ByteplusCdnEdgeFunctionPublishService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnEdgeFunctionPublishService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnEdgeFunctionPublishService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	m["OrderType"] = "create_time"
	return bp.WithPageOffsetQuery(m, "Limit", "Page", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListTicket"

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

		results, err = bp.ObtainSdkValue("Result.Tickets", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Tickets is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusCdnEdgeFunctionPublishService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		result interface{}
		ok     bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf(" Invalid CdnFunctionPublish Id %s ", id)
	}
	ticketId, err := strconv.Atoi(ids[1])
	if err != nil {
		return data, fmt.Errorf(" TicketId cannot convert to int: %v ", ids[1])
	}

	action := "GetTicket"
	req := map[string]interface{}{
		"FunctionId": ids[0],
		"TicketId":   ticketId,
	}
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, action, resp)
	result, err = bp.ObtainSdkValue("Result.Ticket", *resp)
	if err != nil {
		return data, err
	}
	if data, ok = result.(map[string]interface{}); !ok {
		return data, errors.New("Value is not map ")
	}
	if len(data) == 0 {
		return data, fmt.Errorf("edge_function_publish %s is not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnEdgeFunctionPublishService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			status, err = bp.ObtainSdkValue("Status", d)
			if err != nil {
				return nil, "", err
			}

			return d, strconv.Itoa(int(status.(float64))), err
		},
	}
}

func (ByteplusCdnEdgeFunctionPublishService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnEdgeFunctionPublishService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var (
		action       string
		targetStatus []string
	)
	publishAction := resourceData.Get("publish_action").(string)
	if publishAction == "FullPublish" {
		action = "FullPublish"
		targetStatus = []string{"200"}
	} else if publishAction == "CanaryPublish" {
		action = "CanaryPublish"
		targetStatus = []string{"100"}
	} else {
		action = "SnapshotPublish"
		publishType := resourceData.Get("publish_type")
		if publishType == 100 {
			targetStatus = []string{"100"}
		} else if publishType == 200 {
			targetStatus = []string{"200"}
		}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				ticketId, _ := bp.ObtainSdkValue("Result.TicketId", *resp)
				functionId := d.Get("function_id")
				d.SetId(functionId.(string) + ":" + strconv.Itoa(int(ticketId.(float64))))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  targetStatus,
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnEdgeFunctionPublishService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCdnEdgeFunctionPublishService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCdnEdgeFunctionPublishService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "tickets",
		ResponseConverts: map[string]bp.ResponseConvert{
			"Id": {
				TargetField: "ticket_id",
			},
		},
	}
}

func (s *ByteplusCdnEdgeFunctionPublishService) ReadResourceId(id string) string {
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
