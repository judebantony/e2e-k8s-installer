# Enterprise-Grade E2E Kubernetes Installer

A comprehensive, production-ready Go-based CLI tool for deploying and managing Kubernetes clusters across multiple cloud environments with enterprise-grade security, monitoring, and validation capabilities.

## 🎯 Project Overview

This installer provides a unified approach to Kubernetes cluster deployment with:

- **CLI-First Design**: Built with Go and Cobra framework for robust command-line operations
- **Multi-Cloud Native**: Seamless deployment across AWS, Azure, GCP, and on-premises
- **Security-First**: Integrated security scanning, RBAC, and policy enforcement
- **Enterprise Ready**: Production-grade monitoring, logging, and operational tools
- **Validation-Driven**: Comprehensive pre-flight and post-deployment validation

## 🚀 Features & Core Capabilities

### 🎛️ **Installation & Deployment**

| Feature | Description | Benefits |
|---------|-------------|----------|
| **Multi-Phase Installation** | Phased deployment with checkpoint validation | Reliable rollback, progress tracking |
| **Interactive Mode** | Guided setup with intelligent defaults | User-friendly, reduces configuration errors |
| **Dry-Run Support** | Preview changes before execution | Risk-free planning, validation |
| **Resume & Rollback** | Continue from failed phases or revert | Resilient deployments, quick recovery |
| **Configuration Templates** | Pre-built configs for common scenarios | Faster setup, best practices included |

### ☁️ **Multi-Cloud Support**

| Cloud Provider | Features | Supported Services |
|----------------|----------|-------------------|
| **AWS EKS** | ✅ Full automation, VPC setup, IAM roles | EKS, EC2, VPC, ALB, Route53 |
| **Azure AKS** | ✅ Resource group management, RBAC | AKS, Virtual Networks, Load Balancer |
| **Google GKE** | ✅ Project setup, service accounts | GKE, Compute Engine, Cloud Load Balancing |
| **On-Premises** | ✅ Kubeadm, custom networking | Bare metal, VMware, OpenStack |

### 🛡️ **Security Framework**

| Security Component | Capability | Implementation |
|-------------------|------------|----------------|
| **RBAC Management** | Role-based access control | Custom roles, service accounts |
| **Network Policies** | Pod-to-pod communication control | Calico, Cilium integration |
| **Security Scanning** | Vulnerability detection | Trivy, Aqua Security |
| **Runtime Security** | Real-time threat detection | Falco, Sysdig integration |
| **Policy Enforcement** | Admission controllers | OPA Gatekeeper policies |
| **Secret Management** | Encrypted secret storage | External Secrets, Sealed Secrets |

### 📊 **Monitoring & Observability**

| Component | Purpose | Features |
|-----------|---------|----------|
| **Prometheus** | Metrics collection & alerting | Custom metrics, alert rules, federation |
| **Grafana** | Visualization & dashboards | Pre-built dashboards, custom panels |
| **ELK Stack** | Centralized logging | Log aggregation, search, analysis |
| **Jaeger** | Distributed tracing | Request tracing, performance analysis |
| **Falco** | Security monitoring | Runtime security events, compliance |

### 🔧 **Operational Excellence**

| Feature | Description | Use Cases |
|---------|-------------|-----------|
| **Health Checks** | Comprehensive cluster validation | Post-deployment verification |
| **Backup & Recovery** | Automated state management | Disaster recovery, data protection |
| **Upgrade Management** | Rolling cluster upgrades | Zero-downtime updates |
| **Scaling Operations** | Auto & manual scaling | Load adaptation, cost optimization |
| **Certificate Management** | Automated TLS provisioning | Security compliance, cert rotation |

### 🌐 **Networking & Service Mesh**

| Technology | Integration | Benefits |
|------------|-------------|----------|
| **CNI Plugins** | Calico, Flannel, Cilium | Network isolation, policy enforcement |
| **Ingress Controllers** | NGINX, Traefik, Istio Gateway | Traffic routing, SSL termination |
| **Service Mesh** | Istio, Linkerd | Traffic management, security, observability |
| **Load Balancing** | Cloud-native LB integration | High availability, traffic distribution |

### 🔄 **CI/CD Integration**

| Platform | Support | Features |
|----------|---------|----------|
| **GitHub Actions** | ✅ Native workflows | Automated testing, deployment |
| **GitLab CI** | ✅ Pipeline integration | Container registry, deployment stages |
| **Jenkins** | ✅ Plugin support | Custom pipelines, artifact management |
| **ArgoCD** | ✅ GitOps workflows | Declarative deployments, sync policies |

### 📝 **Configuration Management**

| Feature | Capability | Benefits |
|---------|------------|----------|
| **YAML Configuration** | Declarative cluster definition | Version control, reproducibility |
| **Environment Profiles** | Dev/staging/prod templates | Consistent deployments across environments |
| **Variable Substitution** | Dynamic configuration values | Environment-specific customization |
| **Configuration Validation** | Schema-based validation | Early error detection, compliance |
| **Secret Management** | Encrypted sensitive data | Security best practices |

### 🧪 **Validation & Testing**

| Validation Type | Scope | Checks |
|----------------|-------|---------|
| **Pre-flight** | Environment readiness | Dependencies, permissions, connectivity |
| **Configuration** | YAML validation | Schema compliance, resource limits |
| **Post-deployment** | Cluster health | Service status, networking, security |
| **Compliance** | Security standards | CIS benchmarks, best practices |
| **Performance** | Resource utilization | CPU, memory, network performance |

### 🔌 **Extensibility**

| Extension Point | Capability | Examples |
|----------------|------------|----------|
| **Custom Providers** | Plugin architecture | Private cloud integrations |
| **Hook System** | Pre/post deployment scripts | Custom validations, notifications |
| **Template Engine** | Custom resource templates | Organization-specific resources |
| **API Integration** | REST API for automation | External tool integration |

### 📈 **Enterprise Features**

| Feature | Description | Enterprise Value |
|---------|-------------|------------------|
| **Multi-tenancy** | Namespace isolation & quotas | Resource governance, cost allocation |
| **Compliance Reporting** | Automated compliance checks | Audit trails, regulatory compliance |
| **Cost Management** | Resource usage tracking | Budget control, optimization insights |
| **Support Bundle** | Diagnostic data collection | Faster troubleshooting, support |
| **Air-gapped Support** | Offline installation capability | Secure environments, compliance |

## 🏗️ Architecture

### Core Components
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Layer     │    │  Validation      │    │  Installation   │
│   (Cobra)       │───▶│  Engine          │───▶│  Engine         │
│                 │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│ Configuration   │    │  Cloud Provider  │    │  Security       │
│ Management      │    │  Modules         │    │  Framework      │
│ (Viper/YAML)    │    │  (AWS/Azure/GCP) │    │  (RBAC/Policies)│
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### Key Features
- 🔧 **Multi-Phase Installation**: Phased deployment with rollback capabilities
- 🛡️ **Integrated Security**: Trivy, Falco, OPA Gatekeeper integration
- 📊 **Built-in Monitoring**: Prometheus, Grafana, ELK stack deployment
- ☁️ **Cloud Agnostic**: Terraform-based infrastructure provisioning
- 🔍 **Comprehensive Validation**: Pre-flight checks and environment validation
- 📝 **Configuration-Driven**: YAML-based declarative configuration

## 📋 Requirements

### System Prerequisites

| Component | Minimum | Recommended | Notes |
|-----------|---------|-------------|-------|
| **Operating System** | Linux, macOS, Windows | Linux (Ubuntu 20.04+) | Cross-platform Go binary |
| **CPU** | 2 cores | 4+ cores | For build and deployment |
| **Memory** | 4GB RAM | 8GB+ RAM | Depends on cluster size |
| **Storage** | 20GB free | 50GB+ free | For images and logs |
| **Network** | Internet access | High bandwidth | For container images |

### Required Dependencies

#### Core Tools
```bash
# Go Development (for building from source)
Go 1.21+                    # Language runtime
Git 2.30+                  # Version control

# Kubernetes Tools
kubectl 1.28+              # Kubernetes CLI
helm 3.8+                  # Package manager

# Infrastructure Tools
terraform 1.5+             # Infrastructure as Code
docker 20.10+              # Container runtime (optional for build)
```

#### Cloud Provider Tools (Choose based on target)
```bash
# AWS
aws-cli 2.0+               # AWS Command Line Interface
aws-iam-authenticator      # AWS IAM authentication

# Azure
azure-cli 2.30+            # Azure Command Line Interface

# Google Cloud
gcloud 400.0+              # Google Cloud SDK
```

### Optional Tools
```bash
# Security Scanning
trivy latest               # Vulnerability scanner
falco latest               # Runtime security

# Monitoring (auto-installed by installer)
prometheus                 # Metrics collection
grafana                   # Visualization
elasticsearch             # Log storage
```

## � Quick Start

### Installation Options

#### Option 1: Download Pre-built Binary (Recommended)
```bash
# Download latest release
curl -L https://github.com/judebantony/e2e-k8s-installer/releases/latest/download/k8s-installer-$(uname -s)-$(uname -m) -o k8s-installer
chmod +x k8s-installer
sudo mv k8s-installer /usr/local/bin/
```

#### Option 2: Build from Source
```bash
# Clone repository
git clone https://github.com/judebantony/e2e-k8s-installer.git
cd e2e-k8s-installer

# Build binary
go build -o bin/k8s-installer .

# Optional: Install globally
sudo mv bin/k8s-installer /usr/local/bin/
```

### First Run Validation

```bash
# Check installer version
k8s-installer version

# Validate system requirements
k8s-installer validate system

# Generate sample configuration
k8s-installer config generate > sample-config.yaml
```

## ⚙️ Configuration

### Configuration Structure

The installer uses YAML configuration with the following structure:

```yaml
# cluster: Core cluster configuration
cluster:
  name: "my-k8s-cluster"           # Cluster identifier
  version: "1.28"                  # Kubernetes version
  provider: "aws"                  # Cloud provider (aws/azure/gcp/onprem)
  region: "us-west-2"             # Target region

# cloud: Provider-specific settings
cloud:
  aws:
    vpc_cidr: "10.0.0.0/16"       # VPC network range
    instance_type: "t3.medium"     # Node instance size
    node_count: 3                  # Initial node count
    zones: ["us-west-2a", "us-west-2b"]  # Availability zones

# kubernetes: K8s-specific configuration  
kubernetes:
  network_plugin: "calico"         # CNI plugin
  service_cidr: "10.96.0.0/12"    # Service network
  pod_cidr: "10.244.0.0/16"       # Pod network
  
# monitoring: Observability stack
monitoring:
  prometheus:
    enabled: true                  # Enable metrics collection
    retention: "30d"               # Data retention period
  grafana:
    enabled: true                  # Enable dashboards
    
# security: Security configurations
security:
  rbac: true                       # Enable RBAC
  network_policies: true           # Enable network policies
  scanning:
    enabled: true                  # Enable vulnerability scanning
    tools: ["trivy", "falco"]      # Security tools
```

### Configuration Examples

#### AWS EKS Configuration
```yaml
cluster:
  name: "production-eks"
  version: "1.28"
  provider: "aws"
  region: "us-west-2"

cloud:
  aws:
    vpc_cidr: "10.0.0.0/16"
    instance_type: "m5.large"
    node_count: 3
    zones: ["us-west-2a", "us-west-2b", "us-west-2c"]
    
monitoring:
  prometheus:
    enabled: true
    storage: "100Gi"
  grafana:
    enabled: true
```

#### Azure AKS Configuration  
```yaml
cluster:
  name: "production-aks"
  version: "1.28"
  provider: "azure"
  region: "East US"

cloud:
  azure:
    resource_group: "k8s-rg"
    vm_size: "Standard_D2s_v3"
    node_count: 3
    
security:
  rbac: true
  network_policies: true
```

## 🎮 Usage Guide

### Available Commands

```bash
# Core Commands
k8s-installer install     # Install new cluster
k8s-installer validate    # Validate environment/config
k8s-installer status      # Check cluster status
k8s-installer upgrade     # Upgrade cluster
k8s-installer config      # Configuration management
k8s-installer version     # Show version info

# Advanced Commands  
k8s-installer backup      # Backup cluster state
k8s-installer restore     # Restore from backup
k8s-installer scale       # Scale cluster resources
k8s-installer monitor     # Monitor cluster health
```

### Installation Workflows

#### Interactive Installation
```bash
# Start interactive wizard
k8s-installer install --interactive

# Follow prompts for:
# - Cloud provider selection
# - Cluster sizing
# - Security preferences
# - Monitoring setup
```

#### Configuration-Based Installation
```bash
# Generate base configuration
k8s-installer config generate --provider aws > aws-cluster.yaml

# Edit configuration
vim aws-cluster.yaml

# Validate configuration
k8s-installer validate --config aws-cluster.yaml

# Install with dry-run (recommended)
k8s-installer install --config aws-cluster.yaml --dry-run

# Actual installation
k8s-installer install --config aws-cluster.yaml
```

#### Multi-Environment Installation
```bash
# Development environment
k8s-installer install --config dev-config.yaml --environment dev

# Staging environment  
k8s-installer install --config staging-config.yaml --environment staging

# Production environment
k8s-installer install --config prod-config.yaml --environment production
```

### Validation Commands

```bash
# System validation
k8s-installer validate system              # Check prerequisites
k8s-installer validate network             # Test connectivity
k8s-installer validate config              # Validate configuration
k8s-installer validate cloud --provider aws # Check cloud access

# Cluster validation
k8s-installer validate cluster             # Post-install validation
k8s-installer validate security            # Security compliance check
k8s-installer validate monitoring          # Monitoring stack check
```

## � Project Structure

```text
e2e-k8s-installer/
├── cmd/                    # CLI command implementations
│   ├── root.go            # Root command and global flags
│   ├── install.go         # Installation command logic
│   ├── validate.go        # Validation command logic
│   ├── status.go          # Status checking command
│   ├── upgrade.go         # Cluster upgrade command
│   ├── config.go          # Configuration management
│   └── version.go         # Version information
├── pkg/                   # Core business logic packages
│   ├── config/           # Configuration parsing and validation
│   ├── installer/        # Core installation orchestration
│   ├── validation/       # Environment and prerequisite validation
│   ├── cloud/           # Cloud provider implementations
│   │   ├── aws/         # AWS-specific operations
│   │   ├── azure/       # Azure-specific operations
│   │   ├── gcp/         # GCP-specific operations
│   │   └── onprem/      # On-premises deployment
│   ├── k8s/             # Kubernetes cluster management
│   ├── monitoring/      # Monitoring stack deployment
│   ├── security/        # Security framework and policies
│   └── utils/           # Shared utilities and helpers
├── terraform/            # Infrastructure as Code modules
│   ├── aws/             # AWS infrastructure templates
│   ├── azure/           # Azure infrastructure templates
│   └── gcp/             # GCP infrastructure templates
├── charts/              # Helm charts for applications
│   ├── monitoring/      # Prometheus, Grafana charts
│   ├── security/        # Security tool charts
│   └── logging/         # ELK stack charts
├── scripts/             # Utility and automation scripts
├── docs/               # Comprehensive documentation
├── examples/           # Sample configurations
├── tests/              # Test suites
├── .github/            # GitHub workflows and templates
└── bin/                # Built binaries (gitignored)
```

## 🔒 Security Features

### Security Framework

The installer implements a comprehensive security framework:

#### RBAC (Role-Based Access Control)
```yaml
security:
  rbac:
    enabled: true
    custom_roles:
      - name: "developer"
        rules:
          - apiGroups: ["apps"]
            resources: ["deployments"]
            verbs: ["get", "list", "create", "update"]
```

#### Network Policies
```yaml
security:
  network_policies:
    enabled: true
    default_deny: true
    policies:
      - name: "allow-ingress"
        spec:
          podSelector: {}
          policyTypes: ["Ingress"]
```

#### Security Scanning
- **Container Scanning**: Trivy integration for vulnerability detection
- **Runtime Security**: Falco for runtime threat detection  
- **Policy Enforcement**: OPA Gatekeeper for admission control
- **Compliance**: CIS Kubernetes Benchmark validation

### Security Best Practices

1. **Least Privilege**: Minimal required permissions
2. **Network Segmentation**: Default deny network policies
3. **Image Security**: Signed and scanned container images
4. **Secret Management**: Encrypted secret storage
5. **Audit Logging**: Comprehensive audit trail

## 📊 Monitoring & Observability

### Monitoring Stack Components

#### Prometheus Stack
- **Prometheus Server**: Metrics collection and storage
- **Alertmanager**: Alert routing and management
- **Node Exporter**: System metrics collection
- **kube-state-metrics**: Kubernetes object metrics

#### Grafana Dashboards
- **Cluster Overview**: High-level cluster health
- **Node Metrics**: System resource utilization
- **Pod Metrics**: Application performance
- **Security Dashboard**: Security event monitoring

#### Logging Infrastructure
- **Elasticsearch**: Log storage and indexing
- **Logstash**: Log processing and enrichment
- **Kibana**: Log visualization and analysis
- **Fluentd**: Log collection and forwarding

### Custom Monitoring Configuration

```yaml
monitoring:
  prometheus:
    enabled: true
    retention: "30d"
    storage: "100Gi"
    resources:
      requests:
        cpu: "500m"
        memory: "1Gi"
    alerting:
      enabled: true
      rules:
        - alert: "HighCPUUsage"
          expr: "cpu_usage > 80"
          for: "5m"
  grafana:
    enabled: true
    admin_password: "secure_password"
    persistence:
      enabled: true
      size: "10Gi"
  logging:
    elasticsearch:
      enabled: true
      replicas: 3
      storage: "50Gi"
```

## 🚀 Cloud Provider Support

### AWS EKS Integration

```yaml
# AWS-specific configuration
cloud:
  aws:
    region: "us-west-2"
    vpc_cidr: "10.0.0.0/16"
    instance_types:
      - "t3.medium"    # General purpose
      - "c5.large"     # Compute optimized
    node_groups:
      - name: "general"
        instance_type: "t3.medium" 
        min_size: 1
        max_size: 10
        desired_size: 3
    addons:
      - "vpc-cni"
      - "coredns" 
      - "kube-proxy"
```

### Azure AKS Integration

```yaml
# Azure-specific configuration  
cloud:
  azure:
    location: "East US"
    resource_group: "k8s-cluster-rg"
    vm_size: "Standard_D2s_v3"
    node_pools:
      - name: "system"
        vm_size: "Standard_D2s_v3"
        node_count: 3
        mode: "System"
      - name: "user"  
        vm_size: "Standard_D4s_v3"
        node_count: 2
        mode: "User"
```

### GCP GKE Integration

```yaml
# GCP-specific configuration
cloud:
  gcp:
    project: "my-gcp-project"
    region: "us-central1"
    zone: "us-central1-a"
    machine_type: "e2-standard-4"
    node_pools:
      - name: "default-pool"
        machine_type: "e2-standard-4"
        disk_size: "100"
        node_count: 3
```

## 🔧 Troubleshooting & Support

### Debug Mode

```bash
# Enable verbose logging
k8s-installer install --config config.yaml --debug --verbose

# Save logs to file
k8s-installer install --config config.yaml --log-file install.log

# Dry run mode
k8s-installer install --config config.yaml --dry-run
```

### Common Issues & Solutions

#### Issue: Cloud Authentication Failures
```bash
# AWS
aws configure list
aws sts get-caller-identity

# Azure  
az account show
az ad signed-in-user show

# GCP
gcloud auth list
gcloud config list
```

#### Issue: Network Connectivity Problems
```bash
# Test connectivity
k8s-installer validate network

# Check DNS resolution
nslookup kubernetes.default.svc.cluster.local

# Verify firewall rules
kubectl get networkpolicies --all-namespaces
```

#### Issue: Resource Quota Exceeded
```bash
# Check resource usage
kubectl describe nodes
kubectl top nodes
kubectl top pods --all-namespaces

# Check quotas
kubectl describe quota --all-namespaces
```

### Support Bundle Collection

```bash
# Generate comprehensive support bundle
k8s-installer support-bundle \
  --output support-$(date +%Y%m%d-%H%M%S).tar.gz \
  --include-logs \
  --include-config \
  --include-cluster-info
```

## 🧪 Testing & Validation

### Test Suite

```bash
# Run all validation tests
k8s-installer test suite

# Run specific test categories
k8s-installer test --category security
k8s-installer test --category networking  
k8s-installer test --category monitoring

# Performance benchmarking
k8s-installer benchmark --duration 10m
```

### CI/CD Integration

```yaml
# GitHub Actions example
name: K8s Installer Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Build installer
        run: go build -o k8s-installer .
      - name: Run validation
        run: ./k8s-installer validate system
```

## 📚 Documentation

### Additional Resources

- 📖 **[Installation Guide](docs/installation.md)**: Detailed installation instructions
- 🔧 **[Configuration Reference](docs/configuration.md)**: Complete configuration options  
- 🛡️ **[Security Guide](docs/security.md)**: Security best practices
- 📊 **[Monitoring Guide](docs/monitoring.md)**: Observability setup
- 🚀 **[Deployment Examples](docs/examples/)**: Real-world deployment scenarios
- 🐛 **[Troubleshooting](docs/troubleshooting.md)**: Common issues and solutions

### API Reference

```bash
# Generate API documentation
k8s-installer docs generate --format markdown --output docs/api/

# View configuration schema
k8s-installer config schema
```

## 🤝 Contributing

### Development Setup

```bash
# Clone repository
git clone https://github.com/judebantony/e2e-k8s-installer.git
cd e2e-k8s-installer

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build and test
go build -o bin/k8s-installer .
./bin/k8s-installer validate system
```

### Contribution Guidelines

1. **Fork the repository** and create a feature branch
2. **Write tests** for new functionality
3. **Follow Go conventions** and run `gofmt`
4. **Update documentation** for new features
5. **Submit a pull request** with detailed description

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

## 🆘 Support & Community

- 🐛 **Issues**: [GitHub Issues](https://github.com/judebantony/e2e-k8s-installer/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/judebantony/e2e-k8s-installer/discussions)
- 📧 **Email**: support@k8s-installer.io
- 📖 **Documentation**: [docs/](docs/)
- 🌟 **Roadmap**: [ROADMAP.md](ROADMAP.md)

---

## 🎯 Roadmap & Future Features

- [ ] **Air-gapped Installation**: Complete offline deployment support
- [ ] **Service Mesh Integration**: Istio/Linkerd automatic setup
- [ ] **GitOps Integration**: ArgoCD/Flux deployment workflows
- [ ] **Multi-cluster Management**: Cross-cluster networking and policies
- [ ] **Backup & Recovery**: Automated backup and disaster recovery
- [ ] **Cost Optimization**: Resource usage optimization recommendations

---

**Made with ❤️ by Jude Antony**

*Star ⭐ this repository if you find it helpful!*
