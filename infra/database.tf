module "app_db" {
  source               = "git::https://github.com/FIAP-11soat-grupo-21/infra-core.git//modules/RDS?ref=main"
  project_common_tags  = { Project = "orders" }
  app_name             = "order-api-db"
  db_port              = 5432
  db_allocated_storage = 20
  db_storage_type      = "gp2"
  db_engine            = "postgres"
  db_engine_version    = "13"
  db_instance_class    = "db.t3.micro"
  db_username          = "appuser"

  private_subnets = data.terraform_remote_state.infra.outputs.private_subnet_ids
  vpc_id          = data.terraform_remote_state.infra.outputs.vpc_id
}
