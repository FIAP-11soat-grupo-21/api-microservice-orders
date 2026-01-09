data "terraform_remote_state" "infra" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/core/terraform.tfstate"
    region = "us-east-2"
  }
}

data "terraform_remote_state" "kitchen_order_api" {
  backend = "s3"
  config = {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/kitchen-order/terraform.tfstate"
    region = "us-east-2"
  }
}