
# Enterprise-Grade E2E Kubernetes Installer

A comprehensive, production-ready Go-based CLI tool for deploying and managing Kubernetes clusters across multiple cloud environments with enterprise-grade security, monitoring, and validation capabilities.

[![Go Report Card](https://goreportcard.com/badge/github.com/judebantony/e2e-k8s-installer)](https://goreportcard.com/report/github.com/judebantony/e2e-k8s-installer)
[![GitHub Release](https://img.shields.io/github/release/judebantony/e2e-k8s-installer.svg)](https://github.com/judebantony/e2e-k8s-installer/releases/latest)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## 🎯 Overview

**Enterprise-grade Kubernetes installer** that automates the complete deployment lifecycle across AWS, Azure, GCP, and on-premises environments. Supports both connected and air-gapped deployments with multi-mode infrastructure provisioning.

Design and develop a unified, cross-platform, end-to-end (E2E) installer for deploying and managing a Kubernetes-based application across startup, mid-size, enterprise, and air-gapped environments.
The installer automates provisioning, configuration, deployment, validation, and lifecycle management across Azure, AWS, GCP, and on-premises (OpenShift, Rancher, etc.) infrastructures while ensuring compliance, security, and resilience.

### Key Features

- **🏗️ Multi-Mode Infrastructure**: Terraform, Makefile, and Hybrid provisioning modes
- **📦 Artifact Management**: OCI images, Helm charts, and Terraform modules synchronization
- **☁️ Multi-Cloud Ready**: AWS EKS, Azure AKS, GCP GKE, and on-premises support
- **🔒 Security-First**: Enterprise authentication, RBAC, and compliance scanning
- **📊 Enterprise Observability**: Structured logging, progress tracking, and comprehensive reporting
- **🔄 Air-gapped Support**: Complete offline installation capabilities

## 🎯 Core Objectives

- **Simplify Deployment**: Streamline application deployment and lifecycle management in Kubernetes clusters
- **Air-gapped Support**: Provide offline installation and upgrade capabilities for secure environments
- **Compliance & Auditability**: Ensure compliance, auditability, idempotency, and cross-cloud portability
- **Flexible Installation**: Offer both interactive and non-interactive (config-driven) installation modes
- **Self-contained Delivery**: Enable complete deployment without direct vendor package access

## 📋 Scope and Functional Requirements

### Artifact preparation & shipping into the client environment

```mermaid
flowchart TD
    subgraph "🏭 Vendor Environment"
        VR1[📦 GitHub Packages]
        VR2[📦 DockerHub] 
        VR3[📦 Azure ACR]
        VR4[📦 AWS ECR]
        VR5[📦 GCP Artifact Registry]
        
        VG1[📚 Vendor GitHub - Helm Charts]
        VG2[📚 Vendor GitHub - Terraform Modules]
        VG3[📚 Vendor GitHub - DB Scripts]
        
        VR1 & VR2 & VR3 & VR4 & VR5 --> |OCI Images| SCAN[🔍 Security Scan & Validation]
        VG1 & VG2 & VG3 --> |Git Artifacts| VERIFY[✅ Version & Tag Verification]
    end

    subgraph "🚀 Transfer Process"
        SCAN --> PULL[📥 Pull & Package]
        VERIFY --> PULL
        PULL --> MIRROR[🔄 Mirror/Transfer Decision]
        MIRROR --> |Mirror Enabled| PUSH[📤 Push to Client Registry]
        MIRROR --> |Local Checkout| LOCAL[💾 Local Storage]
    end

    subgraph "🏢 Client Environment" 
        CR1[🏬 Harbor Registry]
        CR2[🏬 Nexus Registry]
        CR3[🏬 JFrog Artifactory]
        
        CG1[📚 Client GitHub - Helm Charts]
        CG2[📚 Client GitHub - Terraform Modules] 
        CG3[📚 Client GitHub - DB Scripts]
        
        PUSH --> CR1 & CR2 & CR3
        LOCAL --> |Git Mirror| CG1 & CG2 & CG3
        LOCAL --> |Local Files| WORKSPACE[💼 Local Workspace]
    end

    subgraph "📊 Validation & Reporting"
        CR1 & CR2 & CR3 --> HEALTH[🏥 Health Checks]
        CG1 & CG2 & CG3 --> HEALTH
        WORKSPACE --> HEALTH
        HEALTH --> REPORT[📋 Package-Pull Report]
        REPORT --> |Success| READY[✅ Ready for Deployment]
        REPORT --> |Issues| ALERT[⚠️ Validation Alerts]
    end

    style VR1 fill:#E3F2FD,stroke:#1976D2,stroke-width:2px
    style VR2 fill:#E3F2FD,stroke:#1976D2,stroke-width:2px
    style VR3 fill:#E1F5FE,stroke:#0288D1,stroke-width:2px
    style VR4 fill:#FFF3E0,stroke:#F57C00,stroke-width:2px
    style VR5 fill:#E8F5E8,stroke:#388E3C,stroke-width:2px
    style SCAN fill:#FFF3E0,stroke:#FF8F00,stroke-width:2px
    style READY fill:#E8F5E8,stroke:#4CAF50,stroke-width:2px
    style ALERT fill:#FFEBEE,stroke:#F44336,stroke-width:2px
```

**Detailed Artifact Flow:**

- **🖼️ OCI Container Images**: Transfer from vendor registries (GitHub Packages, DockerHub, Azure ACR, AWS ECR, GCP Artifact Registry) → client's private registry (Harbor, Nexus, JFrog Artifactory) with security scanning
- **📊 Helm Charts**: Migration from vendor GitHub → client GitHub (or maintain local checkout if mirroring is disabled). Charts are versioned and tagged in vendor repositories
- **🏗️ Terraform Modules**: Transfer from vendor GitHub → client GitHub (or local checkout), with version control and tagging
- **🗃️ Database Migration Scripts**: Transfer from vendor GitHub → client GitHub (or local checkout), including repair and migration scripts
- **🔍 Health Checks & Validation**: Comprehensive verification, scanning, and readiness validation with detailed JSON reporting

### Full E2E installation once artifacts are in client environment

- ✅ **Infrastructure Provisioning**: Terraform-based deployment (K8s clusters, managed DBs/services)
- ✅ **Database Migrations & Repair**: Flyway/Liquibase execution as Kubernetes Job/init container
- ✅ **Application Deployment**: Helm-based deployment in configured order with pod readiness + health URL checks
- ✅ **Post-installation Validation**: Ensure all components are correctly installed and configured
- ✅ **Configuration Drift Detection**: Monitor and report any changes to the deployed environment
- ✅ **Post-validation**: Comprehensive checks, housekeeping, and E2E smoke tests
- ✅ **Reporting & Audit Logs**: Detailed JSON reports and structured logs for every step

## 🚀 Current Status

| Component | Status | Description |
|-----------|---------|-------------|
| **🔧 setup Command** | ✅ Complete | Workspace initialization and prerequisite validation |
| **📦 package-pull Command** | ✅ Complete | Artifact synchronization (OCI/Helm/Terraform) |
| **☁️ provision-infra Command** | ✅ Complete | Multi-mode infrastructure provisioning |
| **🗄️ db-migrate Command** | 🚧 In Progress | Database migration framework |
| **🚀 deploy Command** | 🚧 In Progress | Helm-based application deployment |
| **✅ post-validate & e2e-test** | 🔄 Planned | Validation and testing framework |

## 🧭 Flow Diagram

```mermaid
flowchart TD
 subgraph subGraph0["💡 Each Step is Idempotent & Re-runnable"]
        A["🏁 set-up"]
        B["📦 package-pull"]
        C["☁️ provision-infra"]
        D["🧩 db-migrate"]
        E["🚀 deploy"]
        F["🔍 post-validate"]
        G["✅ e2e-test"]
  end
    A --> B
    B --> C
    C --> D
    D --> E
    E --> F
    F --> G
    G --> H["📊 install-summary.json"]
    style A fill:#C6E2FF,stroke:#0366d6,stroke-width:2px
    style B fill:#E0F7FA,stroke:#00ACC1,stroke-width:2px
    style C fill:#E8F5E9,stroke:#2E7D32,stroke-width:2px
    style D fill:#FFF3E0,stroke:#EF6C00,stroke-width:2px
    style E fill:#E8EAF6,stroke:#3949AB,stroke-width:2px
    style F fill:#F3E5F5,stroke:#8E24AA,stroke-width:2px
    style G fill:#FBE9E7,stroke:#D84315,stroke-width:2px
    style H fill:#FFF9C4,stroke:#F9A825,stroke-width:2px
```

## 🧱 Layered Architecture Diagram

```mermaid
graph LR
    subgraph CLI["🧰 Installer CLI"]
        A1[set-up]
        A2[package-pull]
        A3[provision-infra]
        A4[db-migrate]
        A5[deploy]
        A6[post-validate]
        A7[e2e-test]
        A8[install or orchestrator]
    end

    subgraph Core["⚙️ Core Framework"]
        B1[Config Loader - JSON Schema, Validator]
        B2[Logger & ProgressBar - JSONL, pterm]
        B3[Report Generator - per-step summaries]
        B4[Error & Retry Handler - Idempotent logic]
    end

    subgraph Integration["🔗 Integration Modules"]
        C1[OCI Manager - Pull/Push/Scan Images]
        C2[Git Manager - Checkout/Mirror Repos]
        C3[Terraform Runner - Infra Provisioning]
        C4[Helm Deployer - Charts & Values]
        C5[K8s Client - Health Checks & Jobs]
        C6[DB Migration - Flyway/Liquibase]
    end

    subgraph External["🌐 External Systems"]
        D1[Vendor OCI Registries - GitHub, DockerHub, Azure]
        D2[Client Private Registries - Harbor, Artifactory]
        D3[Vendor GitHub Repos - Helm, Terraform, DB]
        D4[Client GitHub Repos - Mirrors]
        D5[Cloud Providers - AWS, Azure, GCP]
        D6[Kubernetes Cluster - OpenShift, Rancher, etc.]
    end

    %% connections
    A1 --> B1
    A2 --> C1
    A2 --> C2
    A3 --> C3
    A4 --> C6
    A5 --> C4
    A5 --> C5
    A6 --> C5
    A7 --> C5
    A8 --> A2
    A8 --> A3
    A8 --> A4
    A8 --> A5
    A8 --> A6
    A8 --> A7

    C1 --> D1
    C1 --> D2
    C2 --> D3
    C2 --> D4
    C3 --> D5
    C4 --> D6
    C5 --> D6
    C6 --> D6

    B1 --> B2
    B2 --> B3
    B3 --> B4

```

## 🧱 Prerequisite & Dependency Flow

```mermaid
sequenceDiagram
    participant U as 👤 User
    participant CLI as 🧰 Installer CLI
    participant CFG as ⚙️ Config File
    participant REG as 🏗️ Registries
    participant GIT as 🧬 GitHub Repos
    participant CLOUD as ☁️ Cloud Providers
    participant K8S as 🌀 Kubernetes Cluster

    U->>CLI: Run "installer install --config installer.json"
    CLI->>CFG: Validate configuration & credentials
    CLI->>REG: Authenticate to Vendor & Client registries
    CLI->>GIT: Clone Helm, Terraform, DB repos
    CLI->>CLOUD: Provision infra using Terraform
    CLI->>K8S: Create/validate cluster connectivity
    CLI->>REG: Pull/mirror OCI images (package-pull)
    CLI->>K8S: Run DB migrations (Flyway/Liquibase)
    CLI->>K8S: Deploy Helm charts in sequence
    CLI->>K8S: Run health checks & smoke tests
    CLI->>CLI: Generate structured JSON reports
    CLI-->>U: Display progress bars & final summary
```

## 🏗️ Architecture

### Multi-Mode Infrastructure Provisioning

The installer supports three distinct infrastructure provisioning modes:

#### **Terraform Mode** 🏗️

Pure Terraform-based infrastructure provisioning for cloud-native deployments.

#### **Makefile Mode** ⚙️

Makefile-based workflows for custom provisioning scripts and legacy systems.

#### **Hybrid Mode** 🔄

Combined approach where Makefiles orchestrate Terraform modules internally.

### System Architecture

```plaintext
┌─────────────────────────────────────────────────────────────────────┐
│                      CLI Application (main.go)                      │
└─────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────┐
│                      Command Layer (cmd/)                          │
│  setup | package-pull | provision-infra | deploy | db-migrate      │
└─────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────┐
│                   Business Logic Layer (pkg/)                      │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐       │
│  │ infrastructure/ │ │   artifacts/    │ │    config/      │       │
│  │ Multi-Mode Mgr  │ │ OCI/Helm/Git    │ │ JSON Validation │       │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘       │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐       │
│  │   terraform/    │ │   makefile/     │ │ logger/progress │       │
│  │   Operations    │ │   Execution     │ │ UI Components   │       │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘       │
└─────────────────────────────────────────────────────────────────────┘
                                     │
┌─────────────────────────────────────────────────────────────────────┐
│              External Integrations                                  │
│  Cloud Providers | OCI Registries | Git Repositories               │
└─────────────────────────────────────────────────────────────────────┘
```

### Package Structure

```plaintext
e2e-k8s-installer/
├── main.go                         # Application entry point
├── cmd/                           # Command implementations
│   ├── setup.go                   # Workspace initialization
│   ├── package_pull.go            # Artifact synchronization
│   ├── provision_infra.go         # Multi-mode infrastructure
│   ├── deploy.go                  # Application deployment
│   └── install.go                 # Full workflow orchestration
├── pkg/                          # Business logic
│   ├── infrastructure/           # Multi-mode infrastructure manager
│   ├── terraform/               # Terraform operations
│   ├── makefile/                # Makefile execution
│   ├── artifacts/               # OCI/Helm/Git management
│   ├── config/                  # Configuration & validation
│   └── logger/                  # Structured logging & progress
└── configs/                     # Sample configurations
```

## ⚙️ Configuration

### Multi-Mode Infrastructure Examples

**Terraform Mode:**

```json
{
  "infrastructure": {
    "provisionMode": "terraform",
    "provider": "aws",
    "region": "us-west-2",
    "terraform": {
      "enabled": true,
      "workingDir": "./terraform",
      "varsFile": "terraform.tfvars"
    }
  }
}
```

**Makefile Mode:**

```json
{
  "infrastructure": {
    "provisionMode": "makefile",
    "makefile": {
      "enabled": true,
      "makefilePath": "./Makefile",
      "targets": ["init", "plan", "apply"],
      "timeout": "30m"
    }
  }
}
```

**Hybrid Mode:**

```json
{
  "infrastructure": {
    "provisionMode": "hybrid",
    "terraform": { "enabled": true },
    "makefile": { "enabled": true }
  }
}
```

## 🎮 Usage

### Quick Start

```bash
# 1. Initialize workspace
./e2e-k8s-installer setup --workspace ./my-k8s-project

# 2. Navigate and configure
cd my-k8s-project
vim installer-config.json

# 3. Synchronize artifacts
./e2e-k8s-installer package-pull --config installer-config.json

# 4. Provision infrastructure
./e2e-k8s-installer provision-infra --config installer-config.json

# 5. Deploy applications (coming soon)
# ./e2e-k8s-installer deploy --config installer-config.json
```

### Available Commands

| Command | Status | Description |
|---------|---------|-------------|
| `setup` | ✅ Ready | Initialize workspace and validate prerequisites |
| `package-pull` | ✅ Ready | Synchronize OCI images, Helm charts, Terraform modules |
| `provision-infra` | ✅ Ready | Deploy infrastructure (terraform/makefile/hybrid modes) |
| `deploy` | 🚧 In Progress | Deploy applications with Helm |
| `db-migrate` | 🚧 In Progress | Run database migrations |
| `install` | 🔄 Planned | Complete workflow orchestration |

### Command Examples

**Setup workspace:**

```bash
./e2e-k8s-installer setup --workspace ./project --config-file custom.json
```

**Pull artifacts:**

```bash
./e2e-k8s-installer package-pull --config config.json --images-only
```

**Provision infrastructure:**

```bash
# Terraform mode
./e2e-k8s-installer provision-infra --config terraform-config.json

# Makefile mode  
./e2e-k8s-installer provision-infra --config makefile-config.json

# Plan only (dry run)
./e2e-k8s-installer provision-infra --config config.json --plan-only
```

## 📋 Requirements

### System Prerequisites

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **OS** | Linux, macOS, Windows | Linux (Ubuntu 20.04+) |
| **CPU** | 2 cores | 4+ cores |
| **Memory** | 4GB RAM | 8GB+ RAM |
| **Storage** | 10GB free | 20GB+ free |

### Required Dependencies

```bash
# Essential tools (required)
kubectl 1.28+     # Kubernetes CLI
helm 3.8+         # Package manager  
terraform 1.5+    # Infrastructure as Code
git 2.30+         # Version control

# Cloud provider tools (choose based on target)
aws-cli 2.0+      # For AWS deployments
azure-cli 2.30+   # For Azure deployments
gcloud 400.0+     # For GCP deployments
```

## 🚀 Installation

### Option 1: Download Binary (Recommended)

```bash
# Download latest release
curl -L https://github.com/judebantony/e2e-k8s-installer/releases/latest/download/e2e-k8s-installer-$(uname -s)-$(uname -m) -o e2e-k8s-installer
chmod +x e2e-k8s-installer
sudo mv e2e-k8s-installer /usr/local/bin/
```

### Option 2: Build from Source

```bash
git clone https://github.com/judebantony/e2e-k8s-installer.git
cd e2e-k8s-installer
go build -o e2e-k8s-installer .
sudo mv e2e-k8s-installer /usr/local/bin/
```

### Verify Installation

```bash
e2e-k8s-installer --help
e2e-k8s-installer setup --help
```

## 🔒 Security & Enterprise Features

### Current Security Implementation

- **✅ Input Validation**: Comprehensive configuration validation
- **✅ Credential Security**: Secure handling of registry credentials and tokens
- **✅ Network Security**: TLS-enabled communications
- **✅ Audit Logging**: Complete audit trail with structured logging

### Planned Security Features

- **🔄 RBAC Integration**: Role-based access control for Kubernetes
- **🔄 Network Policies**: Pod-to-pod communication control
- **🔄 Security Scanning**: Container vulnerability detection
- **🔄 Compliance**: CIS Kubernetes benchmark validation

## 📊 Monitoring & Observability

### Built-in Features

- **Structured Logging**: High-performance JSON logging with zerolog
- **Progress Tracking**: Real-time progress bars with pterm
- **Command Auditing**: Complete audit trail of operations
- **Performance Metrics**: Command timing and resource usage

### Planned Monitoring Stack

- **Prometheus**: Metrics collection and alerting
- **Grafana**: Visualization dashboards
- **ELK Stack**: Centralized log management
- **Jaeger**: Distributed tracing

## 🔧 Troubleshooting

### Common Issues

**Prerequisites validation failed:**

```bash
# Check required tools
which kubectl helm terraform git

# Install missing tools
brew install kubectl helm terraform git  # macOS
```

**Registry authentication failed:**

```bash
# Test registry connectivity
curl -v https://your-registry.io/v2/

# Verify credentials in config
cat installer-config.json | jq .artifacts.images.vendor.auth
```

**Configuration validation errors:**

```bash
# Validate JSON syntax
jq . installer-config.json

# Dry run to check configuration
./e2e-k8s-installer provision-infra --config config.json --dry-run
```

### Debug Mode

```bash
# Enable verbose logging
./e2e-k8s-installer setup --workspace ./test --verbose

# Dry run mode
./e2e-k8s-installer package-pull --config config.json --dry-run
```

## 🤝 Contributing

### Development Setup

```bash
git clone https://github.com/judebantony/e2e-k8s-installer.git
cd e2e-k8s-installer
go mod tidy
go test ./...
go build -o bin/e2e-k8s-installer .
```

### Contribution Guidelines

1. Fork the repository and create a feature branch
2. Write tests for new functionality
3. Follow Go conventions and run `gofmt`
4. Update documentation for new features
5. Submit a pull request with detailed description

## 📚 Documentation & Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/judebantony/e2e-k8s-installer/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/judebantony/e2e-k8s-installer/discussions)
- 📖 **Documentation**: [docs/](docs/)
- 🌟 **Roadmap**: [ROADMAP.md](ROADMAP.md)

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## 🎯 Roadmap

- **✅ Phase 1 (v1.1.0)**: Multi-mode infrastructure provisioning - **COMPLETED**
- **🚧 Phase 2 (v1.2.0)**: Database migrations and application deployment - **IN PROGRESS**  
- **🔄 Phase 3 (v1.3.0)**: Validation and testing framework - **PLANNED**
- **🔄 Phase 4 (v1.4.0)**: Complete workflow orchestration - **PLANNED**

---

***Made with ❤️ by Jude Antony***
