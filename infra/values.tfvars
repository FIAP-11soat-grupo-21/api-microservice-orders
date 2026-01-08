application_name = "order-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8083
app_path_pattern = ["/v1/orders/*"]

# =======================================================
# Configurações do ECS Service
# =======================================================
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8083"
  API_HOST : "0.0.0.0"
  AWS_REGION : "us-east-2"

  DB_RUN_MIGRATIONS : "true"
  DB_NAME : "postgres"
  DB_PORT : "5432"
  DB_USERNAME : "adminuser"

  MESSAGE_BROKER_TYPE : "sqs"
}

container_secrets = {}

health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonRDSFullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess",
  "arn:aws:iam::aws:policy/AmazonCognitoPowerUser",
]
alb_is_internal = true

# =======================================================
# Configurações do API Gateway
# =======================================================
apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

authorization_name = "CognitoAuthorizer"