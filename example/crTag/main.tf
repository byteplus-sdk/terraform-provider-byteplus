# Tag cannot be created,please import by command `terraform import byteplus_cr_tag.default registry:namespace:repository:tag`
resource "byteplus_cr_tag" "default" {
  registry   = "enterprise-1"
  namespace  = "langyu"
  repository = "repo"
  name       = "v2"
}