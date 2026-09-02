variable "aws_region" {
  type    = string
  default = "eu-west-1"
}
variable "project_name" {
  type    = string
  default = "portfolio-compliance-demo"
}
variable "container_image_tag" {
  type    = string
  default = "latest"
}
variable "desired_task_count" {
  type    = number
  default = 1
  validation {
    condition     = var.desired_task_count >= 0
    error_message = "desired_task_count must be non-negative."
  }
}
variable "environment_name" {
  type    = string
  default = "demo"
}

locals {
  name = "${var.project_name}-${var.environment_name}"
  tags = {
    Project = var.project_name, Environment = var.environment_name, ManagedBy = "Terraform"
  }
}
