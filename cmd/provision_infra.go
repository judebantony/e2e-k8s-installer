package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
	"github.com/judebantony/e2e-k8s-installer/pkg/progress"
	"github.com/judebantony/e2e-k8s-installer/pkg/terraform"
)

// provisionInfraCmd represents the provision-infra command
var provisionInfraCmd = &cobra.Command{
	Use:   "provision-infra",
	Short: "Provision Kubernetes cluster and managed services using Terraform",
	Long: `The provision-infra command deploys infrastructure using Terraform modules.
This includes Kubernetes clusters, managed databases, storage, and networking 
components based on your configuration.

This command will:
1. Initialize Terraform backend
2. Plan infrastructure changes
3. Apply Terraform configuration
4. Run embedded health checks
5. Generate infrastructure report

Example:
  k8s-installer provision-infra --config installer-config.json
  k8s-installer provision-infra --config config.json --plan-only`,
	RunE: runProvisionInfra,
}

var (
	provisionConfigFile string
	provisionPlanOnly   bool
	provisionDestroy    bool
	provisionAutoApprove bool
	provisionVarsFile   string
)

func init() {
	provisionInfraCmd.Flags().StringVarP(&provisionConfigFile, "config", "c", "installer-config.json", "Configuration file path")
	provisionInfraCmd.Flags().BoolVar(&provisionPlanOnly, "plan-only", false, "Only show the Terraform plan without applying")
	provisionInfraCmd.Flags().BoolVar(&provisionDestroy, "destroy", false, "Destroy infrastructure instead of creating")
	provisionInfraCmd.Flags().BoolVar(&provisionAutoApprove, "auto-approve", false, "Skip interactive approval of plan")
	provisionInfraCmd.Flags().StringVar(&provisionVarsFile, "vars-file", "", "Additional Terraform variables file")
}

func runProvisionInfra(cmd *cobra.Command, args []string) error {
	// Initialize progress manager
	progress.InitGlobalProgressManager()
	pm := progress.GetProgressManager()
	
	// Show banner
	progress.ShowBanner("1.0.0")
	
	// Start overall progress area
	pm.StartArea("provision-infra")
	
	steps := []string{
		"Load configuration",
		"Initialize Terraform",
		"Plan infrastructure",
		"Apply infrastructure",
		"Run health checks",
		"Generate report",
		"Complete",
	}
	
	if provisionPlanOnly {
		steps = []string{
			"Load configuration",
			"Initialize Terraform", 
			"Plan infrastructure",
			"Complete",
		}
	}
	
	if provisionDestroy {
		steps = []string{
			"Load configuration",
			"Initialize Terraform",
			"Plan destruction",
			"Destroy infrastructure",
			"Generate report",
			"Complete",
		}
	}
	
	currentStep := 0
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 1: Load configuration
	pm.StartSpinner("config", "Loading configuration...")
	logger.StepStart("load-config")
	
	cfg := config.GenerateDefaultConfig()
	// TODO: Implement actual file loading
	// cfg, err := loadConfigFromFile(provisionConfigFile)
	var err error
	if provisionConfigFile != "" {
		// For now, use default config even if file is specified
		// TODO: Implement file parsing
	}
	
	if err := cfg.ValidateConfig(); err != nil {
		pm.FailSpinner("config", "Configuration validation failed")
		logger.StepFailed("load-config", err)
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	pm.SuccessSpinner("config", "Configuration loaded and validated")
	logger.StepComplete("load-config", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 2: Initialize Terraform
	pm.StartSpinner("init", "Initializing Terraform...")
	logger.StepStart("terraform-init")
	
	tfManager, err := terraform.NewManager(&cfg.Infrastructure)
	if err != nil {
		pm.FailSpinner("init", "Terraform initialization failed")
		logger.StepFailed("terraform-init", err)
		return fmt.Errorf("failed to create Terraform manager: %w", err)
	}
	
	if err := tfManager.Init(); err != nil {
		pm.FailSpinner("init", "Terraform init failed")
		logger.StepFailed("terraform-init", err)
		return fmt.Errorf("terraform init failed: %w", err)
	}
	
	pm.SuccessSpinner("init", "Terraform initialized successfully")
	logger.StepComplete("terraform-init", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 3: Plan infrastructure
	pm.StartSpinner("plan", "Planning infrastructure changes...")
	logger.StepStart("terraform-plan")
	
	planOutput, err := tfManager.Plan(provisionDestroy)
	if err != nil {
		pm.FailSpinner("plan", "Terraform planning failed")
		logger.StepFailed("terraform-plan", err)
		return fmt.Errorf("terraform plan failed: %w", err)
	}
	
	pm.SuccessSpinner("plan", "Infrastructure plan completed")
	logger.StepComplete("terraform-plan", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Show plan output
	fmt.Println("\n📋 Terraform Plan:")
	fmt.Println(planOutput)
	
	// If plan-only, stop here
	if provisionPlanOnly {
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
		pm.StopArea("provision-infra")
		progress.ShowSuccess("🎉 Infrastructure planning completed!")
		return nil
	}
	
	// Ask for approval if not auto-approve
	if !provisionAutoApprove && !viper.GetBool("dry-run") {
		// Simple approval prompt
		approved := true // TODO: Implement actual user prompt
		if !approved {
			fmt.Println("Operation cancelled by user")
			return nil
		}
	}
	
	// Step 4: Apply infrastructure
	action := "Applying infrastructure changes..."
	if provisionDestroy {
		action = "Destroying infrastructure..."
	}
	
	pm.StartSpinner("apply", action)
	logger.StepStart("terraform-apply")
	
	if viper.GetBool("dry-run") {
		pm.SuccessSpinner("apply", "Dry run: Infrastructure changes would be applied")
		logger.StepComplete("terraform-apply", 0)
	} else {
		if err := tfManager.Apply(provisionDestroy); err != nil {
			pm.FailSpinner("apply", "Infrastructure operation failed")
			logger.StepFailed("terraform-apply", err)
			return fmt.Errorf("terraform apply failed: %w", err)
		}
		
		pm.SuccessSpinner("apply", "Infrastructure operation completed")
		logger.StepComplete("terraform-apply", 0)
	}
	
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Skip health checks and report for destroy
	if provisionDestroy {
		// Step 5: Generate report for destroy
		pm.StartSpinner("report", "Generating destruction report...")
		logger.StepStart("generate-report")
		
		reportPath, err := generateInfraReport(cfg, tfManager, true)
		if err != nil {
			pm.FailSpinner("report", "Report generation failed")
			logger.StepFailed("generate-report", err)
			return fmt.Errorf("report generation failed: %w", err)
		}
		
		pm.SuccessSpinner("report", "Destruction report generated")
		logger.StepComplete("generate-report", 0)
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
		
		// Complete
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
		pm.StopArea("provision-infra")
		
		progress.ShowSuccess("🎉 Infrastructure destruction completed!")
		fmt.Printf("📄 Report: %s\n", reportPath)
		return nil
	}
	
	// Step 5: Run health checks
	pm.StartSpinner("health", "Running infrastructure health checks...")
	logger.StepStart("health-checks")
	
	if viper.GetBool("dry-run") {
		pm.SuccessSpinner("health", "Dry run: Health checks would be performed")
		logger.StepComplete("health-checks", 0)
	} else {
		if err := tfManager.RunHealthChecks(); err != nil {
			pm.FailSpinner("health", "Health checks failed")
			logger.StepFailed("health-checks", err)
			logger.Warn("Infrastructure health checks failed, but infrastructure was created").Err(err).Send()
			// Don't return error here as infrastructure was successfully created
		} else {
			pm.SuccessSpinner("health", "Health checks passed")
			logger.StepComplete("health-checks", 0)
		}
	}
	
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Step 6: Generate report
	pm.StartSpinner("report", "Generating infrastructure report...")
	logger.StepStart("generate-report")
	
	reportPath, err := generateInfraReport(cfg, tfManager, false)
	if err != nil {
		pm.FailSpinner("report", "Report generation failed")
		logger.StepFailed("generate-report", err)
		return fmt.Errorf("report generation failed: %w", err)
	}
	
	pm.SuccessSpinner("report", "Infrastructure report generated")
	logger.StepComplete("generate-report", 0)
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Complete
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	pm.StopArea("provision-infra")
	
	// Show success message
	progress.ShowSuccess("🎉 Infrastructure provisioning completed!")
	fmt.Printf("📄 Report: %s\n", reportPath)
	
	// Show next steps
	fmt.Println("\n📝 Next steps:")
	fmt.Println("   1. Run 'k8s-installer db-migrate' to initialize databases")
	fmt.Println("   2. Run 'k8s-installer deploy' to deploy applications")
	
	return nil
}

func generateInfraReport(cfg *config.Config, tfManager *terraform.Manager, isDestroy bool) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	reportDir := filepath.Join(cfg.Installer.Workspace, "reports")
	
	if err := os.MkdirAll(reportDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create reports directory: %w", err)
	}
	
	var reportName string
	if isDestroy {
		reportName = fmt.Sprintf("infra-destroy-report_%s.json", timestamp)
	} else {
		reportName = fmt.Sprintf("infra-report_%s.json", timestamp)
	}
	
	reportPath := filepath.Join(reportDir, reportName)
	
	// Get Terraform outputs
	outputs, err := tfManager.GetOutputs()
	if err != nil {
		logger.Warn("Failed to get Terraform outputs").Err(err).Send()
		outputs = make(map[string]interface{})
	}
	
	// Create report
	report := map[string]interface{}{
		"timestamp": timestamp,
		"operation": map[string]interface{}{
			"type":      "provision-infra",
			"destroy":   isDestroy,
			"workspace": cfg.Installer.Workspace,
			"provider":  cfg.Cloud.Provider,
			"region":    cfg.Cloud.Region,
		},
		"terraform": map[string]interface{}{
			"outputs": outputs,
		},
		"status": "completed",
	}
	
	// Convert to JSON
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %w", err)
	}
	
	// Write report
	if err := os.WriteFile(reportPath, jsonData, 0644); err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}
	
	// Create symlink to latest
	latestPath := filepath.Join(reportDir, "infra-report-latest.json")
	if err := os.Remove(latestPath); err != nil && !os.IsNotExist(err) {
		logger.Warn("Failed to remove existing latest report link").Err(err).Send()
	}
	
	if err := os.Symlink(reportName, latestPath); err != nil {
		logger.Warn("Failed to create latest report link").Err(err).Send()
	}
	
	logger.Info("Infrastructure report generated").
		Str("path", reportPath).
		Bool("destroy", isDestroy).
		Send()
	
	return reportPath, nil
}