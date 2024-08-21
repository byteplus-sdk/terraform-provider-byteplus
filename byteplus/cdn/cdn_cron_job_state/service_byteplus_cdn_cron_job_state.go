package cdn_cron_job_state

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

type ByteplusCdnCronJobStateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewCdnCronJobStateService(c *bp.SdkClient) *ByteplusCdnCronJobStateService {
	return &ByteplusCdnCronJobStateService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusCdnCronJobStateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusCdnCronJobStateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithPageOffsetQuery(m, "Limit", "Page", 50, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "ListCronJob"

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

		results, err = bp.ObtainSdkValue("Result.Jobs", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Jobs is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusCdnCronJobStateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 3 {
		return data, fmt.Errorf(" Invalid CdnCronJobState Id %s ", id)
	}

	req := map[string]interface{}{
		"FunctionId": ids[1],
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		var job map[string]interface{}
		if job, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if job["JobName"].(string) == ids[2] {
			data = job
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("cdn_cron_job %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusCdnCronJobStateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
			status, err = bp.ObtainSdkValue("JobState", d)
			if err != nil {
				return nil, "", err
			}

			return d, strconv.Itoa(int(status.(float64))), err
		},
	}
}

func (ByteplusCdnCronJobStateService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusCdnCronJobStateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	targetStatus := []string{"1"}
	jobAction := resourceData.Get("action").(string)
	if jobAction == "Start" {
		action = "StartCronJob"
	} else {
		action = "StopCronJob"
		targetStatus = []string{"3"}
	}

	// 根据 job 当前状态判断是否执行操作
	update, err := s.describeCurrentJobStatus(resourceData, targetStatus)
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}
	if !update {
		resourceData.SetId(fmt.Sprintf("state:%v:%v", resourceData.Get("function_id"), resourceData.Get("job_name")))
		return []bp.Callback{}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				functionId := resourceData.Get("function_id").(string)
				jobName := resourceData.Get("job_name").(string)
				d.SetId(fmt.Sprintf("state:%v:%v", functionId, jobName))
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

func (s *ByteplusCdnCronJobStateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	var action string
	targetStatus := []string{"1"}
	jobAction := resourceData.Get("action").(string)
	if jobAction == "Start" {
		action = "StartCronJob"
	} else {
		action = "StopCronJob"
		targetStatus = []string{"3"}
	}

	// 根据 job 当前状态判断是否执行操作
	update, err := s.describeCurrentJobStatus(resourceData, targetStatus)
	if err != nil {
		return []bp.Callback{{
			Err: err,
		}}
	}
	if !update {
		resourceData.SetId(fmt.Sprintf("state:%v:%v", resourceData.Get("function_id"), resourceData.Get("job_name")))
		return []bp.Callback{}
	}

	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      action,
			ConvertMode: bp.RequestConvertInConvert,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"action": {
					Ignore: true,
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				ids := strings.Split(d.Id(), ":")
				if len(ids) != 3 {
					return false, fmt.Errorf(" Invalid CdnCronJobState Id %s ", d.Id())
				}

				(*call.SdkParam)["FunctionId"] = ids[1]
				(*call.SdkParam)["JobName"] = ids[2]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getPostUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			Refresh: &bp.StateRefresh{
				Target:  targetStatus,
				Timeout: resourceData.Timeout(schema.TimeoutUpdate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusCdnCronJobStateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusCdnCronJobStateService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusCdnCronJobStateService) ReadResourceId(id string) string {
	return id
}

func (s *ByteplusCdnCronJobStateService) describeCurrentJobStatus(resourceData *schema.ResourceData, targetStatus []string) (bool, error) {
	functionId := resourceData.Get("function_id").(string)
	jobName := resourceData.Get("job_name").(string)
	data, err := s.ReadResource(resourceData, "state:"+functionId+":"+jobName)
	if err != nil {
		return false, err
	}
	status, err := bp.ObtainSdkValue("JobState", data)
	if err != nil {
		return false, err
	}
	for _, v := range targetStatus {
		// 目标状态和当前状态相同时，不执行操作
		if v == strconv.Itoa(int(status.(float64))) {
			return false, nil
		}
	}
	return true, nil
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
