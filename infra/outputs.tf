output "dynamodb_table_name" {
  description = "Nome da tabela DynamoDB do Orders"
  value       = module.dynamodb_table.table_name
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Orders"
  value       = module.order_api.service_id
}
