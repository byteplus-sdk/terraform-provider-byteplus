package iam_policy_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamPolicyCreateConfig = `
resource "byteplus_iam_policy" "foo" {
    policy_name = "acc-test-policy"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}
`

func TestAccByteplusIamPolicyResource_Basic(t *testing.T) {
	resourceName := "byteplus_iam_policy.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.ByteplusIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccByteplusIamPolicyUpdateConfig = `
resource "byteplus_iam_policy" "foo" {
    policy_name = "acc-test-policy-new"
	description = "acc-test-new"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"
}
`

func TestAccByteplusIamPolicyResource_Update(t *testing.T) {
	resourceName := "byteplus_iam_policy.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.ByteplusIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers:    byteplus.GetTestAccProviders(),
		CheckDestroy: byteplus.AccTestCheckResourceRemove(acc),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamPolicyCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
				),
			},
			{
				Config: testAccByteplusIamPolicyUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					byteplus.AccTestCheckResourceExists(acc),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_name", "acc-test-policy-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "description", "acc-test-new"),
					resource.TestCheckResourceAttr(acc.ResourceId, "policy_document", "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_trn"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "policy_type"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "create_date"),
					resource.TestCheckResourceAttrSet(acc.ResourceId, "update_date"),
				),
			},
			{
				Config:             testAccByteplusIamPolicyUpdateConfig,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false, // 修改之后，不应该再产生diff
			},
		},
	})
}
