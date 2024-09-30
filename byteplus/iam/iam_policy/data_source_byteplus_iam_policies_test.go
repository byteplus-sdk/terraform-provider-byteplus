package iam_policy_test

import (
	"testing"

	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus"
	"github.com/byteplus-sdk/terraform-provider-byteplus/byteplus/iam/iam_policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const testAccByteplusIamPoliciesDatasourceConfig = `
resource "byteplus_iam_policy" "foo1" {
    policy_name = "acc-test-policy1"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingGroups\"],\"Resource\":[\"*\"]}]}"
}

resource "byteplus_iam_policy" "foo2" {
    policy_name = "acc-test-policy2"
	description = "acc-test"
	policy_document = "{\"Statement\":[{\"Effect\":\"Allow\",\"Action\":[\"auto_scaling:DescribeScalingConfigurations\"],\"Resource\":[\"*\"]}]}"
}

data "byteplus_iam_policies" "foo"{
    query = "${byteplus_iam_policy.foo1.description}"
}
`

func TestAccByteplusIamPoliciesDatasource_Basic(t *testing.T) {
	resourceName := "data.byteplus_iam_policies.foo"

	acc := &byteplus.AccTestResource{
		ResourceId: resourceName,
		Svc:        &iam_policy.ByteplusIamPolicyService{},
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			byteplus.AccTestPreCheck(t)
		},
		Providers: byteplus.GetTestAccProviders(),
		Steps: []resource.TestStep{
			{
				Config: testAccByteplusIamPoliciesDatasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(acc.ResourceId, "policies.#", "2"),
				),
			},
		},
	})
}
