package allow_list_associate

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

type ByteplusMongodbAllowListAssociateService struct {
	Client *bp.SdkClient
}

const (
	ActionAssociateAllowList      = "AssociateAllowList"
	ActionDisassociateAllowList   = "DisassociateAllowList"
	ActionDescribeAllowListDetail = "DescribeAllowListDetail"
)

func NewMongodbAllowListAssociateService(c *bp.SdkClient) *ByteplusMongodbAllowListAssociateService {
	return &ByteplusMongodbAllowListAssociateService{
		Client: c,
	}
}

func (s *ByteplusMongodbAllowListAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusMongodbAllowListAssociateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusMongodbAllowListAssociateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		targetInstanceId string
		allowListId      string
		output           *map[string]interface{}
		resultsMap       map[string]interface{}
		instanceMap      map[string]interface{}
		results          interface{}
		ok               bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid id")
	}
	targetInstanceId = ids[0]
	allowListId = ids[1]
	req := map[string]interface{}{
		"AllowListId": allowListId,
	}
	logger.Debug(logger.ReqFormat, ActionDescribeAllowListDetail, req)
	output, err = s.Client.UniversalClient.DoCall(getUniversalInfo(ActionDescribeAllowListDetail), &req)
	if err != nil {
		return data, err
	}
	logger.Debug(logger.RespFormat, ActionDescribeAllowListDetail, req, *output)
	results, err = bp.ObtainSdkValue("Result", *output)
	if err != nil {
		return data, err
	}
	if resultsMap, ok = results.(map[string]interface{}); !ok {
		return resultsMap, errors.New("Value is not map ")
	}
	if len(resultsMap) == 0 {
		return resultsMap, fmt.Errorf("MongoDB allowlist %s not exist ", allowListId)
	}
	instances := resultsMap["AssociatedInstances"].([]interface{})
	for _, instance := range instances {
		if instanceMap, ok = instance.(map[string]interface{}); !ok {
			return data, errors.New("instance is not map ")
		}
		if instanceMap["InstanceId"].(string) == targetInstanceId {
			data = resultsMap
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("MongoDB allowlist associate %s not associate ", id)
	}
	return data, err
}

func (s *ByteplusMongodbAllowListAssociateService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      1 * time.Second,
		Pending:    []string{},
		Target:     target,
		Timeout:    timeout,
		MinTimeout: 1 * time.Second,

		Refresh: func() (result interface{}, state string, err error) {
			logger.DebugInfo("Refreshing")
			output, err := s.ReadResource(resourceData, id)
			if err != nil {
				if strings.Contains(err.Error(), "not associate") {
					return output, "UnAttached", nil
				}
				return nil, "", err
			}
			return output, "Attached", nil
		},
	}
}

func (s *ByteplusMongodbAllowListAssociateService) WithResourceResponseHandlers(association map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return association, map[string]bp.ResponseConvert{}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusMongodbAllowListAssociateService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionAssociateAllowList,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				var (
					output           *map[string]interface{}
					req              map[string]interface{}
					err              error
					ok               bool
					instanceIdInter  interface{}
					allowListIdInter interface{}
					instanceId       string
					allowListId      string
				)
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				instanceIdInter, ok = (*call.SdkParam)["InstanceId"]
				if !ok {
					return output, fmt.Errorf("please input instance_id")
				}
				instanceId, ok = instanceIdInter.(string)
				if !ok {
					return output, fmt.Errorf("type of instanceIdInter is not string")
				}
				allowListIdInter, ok = (*call.SdkParam)["AllowListId"]
				if !ok {
					return output, fmt.Errorf("please input allow_list_id")
				}
				allowListId, ok = allowListIdInter.(string)
				if !ok {
					return output, fmt.Errorf("type of allowListIdInter is not string")
				}
				req = make(map[string]interface{})
				req["InstanceIds"] = []string{instanceId}
				req["AllowListIds"] = []string{allowListId}
				output, err = s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), &req)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint((*call.SdkParam)["InstanceId"], ":", (*call.SdkParam)["AllowListId"]))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Attached"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongodbAllowListAssociateService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusMongodbAllowListAssociateService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionDisassociateAllowList,
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				id := s.ReadResourceId(d.Id())
				ids := strings.Split(id, ":")
				instanceId := ids[0]
				allowListId := ids[1]
				(*call.SdkParam)["InstanceIds"] = []string{instanceId}
				(*call.SdkParam)["AllowListIds"] = []string{allowListId}
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"UnAttached"},
				Timeout: resourceData.Timeout(schema.TimeoutDelete),
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusMongodbAllowListAssociateService) DatasourceResources(data *schema.ResourceData, resource2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType: bp.ContentTypeJson,
	}
}

func (s *ByteplusMongodbAllowListAssociateService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "mongodb",
		Action:      actionName,
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
	}
}
