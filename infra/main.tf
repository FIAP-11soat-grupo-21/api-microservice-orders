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
      AWS_COGNITO_USER_POOL_ID : data.terraform_remote_state.infra.outputs.cognito_user_pool_id
      AWS_COGNITO_USER_POOL_CLIENT_ID : data.terraform_remote_state.infra.outputs.cognito_user_pool_client_id
      USER_PASSWORD_AUTH : data.terraform_remote_state.infra.outputs.cognito_user_pool_client_secret

      # Database configuration
      DB_HOST : module.app_db.db_connection

      # SQS_PAYMENT_QUEUE_URL : data.terraform_remote_state.kitchen_order_api.outputs.sqs_queue_url,
      SQS_KITCHEN_QUEUE_URL : data.terraform_remote_state.kitchen_order_api.outputs.sqs_queue_url,
  })
  ecs_container_secrets = merge(var.container_secrets,
    {
      DB_PASSWORD : module.app_db.db_secret_password_arn
  })

  private_subnet_ids      = data.terraform_remote_state.infra.outputs.private_subnet_id
  task_execution_role_arn = data.terraform_remote_state.infra.outputs.ecs_task_execution_role_arn
  task_role_policy_arns   = var.task_role_policy_arns
  alb_target_group_arn    = data.terraform_remote_state.infra.outputs.alb_target_group_arn
  alb_security_group_id   = data.terraform_remote_state.infra.outputs.alb_security_group_id

  project_common_tags = data.terraform_remote_state.infra.outputs.project_common_tags
}

module "GetOrderAPIRoute" {
  source     = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/API-Gateway-Routes?ref=main"
  depends_on = [module.order_api]

  api_id       = data.terraform_remote_state.infra.outputs.api_gateway_id
  alb_proxy_id = aws_apigatewayv2_integration.alb_proxy.id

  endpoints = {
    get_order = {
      route_key  = "GET /orders/{id}"
      restricted = false
    },
    get_all_orders = {
      route_key  = "GET /orders"
      restricted = false
    },
    create_order = {
      route_key  = "POST /orders"
      restricted = false
    },
    update_order = {
      route_key  = "PUT /orders/{id}"
      restricted = false
    },
    delete_order = {
      route_key  = "DELETE /orders/{id}"
      restricted = false
    }
  }
}