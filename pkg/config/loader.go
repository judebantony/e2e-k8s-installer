package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// LoadConfig loads and validates configuration from a JSON file
func LoadConfig(path string) (*InstallerConfig, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s", path)
	}

	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse JSON
	var config InstallerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse JSON configuration: %w", err)
	}

	// Validate configuration
	if err := config.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Set default values
	config.setDefaults()

	return &config, nil
}

// ValidateConfig validates the configuration structure
func (c *InstallerConfig) ValidateConfig() error {
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Custom validation logic
	if err := c.validateCustomRules(); err != nil {
		return fmt.Errorf("custom validation failed: %w", err)
	}

	return nil
}

// setDefaults sets default values for optional fields
func (c *InstallerConfig) setDefaults() {
	// Set default log level if not specified
	if c.Installer.LogLevel == "" {
		c.Installer.LogLevel = "info"
	}

	// Set default log format if not specified
	if c.Installer.LogFormat == "" {
		c.Installer.LogFormat = "text"
	}

	// Set default workspace if not specified
	if c.Installer.Workspace == "" {
		c.Installer.Workspace = "./workspace"
	}

	// Set default timeouts
	if c.Artifacts.Images.Vendor.Timeout == "" {
		c.Artifacts.Images.Vendor.Timeout = "30s"
	}
	if c.Artifacts.Images.Client.Timeout == "" {
		c.Artifacts.Images.Client.Timeout = "60s"
	}

	// Set default image pull policy
	for i := range c.Artifacts.Images.Images {
		if c.Artifacts.Images.Images[i].PullPolicy == "" {
			c.Artifacts.Images.Images[i].PullPolicy = "IfNotPresent"
		}
	}

	// Set default Terraform settings
	if c.Infrastructure.Terraform.Enabled {
		if c.Infrastructure.Terraform.Parallelism == 0 {
			c.Infrastructure.Terraform.Parallelism = 10
		}
		if c.Infrastructure.Terraform.Timeout == "" {
			c.Infrastructure.Terraform.Timeout = "30m"
		}
	}

	// Set default database settings
	if c.Database.Enabled {
		if c.Database.Connection.Port == 0 {
			c.Database.Connection.Port = 5432 // Default PostgreSQL port
		}
		if c.Database.Connection.SSLMode == "" {
			c.Database.Connection.SSLMode = "require"
		}
		if c.Database.Connection.Timeout == "" {
			c.Database.Connection.Timeout = "30s"
		}
		if c.Database.Validation.Timeout == "" {
			c.Database.Validation.Timeout = "30s"
		}
		if c.Database.Migration.Timeout == "" {
			c.Database.Migration.Timeout = "10m"
		}
	}

	// Set default deployment settings
	if c.Deployment.Helm.Timeout == "" {
		c.Deployment.Helm.Timeout = "10m"
	}
	if c.Deployment.Kubernetes.Timeout == "" {
		c.Deployment.Kubernetes.Timeout = "5m"
	}
	if c.Deployment.Kubernetes.WaitTimeout == "" {
		c.Deployment.Kubernetes.WaitTimeout = "10m"
	}

	// Set default health check settings
	for i := range c.Deployment.Helm.Charts {
		chart := &c.Deployment.Helm.Charts[i]
		if chart.HealthCheck.Method == "" {
			chart.HealthCheck.Method = "GET"
		}
		if chart.HealthCheck.ExpectedStatus == 0 {
			chart.HealthCheck.ExpectedStatus = 200
		}
		if chart.HealthCheck.Timeout == "" {
			chart.HealthCheck.Timeout = "30s"
		}
		if chart.HealthCheck.Retries == 0 {
			chart.HealthCheck.Retries = 3
		}
		if chart.HealthCheck.Interval == "" {
			chart.HealthCheck.Interval = "10s"
		}
	}

	// Set default validation settings
	if c.Validation.Post.Timeout == "" {
		c.Validation.Post.Timeout = "15m"
	}
	if c.Validation.E2E.Timeout == "" {
		c.Validation.E2E.Timeout = "30m"
	}

	// Set default script settings
	for i := range c.Validation.Post.Scripts {
		script := &c.Validation.Post.Scripts[i]
		if script.Shell == "" {
			script.Shell = "bash"
		}
		if script.Timeout == "" {
			script.Timeout = "5m"
		}
	}
}

// validateCustomRules performs custom validation logic
func (c *InstallerConfig) validateCustomRules() error {
	// Validate workspace path
	if !filepath.IsAbs(c.Installer.Workspace) {
		// Convert relative path to absolute
		absPath, err := filepath.Abs(c.Installer.Workspace)
		if err != nil {
			return fmt.Errorf("invalid workspace path: %w", err)
		}
		c.Installer.Workspace = absPath
	}

	// Validate that at least one registry is configured for images
	if len(c.Artifacts.Images.Images) > 0 {
		if c.Artifacts.Images.Vendor.Registry == "" {
			return fmt.Errorf("vendor registry must be configured when images are specified")
		}
	}

	// Validate Terraform modules if infrastructure is enabled
	if c.Infrastructure.Terraform.Enabled {
		if len(c.Infrastructure.Terraform.Modules) == 0 {
			return fmt.Errorf("terraform modules must be specified when infrastructure is enabled")
		}
	}

	// Validate database connection if database is enabled
	if c.Database.Enabled {
		if c.Database.Connection.Host == "" {
			return fmt.Errorf("database host must be specified when database is enabled")
		}
		if c.Database.Connection.Database == "" {
			return fmt.Errorf("database name must be specified when database is enabled")
		}
		if c.Database.Connection.Username == "" {
			return fmt.Errorf("database username must be specified when database is enabled")
		}
	}

	// Validate deployment charts
	if len(c.Deployment.Helm.Charts) > 0 {
		orders := make(map[int]string)
		for _, chart := range c.Deployment.Helm.Charts {
			if existingChart, exists := orders[chart.Order]; exists {
				return fmt.Errorf("duplicate deployment order %d for charts %s and %s", 
					chart.Order, existingChart, chart.Name)
			}
			orders[chart.Order] = chart.Name
		}
	}

	// Validate E2E test configuration
	if c.Validation.E2E.Enabled {
		if c.Validation.E2E.TestSuite == "" {
			return fmt.Errorf("test suite path must be specified when E2E testing is enabled")
		}
	}

	return nil
}

// SaveConfig saves the configuration to a JSON file
func (c *InstallerConfig) SaveConfig(path string) error {
	// Validate before saving
	if err := c.ValidateConfig(); err != nil {
		return fmt.Errorf("cannot save invalid configuration: %w", err)
	}

	// Convert to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to file
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// GetWorkspaceConfig returns workspace-specific configuration
func (c *InstallerConfig) GetWorkspaceConfig() WorkspaceConfig {
	return WorkspaceConfig{
		Root:           c.Installer.Workspace,
		ArtifactsDir:   filepath.Join(c.Installer.Workspace, "artifacts"),
		LogsDir:        filepath.Join(c.Installer.Workspace, "logs"),
		ReportsDir:     filepath.Join(c.Installer.Workspace, "reports"),
		StateDir:       filepath.Join(c.Installer.Workspace, "state"),
		ScriptsDir:     filepath.Join(c.Installer.Workspace, "scripts"),
		ChartsDir:      filepath.Join(c.Installer.Workspace, "charts"),
		TerraformDir:   filepath.Join(c.Installer.Workspace, "terraform"),
		TestsDir:       filepath.Join(c.Installer.Workspace, "tests"),
	}
}

// WorkspaceConfig contains workspace directory paths
type WorkspaceConfig struct {
	Root         string
	ArtifactsDir string
	LogsDir      string
	ReportsDir   string
	StateDir     string
	ScriptsDir   string
	ChartsDir    string
	TerraformDir string
	TestsDir     string
}

// EnsureDirectories creates all workspace directories
func (w *WorkspaceConfig) EnsureDirectories() error {
	dirs := []string{
		w.Root,
		w.ArtifactsDir,
		w.LogsDir,
		w.ReportsDir,
		w.StateDir,
		w.ScriptsDir,
		w.ChartsDir,
		w.TerraformDir,
		w.TestsDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// ValidateWorkspace ensures the workspace is properly set up
func (c *InstallerConfig) ValidateWorkspace() error {
	workspace := c.GetWorkspaceConfig()
	
	// Check if workspace exists
	if _, err := os.Stat(workspace.Root); os.IsNotExist(err) {
		return fmt.Errorf("workspace directory does not exist: %s", workspace.Root)
	}

	// Ensure all required directories exist
	if err := workspace.EnsureDirectories(); err != nil {
		return fmt.Errorf("failed to ensure workspace directories: %w", err)
	}

	return nil
}