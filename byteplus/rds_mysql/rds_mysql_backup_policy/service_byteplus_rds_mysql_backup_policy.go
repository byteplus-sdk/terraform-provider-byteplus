package rds_mysql_backup_policy

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

type ByteplusRdsMysqlBackupPolicyService struct {
	Client     *bp.SdkClient
	Dispatcher *bp.Dispatcher
}

func NewRdsMysqlBackupPolicyService(c *bp.SdkClient) *ByteplusRdsMysqlBackupPolicyService {
	return &ByteplusRdsMysqlBackupPolicyService{
		Client:     c,
		Dispatcher: &bp.Dispatcher{},
	}
}

func (s *ByteplusRdsMysqlBackupPolicyService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRdsMysqlBackupPolicyService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithSimpleQuery(m, func(condition map[string]interface{}) ([]interface{}, error) {
		action := "DescribeBackupPolicy"

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
		results, err = bp.ObtainSdkValue("Result", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		results = []interface{}{results}
		if data, ok = results.([]interface{}); !ok {
			return data, errors.New("Result.BackupPolicy is not Slice")
		}
		return data, err
	})
}

func (s *ByteplusRdsMysqlBackupPolicyService) ReadResource(resourceData *schema.ResourceData, id string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if id == "" {
		id = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(id, ":")
	req := map[string]interface{}{
		"InstanceId": ids[0],
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
		return data, fmt.Errorf("rds_mysql_backup_policy %s not exist ", id)
	}
	return data, err
}

func (s *ByteplusRdsMysqlBackupPolicyService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
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
					return nil, "", fmt.Errorf("rds_mysql_backup_policy status error, status: %s", status.(string))
				}
			}
			return d, status.(string), err
		},
	}
}

func (s *ByteplusRdsMysqlBackupPolicyService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyBackupPolicy",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"data_full_backup_periods": {
					TargetField: "DataFullBackupPeriods",
					ConvertType: bp.ConvertJsonArray,
				},
				"lock_ddl_time": {
					TargetField: "LockDDLTime",
				},
				"data_full_backup_start_utc_hour": {
					TargetField: "DataFullBackupStartUTCHour",
				},
				"data_incr_backup_periods": {
					TargetField: "DataIncrBackupPeriods",
					ConvertType: bp.ConvertJsonArray,
				},
			},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, call.SdkParam)
				resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, resp, err)
				return resp, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := d.Get("instance_id").(string)
				d.SetId(instanceId + ":" + "backupPolicy")
				return nil
			},
		},
	}
	return []bp.Callback{callback}
}

func (ByteplusRdsMysqlBackupPolicyService) WithResourceResponseHandlers(d map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return d, map[string]bp.ResponseConvert{
			"LockDDLTime": {
				TargetField: "lock_ddl_time",
			},
			"DataFullBackupStartUTCHour": {
				TargetField: "data_full_backup_start_utc_hour",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRdsMysqlBackupPolicyService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      "ModifyBackupPolicy",
			ConvertMode: bp.RequestConvertAll,
			ContentType: bp.ContentTypeJson,
			Convert: map[string]bp.RequestConvert{
				"data_full_backup_periods": {
					TargetField: "DataFullBackupPeriods",
					ConvertType: bp.ConvertJsonArray,
				},
				"data_backup_retention_day": {
					TargetField: "DataBackupRetentionDay",
				},
				"data_full_backup_time": {
					TargetField: "DataFullBackupTime",
				},
				"data_incr_backup_periods": {
					TargetField: "DataIncrBackupPeriods",
					ConvertType: bp.ConvertJsonArray,
				},
				"binlog_file_counts_enable": {
					TargetField: "BinlogFileCountsEnable",
				},
				"binlog_limit_count": {
					TargetField: "BinlogLimitCount",
				},
				"binlog_local_retention_hour": {
					TargetField: "BinlogLocalRetentionHour",
				},
				"binlog_space_limit_enable": {
					TargetField: "BinlogSpaceLimitEnable",
				},
				"binlog_storage_percentage": {
					TargetField: "BinlogStoragePercentage",
				},
				"log_backup_retention_day": {
					TargetField: "LogBackupRetentionDay",
				},
				"log_ddl_time": {
					TargetField: "LogDDLTime",
				},
				"data_full_backup_start_utc_hour": {
					TargetField: "DataFullBackupStartUTCHour",
				},
				"hourly_incr_backup_enable": {
					TargetField: "HourlyIncrBackupEnable",
				},
				"incr_backup_hour_period": {
					TargetField: "IncrBackupHourPeriod",
				},
				"data_backup_encryption_enabled": {
					TargetField: "DataBackupEncryptionEnabled",
				},
				"binlog_backup_encryption_enabled": {
					TargetField: "BinlogBackupEncryptionEnabled",
				},
				"data_keep_policy_after_released": {
					TargetField: "DataKeepPolicyAfterReleased",
				},
				"data_keep_days_after_released": {
					TargetField: "DataKeepDaysAfterReleased",
				},
				"data_backup_all_retention": {
					TargetField: "DataBackupAllRetention",
				},
				"binlog_backup_all_retention": {
					TargetField: "BinlogBackupAllRetention",
				},
				"binlog_backup_enabled": {
					TargetField: "BinlogBackupEnabled",
				},
				"retention_policy_synced": {
					TargetField: "RetentionPolicySynced",
				},
			},
			BeforeCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (bool, error) {
				(*call.SdkParam)["InstanceId"] = d.Get("instance_id").(string)
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

func (s *ByteplusRdsMysqlBackupPolicyService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRdsMysqlBackupPolicyService) DatasourceResources(*schema.ResourceData, *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{}
}

func (s *ByteplusRdsMysqlBackupPolicyService) ReadResourceId(id string) string {
	return id
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
