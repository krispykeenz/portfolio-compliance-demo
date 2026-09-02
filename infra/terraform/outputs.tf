output "frontend_url" {
  value = "https://${aws_cloudfront_distribution.main.domain_name}"
}
output "ecr_repository_url" {
  value = aws_ecr_repository.backend.repository_url
}
output "frontend_bucket" {
  value = aws_s3_bucket.frontend.id
}
output "ecs_cluster_name" {
  value = aws_ecs_cluster.main.name
}
output "ecs_service_name" {
  value = aws_ecs_service.backend.name
}
output "cloudfront_distribution_id" {
  value = aws_cloudfront_distribution.main.id
}
