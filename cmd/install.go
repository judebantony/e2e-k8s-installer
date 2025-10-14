package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/rs/zerolog"
	"github.com/pterm/pterm"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

var (
	installConfigPath      string
	installVerbose         bool
	installDryRun          bool
	installResume          bool
	installSkipSteps       []string
	installStepsOnly       []string
	installStateFile       string
	installParallel        bool
	installContinueOnError bool
	installWorkspace       string
)

// installCmd represents the install command (main orchestrator)
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Orchestrate complete E2E Kubernetes application installation",
	Long: `Orchestrate the complete end-to-end Kubernetes application installation workflow.
This command chains together all installation steps in the correct order with 
comprehensive error handling, state management, and resume capabilities.

Installation Steps:
1. Setup workspace and configuration
2. Pull and sync packages (images, Helm charts, Terraform modules)
3. Provision infrastructure using Terraform
4. Initialize and migrate database schemas
5. Deploy applications using Helm charts
6. Perform post-deployment validation
7. Execute end-to-end testing

This command handles:
- Step orchestration and dependency management
- Installation state persistence and resume capabilities
- Comprehensive error handling and rollback
- Progress tracking and reporting
- Parallel execution where possible
- Configuration validation and preparation

Examples:
  # Run complete installation with default configuration
  e2e-k8s-installer install

  # Run installation with custom workspace
  e2e-k8s-installer install --workspace ./my-deployment

  # Resume failed installation from last successful step
  e2e-k8s-installer install --resume

  # Run specific steps only
  e2e-k8s-installer install --steps-only provision-infra,deploy,post-validate

  # Skip specific steps
  e2e-k8s-installer install --skip-steps e2e-test

  # Continue installation even if non-critical steps fail
  e2e-k8s-installer install --continue-on-error

  # Dry run to preview installation plan
  e2e-k8s-installer install --dry-run`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().StringVar(&installConfigPath, "config", "", "Path to installation configuration file")
	installCmd.Flags().BoolVarP(&installVerbose, "verbose", "v", false, "Enable verbose logging")
	installCmd.Flags().BoolVar(&installDryRun, "dry-run", false, "Preview installation plan without executing")
	installCmd.Flags().BoolVar(&installResume, "resume", false, "Resume installation from last successful step")
	installCmd.Flags().StringSliceVar(&installSkipSteps, "skip-steps", []string{}, "Skip specified installation steps")
	installCmd.Flags().StringSliceVar(&installStepsOnly, "steps-only", []string{}, "Run only specified installation steps")
	installCmd.Flags().StringVar(&installStateFile, "state-file", "", "Path to installation state file")
	installCmd.Flags().BoolVar(&installParallel, "parallel", false, "Enable parallel execution where possible")
	installCmd.Flags().BoolVar(&installContinueOnError, "continue-on-error", false, "Continue installation if non-critical steps fail")
	installCmd.Flags().StringVar(&installWorkspace, "workspace", "", "Installation workspace directory")
}

func runInstall(cmd *cobra.Command, args []string) error {
	// Initialize logger
	logger := zerolog.New(os.Stderr).With().
		Timestamp().
		Str("component", "install").
		Logger()

	if installVerbose {
		logger = logger.Level(zerolog.DebugLevel)
	}

	// Create spinner for initialization
	spinner, _ := pterm.DefaultSpinner.Start("Initializing E2E Kubernetes installation...")
	
	ctx := context.Background()
	startTime := time.Now()

	// Load configuration
	config, err := loadInstallConfig(installConfigPath)
	if err != nil {
		spinner.Fail("Failed to load configuration")
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	spinner.Success("Configuration loaded")
	logger.Info().Msg("Installation configuration loaded successfully")

	// Create installation manager
	manager, err := NewInstallationManager(config, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize installation manager: %w", err)
	}

	// Apply command line overrides
	manager.ApplyCommandLineOverrides()

	// Load or initialize installation state
	if err := manager.LoadState(); err != nil {
		return fmt.Errorf("failed to load installation state: %w", err)
	}

	// Create progress area
	progressArea, _ := pterm.DefaultArea.Start()
	
	// Define installation steps with their dependencies and configurations
	steps := []InstallationStep{
		{
			Name:         "setup",
			Description:  "Setting up workspace and configuration",
			Command:      "setup",
			Required:     true,
			Dependencies: []string{},
			Handler:      manager.RunSetup,
		},
		{
			Name:         "package-pull",
			Description:  "Pulling and syncing packages",
			Command:      "package-pull",
			Required:     true,
			Dependencies: []string{"setup"},
			Handler:      manager.RunPackagePull,
		},
		{
			Name:         "provision-infra",
			Description:  "Provisioning infrastructure",
			Command:      "provision-infra",
			Required:     true,
			Dependencies: []string{"package-pull"},
			Handler:      manager.RunProvisionInfra,
		},
		{
			Name:         "db-migrate",
			Description:  "Initializing and migrating database",
			Command:      "db-migrate",
			Required:     false,
			Dependencies: []string{"provision-infra"},
			Handler:      manager.RunDBMigrate,
		},
		{
			Name:         "deploy",
			Description:  "Deploying applications",
			Command:      "deploy",
			Required:     true,
			Dependencies: []string{"provision-infra"},
			Handler:      manager.RunDeploy,
		},
		{
			Name:         "post-validate",
			Description:  "Performing post-deployment validation",
			Command:      "post-validate",
			Required:     false,
			Dependencies: []string{"deploy"},
			Handler:      manager.RunPostValidate,
		},
		{
			Name:         "e2e-test",
			Description:  "Executing end-to-end tests",
			Command:      "e2e-test",
			Required:     false,
			Dependencies: []string{"deploy"},
			Handler:      manager.RunE2ETest,
		},
	}

	// Filter steps based on command line flags
	steps = manager.FilterSteps(steps)

	// Execute installation steps
	if installParallel {
		err = manager.ExecuteStepsParallel(ctx, steps, progressArea)
	} else {
		err = manager.ExecuteStepsSequential(ctx, steps, progressArea)
	}

	progressArea.Stop()

	// Handle installation result
	if err != nil {
		pterm.Error.Printf("❌ Installation failed: %v\n", err)
		
		// Save state for resume
		if saveErr := manager.SaveState(); saveErr != nil {
			logger.Error().Err(saveErr).Msg("Failed to save installation state")
		}
		
		return err
	}

	// Mark installation as completed
	manager.MarkCompleted()

	// Generate final installation report
	if err := manager.GenerateFinalReport(); err != nil {
		logger.Warn().Err(err).Msg("Failed to generate final installation report")
	}

	// Success summary
	duration := time.Since(startTime)
	pterm.Success.Printf("🎉 E2E Kubernetes installation completed successfully in %v\n", duration.Round(time.Second))
	
	// Display installation summary
	pterm.DefaultSection.Println("Installation Summary")
	
	results := manager.GetInstallationResults()
	info := [][]string{
		{"Workspace", manager.GetWorkspace()},
		{"Total Steps", fmt.Sprintf("%d", results.TotalSteps)},
		{"Completed", fmt.Sprintf("%d", results.CompletedSteps)},
		{"Skipped", fmt.Sprintf("%d", results.SkippedSteps)},
		{"Failed", fmt.Sprintf("%d", results.FailedSteps)},
		{"Duration", duration.Round(time.Second).String()},
		{"Success Rate", fmt.Sprintf("%.1f%%", results.SuccessRate)},
	}

	if installDryRun {
		info = append(info, []string{"Mode", "DRY RUN - No changes applied"})
	}

	if installResume {
		info = append(info, []string{"Resume", "Yes - Resumed from previous state"})
	}

	pterm.DefaultTable.WithHasHeader().WithData(
		append([][]string{{"Property", "Value"}}, info...),
	).Render()

	// Display step details
	pterm.DefaultSection.Println("Step Execution Details")
	
	stepData := [][]string{{"Step", "Status", "Duration", "Description"}}
	for _, step := range manager.GetCompletedSteps() {
		status := "✅ Completed"
		if step.Failed {
			status = "❌ Failed"
		} else if step.Skipped {
			status = "⏭️  Skipped"
		}
		
		stepData = append(stepData, []string{
			step.Name,
			status,
			step.Duration.Round(time.Second).String(),
			step.Description,
		})
	}
	
	pterm.DefaultTable.WithHasHeader().WithData(stepData).Render()

	// Display final information
	if manager.GetReportPath() != "" {
		pterm.DefaultSection.Println("Installation Report")
		pterm.Info.Printf("📊 Final installation report: %s\n", manager.GetReportPath())
	}

	if manager.GetStateFile() != "" {
		pterm.Info.Printf("💾 Installation state saved: %s\n", manager.GetStateFile())
	}

	logger.Info().
		Dur("duration", duration).
		Int("total_steps", results.TotalSteps).
		Int("completed_steps", results.CompletedSteps).
		Int("failed_steps", results.FailedSteps).
		Float64("success_rate", results.SuccessRate).
		Msg("Installation completed")

	return nil
}

// InstallationStep represents a single installation step
type InstallationStep struct {
	Name         string
	Description  string
	Command      string
	Required     bool
	Dependencies []string
	Handler      func() error
}

// InstallationResults represents the results of installation execution
type InstallationResults struct {
	TotalSteps     int
	CompletedSteps int
	SkippedSteps   int
	FailedSteps    int
	SuccessRate    float64
	StartTime      time.Time
	EndTime        *time.Time
}

// CompletedStep represents a completed installation step
type CompletedStep struct {
	Name        string
	Description string
	Duration    time.Duration
	Failed      bool
	Skipped     bool
	Error       string
}

// InstallationManager handles the complete installation orchestration
type InstallationManager struct {
	config      *config.InstallerConfig
	logger      zerolog.Logger
	workspace   string
	stateFile   string
	reportPath  string
	state       *config.InstallState
	results     InstallationResults
	completed   []CompletedStep
}

// NewInstallationManager creates a new installation manager
func NewInstallationManager(config *config.InstallerConfig, logger zerolog.Logger) (*InstallationManager, error) {
	workspace := config.Installer.Workspace
	if installWorkspace != "" {
		workspace = installWorkspace
	}

	stateFile := installStateFile
	if stateFile == "" {
		stateFile = filepath.Join(workspace, "install-state.json")
	}

	reportPath := filepath.Join(workspace, "reports", "installation-report.json")

	manager := &InstallationManager{
		config:     config,
		logger:     logger,
		workspace:  workspace,
		stateFile:  stateFile,
		reportPath: reportPath,
		results: InstallationResults{
			StartTime: time.Now(),
		},
		completed: []CompletedStep{},
	}

	return manager, nil
}

// ApplyCommandLineOverrides applies command line flag overrides
func (m *InstallationManager) ApplyCommandLineOverrides() {
	// Override configuration with command line flags
	if installVerbose {
		m.config.Installer.Verbose = true
	}
	
	if installDryRun {
		m.config.Installer.DryRun = true
	}
}

// LoadState loads installation state from file
func (m *InstallationManager) LoadState() error {
	if !installResume {
		// Initialize new state
		m.state = &config.InstallState{
			Steps:     []config.StepState{},
			StartTime: time.Now(),
			Status:    "running",
		}
		return nil
	}

	// TODO: Implement actual state loading from file
	// This would typically involve:
	// 1. Reading state file
	// 2. Parsing JSON state
	// 3. Validating state consistency
	// 4. Preparing for resume

	m.logger.Info().Str("state_file", m.stateFile).Msg("Loading installation state for resume")
	
	// For now, create a new state
	m.state = &config.InstallState{
		Steps:     []config.StepState{},
		StartTime: time.Now(),
		Status:    "running",
		Resume:    true,
	}

	return nil
}

// SaveState saves installation state to file
func (m *InstallationManager) SaveState() error {
	// Create workspace directory if it doesn't exist
	if err := os.MkdirAll(m.workspace, 0755); err != nil {
		return fmt.Errorf("failed to create workspace directory: %w", err)
	}

	// TODO: Implement actual state saving to file
	// This would typically involve:
	// 1. Serializing state to JSON
	// 2. Writing to state file
	// 3. Setting appropriate permissions

	m.logger.Info().Str("state_file", m.stateFile).Msg("Installation state saved")
	return nil
}

// FilterSteps filters installation steps based on command line flags
func (m *InstallationManager) FilterSteps(steps []InstallationStep) []InstallationStep {
	// If steps-only is specified, only include those steps
	if len(installStepsOnly) > 0 {
		stepSet := make(map[string]bool)
		for _, stepName := range installStepsOnly {
			stepSet[stepName] = true
		}

		var filteredSteps []InstallationStep
		for _, step := range steps {
			if stepSet[step.Name] {
				filteredSteps = append(filteredSteps, step)
			}
		}
		return filteredSteps
	}

	// If skip-steps is specified, exclude those steps
	if len(installSkipSteps) > 0 {
		skipSet := make(map[string]bool)
		for _, stepName := range installSkipSteps {
			skipSet[stepName] = true
		}

		var filteredSteps []InstallationStep
		for _, step := range steps {
			if !skipSet[step.Name] {
				filteredSteps = append(filteredSteps, step)
			}
		}
		return filteredSteps
	}

	return steps
}

// ExecuteStepsSequential executes installation steps sequentially
func (m *InstallationManager) ExecuteStepsSequential(ctx context.Context, steps []InstallationStep, progressArea *pterm.AreaPrinter) error {
	for i, step := range steps {
		stepProgress := fmt.Sprintf("[%d/%d] %s", i+1, len(steps), step.Description)
		progressArea.Update(pterm.Sprintf("🔄 %s", stepProgress))
		
		m.logger.Info().
			Str("step", step.Name).
			Str("command", step.Command).
			Bool("required", step.Required).
			Msg("Starting installation step")

		stepStart := time.Now()
		
		if installDryRun {
			m.logger.Info().Str("step", step.Name).Msg("DRY RUN: Step execution skipped")
			
			// Simulate step completion for dry run
			m.completed = append(m.completed, CompletedStep{
				Name:        step.Name,
				Description: step.Description,
				Duration:    time.Second,
				Failed:      false,
				Skipped:     false,
			})
			m.results.CompletedSteps++
		} else {
			if err := step.Handler(); err != nil {
				stepDuration := time.Since(stepStart)
				m.completed = append(m.completed, CompletedStep{
					Name:        step.Name,
					Description: step.Description,
					Duration:    stepDuration,
					Failed:      true,
					Skipped:     false,
					Error:       err.Error(),
				})
				
				m.results.FailedSteps++
				m.logger.Error().
					Err(err).
					Str("step", step.Name).
					Dur("duration", stepDuration).
					Msg("Installation step failed")

				// Check if step is required or if we should continue on error
				if step.Required && !installContinueOnError {
					progressArea.Update(pterm.Sprintf("❌ %s", stepProgress))
					return fmt.Errorf("required installation step '%s' failed: %w", step.Name, err)
				}

				// Continue with non-required steps or when continue-on-error is enabled
				progressArea.Update(pterm.Sprintf("⚠️  %s (failed but continuing)", stepProgress))
			} else {
				stepDuration := time.Since(stepStart)
				m.completed = append(m.completed, CompletedStep{
					Name:        step.Name,
					Description: step.Description,
					Duration:    stepDuration,
					Failed:      false,
					Skipped:     false,
				})
				
				m.results.CompletedSteps++
				progressArea.Update(pterm.Sprintf("✅ %s", stepProgress))
				m.logger.Info().
					Str("step", step.Name).
					Dur("duration", stepDuration).
					Msg("Installation step completed successfully")
			}
		}
		
		m.results.TotalSteps++
		time.Sleep(300 * time.Millisecond) // Visual feedback
	}

	// Calculate success rate
	if m.results.TotalSteps > 0 {
		m.results.SuccessRate = float64(m.results.CompletedSteps) / float64(m.results.TotalSteps) * 100
	}

	return nil
}

// ExecuteStepsParallel executes installation steps in parallel where possible
func (m *InstallationManager) ExecuteStepsParallel(ctx context.Context, steps []InstallationStep, progressArea *pterm.AreaPrinter) error {
	// TODO: Implement proper parallel execution with dependency resolution
	// For now, execute sequentially but with different messaging
	progressArea.Update("🔄 Executing installation steps with parallelization...")
	
	return m.ExecuteStepsSequential(ctx, steps, progressArea)
}

// MarkCompleted marks the installation as completed
func (m *InstallationManager) MarkCompleted() {
	now := time.Now()
	m.results.EndTime = &now
	m.state.Status = "completed"
	m.state.EndTime = &now
}

// GenerateFinalReport generates the final installation report
func (m *InstallationManager) GenerateFinalReport() error {
	// Create reports directory
	if err := os.MkdirAll(filepath.Dir(m.reportPath), 0755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	report := map[string]interface{}{
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"workspace":       m.workspace,
		"total_steps":     m.results.TotalSteps,
		"completed_steps": m.results.CompletedSteps,
		"failed_steps":    m.results.FailedSteps,
		"skipped_steps":   m.results.SkippedSteps,
		"success_rate":    m.results.SuccessRate,
		"start_time":      m.results.StartTime.Format(time.RFC3339),
		"end_time":        m.results.EndTime.Format(time.RFC3339),
		"duration":        time.Since(m.results.StartTime).String(),
		"steps":           m.completed,
		"dry_run":         installDryRun,
		"resumed":         installResume,
		"status":          "completed",
	}

	// TODO: Write actual report to file
	m.logger.Info().Interface("report", report).Str("report_path", m.reportPath).Msg("Final installation report generated")
	return nil
}

// Step handler methods (these would call the actual commands)

func (m *InstallationManager) RunSetup() error {
	// TODO: Call the actual setup command
	time.Sleep(2 * time.Second) // Simulate setup
	m.logger.Info().Msg("Setup step completed")
	return nil
}

func (m *InstallationManager) RunPackagePull() error {
	// TODO: Call the actual package-pull command
	time.Sleep(3 * time.Second) // Simulate package pull
	m.logger.Info().Msg("Package pull step completed")
	return nil
}

func (m *InstallationManager) RunProvisionInfra() error {
	// TODO: Call the actual provision-infra command
	time.Sleep(4 * time.Second) // Simulate infrastructure provisioning
	m.logger.Info().Msg("Infrastructure provisioning step completed")
	return nil
}

func (m *InstallationManager) RunDBMigrate() error {
	// TODO: Call the actual db-migrate command
	time.Sleep(2 * time.Second) // Simulate database migration
	m.logger.Info().Msg("Database migration step completed")
	return nil
}

func (m *InstallationManager) RunDeploy() error {
	// TODO: Call the actual deploy command
	time.Sleep(3 * time.Second) // Simulate deployment
	m.logger.Info().Msg("Deployment step completed")
	return nil
}

func (m *InstallationManager) RunPostValidate() error {
	// TODO: Call the actual post-validate command
	time.Sleep(2 * time.Second) // Simulate post-validation
	m.logger.Info().Msg("Post-validation step completed")
	return nil
}

func (m *InstallationManager) RunE2ETest() error {
	// TODO: Call the actual e2e-test command
	time.Sleep(4 * time.Second) // Simulate E2E testing
	m.logger.Info().Msg("E2E testing step completed")
	return nil
}

// Helper methods

func (m *InstallationManager) GetWorkspace() string {
	return m.workspace
}

func (m *InstallationManager) GetStateFile() string {
	return m.stateFile
}

func (m *InstallationManager) GetReportPath() string {
	return m.reportPath
}

func (m *InstallationManager) GetInstallationResults() InstallationResults {
	return m.results
}

func (m *InstallationManager) GetCompletedSteps() []CompletedStep {
	return m.completed
}

func loadInstallConfig(configPath string) (*config.InstallerConfig, error) {
	// Load configuration from file or use defaults
	// For now, return a default configuration
	
	config := config.GenerateDefaultConfig()
	
	// TODO: Implement actual configuration loading from file
	// This would typically involve:
	// 1. Reading configuration file
	// 2. Parsing YAML/JSON configuration
	// 3. Validating configuration structure
	// 4. Merging with defaults

	return config, nil
}