module "ALB" {
  source             = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ALB?ref=main"
  loadbalancer_name  = var.application_name
  health_check_path  = var.health_check_path
  app_port           = var.image_port
  is_internal        = true
  private_subnet_ids = data.terraform_remote_state.infra.outputs.private_subnet_id
  vpc_id             = data.terraform_remote_state.infra.outputs.vpc_id
  vpc_cidr_blocks    = [data.terraform_remote_state.infra.outputs.vpc_cdir_block]

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "order_api" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ECS-Service?ref=main"

  cluster_id            = data.terraform_remote_state.infra.outputs.ecs_cluster_id
  ecs_security_group_id = data.terraform_remote_state.infra.outputs.ecs_security_group_id

  cloudwatch_log_group     = data.terraform_remote_state.infra.outputs.ecs_cloudwatch_log_group
  ecs_container_image      = var.image_name
  ecs_container_name       = var.application_name
  ecs_container_port       = var.image_port
  ecs_service_name         = var.application_name
  ecs_desired_count        = var.desired_count
  registry_credentials_arn = data.terraform_remote_state.infra.outputs.ecr_registry_credentials_arn

  ecs_container_environment_variables = var.container_environment_variables

  private_subnet_ids      = data.terraform_remote_state.infra.outputs.private_subnet_id
  task_execution_role_arn = data.terraform_remote_state.infra.outputs.ecs_task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = module.ALB.target_group_arn
  alb_security_group_id   = module.ALB.alb_security_group_id

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "api_gateway_routes" {
  source     = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.order_api]

  api_id            = data.terraform_remote_state.infra.outputs.api_gateway_id
  vpc_link_id       = data.terraform_remote_state.infra.outputs.api_gateway_vpc_link_id
  alb_listener_arn  = module.ALB.listener_arn
  gwapi_route_key   = "GET /{proxy+}"
  gwapi_auto_deploy = true
  stage_name        = data.terraform_remote_state.infra.outputs.api_gateway_stage_name

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
  api_gw_logs_arn     = data.terraform_remote_state.infra.outputs.api_gateway_logs_arn
}

module "dynamodb_table" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/Dynamo?ref=main"

  name          = "${var.application_name}-table"
  hash_key      = "id"
  hash_key_type = "S"
  billing_mode  = "PAY_PER_REQUEST"
}