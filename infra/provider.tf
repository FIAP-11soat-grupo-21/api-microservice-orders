provider "aws" {
  region = "us-east-2"
}

terraform {
  backend "s3" {
    bucket = "fiap-tc-terraform-846874"
    key    = "tech-challenge-project/order/terraform.tfstate"
    region = "us-east-2"
  }
}