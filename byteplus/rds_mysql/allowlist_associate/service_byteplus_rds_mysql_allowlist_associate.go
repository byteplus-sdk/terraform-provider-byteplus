package allowlist_associate

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/rds_mysql/rds_mysql_instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRdsMysqlAllowListAssociateService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func (s *ByteplusRdsMysqlAllowListAssociateService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlAllowListAssociateService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusRdsMysqlAllowListAssociateService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results     interface{}
		resultsMap  map[string]interface{}
		instanceMap map[string]interface{}
		instances   []interface{}
		ok          bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return data, err
	}
	req := map[string]interface{}{
		"AllowListId": ids[1],
	}
	action := "DescribeAllowListDetail"
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &req)
	if err != nil {
		return data, err
	}
	results, err = bp.ObtainSdkValue("Result", *resp)
	if err != nil {
		return data, err
	}
	if resultsMap, ok = results.(map[string]interface{}); !ok {
		return resultsMap, errors.New("Value is not map ")
	}
	if len(resultsMap) == 0 {
		return resultsMap, fmt.Errorf("Rds allowlist %s not exist ", ids[1])
	}
	logger.Debug(logger.ReqFormat, action, resultsMap)
	instances = resultsMap["AssociatedInstances"].([]interface{})
	logger.Debug(logger.ReqFormat, action, instances)
	for _, instance := range instances {
		if instanceMap, ok = instance.(map[string]interface{}); !ok {
			return data, errors.New("instance is not map ")
		}
		if len(instanceMap) == 0 {
			continue
		}
		if instanceMap["InstanceId"].(string) == ids[0] {
			data = resultsMap
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Rds allowlist associate %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusRdsMysqlAllowListAssociateService) RefreshResourceState(data *schema.ResourceData, strings []string, duration time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, map[string]bp.ResponseConvert{}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := data.Get("instance_id").(string)
	allowListId := data.Get("allow_list_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "AssociateAllowList",
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceIds":  []string{instanceId},
				"AllowListIds": []string{allowListId},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprint(instanceId, ":", allowListId))
				return nil
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    data.Timeout(schema.TimeoutCreate),
					ResourceId: instanceId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	instanceId := data.Get("instance_id").(string)
	allowListId := data.Get("allow_list_id").(string)
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DisassociateAllowList",
			ContentType: bp.ContentTypeJson,
			SdkParam: &map[string]interface{}{
				"InstanceIds":  []string{instanceId},
				"AllowListIds": []string{allowListId},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				err := bp.CheckResourceUtilRemoved(d, s.ReadResource, 10*time.Minute)
				return err
			},
			LockId: func(d *schema.ResourceData) string {
				return instanceId
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				rds_mysql_instance.NewRdsMysqlInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    data.Timeout(schema.TimeoutDelete),
					ResourceId: instanceId,
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusRdsMysqlAllowListAssociateService) ReadResourceId(id string) string {
	return id
}

func NewRdsMysqlAllowListAssociateService(client *bp.SdkClient) *ByteplusRdsMysqlAllowListAssociateService {
	return &ByteplusRdsMysqlAllowListAssociateService{
		Client:     client,
		Dispatcher: &bp.Dispatcher{},
	}
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_mysql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
