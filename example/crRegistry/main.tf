# create cr registry
resource "byteplus_cr_registry" "foo" {
  name               = "acc-test-cr"
  delete_immediately = false
  password           = "1qaz!QAZ"
  project            = "default"
}

# create cr namespace
resource "byteplus_cr_namespace" "foo" {
  registry = byteplus_cr_registry.foo.id
  name     = "acc-test-namespace"
  project  = "default"
}

# create cr repository
resource "byteplus_cr_repository" "foo" {
  registry     = byteplus_cr_registry.foo.id
  namespace    = byteplus_cr_namespace.foo.name
  name         = "acc-test-repository"
  description  = "A test repository created by terraform."
  access_level = "Public"
}
