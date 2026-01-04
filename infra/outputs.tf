output "dynamodb_table_name" {
  description = "Nome da tabela DynamoDB do Orders"
  value       = module.dynamodb_table.table_name
}

output "alb_dns_name" {
  description = "DNS name do Application Load Balancer do Orders"
  value       = module.ALB.alb_dns_name
}

output "cognito_user_pool_id" {
  description = "ID do Cognito User Pool do Orders"
  value       = module.cognito.user_pool_id
}

output "cognito_user_pool_client_id" {
  description = "ID do Cognito User Pool Client do Orders"
  value       = module.cognito.user_pool_client_id
}

# Outputs relacionados ao SQS
output "sqs_kitchen_queue_url" {
  description = "URL da fila SQS para enviar pedidos para a cozinha"
  value       = module.sqs_kitchen_orders.sqs_queue_url
}

output "sqs_payment_queue_url" {
  description = "URL da fila SQS para receber confirmações de pagamento"
  value       = data.terraform_remote_state.kitchen_order.outputs.sqs_queue_url
}
