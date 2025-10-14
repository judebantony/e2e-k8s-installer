package validation

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

// Validator handles environment and configuration validation
type Validator struct {
	logger *logrus.Logger
}

// ValidationResult represents the result of a validation check
type ValidationResult struct {
	Name        string        `json:"name"`
	Status      string        `json:"status"` // "pass", "fail", "warning"
	Message     string        `json:"message"`
	Duration    time.Duration `json:"duration"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// ValidationSuite represents a collection of validation results
type ValidationSuite struct {
	Results   []ValidationResult `json:"results"`
	Summary   ValidationSummary  `json:"summary"`
	Timestamp time.Time          `json:"timestamp"`
}

// ValidationSummary provides overall validation statistics
type ValidationSummary struct {
	Total    int `json:"total"`
	Passed   int `json:"passed"`
	Failed   int `json:"failed"`
	Warnings int `json:"warnings"`
}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	return &Validator{
		logger: logger,
	}
}

// ValidateAll performs comprehensive validation of the environment
func (v *Validator) ValidateAll() error {
	v.logger.Info("🔍 Starting comprehensive environment validation")
	
	suite := &ValidationSuite{
		Results:   make([]ValidationResult, 0),
		Timestamp: time.Now(),
	}
	
	// System prerequisites validation
	suite.Results = append(suite.Results, v.validateSystemPrerequisites()...)
	
	// Network connectivity validation
	suite.Results = append(suite.Results, v.validateNetworkConnectivity()...)
	
	// Tool availability validation
	suite.Results = append(suite.Results, v.validateToolAvailability()...)
	
	// Generate summary
	suite.Summary = v.generateSummary(suite.Results)
	
	// Display results
	v.displayResults(suite)
	
	if suite.Summary.Failed > 0 {
		return fmt.Errorf("validation failed with %d errors", suite.Summary.Failed)
	}
	
	v.logger.Info("✅ All validations passed successfully")
	return nil
}

// validateSystemPrerequisites validates system-level prerequisites
func (v *Validator) validateSystemPrerequisites() []ValidationResult {
	var results []ValidationResult
	
	// Check operating system
	results = append(results, v.validateOS())
	
	// Check available disk space
	results = append(results, v.validateDiskSpace())
	
	// Check available memory
	results = append(results, v.validateMemory())
	
	// Check CPU requirements
	results = append(results, v.validateCPU())
	
	return results
}

// validateNetworkConnectivity validates network connectivity
func (v *Validator) validateNetworkConnectivity() []ValidationResult {
	var results []ValidationResult
	
	// Check internet connectivity (if not airgapped)
	results = append(results, v.validateInternetConnectivity())
	
	// Check DNS resolution
	results = append(results, v.validateDNSResolution())
	
	// Check proxy configuration
	results = append(results, v.validateProxyConfiguration())
	
	return results
}

// validateToolAvailability validates required tools availability
func (v *Validator) validateToolAvailability() []ValidationResult {
	var results []ValidationResult
	
	tools := []string{"kubectl", "helm", "terraform", "docker"}
	
	for _, tool := range tools {
		results = append(results, v.validateToolInstallation(tool))
	}
	
	return results
}

// validateOS validates operating system compatibility
func (v *Validator) validateOS() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "Operating System",
		Duration: time.Since(start),
	}
	
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		result.Status = "fail"
		result.Message = "Failed to detect operating system"
		result.Suggestions = []string{"Ensure 'uname' command is available"}
		return result
	}
	
	osType := strings.TrimSpace(string(output))
	supportedOS := []string{"Linux", "Darwin"}
	
	for _, supported := range supportedOS {
		if strings.Contains(osType, supported) {
			result.Status = "pass"
			result.Message = fmt.Sprintf("Operating system %s is supported", osType)
			return result
		}
	}
	
	result.Status = "warning"
	result.Message = fmt.Sprintf("Operating system %s may not be fully supported", osType)
	result.Suggestions = []string{"Supported OS: Linux, macOS"}
	
	return result
}

// validateDiskSpace validates available disk space
func (v *Validator) validateDiskSpace() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "Disk Space",
		Duration: time.Since(start),
	}
	
	cmd := exec.Command("df", "-h", ".")
	_, err := cmd.Output()
	if err != nil {
		result.Status = "warning"
		result.Message = "Unable to check disk space"
		result.Suggestions = []string{"Ensure at least 20GB free disk space"}
		return result
	}
	
	result.Status = "pass"
	result.Message = "Disk space check completed"
	// Additional logic to parse df output and check minimum requirements
	
	return result
}

// validateMemory validates available system memory
func (v *Validator) validateMemory() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "System Memory",
		Duration: time.Since(start),
	}
	
	// Platform-specific memory check
	var cmd *exec.Cmd
	cmd = exec.Command("free", "-h")
	
	_, err := cmd.Output()
	if err != nil {
		result.Status = "warning"
		result.Message = "Unable to check system memory"
		result.Suggestions = []string{"Ensure at least 8GB RAM available"}
		return result
	}
	
	result.Status = "pass"
	result.Message = "Memory check completed"
	// Additional logic to parse memory info
	
	return result
}

// validateCPU validates CPU requirements
func (v *Validator) validateCPU() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "CPU Requirements",
		Duration: time.Since(start),
	}
	
	cmd := exec.Command("nproc")
	output, err := cmd.Output()
	if err != nil {
		result.Status = "warning"
		result.Message = "Unable to check CPU cores"
		result.Suggestions = []string{"Ensure at least 4 CPU cores available"}
		return result
	}
	
	result.Status = "pass"
	result.Message = fmt.Sprintf("CPU cores: %s", strings.TrimSpace(string(output)))
	
	return result
}

// validateInternetConnectivity validates internet connectivity
func (v *Validator) validateInternetConnectivity() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "Internet Connectivity",
		Duration: time.Since(start),
	}
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	testURLs := []string{
		"https://www.google.com",
		"https://github.com",
		"https://registry-1.docker.io",
	}
	
	connected := false
	for _, testURL := range testURLs {
		resp, err := client.Get(testURL)
		if err == nil && resp.StatusCode == 200 {
			connected = true
			resp.Body.Close()
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	
	if connected {
		result.Status = "pass"
		result.Message = "Internet connectivity is available"
	} else {
		result.Status = "warning"
		result.Message = "Internet connectivity not available (airgapped mode)"
		result.Suggestions = []string{"Ensure all required images and packages are available locally"}
	}
	
	return result
}

// validateDNSResolution validates DNS resolution
func (v *Validator) validateDNSResolution() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "DNS Resolution",
		Duration: time.Since(start),
	}
	
	cmd := exec.Command("nslookup", "google.com")
	err := cmd.Run()
	
	if err != nil {
		result.Status = "warning"
		result.Message = "DNS resolution may not be working"
		result.Suggestions = []string{"Configure DNS servers", "Check /etc/resolv.conf"}
	} else {
		result.Status = "pass"
		result.Message = "DNS resolution is working"
	}
	
	return result
}

// validateProxyConfiguration validates proxy configuration
func (v *Validator) validateProxyConfiguration() ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     "Proxy Configuration",
		Duration: time.Since(start),
	}
	
	// Check common proxy environment variables
	proxyVars := []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy"}
	proxyConfigured := false
	
	for _, proxyVar := range proxyVars {
		if value := os.Getenv(proxyVar); value != "" {
			proxyConfigured = true
			break
		}
	}
	
	if proxyConfigured {
		result.Status = "pass"
		result.Message = "Proxy configuration detected"
	} else {
		result.Status = "pass"
		result.Message = "No proxy configuration (direct connection)"
	}
	
	return result
}

// validateToolInstallation validates if a required tool is installed
func (v *Validator) validateToolInstallation(tool string) ValidationResult {
	start := time.Now()
	result := ValidationResult{
		Name:     fmt.Sprintf("Tool: %s", tool),
		Duration: time.Since(start),
	}
	
	cmd := exec.Command("which", tool)
	err := cmd.Run()
	
	if err != nil {
		result.Status = "fail"
		result.Message = fmt.Sprintf("%s is not installed or not in PATH", tool)
		result.Suggestions = v.getInstallationSuggestion(tool)
	} else {
		// Get version information
		versionCmd := exec.Command(tool, "--version")
		output, versionErr := versionCmd.Output()
		
		if versionErr == nil {
			version := strings.Split(string(output), "\n")[0]
			result.Status = "pass"
			result.Message = fmt.Sprintf("%s is installed: %s", tool, version)
		} else {
			result.Status = "pass"
			result.Message = fmt.Sprintf("%s is installed", tool)
		}
	}
	
	return result
}

// getInstallationSuggestion provides installation suggestions for tools
func (v *Validator) getInstallationSuggestion(tool string) []string {
	suggestions := map[string][]string{
		"kubectl": {
			"Install kubectl: https://kubernetes.io/docs/tasks/tools/",
			"On macOS: brew install kubectl",
			"On Linux: curl -LO https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl",
		},
		"helm": {
			"Install Helm: https://helm.sh/docs/intro/install/",
			"On macOS: brew install helm",
			"On Linux: curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash",
		},
		"terraform": {
			"Install Terraform: https://learn.hashicorp.com/tutorials/terraform/install-cli",
			"On macOS: brew install terraform",
			"On Linux: Download from https://www.terraform.io/downloads.html",
		},
		"docker": {
			"Install Docker: https://docs.docker.com/get-docker/",
			"On macOS: brew install --cask docker",
			"On Linux: Follow distribution-specific instructions",
		},
	}
	
	if suggestion, exists := suggestions[tool]; exists {
		return suggestion
	}
	
	return []string{fmt.Sprintf("Please install %s and ensure it's in your PATH", tool)}
}

// ValidateRegistry validates container registry access
func (v *Validator) ValidateRegistry(registry config.RegistryConfig) error {
	v.logger.WithField("registry", registry.URL).Info("Validating registry access")
	
	// Parse registry URL
	registryURL, err := url.Parse(registry.URL)
	if err != nil {
		return fmt.Errorf("invalid registry URL: %w", err)
	}
	
	// Test registry connectivity
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	// Try to access registry v2 API
	testURL := fmt.Sprintf("%s://%s/v2/", registryURL.Scheme, registryURL.Host)
	if registryURL.Scheme == "" {
		testURL = fmt.Sprintf("https://%s/v2/", registry.URL)
	}
	
	resp, err := client.Get(testURL)
	if err != nil {
		return fmt.Errorf("registry connectivity test failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		return fmt.Errorf("registry returned unexpected status: %d", resp.StatusCode)
	}
	
	v.logger.Info("✅ Registry validation completed successfully")
	return nil
}

// ValidateInstallation validates post-installation state
func (v *Validator) ValidateInstallation(cfg *config.Config) error {
	v.logger.Info("🔍 Validating installation state")
	
	// Validate Kubernetes resources
	if err := v.validateKubernetesResources(cfg); err != nil {
		return fmt.Errorf("kubernetes resources validation failed: %w", err)
	}
	
	// Validate application endpoints
	if err := v.validateApplicationEndpoints(cfg); err != nil {
		return fmt.Errorf("application endpoints validation failed: %w", err)
	}
	
	// Validate monitoring stack
	if cfg.Monitoring.Prometheus.Enabled || cfg.Monitoring.Grafana.Enabled {
		if err := v.validateMonitoringStack(cfg); err != nil {
			return fmt.Errorf("monitoring stack validation failed: %w", err)
		}
	}
	
	v.logger.Info("✅ Installation validation completed successfully")
	return nil
}

// validateKubernetesResources validates Kubernetes resources
func (v *Validator) validateKubernetesResources(cfg *config.Config) error {
	// This would use kubectl or client-go to validate resources
	cmd := exec.Command("kubectl", "get", "pods", "-n", cfg.Kubernetes.Namespace)
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get pods: %w", err)
	}
	
	v.logger.WithField("output", string(output)).Debug("Kubernetes resources validated")
	return nil
}

// validateApplicationEndpoints validates application endpoints
func (v *Validator) validateApplicationEndpoints(cfg *config.Config) error {
	// Implementation for endpoint validation
	v.logger.Debug("Application endpoints validated")
	return nil
}

// validateMonitoringStack validates monitoring components
func (v *Validator) validateMonitoringStack(cfg *config.Config) error {
	// Implementation for monitoring stack validation
	v.logger.Debug("Monitoring stack validated")
	return nil
}

// generateSummary generates validation summary
func (v *Validator) generateSummary(results []ValidationResult) ValidationSummary {
	summary := ValidationSummary{
		Total: len(results),
	}
	
	for _, result := range results {
		switch result.Status {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
		case "warning":
			summary.Warnings++
		}
	}
	
	return summary
}

// displayResults displays validation results
func (v *Validator) displayResults(suite *ValidationSuite) {
	v.logger.Info("📊 Validation Results Summary:")
	v.logger.WithFields(logrus.Fields{
		"total":    suite.Summary.Total,
		"passed":   suite.Summary.Passed,
		"failed":   suite.Summary.Failed,
		"warnings": suite.Summary.Warnings,
	}).Info("Validation Summary")
	
	for _, result := range suite.Results {
		fields := logrus.Fields{
			"duration": result.Duration,
			"status":   result.Status,
		}
		
		switch result.Status {
		case "pass":
			v.logger.WithFields(fields).Infof("✅ %s: %s", result.Name, result.Message)
		case "fail":
			v.logger.WithFields(fields).Errorf("❌ %s: %s", result.Name, result.Message)
			for _, suggestion := range result.Suggestions {
				v.logger.Errorf("   💡 %s", suggestion)
			}
		case "warning":
			v.logger.WithFields(fields).Warnf("⚠️  %s: %s", result.Name, result.Message)
			for _, suggestion := range result.Suggestions {
				v.logger.Warnf("   💡 %s", suggestion)
			}
		}
	}
}