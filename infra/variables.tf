variable "application_name" {
  description = "Nome da aplicação ECS"
  type        = string
}

variable "image_name" {
  description = "Nome da imagem do container"
  type        = string
}

variable "image_port" {
  description = "Porta do container"
  type        = number
}

variable "desired_count" {
  description = "Número desejado de tarefas ECS"
  type        = number
  default     = 1
}

variable "container_environment_variables" {
  description = "Variáveis de ambiente do container"
  type        = map(string)
  default     = {}
}

variable "container_secrets" {
  description = "Segredos do container"
  type        = map(string)
  default     = {}
}

variable "health_check_path" {
  description = "Caminho de verificação de integridade do serviço"
  type        = string
  default     = "/health"
}

variable "task_role_policy_arns" {
  description = "Lista de ARNs de políticas para anexar à função da tarefa ECS"
  type        = list(string)
  default     = []
}

variable "api_endpoints" {
  description = "Lista de endpoints da API Gateway"
  type = map(object({
    route_key           = string
    target              = optional(string)
    restricted          = optional(bool, false)
    auth_integration_id = optional(string)
  }))
}

variable "alb_is_internal" {
  description = "Se o ALB é interno"
  type        = bool
  default     = true
}


#########################################################
################# Variáveis do DynamoDB #################
#########################################################

variable "dynamodb_secondary_indexes" {
  description = "Lista de índices secundários para a tabela DynamoDB"
  type = list(object({
    name            = string
    hash_key        = string
    range_key       = string
    projection_type = string
  }))
}

variable "dynamodb_hash_key" {
  description = "Hash key da tabela DynamoDB"
  type        = string
  default     = "id"
}

variable "dynamodb_hash_key_type" {
  description = "Tipo da hash key da tabela DynamoDB"
  type        = string
  default     = "S"
}

variable "dynamodb_billing_mode" {
  description = "Billing mode da tabela DynamoDB"
  type        = string
  default     = "PAY_PER_REQUEST"
}

variable "dynamodb_range_keys" {
  description = "Lista de range keys (nome, tipo) para a tabela DynamoDB"
  type = list(object({
    name = string
    type = string
  }))
  default = []
}

#########################################################
############### Variáveis do API Gateway ################
#########################################################

variable "apigw_integration_type" {
  description = "Tipo de integração do API Gateway"
  type        = string
  default     = "HTTP_PROXY"
}

variable "apigw_integration_method" {
  description = "Método de integração do API Gateway"
  type        = string
  default     = "ANY"
}

variable "apigw_payload_format_version" {
  description = "Versão do payload do API Gateway"
  type        = string
  default     = "1.0"
}

variable "apigw_connection_type" {
  description = "Tipo de conexão do API Gateway"
  type        = string
  default     = "VPC_LINK"
}


##########################################################
############ Variáveis do Cognito User Pool ##############
##########################################################

variable "cognito_user_pool_name" {
  description = "Nome do Cognito User Pool"
  type        = string
  default     = "tech-challenge-user-pool"
}

variable "allow_admin_create_user_only" {
  description = "Permitir apenas criação de usuário por admin"
  type        = bool
  default     = false
}

variable "auto_verified_attributes" {
  description = "Atributos auto verificados"
  type        = list(string)
  default     = ["email"]
}

variable "username_attributes" {
  description = "Atributos usados como username"
  type        = list(string)
  default     = []
}

variable "email_required" {
  description = "E-mail é requerido"
  type        = bool
  default     = true
}

variable "name_required" {
  description = "Nome é requerido"
  type        = bool
  default     = true
}

variable "generate_secret" {
  description = "Gerar secret para o client"
  type        = bool
  default     = true
}

variable "access_token_validity" {
  description = "Validade do access token em minutos"
  type        = number
  default     = 60
}

variable "id_token_validity" {
  description = "Validade do id token em minutos"
  type        = number
  default     = 60
}

variable "refresh_token_validity" {
  description = "Validade do refresh token em dias"
  type        = number
  default     = 30
}