package image

import (
	"errors"
	"fmt"
	"time"

	bp "github.com/byteplus-sdk/terraform-provider-byteplus/common"
	"github.com/byteplus-sdk/terraform-provider-byteplus/logger"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ByteplusImageService struct {
	Client *bp.SdkClient
}

func NewImageService(c *bp.SdkClient) *ByteplusImageService {
	return &ByteplusImageService{
		Client: c,
	}
}

func (s *ByteplusImageService) GetClient() *bp.SdkClient {
	return s.Client
}

func (s *ByteplusImageService) ReadResources(condition map[string]interface{}) ([]interface{}, error) {
	var (
		resp    *map[string]interface{}
		results interface{}
		ok      bool
	)
	return bp.WithNextTokenQuery(condition, "MaxResults", "NextToken", 20, nil, func(m map[string]interface{}) (data []interface{}, next string, err error) {
		action := "DescribeImages"
		logger.Debug(logger.ReqFormat, action, condition)
		if condition == nil {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), nil)
			if err != nil {
				return data, next, err
			}
		} else {
			resp, err = s.Client.UniversalClient.DoCall(getUniversalInfo(action), &condition)
			if err != nil {
				return data, next, err
			}
		}
		logger.Debug(logger.RespFormat, action, condition, *resp)

		results, err = bp.ObtainSdkValue("Result.Images", *resp)
		if err != nil {
			return data, next, err
		}
		nextToken, err := bp.ObtainSdkValue("Result.NextToken", *resp)
		if err != nil {
			return data, next, err
		}
		next = nextToken.(string)
		if results == nil {
			results = []interface{}{}
		}

		if _, ok = results.([]interface{}); !ok {
			return data, next, errors.New("Result.Images is not Slice")
		}

		return results.([]interface{}), next, err
	})
}

func (s *ByteplusImageService) ReadResource(resourceData *schema.ResourceData, imageId string) (data map[string]interface{}, err error) {
	var (
		results []interface{}
		ok      bool
	)
	if imageId == "" {
		imageId = s.ReadResourceId(resourceData.Id())
	}
	req := map[string]interface{}{
		"ImageIds.1": imageId,
	}

	results, err = s.ReadResources(req)
	if err != nil {
		return nil, err
	}
	for _, v := range results {
		if data, ok = v.(map[string]interface{}); !ok {
			return nil, errors.New("Value is not map ")
		}
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("Image %s not exist ", imageId)
	}
	return data, err
}

func (s *ByteplusImageService) RefreshResourceState(resourceData *schema.ResourceData, target []string, timeout time.Duration, id string) *resource.StateChangeConf {
	return nil
}

func (s *ByteplusImageService) WithResourceResponseHandlers(image map[string]interface{}) []bp.ResourceResponseHandler {
	handler := func() (map[string]interface{}, map[string]bp.ResponseConvert, error) {
		return image, nil, nil
	}
	return []bp.ResourceResponseHandler{handler}
}

func (s *ByteplusImageService) CreateResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusImageService) ModifyResource(resourceData *schema.ResourceData, resource *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusImageService) RemoveResource(resourceData *schema.ResourceData, r *schema.Resource) []bp.Callback {
	return []bp.Callback{}
}

func (s *ByteplusImageService) DatasourceResources(data *schema.ResourceData, resource *schema.Resource) bp.DataSourceInfo {
	return bp.DataSourceInfo{
		RequestConverts: map[string]bp.RequestConvert{
			"ids": {
				TargetField: "ImageIds",
				ConvertType: bp.ConvertWithN,
			},
		},
		NameField:    "ImageName",
		IdField:      "ImageId",
		CollectField: "images",
	}
}

func (s *ByteplusImageService) ReadResourceId(id string) string {
	return id
}

func getUniversalInfo(actionName string) bp.UniversalInfo {
	return bp.UniversalInfo{
		ServiceName: "ecs",
		Version:     "2020-04-01",
		HttpMethod:  bp.GET,
		Action:      actionName,
	}
}
