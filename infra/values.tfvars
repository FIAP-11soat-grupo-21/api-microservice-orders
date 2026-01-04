application_name = "order-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8080

# =======================================================
# Configurações do ECS Service
# =======================================================
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8080"
  API_HOST : "0.0.0.0"
  AWS_REGION : "us-east-2"
  AWS_DYNAMO_TABLE_NAME : "order-api-table"
  
  MESSAGE_BROKER_TYPE : "sqs"
}

container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
  "arn:aws:iam::aws:policy/AmazonSQSFullAccess"
]
alb_is_internal = true

# =======================================================
# Configurações do API Gateaway
# =======================================================
# API Gateway
apigw_integration_type       = "HTTP_PROXY"
apigw_integration_method     = "ANY"
apigw_payload_format_version = "1.0"
apigw_connection_type        = "VPC_LINK"

# Definição dos endpoints da API
api_endpoints = {
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

