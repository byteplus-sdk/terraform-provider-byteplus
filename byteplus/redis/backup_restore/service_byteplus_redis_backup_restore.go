package backup_restore

import (
	"fmt"
	"strings"
	"time"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/redis/instance"
	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusRedisBackupRestoreService struct {
	Client *bp.SdkClient
}

const (
	ActionRestoreDBInstance = "RestoreDBInstance"
)

func NewRedisBackupRestoreService(c *bp.SdkClient) *ByteplusRedisBackupRestoreService {
	return &ByteplusRedisBackupRestoreService{
		Client: c,
	}
}

func (s *ByteplusRedisBackupRestoreService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusRedisBackupRestoreService) ReadResources(m map[string]interface{}) (data []interface{}, err error) {
	return nil, nil
}

func (s *ByteplusRedisBackupRestoreService) ReadResource(resourceData *schema.ResourceData, tmpId string) (data map[string]interface{}, err error) {
	var (
		ids []string
	)
	if tmpId == "" {
		tmpId = s.ReadResourceId(resourceData.Id())
	}
	// 无法读取出来
	ids = strings.Split(tmpId, ":")
	if len(ids) != 2 {
		return data, fmt.Errorf("invalid resource id")
	}
	return data, nil
}

func (s *ByteplusRedisBackupRestoreService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Delay:      time.Second,
		Pending:    []string{},
		Target:     target,
		Timeout:    timeout,
		MinTimeout: time.Second,
		Refresh:    nil,
	}
}

func (s *ByteplusRedisBackupRestoreService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	return []bp.ResourceResponseHandler{}
}

func (s *ByteplusRedisBackupRestoreService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	callback := bp.Callback{
		Call: bp.SdkCall{
			Action:      ActionRestoreDBInstance,
			ContentType: bp.ContentTypeJson,
			Convert:     map[string]bp.RequestConvert{},
			ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
				logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
				output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
				logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
				return output, err
			},
			AfterCall: func(d *schema.ResourceData, client *bp.SdkClient, resp *map[string]interface{}, call bp.SdkCall) error {
				d.SetId(fmt.Sprintf("restore:%s", d.Get("instance_id")))
				return nil
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

func (s *ByteplusRedisBackupRestoreService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) (callbacks []bp.Callback) {
	if resourceData.HasChanges("time_point", "backup_point_id") {
		callback := bp.Callback{
			Call: bp.SdkCall{
				Action:      ActionRestoreDBInstance,
				ContentType: bp.ContentTypeJson,
				ConvertMode: bp.RequestConvertInConvert,
				Convert: map[string]bp.RequestConvert{
					"instance_id": {
						ForceGet: true,
					},
					"backup_type": {
						ForceGet: true,
					},
					"time_point": {
						TargetField: "TimePoint",
					},
					"backup_point_id": {
						TargetField: "BackupPointId",
					},
				},
				ExecuteCall: func(d *schema.ResourceData, client *bp.SdkClient, call bp.SdkCall) (*map[string]interface{}, error) {
					logger.Debug(logger.ReqFormat, call.Action, *call.SdkParam)
					output, err := s.Client.UniversalClient.DoCall(getUniversalInfo(call.Action), call.SdkParam)
					logger.Debug(logger.RespFormat, call.Action, *call.SdkParam, *output)
					return output, err
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
		callbacks = append(callbacks, callback)
	}
	return callbacks
}

func (s *ByteplusRedisBackupRestoreService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusRedisBackupRestoreService) DatasourceResources(data *schema.ResourceData, resource2 *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		ContentType: bp.ContentTypeJson,
	}
}

func (s *ByteplusRedisBackupRestoreService) ReadResourceId(id string) string {
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
