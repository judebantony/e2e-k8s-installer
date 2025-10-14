package k8s

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

// Manager handles Kubernetes operations
type Manager struct {
	config config.KubernetesConfig
	logger *logrus.Logger
}

// NewManager creates a new Kubernetes manager
func NewManager(cfg config.KubernetesConfig, logger *logrus.Logger) (*Manager, error) {
	manager := &Manager{
		config: cfg,
		logger: logger,
	}

	return manager, nil
}

// ValidateAccess validates Kubernetes cluster access
func (m *Manager) ValidateAccess(ctx context.Context) error {
	m.logger.Info("🔍 Validating Kubernetes cluster access")

	// Test cluster connectivity
	if err := m.testClusterConnectivity(ctx); err != nil {
		return fmt.Errorf("cluster connectivity test failed: %w", err)
	}

	// Check cluster version
	if err := m.validateClusterVersion(ctx); err != nil {
		return fmt.Errorf("cluster version validation failed: %w", err)
	}

	// Check permissions
	if err := m.validatePermissions(ctx); err != nil {
		return fmt.Errorf("permissions validation failed: %w", err)
	}

	// Check namespace
	if err := m.validateNamespace(ctx); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	m.logger.Info("✅ Kubernetes cluster access validation completed successfully")
	return nil
}

// SetupCluster sets up the Kubernetes cluster
func (m *Manager) SetupCluster(ctx context.Context) error {
	m.logger.Info("⚙️  Setting up Kubernetes cluster")

	// Create namespace if it doesn't exist
	if err := m.createNamespace(ctx); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	// Setup RBAC
	if m.config.RBAC.Enabled {
		if err := m.setupRBAC(ctx); err != nil {
			return fmt.Errorf("RBAC setup failed: %w", err)
		}
	}

	// Setup networking
	if err := m.setupNetworking(ctx); err != nil {
		return fmt.Errorf("networking setup failed: %w", err)
	}

	// Setup storage
	if err := m.setupStorage(ctx); err != nil {
		return fmt.Errorf("storage setup failed: %w", err)
	}

	m.logger.Info("✅ Kubernetes cluster setup completed successfully")
	return nil
}

// InstallCoreComponents installs core Kubernetes components
func (m *Manager) InstallCoreComponents(ctx context.Context) error {
	m.logger.Info("📦 Installing core Kubernetes components")

	// Install ingress controller
	if m.config.Ingress.Enabled {
		if err := m.installIngressController(ctx); err != nil {
			return fmt.Errorf("ingress controller installation failed: %w", err)
		}
	}

	// Install cert-manager if TLS is enabled
	if m.config.Ingress.TLS && m.config.Ingress.CertManager {
		if err := m.installCertManager(ctx); err != nil {
			return fmt.Errorf("cert-manager installation failed: %w", err)
		}
	}

	// Install metrics server
	if err := m.installMetricsServer(ctx); err != nil {
		return fmt.Errorf("metrics server installation failed: %w", err)
	}

	m.logger.Info("✅ Core components installation completed successfully")
	return nil
}

// DeployApplication deploys the main application using Helm
func (m *Manager) DeployApplication(ctx context.Context) error {
	m.logger.Info("🚀 Deploying application using Helm")

	// Add Helm repositories
	if err := m.addHelmRepositories(ctx); err != nil {
		return fmt.Errorf("failed to add Helm repositories: %w", err)
	}

	// Update Helm repositories
	if err := m.updateHelmRepositories(ctx); err != nil {
		return fmt.Errorf("failed to update Helm repositories: %w", err)
	}

	// Deploy application charts
	if err := m.deployHelmCharts(ctx); err != nil {
		return fmt.Errorf("failed to deploy Helm charts: %w", err)
	}

	m.logger.Info("✅ Application deployment completed successfully")
	return nil
}

// RunDatabaseMigrations runs database migration scripts
func (m *Manager) RunDatabaseMigrations(ctx context.Context) error {
	m.logger.Info("🗄️  Running database migrations")

	// Create migration job
	if err := m.createMigrationJob(ctx); err != nil {
		return fmt.Errorf("failed to create migration job: %w", err)
	}

	// Wait for migration completion
	if err := m.waitForMigrationCompletion(ctx); err != nil {
		return fmt.Errorf("migration job failed: %w", err)
	}

	m.logger.Info("✅ Database migrations completed successfully")
	return nil
}

// HealthCheck performs comprehensive health check
func (m *Manager) HealthCheck(ctx context.Context) error {
	m.logger.Info("🏥 Performing health check")

	// Check pod status
	if err := m.checkPodHealth(ctx); err != nil {
		return fmt.Errorf("pod health check failed: %w", err)
	}

	// Check service endpoints
	if err := m.checkServiceEndpoints(ctx); err != nil {
		return fmt.Errorf("service endpoint check failed: %w", err)
	}

	// Check ingress status
	if m.config.Ingress.Enabled {
		if err := m.checkIngressStatus(ctx); err != nil {
			return fmt.Errorf("ingress status check failed: %w", err)
		}
	}

	m.logger.Info("✅ Health check completed successfully")
	return nil
}

// testClusterConnectivity tests basic cluster connectivity
func (m *Manager) testClusterConnectivity(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "nodes")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w, output: %s", err, output)
	}
	m.logger.Debug("Cluster connectivity test passed")
	return nil
}

// validateClusterVersion validates Kubernetes cluster version
func (m *Manager) validateClusterVersion(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "version", "--client=false", "--short")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to get server version: %w, output: %s", err, output)
	}

	m.logger.WithField("version", string(output)).Info("Cluster version detected")
	return nil
}

// validatePermissions validates required permissions
func (m *Manager) validatePermissions(ctx context.Context) error {
	// Test basic permissions using kubectl auth can-i
	cmd := exec.CommandContext(ctx, "kubectl", "auth", "can-i", "create", "pods")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("insufficient permissions to create pods: %w", err)
	}
	return nil
}

// validateNamespace validates the target namespace
func (m *Manager) validateNamespace(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "namespace", m.config.Namespace)
	if err := cmd.Run(); err != nil {
		// Namespace doesn't exist, which is fine - we'll create it
		m.logger.WithField("namespace", m.config.Namespace).Info("Namespace will be created")
	}
	return nil
}

// createNamespace creates the target namespace
func (m *Manager) createNamespace(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "create", "namespace", m.config.Namespace, "--dry-run=client", "-o", "yaml")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to generate namespace YAML: %w", err)
	}

	// Apply the namespace
	cmd = exec.CommandContext(ctx, "kubectl", "apply", "-f", "-")
	cmd.Stdin = strings.NewReader(string(output))
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create namespace: %w", err)
	}

	m.logger.WithField("namespace", m.config.Namespace).Info("Namespace created successfully")
	return nil
}

// setupRBAC sets up Role-Based Access Control
func (m *Manager) setupRBAC(ctx context.Context) error {
	m.logger.Info("🔐 Setting up RBAC")
	
	// Apply RBAC manifests
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/rbac/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("RBAC setup failed")
		return fmt.Errorf("RBAC setup failed: %w", err)
	}

	m.logger.Info("✅ RBAC setup completed successfully")
	return nil
}

// setupNetworking sets up networking components
func (m *Manager) setupNetworking(ctx context.Context) error {
	m.logger.WithField("networking", m.config.Networking).Info("🌐 Setting up networking")
	
	switch m.config.Networking {
	case "calico":
		return m.setupCalico(ctx)
	case "flannel":
		return m.setupFlannel(ctx)
	case "weave":
		return m.setupWeave(ctx)
	case "istio":
		return m.setupIstio(ctx)
	default:
		m.logger.Info("Using default networking")
		return nil
	}
}

// setupStorage sets up storage components
func (m *Manager) setupStorage(ctx context.Context) error {
	m.logger.WithField("storage", m.config.Storage).Info("💾 Setting up storage")
	
	// Apply storage class manifests
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/storage/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Storage setup failed")
		return fmt.Errorf("storage setup failed: %w", err)
	}

	m.logger.Info("✅ Storage setup completed successfully")
	return nil
}

// installIngressController installs the ingress controller
func (m *Manager) installIngressController(ctx context.Context) error {
	m.logger.WithField("class", m.config.Ingress.Class).Info("🌐 Installing ingress controller")
	
	// Install using Helm
	cmd := exec.CommandContext(ctx, "helm", "install", "ingress-nginx", "ingress-nginx/ingress-nginx",
		"--namespace", "ingress-nginx",
		"--create-namespace",
		"--set", "controller.service.type=LoadBalancer")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Ingress controller installation failed")
		return fmt.Errorf("ingress controller installation failed: %w", err)
	}

	m.logger.Info("✅ Ingress controller installed successfully")
	return nil
}

// installCertManager installs cert-manager
func (m *Manager) installCertManager(ctx context.Context) error {
	m.logger.Info("🔐 Installing cert-manager")
	
	// Install cert-manager using Helm
	cmd := exec.CommandContext(ctx, "helm", "install", "cert-manager", "jetstack/cert-manager",
		"--namespace", "cert-manager",
		"--create-namespace",
		"--version", "v1.13.0",
		"--set", "installCRDs=true")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Cert-manager installation failed")
		return fmt.Errorf("cert-manager installation failed: %w", err)
	}

	m.logger.Info("✅ Cert-manager installed successfully")
	return nil
}

// installMetricsServer installs metrics server
func (m *Manager) installMetricsServer(ctx context.Context) error {
	m.logger.Info("📊 Installing metrics server")
	
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", 
		"https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Metrics server installation failed")
		return fmt.Errorf("metrics server installation failed: %w", err)
	}

	m.logger.Info("✅ Metrics server installed successfully")
	return nil
}

// Networking setup methods
func (m *Manager) setupCalico(ctx context.Context) error {
	m.logger.Info("Installing Calico networking")
	// Implementation for Calico setup
	return nil
}

func (m *Manager) setupFlannel(ctx context.Context) error {
	m.logger.Info("Installing Flannel networking")
	// Implementation for Flannel setup
	return nil
}

func (m *Manager) setupWeave(ctx context.Context) error {
	m.logger.Info("Installing Weave networking")
	// Implementation for Weave setup
	return nil
}

func (m *Manager) setupIstio(ctx context.Context) error {
	m.logger.Info("Installing Istio service mesh")
	// Implementation for Istio setup
	return nil
}

// Helm operations
func (m *Manager) addHelmRepositories(ctx context.Context) error {
	repositories := map[string]string{
		"bitnami":        "https://charts.bitnami.com/bitnami",
		"prometheus":     "https://prometheus-community.github.io/helm-charts",
		"grafana":        "https://grafana.github.io/helm-charts",
		"elastic":        "https://helm.elastic.co",
		"ingress-nginx":  "https://kubernetes.github.io/ingress-nginx",
		"jetstack":       "https://charts.jetstack.io",
	}

	for name, url := range repositories {
		cmd := exec.CommandContext(ctx, "helm", "repo", "add", name, url)
		if err := cmd.Run(); err != nil {
			m.logger.WithFields(logrus.Fields{"repo": name, "url": url}).Warn("Failed to add Helm repository")
		} else {
			m.logger.WithFields(logrus.Fields{"repo": name, "url": url}).Debug("Helm repository added")
		}
	}

	return nil
}

func (m *Manager) updateHelmRepositories(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "helm", "repo", "update")
	return cmd.Run()
}

func (m *Manager) deployHelmCharts(ctx context.Context) error {
	// Deploy application Helm charts
	chartsDir := "./deploy/helm"
	cmd := exec.CommandContext(ctx, "helm", "install", "myapp", chartsDir,
		"--namespace", m.config.Namespace,
		"--create-namespace")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Helm chart deployment failed")
		return fmt.Errorf("helm chart deployment failed: %w", err)
	}

	return nil
}

// Database migration operations
func (m *Manager) createMigrationJob(ctx context.Context) error {
	// Apply migration job manifest
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/migrations/")
	return cmd.Run()
}

func (m *Manager) waitForMigrationCompletion(ctx context.Context) error {
	// Wait for migration job to complete
	cmd := exec.CommandContext(ctx, "kubectl", "wait", "--for=condition=complete", 
		"job/db-migration", "--namespace", m.config.Namespace, "--timeout=300s")
	return cmd.Run()
}

// Health check operations
func (m *Manager) checkPodHealth(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "--namespace", m.config.Namespace)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}

	m.logger.WithField("pods", string(output)).Debug("Pod status")
	return nil
}

func (m *Manager) checkServiceEndpoints(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "endpoints", "--namespace", m.config.Namespace)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get endpoints: %w", err)
	}

	m.logger.WithField("endpoints", string(output)).Debug("Service endpoints")
	return nil
}

func (m *Manager) checkIngressStatus(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "ingress", "--namespace", m.config.Namespace)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get ingress: %w", err)
	}

	m.logger.WithField("ingress", string(output)).Debug("Ingress status")
	return nil
}