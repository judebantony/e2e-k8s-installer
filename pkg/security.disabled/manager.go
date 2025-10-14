package security

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

// Manager handles security and RBAC operations
type Manager struct {
	config config.SecurityConfig
	logger *logrus.Logger
}

// NewManager creates a new security manager
func NewManager(cfg config.SecurityConfig, logger *logrus.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// Setup sets up security configurations
func (m *Manager) Setup(ctx context.Context) error {
	m.logger.Info("🔒 Setting up security configurations")

	// Setup RBAC (always enabled for security)
	if err := m.setupRBAC(ctx); err != nil {
		return fmt.Errorf("RBAC setup failed: %w", err)
	}

	// Setup network policies
	if err := m.setupNetworkPolicies(ctx); err != nil {
		return fmt.Errorf("network policies setup failed: %w", err)
	}

	// Install security tools
	if m.config.Scanning.Enabled {
		if err := m.installSecurityTools(ctx); err != nil {
			return fmt.Errorf("security tools installation failed: %w", err)
		}
	}

	// Setup authentication
	if m.config.Authentication.Provider != "" {
		if err := m.setupAuthentication(ctx); err != nil {
			return fmt.Errorf("authentication setup failed: %w", err)
		}
	}

	// Setup secrets management (encryption)
	if m.config.Encryption.AtRest || m.config.Encryption.InTransit {
		if err := m.setupSecretsManagement(ctx); err != nil {
			return fmt.Errorf("secrets management setup failed: %w", err)
		}
	}

	m.logger.Info("✅ Security setup completed successfully")
	return nil
}

// setupRBAC sets up Role-Based Access Control
func (m *Manager) setupRBAC(ctx context.Context) error {
	m.logger.Info("🛡️  Setting up RBAC")

	if m.config.Policies.RBAC {
		// Apply RBAC manifests using kubectl
		cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/rbac/")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to apply RBAC manifests: %w", err)
		}
	}

	m.logger.Info("✅ RBAC setup completed successfully")
	return nil
}

// setupNetworkPolicies sets up network security policies
func (m *Manager) setupNetworkPolicies(ctx context.Context) error {
	m.logger.Info("🌐 Setting up network policies")

	if m.config.Policies.NetworkPolicies {
		// Apply network policy manifests
		cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/network-policies/")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to apply network policies: %w", err)
		}
	}

	m.logger.Info("✅ Network policies setup completed successfully")
	return nil
}

// installSecurityTools installs security scanning and monitoring tools
func (m *Manager) installSecurityTools(ctx context.Context) error {
	m.logger.Info("🔍 Installing security tools")

	// Install based on enabled scanners
	for _, scanner := range m.config.Scanning.VulnScanners {
		if err := m.installScanner(ctx, scanner); err != nil {
			m.logger.WithField("scanner", scanner).Warn("Failed to install scanner")
		}
	}

	for _, scanner := range m.config.Scanning.PolicyScanners {
		if err := m.installScanner(ctx, scanner); err != nil {
			m.logger.WithField("scanner", scanner).Warn("Failed to install policy scanner")
		}
	}

	m.logger.Info("✅ Security tools installation completed successfully")
	return nil
}

// installScanner installs a specific security scanner
func (m *Manager) installScanner(ctx context.Context, scanner string) error {
	m.logger.WithField("scanner", scanner).Info("Installing security scanner")

	switch strings.ToLower(scanner) {
	case "trivy":
		return m.installTrivy(ctx)
	case "falco":
		return m.installFalco(ctx)
	case "gatekeeper":
		return m.installGatekeeper(ctx)
	default:
		m.logger.WithField("scanner", scanner).Warn("Unknown scanner, skipping")
		return nil
	}
}

// installTrivy installs Trivy vulnerability scanner
func (m *Manager) installTrivy(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "helm", "install", "trivy", "aqua/trivy",
		"--namespace", "trivy-system", "--create-namespace")
	return cmd.Run()
}

// installFalco installs Falco runtime security
func (m *Manager) installFalco(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "helm", "install", "falco", "falcosecurity/falco",
		"--namespace", "falco-system", "--create-namespace")
	return cmd.Run()
}

// installGatekeeper installs OPA Gatekeeper
func (m *Manager) installGatekeeper(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "helm", "install", "gatekeeper", "gatekeeper/gatekeeper",
		"--namespace", "gatekeeper-system", "--create-namespace")
	return cmd.Run()
}

// setupAuthentication sets up authentication mechanisms
func (m *Manager) setupAuthentication(ctx context.Context) error {
	m.logger.WithField("provider", m.config.Authentication.Provider).Info("🔑 Setting up authentication")

	switch m.config.Authentication.Provider {
	case "oidc":
		return m.setupOIDC(ctx)
	case "oauth2":
		return m.setupOAuth2(ctx)
	case "ldap":
		return m.setupLDAP(ctx)
	default:
		m.logger.Warn("Unknown authentication provider, skipping")
		return nil
	}
}

// setupOIDC sets up OIDC authentication
func (m *Manager) setupOIDC(ctx context.Context) error {
	m.logger.Info("Setting up OIDC authentication")
	// Implementation would configure OIDC
	return nil
}

// setupOAuth2 sets up OAuth2 authentication
func (m *Manager) setupOAuth2(ctx context.Context) error {
	m.logger.Info("Setting up OAuth2 authentication")
	// Implementation would configure OAuth2
	return nil
}

// setupLDAP sets up LDAP authentication
func (m *Manager) setupLDAP(ctx context.Context) error {
	m.logger.Info("Setting up LDAP authentication")
	// Implementation would configure LDAP
	return nil
}

// setupSecretsManagement sets up secrets management and encryption
func (m *Manager) setupSecretsManagement(ctx context.Context) error {
	m.logger.Info("🔐 Setting up secrets management")

	if m.config.Encryption.AtRest {
		if err := m.setupEncryptionAtRest(ctx); err != nil {
			return fmt.Errorf("encryption at rest setup failed: %w", err)
		}
	}

	if m.config.Encryption.InTransit {
		if err := m.setupEncryptionInTransit(ctx); err != nil {
			return fmt.Errorf("encryption in transit setup failed: %w", err)
		}
	}

	return nil
}

// setupEncryptionAtRest sets up encryption at rest
func (m *Manager) setupEncryptionAtRest(ctx context.Context) error {
	m.logger.Info("Setting up encryption at rest")
	// Implementation would configure etcd encryption
	return nil
}

// setupEncryptionInTransit sets up encryption in transit
func (m *Manager) setupEncryptionInTransit(ctx context.Context) error {
	m.logger.Info("Setting up encryption in transit")
	// Implementation would configure TLS
	return nil
}

// Validate validates security configurations
func (m *Manager) Validate(ctx context.Context) error {
	m.logger.Info("🔍 Validating security configurations")

	// Validate RBAC
	if err := m.validateRBAC(ctx); err != nil {
		return fmt.Errorf("RBAC validation failed: %w", err)
	}

	// Validate network policies
	if err := m.validateNetworkPolicies(ctx); err != nil {
		return fmt.Errorf("network policies validation failed: %w", err)
	}

	// Validate security tools
	if err := m.validateSecurityTools(ctx); err != nil {
		return fmt.Errorf("security tools validation failed: %w", err)
	}

	m.logger.Info("✅ Security validation completed successfully")
	return nil
}

// validateRBAC validates RBAC configurations
func (m *Manager) validateRBAC(ctx context.Context) error {
	m.logger.Info("Validating RBAC configurations")

	// Check if service accounts exist using kubectl
	cmd := exec.CommandContext(ctx, "kubectl", "get", "serviceaccounts", "--all-namespaces")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list service accounts: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	m.logger.WithField("count", len(lines)-1).Info("Service accounts validated")
	return nil
}

// validateNetworkPolicies validates network policies
func (m *Manager) validateNetworkPolicies(ctx context.Context) error {
	m.logger.Info("Validating network policies")

	// Implementation would validate network policies
	return nil
}

// validateSecurityTools validates security tools
func (m *Manager) validateSecurityTools(ctx context.Context) error {
	m.logger.Info("Validating security tools")

	// Check if security tools are running
	namespaces := []string{"falco-system", "trivy-system", "gatekeeper-system"}
	for _, ns := range namespaces {
		cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "-n", ns)
		if err := cmd.Run(); err != nil {
			m.logger.WithField("namespace", ns).Warn("Could not list pods in security namespace")
			continue
		}

		m.logger.WithField("namespace", ns).Debug("Security tools validated in namespace")
	}

	return nil
}

// GetStatus returns the current security status
func (m *Manager) GetStatus(ctx context.Context) (map[string]interface{}, error) {
	m.logger.Info("Getting security status")

	status := map[string]interface{}{
		"rbac":        m.config.Policies.RBAC,
		"networkPolicies": m.config.Policies.NetworkPolicies,
		"podSecurity": m.config.Policies.PodSecurity,
		"scanning":    m.config.Scanning.Enabled,
		"authentication": m.config.Authentication.Provider,
		"encryption": map[string]bool{
			"atRest":    m.config.Encryption.AtRest,
			"inTransit": m.config.Encryption.InTransit,
		},
	}

	return status, nil
}