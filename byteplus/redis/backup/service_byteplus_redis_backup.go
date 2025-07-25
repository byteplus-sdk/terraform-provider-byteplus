package backup

import (
	"errors"
	"fmt"
	"github.com/byteplus-sdk/byteplus-go-sdk/byteplus/bytepluserr"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/redis/instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ActionDescribeBackups = "DescribeBackups"
	ActionCreateBackup    = "CreateBackup"
)

type ByteplusRedisBackupService struct {
	Client *bp.SdkClient
}

func NewRedisBackupService(c *bp.SdkClient) *ByteplusRedisBackupService {
	return &ByteplusRedisBackupService{
		Client: c,
	}
}

func (s *ByteplusRedisBackupService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisBackupService) readBackupPointDownloadUrls(instanceId, backupPointId string) (instance interface{}, err error) {
	action := "DescribeBackupPointDownloadUrls"
	cond := map[string]interface{}{
		"InstanceId":    instanceId,
		"BackupPointId": backupPointId,
	}
	logger.Debug(logger.RespFormat, action, cond)
	resp, err := s.Client.UniversalClient.DoCall(getUniversalInfo(action), &cond)
	logger.Debug(logger.RespFormat, action, *resp)
	if err != nil {
		return instance, err
	}

	instance, err = bp.ObtainSdkValue("Result.BackupPointDownloadUrls", *resp)
	if err != nil {
		return instance, err
	}
	if instance == nil {
		return instance, fmt.Errorf("instance %s or BackupPointId %s is not exist", instanceId, backupPointId)
	}
	return instance, err
}

func (s *ByteplusRedisBackupService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	var (
		resp    *map[string]interface{}
		results interface{}
	)

	return bp.WithPageNumberQuery(m, "PageSize", "PageNumber", 20, 1, func(condition map[string]interface{}) ([]interface{}, error) {
		logger.Debug(logger.ReqFormat, ActionDescribeBackups, m)
		if m == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(ActionDescribeBackups), nil)
			if err != nil {
				return data, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(ActionDescribeBackups), &m)
			if err != nil {
				return data, err
			}
		}
		if resp == nil {
			return data, fmt.Errorf("can not describe backup")
		}
		logger.Debug(logger.RespFormat, ActionDescribeBackups, m, *resp)
		results, err = bp.ObtainSdkValue("Result.Backups", *resp)
		if err != nil {
			return data, err
		}
		if results == nil {
			results = []interface{}{}
		}
		backupDetails, ok := results.([]interface{})
		if !ok {
			return data, fmt.Errorf("Result.Backups is not slice")
		}

		for _, ele := range backupDetails {
			backupDetail := ele.(map[string]interface{})
			instanceId, ok := backupDetail["InstanceId"].(string)
			if !ok {
				return data, fmt.Errorf("InstanceId is not string")
			}

			backupPointId, ok := backupDetail["BackupPointId"].(string)
			if !ok {
				return data, fmt.Errorf("BackupPointId is not string")
			}

			backupPointDownloadUrls, err := s.readBackupPointDownloadUrls(instanceId, backupPointId)
			if err != nil {
				return data, err
			}
			backupDetail["BackupPointDownloadUrls"] = backupPointDownloadUrls

			data = append(data, backupDetail)
		}
		return data, nil
	})
}

func (s *ByteplusRedisBackupService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)

	if tmpId == "" {
		tmpId = s.ReadResourceId(resourceData.Id())
	}
	ids := strings.Split(tmpId, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid id format")
	}
	req := map[string]interface{}{
		"InstanceId": ids[0],
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return data, err
	}
	if len(results) == 0 {
		return data, errors.New("backup not exist")
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); ok {
			backupPointId, ok := data["BackupPointId"].(string)
			if !ok {
				return data, fmt.Errorf("BackupPointId is not string")
			}
			if backupPointId == ids[1] {
				return data, nil
			}
		}
	}

	data["InstanceDetail"] = []interface{}{}
	return data, errors.New("backup not exist")
}

func (s *ByteplusRedisBackupService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			var (
				demo       map[string]interface{}
				status     interface{}
				failStates []string
			)
			failStates = append(failStates, "Error", "Unavailable", "Deleting")

			// 可能查询不到
			if err = resource.Retry(20*time.Minute, func() *resource.RetryError {
				demo, err = s.ReadResource(resourceData, id)
				if err != nil {
					if s.ResourceNotFoundError(err) {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				return nil
			}); err != nil {
				return nil, "", err
			}

			demo, err = s.ReadResource(resourceData, id)
			if err != nil {
				return nil, "", err
			}
			status, err = bp.ObtainSdkValue("Status", demo)
			if err != nil {
				return nil, "", err
			}
			for _, v := range failStates {
				if v == status.(string) {
					return nil, "", fmt.Errorf("Vpc  status  error, status:%s", status.(string))
				}
			}
			//注意 返回的第一个参数不能为空 否则会一直等下去
			return demo, status.(string), err
		},
	}
}

func (s *ByteplusRedisBackupService) WithResourceResponseHandlers(backup map[string]interface{}) []bp.ResourceResponseHandler {
	detail, ok := backup["InstanceDetail"].(map[string]interface{})
	if ok {
		vpcInfo, vpcInfoOk := detail["VpcInfo"].(map[string]interface{})
		if vpcInfoOk {
			vpcInfo["Id"] = vpcInfo["ID"] // id change
			detail["VpcInfo"] = []interface{}{vpcInfo}
			backup["InstanceDetail"] = detail
		}
	}

	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return backup, map[string]bp.ResponseConvert{
			"TTL": {
				TargetField: "ttl",
			},
		}, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusRedisBackupService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionCreateBackup,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam)
				output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				instanceId := (*call.SdkParam)["InstanceId"]
				id, err := bp.ObtainSdkValue("Result.BackupPointId", *resp)
				if err != nil {
					return err
				}
				d.SetId(fmt.Sprintf("%s:%s", instanceId, id))
				return nil
			},
			Refresh: &bp.StateRefresh{
				Target:  []string{"Available"},
				Timeout: resourceData.Timeout(schema.TimeoutCreate),
			},
			LockId: func(d *schema.ResourceData) string {
				return d.Get("instance_id").(string)
			},
			ExtraRefresh: map[bp.ResourceService]*bp.StateRefresh{
				instance.NewRedisDbInstanceService(s.Client): {
					Target:     []string{"Running"},
					Timeout:    resourceData.Timeout(schema.TimeoutCreate),
					ResourceId: resourceData.Get("instance_id").(string),
				},
			},
		},
	}
	return []bp.Callback{callback}
}

func (s *ByteplusRedisBackupService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRedisBackupService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRedisBackupService) DatasourceResources(data *schema.ResourceData, resource2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType:  bp.ContentTypeJson,
		CollectField: "backups",
		RequestConverts: map[string]bp.RequestConvert{
			"backup_strategy_list": {
				TargetField: "BackupStrategyList",
				ConvertType: bp.ConvertJsonArray,
			},
		},
		ResponseConverts: map[string]bp.ResponseConvert{
			"ID": {
				TargetField: "id",
			},
			"TTL": {
				TargetField: "ttl",
			},
		},
	}
}

func (s *ByteplusRedisBackupService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "Redis",
		Version:     "2020-12-07",
		HttpMethod:  bp.POST,
		ContentType: bp.ApplicationJSON,
		Action:      actionName,
	}
}

func (s *ByteplusRedisBackupService) ResourceNotFoundError(err error) bool {
	if e, ok := err.(bytepluserr.RequestFailure); ok && e.StatusCode() == 404 {
		return true
	}
	errMessage := strings.ToLower(err.Error())
	return strings.Contains(errMessage, "notfound") ||
		strings.Contains(errMessage, "not found") ||
		strings.Contains(errMessage, "not exist") ||
		strings.Contains(errMessage, "not associate") ||
		strings.Contains(errMessage, "invalid") ||
		strings.Contains(errMessage, "not_found") ||
		strings.Contains(errMessage, "notexist") ||
		strings.Contains(errMessage, "unavailable")
}
