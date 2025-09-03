---
subcategory: "TOS"
layout: "byteplus"
page_title: "Byteplus: byteplus_tos_bucket_policy"
sidebar_current: "docs-byteplus-resource-tos_bucket_policy"
description: |-
  Provides a resource to manage tos bucket policy
---
# byteplus_tos_bucket_policy
Provides a resource to manage tos bucket policy
## Example Usage
```hcl
resource "byteplus_tos_bucket_policy" "default" {
  bucket_name = "tf-acc-test-bucket"
  policy = jsonencode({
    Statement = [
      {
        Sid    = "test"
        Effect = "Allow"
        Principal = [
          "AccountId/subUserName"
        ]
        Action = [
          "tos:List*"
        ]
        Resource = [
          "trn:tos:::tf-acc-test-bucket"
        ]
      }
    ]
  })
}
```
## Argument Reference
The following arguments are supported:
* `bucket_name` - (Required, ForceNew) The name of the bucket.
* `policy` - (Required) The policy document. This is a JSON formatted string. For more information about building Byteplus IAM policy documents with Terraform, see the  [Byteplus IAM Policy Document Guide](https://www.byteplus.com/docs/6349/102127).

## Attributes Reference
In addition to all arguments above, the following attributes are exported:
* `id` - ID of the resource.



## Import
Tos Bucket can be imported using the id, e.g.
```
$ terraform import byteplus_tos_bucket_policy.default bucketName:policy
```

