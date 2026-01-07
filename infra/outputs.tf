output "db_address" {
  description = "Endereço do banco de dados RDS compartilhado"
  value       = data.terraform_remote_state.infra.outputs.rds_address
}

output "db_secret_arn" {
  description = "ARN do segredo do banco de dados RDS compartilhado"
  value       = data.terraform_remote_state.infra.outputs.rds_secret_arn
}

output "ecs_service_id" {
  description = "ID do serviço ECS do Orders"
  value       = module.order_api.service_id
}
