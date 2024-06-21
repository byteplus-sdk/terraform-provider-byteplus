package rds_postgresql_database

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

type ByteplusRdsPostgresqlDatabaseService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsPostgresqlDatabaseService(c *bp.SdkClient) *ByteplusRdsPostgresqlDatabaseService {
	return &ByteplusRdsPostgresqlDatabaseService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsPostgresqlDatabaseService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsPostgresqlDatabaseService) ReadResources(m map[string]interface{}) ([]interface{}, error) {
	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 10, 1, func(condition map[string]interface{}) (data []interface{}, err error) {
		var (
			resp    *map[string]interface{}
			results interface{}
			ok      bool
		)
		action := "DescribeDatabases"
		logger.Debug(logger.ReqFormat, action, condition)
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

		results, err = bp.ObtainSdkValue("Result.Databases", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.Databases is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusRdsPostgresqlDatabaseService) ReadResource(resourceData *schema.ResourceData, rdsDatabaseId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if rdsDatabaseId == "" {
		rdsDatabaseId = s.ReadResourceId(resourceData.Id())
	}

	ids := strings.Split(rdsDatabaseId, ":")
	if len(ids) != 2 {
		return map[string]interface{}{}, fmt.Errorf("invalid database id")
	}

	instanceId := ids[0]
	dbName := ids[1]

	req := map[string]interface{}{
		"InstanceId": instanceId,
		"DBName":     dbName,
	}
	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	for _, v := range results {
		var dbMap map[string]interface{}
		if dbMap, ok = v.(map[string]interface{}); !ok {
			return data, errors.New("Value is not map ")
		}
		if dbName == dbMap["DBName"].(string) {
			data = dbMap
			break
		}
	}
	if len(data) == 0 {
		return data, fmt.Errorf("RDS postgresql database %s not exist ", rdsDatabaseId)
	}

	return data, err
}

func (s *ByteplusRdsPostgresqlDatabaseService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{}
}

func (s *ByteplusRdsPostgresqlDatabaseService) WithResourceResponseHandlers(database map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return database, map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
			"DBStatus": {
				TargetField: "db_status",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsPostgresqlDatabaseService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "CreateDatabase",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertAll,
			Convert: map[string]bp.RequestConvert{
				"db_name": {
					TargetField: "DBName",
				},
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				id := fmt.Sprintf("%s:%s", d.Get("instance_id"), d.Get("db_name"))
				d.SetId(id)
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlDatabaseService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRdsPostgresqlDatabaseService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "DeleteDatabase",
			ContentType: bp.ContentTypeJson,
			ConvertMode: bp.RequestConvertIgnore,
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				databaseId := d.Id()
				ids := strings.Split(databaseId, ":")
				if len(ids) != 2 {
					return false, fmt.Errorf("invalid rds postgresql database id")
				}
				(*call.SdkParam)["InstanceId"] = ids[0]
				(*call.SdkParam)["DBName"] = ids[1]
				return true, nil
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, call.SdkParam)
				return s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsPostgresqlDatabaseService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType: bp.ContentTypeJson,
		RequestConverts: map[string]bp.RequestConvert{
			"db_name": {
				TargetField: "DBName",
			},
		},
		NameField:    "DBName",
		IdField:      "DBName",
		CollectField: "databases",
		ResponseConverts: map[string]bp.ResponseConvert{
			"DBName": {
				TargetField: "db_name",
			},
			"DBStatus": {
				TargetField: "db_status",
			},
		},
	}
}

func (s *ByteplusRdsPostgresqlDatabaseService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "rds_postgresql",
		Version:     "2022-01-01",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}
