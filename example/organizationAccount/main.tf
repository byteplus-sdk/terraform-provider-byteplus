resource "byteplus_organization_unit" "foo" {
  name        = "acc-test-org-unit"
  parent_id   = "730671013833632****"
  description = "acc-test"
}

resource "byteplus_organization_account" "foo" {
  account_name             = "acc-test-account"
  show_name                = "acc-test-account"
  description              = "acc-test"
  org_unit_id              = byteplus_organization_unit.foo.id
  verification_relation_id = "210026****"

  tags {
    key   = "k1"
    value = "v1"
  }
}
