# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

# EKS Outputs
output "cluster_id" {
  description = "The name/id of the EKS cluster"
  value       = module.eks.cluster_name
}

output "cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the cluster"
  value       = module.eks.cluster_arn
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_version" {
  description = "The Kubernetes version for the EKS cluster"
  value       = module.eks.cluster_version
}

output "cluster_platform_version" {
  description = "Platform version for the EKS cluster"
  value       = module.eks.cluster_platform_version
}

output "cluster_status" {
  description = "Status of the EKS cluster. One of `CREATING`, `ACTIVE`, `DELETING`, `FAILED`"
  value       = module.eks.cluster_status
}

output "cluster_security_group_id" {
  description = "Security group ID attached to the EKS cluster"
  value       = module.eks.cluster_security_group_id
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks.cluster_certificate_authority_data
}

output "cluster_oidc_issuer_url" {
  description = "The URL on the EKS cluster OIDC Issuer"
  value       = module.eks.cluster_oidc_issuer_url
}

output "cluster_primary_security_group_id" {
  description = "The cluster primary security group ID created by EKS"
  value       = module.eks.cluster_primary_security_group_id
}

# Node Group Outputs
output "node_groups" {
  description = "Map of attribute maps for all EKS managed node groups created"
  value       = module.eks.eks_managed_node_groups
}

output "node_security_group_id" {
  description = "ID of the node shared security group"
  value       = module.eks.node_security_group_id
}

# OIDC Provider Outputs
output "oidc_provider_arn" {
  description = "The ARN of the OIDC Provider if one is created"
  value       = module.eks.oidc_provider_arn
}

# IAM Role Outputs
output "cluster_autoscaler_role_arn" {
  description = "ARN of the cluster autoscaler IAM role"
  value       = aws_iam_role.cluster_autoscaler.arn
}

output "aws_load_balancer_controller_role_arn" {
  description = "ARN of the AWS Load Balancer Controller IAM role"
  value       = aws_iam_role.aws_load_balancer_controller.arn
}

# ALB Outputs
output "alb_dns_name" {
  description = "DNS name of the load balancer"
  value       = var.create_alb ? aws_lb.main[0].dns_name : null
}

output "alb_arn" {
  description = "ARN of the load balancer"
  value       = var.create_alb ? aws_lb.main[0].arn : null
}

output "alb_zone_id" {
  description = "Canonical hosted zone ID of the load balancer"
  value       = var.create_alb ? aws_lb.main[0].zone_id : null
}

# RDS Outputs
output "db_instance_id" {
  description = "The RDS instance ID"
  value       = var.create_rds ? aws_db_instance.main[0].id : null
}

output "db_instance_endpoint" {
  description = "The RDS instance endpoint"
  value       = var.create_rds ? aws_db_instance.main[0].endpoint : null
  sensitive   = true
}

output "db_instance_port" {
  description = "The RDS instance port"
  value       = var.create_rds ? aws_db_instance.main[0].port : null
}

output "db_instance_address" {
  description = "The RDS instance hostname"
  value       = var.create_rds ? aws_db_instance.main[0].address : null
  sensitive   = true
}

# ElastiCache Outputs
output "elasticache_replication_group_id" {
  description = "ID of the ElastiCache replication group"
  value       = var.create_elasticache ? aws_elasticache_replication_group.main[0].id : null
}

output "elasticache_primary_endpoint_address" {
  description = "Address of the replication group configuration endpoint"
  value       = var.create_elasticache ? aws_elasticache_replication_group.main[0].primary_endpoint_address : null
  sensitive   = true
}

output "elasticache_reader_endpoint_address" {
  description = "Address of the endpoint for the reader node in the replication group"
  value       = var.create_elasticache ? aws_elasticache_replication_group.main[0].reader_endpoint_address : null
  sensitive   = true
}

# S3 Outputs
output "s3_bucket_id" {
  description = "The name of the S3 bucket"
  value       = var.create_s3_bucket ? aws_s3_bucket.app_data[0].id : null
}

output "s3_bucket_arn" {
  description = "The ARN of the S3 bucket"
  value       = var.create_s3_bucket ? aws_s3_bucket.app_data[0].arn : null
}

output "s3_bucket_domain_name" {
  description = "The bucket domain name"
  value       = var.create_s3_bucket ? aws_s3_bucket.app_data[0].bucket_domain_name : null
}

# Additional Security Group Output
output "additional_security_group_id" {
  description = "ID of the additional security group"
  value       = aws_security_group.additional.id
}

# Kubeconfig
output "configure_kubectl" {
  description = "Configure kubectl: make sure you're logged in with the correct AWS profile and run the following command to update your kubeconfig"
  value       = "aws eks --region ${var.aws_region} update-kubeconfig --name ${module.eks.cluster_name}"
}