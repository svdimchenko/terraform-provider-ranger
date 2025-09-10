resource "ranger_policy" "example" {
  service    = "service-name"
  name       = "test-policy-name"
  definition = file("${path.module}/policies/test-rls-tf.json")
}
