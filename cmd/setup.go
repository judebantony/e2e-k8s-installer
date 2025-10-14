package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
	"github.com/judebantony/e2e-k8s-installer/pkg/progress"
)

// setupCmd represents the set-up command
var setupCmd = &cobra.Command{
	Use:   "set-up",
	Short: "Initialize workspace and generate sample configuration",
	Long: `The set-up command initializes a new workspace for the K8s installer and generates
a sample configuration file that can be customized for your specific deployment needs.

This command will:
1. Create the workspace directory structure
2. Generate a sample installer-config.json file
3. Validate prerequisites (kubectl, helm, terraform)
4. Create necessary subdirectories for artifacts and logs

Example:
  k8s-installer set-up --workspace ./my-project
  k8s-installer set-up --config-file custom-config.json`,
	RunE: runSetup,
}

var (
	setupWorkspace  string
	setupConfigFile string
	setupForce      bool
)

func init() {
	rootCmd.AddCommand(setupCmd)
	
	setupCmd.Flags().StringVarP(&setupWorkspace, "workspace", "w", "./workspace", "Workspace directory path")
	setupCmd.Flags().StringVarP(&setupConfigFile, "config-file", "c", "installer-config.json", "Configuration file name")
	setupCmd.Flags().BoolVarP(&setupForce, "force", "f", false, "Force overwrite existing files")
}

func runSetup(cmd *cobra.Command, args []string) error {
	// Initialize progress manager
	progress.InitGlobalProgressManager()
	pm := progress.GetProgressManager()
	
	// Show banner
	progress.ShowBanner("1.0.0")
	
	// Start overall progress area
	pm.StartArea("setup")
	
	steps := []string{
		"Validate prerequisites",
		"Create workspace structure", 
		"Generate configuration file",
		"Initialize logging directories",
		"Setup complete",
	}
	
	currentStep := 0
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 1: Validate prerequisites
	pm.StartSpinner("prereq", "Validating prerequisites...")
	logger.StepStart("validate-prerequisites")
	
	if err := validatePrerequisites(); err != nil {
		pm.FailSpinner("prereq", "Prerequisites validation failed")
		logger.StepFailed("validate-prerequisites", err)
		return fmt.Errorf("prerequisite validation failed: %w", err)
	}
	
	pm.SuccessSpinner("prereq", "Prerequisites validated successfully")
	logger.StepComplete("validate-prerequisites", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 2: Create workspace structure
	pm.StartSpinner("workspace", "Creating workspace structure...")
	logger.StepStart("create-workspace")
	
	if err := createWorkspaceStructure(setupWorkspace); err != nil {
		pm.FailSpinner("workspace", "Workspace creation failed")
		logger.StepFailed("create-workspace", err)
		return fmt.Errorf("workspace creation failed: %w", err)
	}
	
	pm.SuccessSpinner("workspace", "Workspace structure created")
	logger.StepComplete("create-workspace", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 3: Generate configuration file
	configPath := filepath.Join(setupWorkspace, setupConfigFile)
	pm.StartSpinner("config", "Generating configuration file...")
	logger.StepStart("generate-config")
	
	if err := generateConfigFile(configPath); err != nil {
		pm.FailSpinner("config", "Configuration generation failed")
		logger.StepFailed("generate-config", err)
		return fmt.Errorf("configuration generation failed: %w", err)
	}
	
	pm.SuccessSpinner("config", "Configuration file generated")
	logger.StepComplete("generate-config", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 4: Initialize logging directories
	pm.StartSpinner("logging", "Initializing logging directories...")
	logger.StepStart("init-logging")
	
	if err := initializeLoggingDirs(setupWorkspace); err != nil {
		pm.FailSpinner("logging", "Logging directory initialization failed")
		logger.StepFailed("init-logging", err)
		return fmt.Errorf("logging directory initialization failed: %w", err)
	}
	
	pm.SuccessSpinner("logging", "Logging directories initialized")
	logger.StepComplete("init-logging", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Complete setup
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Stop progress area
	pm.StopArea("setup")
	
	// Show success message
	progress.ShowSuccess("🎉 Setup completed successfully!")
	fmt.Printf("\n📁 Workspace: %s\n", setupWorkspace)
	fmt.Printf("⚙️  Config file: %s\n", configPath)
	fmt.Println("\n📝 Next steps:")
	fmt.Println("   1. Edit the configuration file to match your environment")
	fmt.Println("   2. Run 'k8s-installer package-pull' to sync artifacts")
	fmt.Println("   3. Run 'k8s-installer install' to start the installation")
	
	return nil
}

func validatePrerequisites() error {
	logger.Info("Validating prerequisites").Send()
	
	// Check for required tools
	tools := []struct {
		name    string
		command string
		version string
	}{
		{"kubectl", "kubectl", "version --client=true"},
		{"helm", "helm", "version"},
		{"terraform", "terraform", "version"},
	}
	
	for _, tool := range tools {
		logger.Debug("Checking tool").Str("tool", tool.name).Send()
		
		// Check if tool exists in PATH
		if _, err := exec.LookPath(tool.command); err != nil {
			return fmt.Errorf("%s not found in PATH - please install %s", tool.command, tool.name)
		}
		
		logger.Info("Tool found").Str("tool", tool.name).Send()
	}
	
	// Check Go version (for building if needed)
	if _, err := exec.LookPath("go"); err != nil {
		logger.Warn("Go not found in PATH - some features may be limited").Send()
	}
	
	return nil
}

func createWorkspaceStructure(workspace string) error {
	logger.Info("Creating workspace structure").Str("workspace", workspace).Send()
	
	// Create main workspace directory
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}
	
	// Create subdirectories
	subdirs := []string{
		"artifacts/images",
		"artifacts/helm",
		"artifacts/terraform",
		"artifacts/db-scripts",
		"logs",
		"reports",
		"state",
		"scripts",
		"charts",
		"terraform",
		"tests",
	}
	
	for _, subdir := range subdirs {
		path := filepath.Join(workspace, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
		logger.Debug("Created directory").Str("path", path).Send()
	}
	
	// Create .gitignore file
	gitignorePath := filepath.Join(workspace, ".gitignore")
	gitignoreContent := `# Installer generated files
logs/
state/
reports/
artifacts/
*.log

# Sensitive files
**/secrets/
**/*-secret.yaml
**/*-secret.yml

# Terraform files
*.tfstate
*.tfstate.*
.terraform/
.terraform.lock.hcl

# Helm files
charts/*/charts/
charts/*/Chart.lock

# OS files
.DS_Store
Thumbs.db
`
	
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	
	return nil
}

func generateConfigFile(configPath string) error {
	logger.Info("Generating configuration file").Str("path", configPath).Send()
	
	// Check if file exists and force flag
	if _, err := os.Stat(configPath); err == nil && !setupForce {
		return fmt.Errorf("configuration file already exists at %s (use --force to overwrite)", configPath)
	}
	
	// Generate default configuration
	defaultConfig := config.GenerateDefaultConfig()
	
	// Update workspace path in config
	defaultConfig.Installer.Workspace = setupWorkspace
	
	// Convert to JSON
	jsonData, err := defaultConfig.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal configuration: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(configPath, []byte(jsonData), 0644); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}
	
	logger.Info("Configuration file generated").Str("path", configPath).Send()
	return nil
}

func initializeLoggingDirs(workspace string) error {
	logger.Info("Initializing logging directories").Str("workspace", workspace).Send()
	
	// Create timestamped log directory
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logDir := filepath.Join(workspace, "logs", timestamp)
	
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}
	
	// Create symlink to latest
	latestLink := filepath.Join(workspace, "logs", "latest")
	if err := os.Remove(latestLink); err != nil && !os.IsNotExist(err) {
		logger.Warn("Failed to remove existing latest symlink").Err(err).Send()
	}
	
	if err := os.Symlink(timestamp, latestLink); err != nil {
		logger.Warn("Failed to create latest symlink").Err(err).Send()
	}
	// Create README file
	readmePath := filepath.Join(workspace, "README.md")
	readmeContent := `# K8s Installer Workspace

This workspace was created by the K8s installer set-up command.

## Directory Structure

- artifacts/ - Downloaded and synchronized artifacts (images, charts, terraform)
- charts/ - Helm charts for deployment
- logs/ - Installation logs and reports
- reports/ - Test and validation reports
- scripts/ - Custom scripts for deployment and validation
- state/ - Installation state and progress tracking
- terraform/ - Terraform modules for infrastructure
- tests/ - End-to-end tests

## Configuration

Edit installer-config.json to customize the installation for your environment.

## Usage

1. **Package Pull**: k8s-installer package-pull
2. **Infrastructure**: k8s-installer provision-infra 
3. **Database**: k8s-installer db-migrate
4. **Deploy**: k8s-installer deploy
5. **Validate**: k8s-installer post-validate
6. **Test**: k8s-installer e2e-test

Or run all steps: k8s-installer install

## Documentation

For more information, see the main project documentation.
`
	
	if err := os.WriteFile(readmePath, []byte(readmeContent), 0644); err != nil {
		logger.Warn("Failed to create README file").Err(err).Send()
	}
	
	return nil
}