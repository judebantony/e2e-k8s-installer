# GitHub Copilot Custom Instructions for E2E Installer for K8S application and any airgapped & high-security environments

## Project Context

You are a software engineer and an architect assisting with the design and development of an E2E Installer for a K8S application designed to function in startup, medium size, highly-secure and compliant, airgapped environments. The installer aims to simplify the deployment and management of applications in Kubernetes clusters without direct access to the vendor package providing enterprise solutions. The project focuses on creating a seamless installation experience, ensuring that all necessary components and dependencies are included in the installer package. And also need to setup the e2e infrastructure in Azure cloud, GCP, AWS and on-prem k8s infra setup like OpenShift.

Need to provision the k8s cluster using terraform and helm charts for deploying the application components. And also need to provision some the cloud managed service of different providers, so need to take parameter from the customer. 

Need to work on the Service Mesh, CI/CD pipeline using GitHub actions and ArgoCD for the deployment. Need to work on the monitoring and logging using ELK stack, Prometheus and Grafana. Need to work on the security aspect using OIDC 2.0 and OAuth2.0 using Auth0. Need to work on the API gateway for microservices architecture. Need to work on the caching using Redis and database using MongoDB and PostgreSQL. Need to work on the circuit breaker pattern for the resilience of the application. Need to work on the data streaming using Apache Kafka, RabbitMQ and Apache Flink. Need to work on the AI/ML integration using OpenAI API for some specific features.

The app images would be in OCI compliant registry like DockerHub, GitHub Package, Azure Container Registry, AWS ECR, GCP Artifact Registry. Need to ship to customer private registry like Harbor, Nexus or JFrog Artifactory and perform image scanning and vulnerability assessment.

The helm charts and terraform scripts would be in the GitHub repository. The installer should be able to pull the required images and helm charts from the private registry and deploy them to the k8s cluster. It should be able to check in to private git repository of the client and perform validation. The installer should also be able to handle the configuration of the application components, including setting up environment variables, secrets, and config maps.

Once the image and helm charts are pulled from the private registry, the installer should be able to deploy the application components to the k8s cluster using helm charts. The installer should also be able to perform post-installation checks to ensure that all components are running correctly and that the application is functioning as expected. And also do a dry-run before the actual installation. All steps should be logged for audit and compliance purpose. All the steps in the installation process should be idempotent, meaning that if the installation is interrupted or fails, it can be resumed from the point of failure without having to start over from the beginning.

Once the infrastructure installation is completed, need to have capability to run some db migration scripts using flyway or liquibase in an idempotent manner. It should be able to handle any errors or issues that may arise during the migration process and provide clear feedback to the user. It should be as init container or job in the k8s cluster.

Once the infrastructure provision and application is deployed, the installer should be able to perform some basic tests to ensure that the application is functioning correctly. This could include checking that all services are running, that the application is accessible, and that any necessary integrations are working as expected. The installer should provide clear feedback to the user about the results of these tests, including any issues that were encountered and how to resolve them. And also do the proper health check of the application and its components and clean up any temporary resources that were created during the installation process. It should have a capability to do the rollback if anything goes wrong during the installation process. And also need to have a capability to do the upgrade of the application in an idempotent manner. It should have capability to do some housekeeping tasks like log rotation, db backup and restore etc and run db insert/update script as post deployment.

Installer should be able to update the secrets and config maps in an idempotent manner. and also have a capability to do the secret management using HashiCorp Vault or Azure Key Vault or AWS Secrets Manager or GCP Secret Manager.

Once the installation is complete, the installer should provide a summary of the installation process, including any errors or issues that were encountered. The installer should also provide guidance on next steps, such as how to access the application, how to perform basic maintenance tasks, and where to find additional documentation or support resources.

The primary goal is to ensure that the installer is capable of handling the complexities associated with airgapped environments, such as managing dependencies, configurations, and network restrictions.

The installer should be user-friendly, robust, and capable of handling various configurations and scenarios that may arise during the installation process. And also able to provide full transparency to the user about the installation progress and any potential issues that may arise. Take the user credentials from the customer for the cloud provider and k8s cluster interactively or config.

The installer should be able to run on different OS like Windows, Linux and MacOS. The installer should be able to handle different versions of Kubernetes and should be compatible with popular distributions such as OpenShift, Rancher, and others. The installer should also be able to handle different networking configurations, such as Istio, Calico, Flannel, and Weave. It should be able to handle different storage configurations, such as NFS, GlusterFS, and Ceph. The installer should also be able to handle different authentication mechanisms, such as LDAP, Active Directory, and OAuth2.0.

The installer should be like to install interactively or in a non-interactive mode using configuration files. And able install like yum, brew and apt-get. The installer should have help and menu support. The installer should be able to provide detailed logs and error messages to help users troubleshoot any issues that may arise during the installation process. The installer should also be able to provide a rollback mechanism in case of any failures during the installation process.

In summary, the E2E Installer for K8S application in airgapped environments aims to provide a seamless and efficient installation experience, ensuring that users can deploy and manage their applications in Kubernetes clusters without internet access, while handling all necessary configurations, dependencies, and potential issues that may arise during the installation process.

## Project Overview

The E2E Installer for K8S application is a comprehensive solution designed to facilitate the deployment of applications in Kubernetes clusters, particularly in airgapped environments where internet access is restricted. The installer package includes all necessary components, dependencies, and configurations required for a successful installation.

## Key Components

1. **Pre-Installation Checks**: The installer performs a series of checks to ensure that the target environment meets all prerequisites for installation, including verifying Kubernetes cluster access and resource availability.

2. **Configuration Input**: Users can provide configuration details through an interactive command-line interface (CLI) or configuration files, specifying parameters such as namespace, resource limits, and any custom settings.

3. **Dependency Management**: The installer automatically resolves and packages all necessary dependencies, ensuring that the application can function correctly in an airgapped environment.

4. **Installation Process**: The installer deploys the application components to the Kubernetes cluster using Helm charts, ensuring that all services are correctly configured and started.

5. **Post-Installation Validation**: After installation, the installer performs validation checks to confirm that all components are running as expected and that the application is accessible.

6. **Logging and Reporting**: The installer provides detailed logs of the installation process, including any errors or warnings encountered, and generates a summary report for user reference.

## Technologies Used

Build this installer using Go or Python for scripting, Helm for Kubernetes package management, and Terraform for infrastructure provisioning. The installer can be distributed as a standalone binary or script that can be executed on various operating systems.

## User Experience

The installer is designed to be user-friendly, with clear prompts and guidance throughout the installation process. Users can choose to run the installer in interactive mode or provide a configuration file for non-interactive installations. The installer also includes a help menu and detailed documentation to assist users.

## Input Requirements

- Target OS details like Windows, Linux and MacOS
- Cloud provider details like Azure, AWS, GCP or on-prem
- K8S cluster details like version, distribution, networking and storage configuration
- Private registry details like Harbor, Nexus or JFrog Artifactory
- Application configuration details like namespace, resource limits, environment variables, secrets and config maps
- User credentials for cloud provider and k8s cluster
- Any specific requirements or customizations needed for the application deployment
- GitHub repository details for pulling helm charts and terraform scripts
- Database migration scripts using flyway or liquibase
- Monitoring and logging configuration details
- Security configuration details like OIDC 2.0 and OAuth2.0 using Auth0

## Expected Output

- Successful deployment of the application in the specified Kubernetes cluster
- Detailed logs of the installation process
- Summary report of the installation, including any errors or warnings encountered
- Post-installation validation results confirming the application is running as expected
- Rollback mechanism in case of installation failure
- Upgrade mechanism for future application updates
- Housekeeping tasks like log rotation, db backup and restore etc
- Database migration scripts execution results
- Secret management using HashiCorp Vault or Azure Key Vault or AWS Secrets Manager or GCP Secret Manager
- Monitoring and alerting setup status
- CI/CD pipeline setup status
- Caching setup status
- Circuit breaker pattern implementation status
- AI/ML integration status

## Architecture Guidelines

### Core Principles
- **Idempotency**: All operations must be idempotent and resumable
- **Observability**: Comprehensive logging, monitoring, and audit trails
- **Security First**: Zero-trust security model with least privilege access
- **Modularity**: Component-based architecture with clear separation of concerns
- **Offline-First**: Designed for airgapped environments with optional online capabilities

### Technical Stack
- **Language**: Go (primary) or Python for cross-platform CLI tools
- **Infrastructure**: Terraform with provider-specific modules
- **Package Management**: Helm charts with OCI registry support
- **Configuration**: Viper for configuration management, Cobra for CLI
- **Security**: RBAC, image scanning, vulnerability assessment
- **Monitoring**: Prometheus, Grafana, ELK stack integration
- **Secret Management**: Integration with HashiCorp Vault, cloud-native secret managers

### Development Standards
- Follow semantic versioning for releases
- Implement comprehensive testing (unit, integration, e2e)
- Maintain backward compatibility across versions
- Use structured logging with correlation IDs
- Implement graceful error handling and recovery
- Support multiple deployment targets (dev, staging, prod)

### Compliance & Security
- Implement audit logging for all operations
- Support compliance frameworks (SOC2, FedRAMP, HIPAA)
- Image vulnerability scanning and policy enforcement
- Network segmentation and service mesh security
- Encryption at rest and in transit
- Regular security assessments and penetration testing

### Airgapped Environment Considerations
- **Offline Package Management**: Bundle all dependencies and OCI images
- **Local Registry Support**: Integrate with Harbor, Nexus, JFrog Artifactory
- **Certificate Management**: Handle custom CA certificates and TLS
- **Network Isolation**: Support proxy configurations and restricted networks
- **Local Identity Providers**: LDAP, Active Directory integration
- **Offline Documentation**: Include comprehensive offline help and documentation

### Installation Modes
- **Interactive Mode**: Guided CLI with prompts and validation
- **Silent Mode**: Configuration file-driven installation
- **Package Manager**: Distribution via yum, brew, apt-get
- **Dry Run**: Validation without actual deployment
- **Resume**: Continue interrupted installations

### Quality Assurance
- **Pre-flight Checks**: Validate environment and prerequisites
- **Health Checks**: Continuous monitoring of deployed components
- **Integration Tests**: Automated testing of all components
- **Rollback Capability**: Safe recovery from failed installations
- **Upgrade Path**: Seamless version upgrades with data migration

### Operational Excellence
- **Housekeeping**: Automated maintenance tasks (log rotation, cleanup)
- **Backup & Restore**: Database and configuration backup strategies
- **Performance Monitoring**: Resource utilization and performance metrics
- **Alerting**: Proactive notifications for issues and thresholds
- **Documentation**: Comprehensive user guides and troubleshooting

### Multi-Cloud & Hybrid Support
- **Cloud Providers**: AWS, Azure, GCP with unified interface
- **On-Premises**: OpenShift, Rancher, vanilla Kubernetes
- **Hybrid Deployments**: Seamless multi-cloud and hybrid configurations
- **Provider Abstraction**: Unified API across different platforms

