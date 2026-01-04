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
