resource "byteplus_ecs_key_pair" "foo" {
  key_pair_name = "acc-test-key-name"
  description   = "acc-test"
}
data "byteplus_ecs_key_pairs" "foo" {
  key_pair_name = byteplus_ecs_key_pair.foo.key_pair_name
}