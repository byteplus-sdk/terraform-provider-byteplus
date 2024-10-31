resource "byteplus_organization" "foo" {

}

data "byteplus_organization_units" "foo" {
  depends_on = [byteplus_organization.foo]
}

resource "byteplus_organization_unit" "foo" {
  name        = "tf-test-unit"
  parent_id   = [for unit in data.byteplus_organization_units.foo.units : unit.id if unit.parent_id == "0"][0]
  description = "tf-test"
}