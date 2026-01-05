output "db_address" {
  description = "Endereço do banco de dados RDS do Orders"
  value       = module.app_db.db_connection
}

output "db_secret_arn" {
  description = "ARN do segredo do banco de dados RDS do Orders"
  value       = module.app_db.db_secret_password_arn
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Orders"
  value       = module.order_api.service_id
}
