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