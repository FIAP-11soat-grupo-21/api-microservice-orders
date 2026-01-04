module "dynamodb_table" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Dynamo?ref=main"

  name          = "${var.application_name}-table"
  hash_key      = var.dynamodb_hash_key
  hash_key_type = var.dynamodb_hash_key_type
  billing_mode  = var.dynamodb_billing_mode

  secondary_indexes = var.dynamodb_secondary_indexes

  range_key = var.dynamodb_range_keys
}
