package zone

import (
	"errors"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusClbZoneService struct {
	Client *bp.SdkClient
}

func (s *ByteplusClbZoneService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
		err     error
		data    []interface{}
	)
	action := "DescribeZones"
	if condition == nil {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
	} else {
		resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
	}
	if err != nil {
		return nil, err
	}
	logger.Debug(logger.RespFormat, action, condition, *resp)
	results, err = bp.ObtainSdkValue("Result.MasterZones", *resp)
	if err != nil {
		return nil, err
	}
	if results == nil {
		results = make([]interface{}, 0)
	}
	if data, ok = results.([]interface{}); !ok {
		return nil, errors.New("Result.MasterZones is not Slice")
	}
	return data, nil
}

func (s *ByteplusClbZoneService) ReadResource(data *schema.ResourceData, s2 string) (map[string]interface{}, error) {
	return nil, nil
}

func (s *ByteplusClbZoneService) RefreshResourceState(data *schema.ResourceData, target []string, timeout time.Duration, s2 string) *resource.StateChangeConf {
	return &resource.StateChangeConf{
		Pending:    []string{},
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
		Target:     target,
		Timeout:    timeout,
		Refresh: func() (result interface{}, state string, err error) {
			return nil, "", err
		},
	}
}

func (s *ByteplusClbZoneService) WithResourceResponseHandlers(m map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return m, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusClbZoneService) CreateResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusClbZoneService) ModifyResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusClbZoneService) RemoveResource(data *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusClbZoneService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		CollectField: "master_zones",
	}
}

func (s *ByteplusClbZoneService) ReadResourceId(id string) string {
	return id
}

func NewClbZoneService(c *bp.SdkClient) *ByteplusClbZoneService {
	return &ByteplusClbZoneService{
		Client: c,
	}
}

func (s *ByteplusClbZoneService) GetClient() *bp.SdkClient {
	return s.Client
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "clb",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		ContentType: bp.Default,
		Action:      actionName,
	}
}
