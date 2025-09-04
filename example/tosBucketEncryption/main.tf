resource "byteplus_tos_bucket" "foo" {
  bucket_name   = "tf-acc-test-bucket"
  public_acl    = "private"
  az_redundancy = "multi-az"
  project_name  = "default"
  tags {
    key   = "k1"
    value = "v1"
  }
}

resource "byteplus_tos_bucket_encryption" "foo" {
  bucket_name = byteplus_tos_bucket.foo.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm       = "kms"
      kms_data_encryption = "AES256"
      kms_master_key_id   = "trn:kms:ap-southeast-1:30000*****:keyrings/acc-test-kms"
    }
  }
}
