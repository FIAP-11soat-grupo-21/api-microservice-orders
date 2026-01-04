module "cognito" {
  source = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/cognito?ref=main"

  user_pool_name               = var.cognito_user_pool_name
  allow_admin_create_user_only = var.allow_admin_create_user_only
  auto_verified_attributes     = var.auto_verified_attributes
  username_attributes          = var.username_attributes
  email_required               = var.email_required
  name_required                = var.name_required
  generate_secret              = var.generate_secret
  access_token_validity        = var.access_token_validity
  id_token_validity            = var.id_token_validity
  refresh_token_validity       = var.refresh_token_validity

  tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "ALB" {
  source             = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/ALB?ref=main"
  loadbalancer_name  = var.application_name
  health_check_path  = var.health_check_path
  app_port           = var.image_port
  is_internal        = var.alb_is_internal
  private_subnet_ids = data.terraform_remote_state.infra.outputs.private_subnet_id
  vpc_id             = data.terraform_remote_state.infra.outputs.vpc_id

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

  ecs_container_environment_variables = merge(var.container_environment_variables,
    {
      AWS_COGNITO_USER_POOL_ID : module.cognito.user_pool_id,
      AWS_COGNITO_USER_POOL_CLIENT_ID : module.cognito.user_pool_client_id,

      # SQS_PAYMENT_QUEUE_URL : data.terraform_remote_state.kitchen_order_api.outputs.sqs_queue_url,
      SQS_KITCHEN_QUEUE_URL : data.terraform_remote_state.kitchen_order_api.outputs.sqs_queue_url,
  })

  private_subnet_ids      = data.terraform_remote_state.infra.outputs.private_subnet_id
  task_execution_role_arn = data.terraform_remote_state.infra.outputs.ecs_task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = module.ALB.target_group_arn
  alb_security_group_id   = module.ALB.alb_security_group_id

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "GetOrderAPIRoute" {
  source     = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.order_api]

  api_id       = data.terraform_remote_state.infra.outputs.api_gateway_id
  alb_proxy_id = aws_apigatewayv2_integration.alb_proxy.id

  endpoints = var.api_endpoints
}