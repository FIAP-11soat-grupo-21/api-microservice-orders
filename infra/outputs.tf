output "db_address" {
  description = "Endereço do banco de dados RDS compartilhado"
  value       = data.terraform_remote_state.infra.outputs.rds_address
}

output "db_secret_arn" {
  description = "ARN do segredo do banco de dados RDS compartilhado"
  value       = data.terraform_remote_state.infra.outputs.rds_secret_arn
}

output "sqs_orders_queue_url" {
  description = "URL da fila SQS do Orders (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_orders_queue_url
}

output "sqs_orders_order_error_queue_url" {
  description = "URL da fila SQS de erro do Orders (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_orders_order_error_queue_url
}

output "sqs_payments_queue_url" {
  description = "URL da fila SQS do Payments (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_payments_queue_url
}

output "sqs_kitchen_orders_queue_url" {
  description = "URL da fila SQS do Kitchen Orders (do infra-core)"
  value       = data.terraform_remote_state.infra.outputs.sqs_kitchen_orders_queue_url
}



output "ecs_service_id" {
  description = "ID do serviço ECS do Orders"
  value       = module.order_api.service_id
}
