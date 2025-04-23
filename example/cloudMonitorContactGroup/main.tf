resource "byteplus_cloud_monitor_contact" "foo1" {
  name  = "acc-test-contact-1"
  email = "test1@163.com"
}

resource "byteplus_cloud_monitor_contact" "foo2" {
  name  = "acc-test-contact-2"
  email = "test2@163.com"
}

resource "byteplus_cloud_monitor_contact_group" "foo" {
  name             = "acc-test-contact-group-new"
  description      = "tf-test-new"
  contacts_id_list = [byteplus_cloud_monitor_contact.foo1.id, byteplus_cloud_monitor_contact.foo2.id]
}
