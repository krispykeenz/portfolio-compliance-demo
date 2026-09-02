terraform {
  required_version = ">= 1.13.0, < 2.0.0"
  required_providers {
    aws = {
      source = "hashicorp/aws", version = "~> 6.13"
    }
    random = {
      source = "hashicorp/random", version = "~> 3.7"
    }

  }
}

provider "aws" {
  region = var.aws_region
}
