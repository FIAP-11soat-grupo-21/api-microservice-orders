application_name = "order-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8083
app_path_pattern = ["/orders*", "/orders/*"]

# =======================================================
# Configurações do ECS Service
# =======================================================
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8083"
  API_HOST : "0.0.0.0"
  AWS_REGION : "us-east-2"
  AWS_DYNAMO_TABLE_NAME : "order-api-table"

  MESSAGE_BROKER_TYPE : "sqs"
}

container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
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

# =======================================================
# Configurações do dynamoDB
# =======================================================
dynamodb_secondary_indexes = [
  {
    name            = "cpf-index"
    hash_key        = "cpf"
    range_key       = "S"
    projection_type = "ALL"
  },
  {
    name            = "email-index"
    hash_key        = "email"
    range_key       = "S"
    projection_type = "ALL"
  }
]

dynamodb_hash_key      = "id"
dynamodb_hash_key_type = "S"
dynamodb_billing_mode  = "PAY_PER_REQUEST"

dynamodb_range_keys = [
  {
    name = "cpf",
    type = "S"
  },
  {
    name = "email",
    type = "S"
  },
]

