application_name = "order-api"
image_name       = "GHCR_IMAGE_TAG"
image_port       = 8080
container_environment_variables = {
  GO_ENV : "production"
  API_PORT : "8080"
  API_HOST : "0.0.0.0"
  AWS_REGION : "us-east-2"
  AWS_DYNAMO_TABLE_NAME : "order-api-table"
}
container_secrets = {}
health_check_path = "/health"
task_role_policy_arns = [
  "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess"
]