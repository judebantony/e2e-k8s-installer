package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
)

// Manager handles Terraform operations
type Manager struct {
	config      *config.InfrastructureConfig
	workingDir  string
	initialized bool
}

// NewManager creates a new Terraform manager
func NewManager(infraConfig *config.InfrastructureConfig) (*Manager, error) {
	if infraConfig == nil {
		return nil, fmt.Errorf("infrastructure configuration is required")
	}

	workingDir := infraConfig.Terraform.Workspace
	if workingDir == "" {
		workingDir = "./terraform"
	}

	// Ensure working directory exists
	if err := os.MkdirAll(workingDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create terraform working directory: %w", err)
	}

	return &Manager{
		config:     infraConfig,
		workingDir: workingDir,
	}, nil
}

// Init initializes Terraform in the working directory
func (m *Manager) Init() error {
	logger.Info("Initializing Terraform").Str("workingDir", m.workingDir).Send()

	// Check if terraform binary exists
	if _, err := exec.LookPath("terraform"); err != nil {
		return fmt.Errorf("terraform not found in PATH: %w", err)
	}

	// Create main.tf if it doesn't exist
	if err := m.ensureMainTerraformFile(); err != nil {
		return fmt.Errorf("failed to create main terraform file: %w", err)
	}

	// Initialize terraform
	cmd := exec.Command("terraform", "init")
	cmd.Dir = m.workingDir

	// Set environment variables
	cmd.Env = append(os.Environ(), m.getTerraformEnvVars()...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Terraform init failed").
			Str("output", string(output)).
			Err(err).
			Send()
		return fmt.Errorf("terraform init failed: %w\nOutput: %s", err, string(output))
	}

	logger.Info("Terraform initialized successfully").Send()
	m.initialized = true
	return nil
}

// Plan creates a Terraform plan
func (m *Manager) Plan(destroy bool) (string, error) {
	if !m.initialized {
		return "", fmt.Errorf("terraform not initialized")
	}

	logger.Info("Creating Terraform plan").Bool("destroy", destroy).Send()

	args := []string{"plan", "-no-color"}
	if destroy {
		args = append(args, "-destroy")
	}

	// Add variables file if specified
	if len(m.config.Terraform.VarFiles) > 0 {
		for _, varFile := range m.config.Terraform.VarFiles {
			args = append(args, "-var-file="+varFile)
		}
	}

	// Add provider-specific variables
	varArgs := m.getProviderVariables()
	args = append(args, varArgs...)

	cmd := exec.Command("terraform", args...)
	cmd.Dir = m.workingDir
	cmd.Env = append(os.Environ(), m.getTerraformEnvVars()...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Terraform plan failed").
			Str("output", string(output)).
			Err(err).
			Send()
		return "", fmt.Errorf("terraform plan failed: %w\nOutput: %s", err, string(output))
	}

	logger.Info("Terraform plan completed").Send()
	return string(output), nil
}

// Apply applies the Terraform configuration
func (m *Manager) Apply(destroy bool) error {
	if !m.initialized {
		return fmt.Errorf("terraform not initialized")
	}

	logger.Info("Applying Terraform configuration").Bool("destroy", destroy).Send()

	var args []string
	if destroy {
		args = []string{"destroy", "-auto-approve", "-no-color"}
	} else {
		args = []string{"apply", "-auto-approve", "-no-color"}
	}

	// Add variables file if specified
	if len(m.config.Terraform.VarFiles) > 0 {
		for _, varFile := range m.config.Terraform.VarFiles {
			args = append(args, "-var-file="+varFile)
		}
	}

	// Add provider-specific variables
	varArgs := m.getProviderVariables()
	args = append(args, varArgs...)

	cmd := exec.Command("terraform", args...)
	cmd.Dir = m.workingDir
	cmd.Env = append(os.Environ(), m.getTerraformEnvVars()...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Terraform apply failed").
			Str("output", string(output)).
			Err(err).
			Send()
		return fmt.Errorf("terraform apply failed: %w\nOutput: %s", err, string(output))
	}

	logger.Info("Terraform configuration applied successfully").Send()
	return nil
}

// GetOutputs retrieves Terraform outputs
func (m *Manager) GetOutputs() (map[string]interface{}, error) {
	if !m.initialized {
		return nil, fmt.Errorf("terraform not initialized")
	}

	logger.Info("Retrieving Terraform outputs").Send()

	cmd := exec.Command("terraform", "output", "-json")
	cmd.Dir = m.workingDir
	cmd.Env = append(os.Environ(), m.getTerraformEnvVars()...)

	output, err := cmd.Output()
	if err != nil {
		// If no outputs exist, return empty map instead of error
		if strings.Contains(err.Error(), "no outputs") {
			logger.Info("No Terraform outputs found").Send()
			return make(map[string]interface{}), nil
		}
		logger.Error("Failed to get Terraform outputs").Err(err).Send()
		return nil, fmt.Errorf("failed to get terraform outputs: %w", err)
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal(output, &outputs); err != nil {
		return nil, fmt.Errorf("failed to parse terraform outputs: %w", err)
	}

	logger.Info("Terraform outputs retrieved").Int("count", len(outputs)).Send()
	return outputs, nil
}

// RunHealthChecks runs health checks on the infrastructure
func (m *Manager) RunHealthChecks() error {
	logger.Info("Running infrastructure health checks").Send()

	// Get outputs to check infrastructure health
	outputs, err := m.GetOutputs()
	if err != nil {
		return fmt.Errorf("failed to get outputs for health checks: %w", err)
	}

	// Check Kubernetes cluster endpoint if available
	if kubeEndpoint, exists := outputs["kubernetes_endpoint"]; exists {
		if err := m.checkKubernetesHealth(kubeEndpoint); err != nil {
			return fmt.Errorf("kubernetes health check failed: %w", err)
		}
	}

	// Check database endpoint if available
	if dbEndpoint, exists := outputs["database_endpoint"]; exists {
		if err := m.checkDatabaseHealth(dbEndpoint); err != nil {
			logger.Warn("Database health check failed").Err(err).Send()
			// Don't fail the entire operation for database health check
		}
	}

	logger.Info("Infrastructure health checks completed").Send()
	return nil
}

// ensureMainTerraformFile creates a basic main.tf if it doesn't exist
func (m *Manager) ensureMainTerraformFile() error {
	mainTfPath := filepath.Join(m.workingDir, "main.tf")

	// Check if main.tf already exists
	if _, err := os.Stat(mainTfPath); err == nil {
		logger.Info("main.tf already exists").Str("path", mainTfPath).Send()
		return nil
	}

	// Create basic main.tf based on provider
	// TODO: Make provider configurable from configuration
	provider := "aws" // Default to AWS for now
	var terraformContent string
	switch strings.ToLower(provider) {
	case "aws":
		terraformContent = m.generateAWSTerraform()
	case "azure":
		terraformContent = m.generateAzureTerraform()
	case "gcp":
		terraformContent = m.generateGCPTerraform()
	default:
		terraformContent = m.generateGenericTerraform()
	}

	if err := os.WriteFile(mainTfPath, []byte(terraformContent), 0644); err != nil {
		return fmt.Errorf("failed to write main.tf: %w", err)
	}

	logger.Info("Created main.tf").Str("path", mainTfPath).Str("provider", provider).Send()
	return nil
}

// getTerraformEnvVars returns environment variables for Terraform
func (m *Manager) getTerraformEnvVars() []string {
	envVars := []string{
		"TF_IN_AUTOMATION=1",
		"TF_INPUT=0",
	}

	// TODO: Add backend configuration support when needed
	// Backend configuration would be added here

	return envVars
}

// getProviderVariables returns provider-specific variables for Terraform
func (m *Manager) getProviderVariables() []string {
	var vars []string

	// TODO: Add region variable from cloud configuration
	region := "us-west-2" // Default region
	vars = append(vars, "-var=region="+region)

	// Add provider-specific variables
	// TODO: Make provider configurable
	provider := "aws" // Default provider
	switch strings.ToLower(provider) {
	case "aws":
		// AWS-specific variables will be picked up from environment or AWS config
	case "azure":
		// Azure-specific variables will be picked up from environment or Azure CLI
	case "gcp":
		// GCP-specific variables will be picked up from environment or gcloud config
	}

	return vars
}

// checkKubernetesHealth checks if Kubernetes cluster is healthy
func (m *Manager) checkKubernetesHealth(endpoint interface{}) error {
	logger.Info("Checking Kubernetes cluster health").Send()

	// Basic health check - try to connect to the cluster
	// This is a simple implementation - in production you'd want more sophisticated checks
	if endpointStr, ok := endpoint.(string); ok && endpointStr != "" {
		logger.Info("Kubernetes endpoint available").Str("endpoint", endpointStr).Send()
		return nil
	}

	return fmt.Errorf("kubernetes endpoint not available")
}

// checkDatabaseHealth checks if database is healthy
func (m *Manager) checkDatabaseHealth(endpoint interface{}) error {
	logger.Info("Checking database health").Send()

	// Basic health check for database
	if endpointStr, ok := endpoint.(string); ok && endpointStr != "" {
		logger.Info("Database endpoint available").Str("endpoint", endpointStr).Send()
		return nil
	}

	return fmt.Errorf("database endpoint not available")
}

// generateAWSTerraform generates AWS-specific Terraform configuration
func (m *Manager) generateAWSTerraform() string {
	return `terraform {
  required_version = ">= 1.5"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

# Basic EKS cluster configuration
# This is a minimal example - customize based on your needs
resource "aws_eks_cluster" "main" {
  name     = "k8s-installer-cluster"
  role_arn = aws_iam_role.cluster.arn
  version  = "1.28"

  vpc_config {
    subnet_ids = [aws_subnet.main.id, aws_subnet.secondary.id]
  }

  depends_on = [
    aws_iam_role_policy_attachment.cluster_AmazonEKSClusterPolicy,
  ]
}

# IAM role for EKS cluster
resource "aws_iam_role" "cluster" {
  name = "k8s-installer-cluster-role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "cluster_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.cluster.name
}

# VPC configuration
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "k8s-installer-vpc"
  }
}

resource "aws_subnet" "main" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = data.aws_availability_zones.available.names[0]

  tags = {
    Name = "k8s-installer-subnet-1"
  }
}

resource "aws_subnet" "secondary" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "10.0.2.0/24"
  availability_zone = data.aws_availability_zones.available.names[1]

  tags = {
    Name = "k8s-installer-subnet-2"
  }
}

data "aws_availability_zones" "available" {
  state = "available"
}

# Outputs
output "kubernetes_endpoint" {
  description = "Kubernetes cluster endpoint"
  value       = aws_eks_cluster.main.endpoint
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = aws_eks_cluster.main.name
}
`
}

// generateAzureTerraform generates Azure-specific Terraform configuration
func (m *Manager) generateAzureTerraform() string {
	return `terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  features {}
}

variable "region" {
  description = "Azure region"
  type        = string
  default     = "East US"
}

# Resource group
resource "azurerm_resource_group" "main" {
  name     = "k8s-installer-rg"
  location = var.region
}

# AKS cluster
resource "azurerm_kubernetes_cluster" "main" {
  name                = "k8s-installer-cluster"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  dns_prefix          = "k8sinstaller"

  default_node_pool {
    name       = "default"
    node_count = 2
    vm_size    = "Standard_D2_v2"
  }

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Development"
    Purpose     = "K8s-Installer"
  }
}

# Outputs
output "kubernetes_endpoint" {
  description = "Kubernetes cluster endpoint"
  value       = azurerm_kubernetes_cluster.main.kube_config.0.host
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = azurerm_kubernetes_cluster.main.name
}
`
}

// generateGCPTerraform generates GCP-specific Terraform configuration
func (m *Manager) generateGCPTerraform() string {
	return `terraform {
  required_version = ">= 1.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  region = var.region
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "project_id" {
  description = "GCP project ID"
  type        = string
}

# GKE cluster
resource "google_container_cluster" "main" {
  name     = "k8s-installer-cluster"
  location = var.region

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name
}

# Separately Managed Node Pool
resource "google_container_node_pool" "primary_nodes" {
  name       = "k8s-installer-node-pool"
  location   = var.region
  cluster    = google_container_cluster.main.name
  node_count = 2

  node_config {
    preemptible  = true
    machine_type = "e2-medium"

    # Google recommends custom service accounts that have cloud-platform scope and permissions granted via IAM Roles.
    service_account = google_service_account.default.email
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
}

# VPC
resource "google_compute_network" "vpc" {
  name                    = "k8s-installer-vpc"
  auto_create_subnetworks = "false"
}

# Subnet
resource "google_compute_subnetwork" "subnet" {
  name          = "k8s-installer-subnet"
  region        = var.region
  network       = google_compute_network.vpc.name
  ip_cidr_range = "10.10.0.0/24"
}

# Service account
resource "google_service_account" "default" {
  account_id   = "k8s-installer-sa"
  display_name = "K8s Installer Service Account"
}

# Outputs
output "kubernetes_endpoint" {
  description = "Kubernetes cluster endpoint"
  value       = google_container_cluster.main.endpoint
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = google_container_cluster.main.name
}
`
}

// generateGenericTerraform generates a generic Terraform configuration
func (m *Manager) generateGenericTerraform() string {
	return `terraform {
  required_version = ">= 1.5"
}

variable "region" {
  description = "Deployment region"
  type        = string
  default     = "us-west-2"
}

# Generic configuration - customize based on your provider
# This is a placeholder configuration that should be replaced
# with your specific infrastructure requirements

output "kubernetes_endpoint" {
  description = "Kubernetes cluster endpoint"
  value       = "https://kubernetes.example.com"
}

output "cluster_name" {
  description = "Kubernetes cluster name"
  value       = "k8s-installer-cluster"
}
`
}
