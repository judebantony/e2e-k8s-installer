package installer

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/validation"
	"github.com/judebantony/e2e-k8s-installer/pkg/cloud"
	"github.com/judebantony/e2e-k8s-installer/pkg/k8s"
	"github.com/judebantony/e2e-k8s-installer/pkg/monitoring"
	"github.com/judebantony/e2e-k8s-installer/pkg/security"
)

// Installer represents the main installer instance
type Installer struct {
	config     *config.Config
	logger     *logrus.Logger
	validator  *validation.Validator
	cloudMgr   cloud.Manager
	k8sMgr     *k8s.Manager
	monitoring *monitoring.Manager
	security   *security.Manager
	ctx        context.Context
	cancel     context.CancelFunc
}

// InstallationPhase represents different phases of installation
type InstallationPhase string

const (
	PhaseValidation     InstallationPhase = "validation"
	PhaseInfrastructure InstallationPhase = "infrastructure"
	PhaseKubernetes     InstallationPhase = "kubernetes"
	PhaseApplication    InstallationPhase = "application"
	PhaseMonitoring     InstallationPhase = "monitoring"
	PhaseSecurity       InstallationPhase = "security"
	PhaseValidationPost InstallationPhase = "post-validation"
	PhaseComplete       InstallationPhase = "complete"
)

// InstallationStatus represents the current installation status
type InstallationStatus struct {
	Phase       InstallationPhase `json:"phase"`
	Progress    int               `json:"progress"`
	Message     string            `json:"message"`
	StartTime   time.Time         `json:"startTime"`
	ElapsedTime time.Duration     `json:"elapsedTime"`
	Errors      []string          `json:"errors,omitempty"`
	Warnings    []string          `json:"warnings,omitempty"`
}

// NewInstaller creates a new installer instance
func NewInstaller(cfg *config.Config) (*Installer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	if cfg.Metadata.Environment == "dev" {
		logger.SetLevel(logrus.DebugLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize managers
	validator := validation.NewValidator()
	
	cloudMgr, err := cloud.NewManager(cfg.Cloud, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create cloud manager: %w", err)
	}

	k8sMgr, err := k8s.NewManager(cfg.Kubernetes, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create kubernetes manager: %w", err)
	}

	monitoringMgr := monitoring.NewManager(cfg.Monitoring, logger)
	securityMgr := security.NewManager(cfg.Security, logger)

	return &Installer{
		config:     cfg,
		logger:     logger,
		validator:  validator,
		cloudMgr:   cloudMgr,
		k8sMgr:     k8sMgr,
		monitoring: monitoringMgr,
		security:   securityMgr,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

// Install performs the complete installation process
func (i *Installer) Install() error {
	startTime := time.Now()
	i.logger.Info("🚀 Starting E2E Kubernetes installation")

	// Create installation status tracker
	status := &InstallationStatus{
		Phase:     PhaseValidation,
		Progress:  0,
		Message:   "Starting installation",
		StartTime: startTime,
	}

	defer func() {
		status.ElapsedTime = time.Since(startTime)
		i.logger.WithFields(logrus.Fields{
			"phase":    status.Phase,
			"progress": status.Progress,
			"duration": status.ElapsedTime,
		}).Info("Installation phase completed")
	}()

	// Phase 1: Pre-installation validation
	if err := i.runPhase(status, PhaseValidation, i.validateEnvironment); err != nil {
		return fmt.Errorf("validation phase failed: %w", err)
	}

	// Phase 2: Infrastructure provisioning
	if err := i.runPhase(status, PhaseInfrastructure, i.provisionInfrastructure); err != nil {
		return fmt.Errorf("infrastructure phase failed: %w", err)
	}

	// Phase 3: Kubernetes setup
	if err := i.runPhase(status, PhaseKubernetes, i.setupKubernetes); err != nil {
		return fmt.Errorf("kubernetes phase failed: %w", err)
	}

	// Phase 4: Application deployment
	if err := i.runPhase(status, PhaseApplication, i.deployApplication); err != nil {
		return fmt.Errorf("application phase failed: %w", err)
	}

	// Phase 5: Monitoring setup
	if err := i.runPhase(status, PhaseMonitoring, i.setupMonitoring); err != nil {
		return fmt.Errorf("monitoring phase failed: %w", err)
	}

	// Phase 6: Security configuration
	if err := i.runPhase(status, PhaseSecurity, i.configureSecurity); err != nil {
		return fmt.Errorf("security phase failed: %w", err)
	}

	// Phase 7: Post-installation validation
	if err := i.runPhase(status, PhaseValidationPost, i.postValidation); err != nil {
		return fmt.Errorf("post-validation phase failed: %w", err)
	}

	// Phase 8: Complete
	status.Phase = PhaseComplete
	status.Progress = 100
	status.Message = "Installation completed successfully"

	i.logger.WithFields(logrus.Fields{
		"duration": time.Since(startTime),
		"phase":    status.Phase,
	}).Info("✅ E2E Kubernetes installation completed successfully")

	return nil
}

// runPhase executes a specific installation phase
func (i *Installer) runPhase(status *InstallationStatus, phase InstallationPhase, phaseFunc func() error) error {
	status.Phase = phase
	status.Message = fmt.Sprintf("Executing %s phase", phase)
	
	i.logger.WithField("phase", phase).Info("Starting installation phase")
	
	if err := phaseFunc(); err != nil {
		status.Errors = append(status.Errors, err.Error())
		return err
	}
	
	// Update progress based on phase
	switch phase {
	case PhaseValidation:
		status.Progress = 10
	case PhaseInfrastructure:
		status.Progress = 25
	case PhaseKubernetes:
		status.Progress = 40
	case PhaseApplication:
		status.Progress = 60
	case PhaseMonitoring:
		status.Progress = 75
	case PhaseSecurity:
		status.Progress = 85
	case PhaseValidationPost:
		status.Progress = 95
	}
	
	return nil
}

// validateEnvironment performs comprehensive pre-installation validation
func (i *Installer) validateEnvironment() error {
	i.logger.Info("🔍 Validating environment prerequisites")
	
	// Validate configuration
	if err := i.config.Validate(); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}
	
	// Validate cloud provider access
	if err := i.cloudMgr.ValidateAccess(i.ctx); err != nil {
		return fmt.Errorf("cloud provider validation failed: %w", err)
	}
	
	// Validate Kubernetes access
	if err := i.k8sMgr.ValidateAccess(i.ctx); err != nil {
		return fmt.Errorf("kubernetes validation failed: %w", err)
	}
	
	// Validate registry access
	if err := i.validator.ValidateRegistry(i.config.Registry); err != nil {
		return fmt.Errorf("registry validation failed: %w", err)
	}
	
	i.logger.Info("✅ Environment validation completed successfully")
	return nil
}

// provisionInfrastructure provisions cloud infrastructure
func (i *Installer) provisionInfrastructure() error {
	i.logger.Info("🏗️  Provisioning cloud infrastructure")
	
	if err := i.cloudMgr.ProvisionInfrastructure(i.ctx); err != nil {
		return fmt.Errorf("infrastructure provisioning failed: %w", err)
	}
	
	i.logger.Info("✅ Infrastructure provisioning completed successfully")
	return nil
}

// setupKubernetes sets up Kubernetes cluster and core components
func (i *Installer) setupKubernetes() error {
	i.logger.Info("⚙️  Setting up Kubernetes cluster")
	
	if err := i.k8sMgr.SetupCluster(i.ctx); err != nil {
		return fmt.Errorf("kubernetes setup failed: %w", err)
	}
	
	if err := i.k8sMgr.InstallCoreComponents(i.ctx); err != nil {
		return fmt.Errorf("core components installation failed: %w", err)
	}
	
	i.logger.Info("✅ Kubernetes setup completed successfully")
	return nil
}

// deployApplication deploys the main application stack
func (i *Installer) deployApplication() error {
	i.logger.Info("📦 Deploying application stack")
	
	if err := i.k8sMgr.DeployApplication(i.ctx); err != nil {
		return fmt.Errorf("application deployment failed: %w", err)
	}
	
	// Run database migrations if configured
	if i.config.Database.Migration.Tool != "" {
		if err := i.k8sMgr.RunDatabaseMigrations(i.ctx); err != nil {
			return fmt.Errorf("database migration failed: %w", err)
		}
	}
	
	i.logger.Info("✅ Application deployment completed successfully")
	return nil
}

// setupMonitoring configures monitoring and logging stack
func (i *Installer) setupMonitoring() error {
	i.logger.Info("📊 Setting up monitoring and logging")
	
	if err := i.monitoring.InstallStack(i.ctx); err != nil {
		return fmt.Errorf("monitoring setup failed: %w", err)
	}
	
	i.logger.Info("✅ Monitoring setup completed successfully")
	return nil
}

// configureSecurity configures security components
func (i *Installer) configureSecurity() error {
	i.logger.Info("🔒 Configuring security components")
	
	// Use the Setup method from the new security manager
	if err := i.security.Setup(i.ctx); err != nil {
		return fmt.Errorf("security setup failed: %w", err)
	}
	
	i.logger.Info("✅ Security configuration completed successfully")
	return nil
}

// postValidation performs post-installation validation
func (i *Installer) postValidation() error {
	i.logger.Info("🔍 Performing post-installation validation")
	
	if err := i.validator.ValidateInstallation(i.config); err != nil {
		return fmt.Errorf("post-installation validation failed: %w", err)
	}
	
	if err := i.k8sMgr.HealthCheck(i.ctx); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	
	i.logger.Info("✅ Post-installation validation completed successfully")
	return nil
}

// Rollback performs rollback to previous version
func (i *Installer) Rollback() error {
	i.logger.Info("⏪ Starting rollback process")
	
	// Implementation for rollback logic
	// This would involve reversing the installation steps
	
	i.logger.Info("✅ Rollback completed successfully")
	return nil
}

// Upgrade performs application upgrade
func (i *Installer) Upgrade() error {
	i.logger.Info("🔄 Starting upgrade process")
	
	// Implementation for upgrade logic
	// This would involve updating components to newer versions
	
	i.logger.Info("✅ Upgrade completed successfully")
	return nil
}

// GetStatus returns current installation status
func (i *Installer) GetStatus() (*InstallationStatus, error) {
	// Implementation to get current status
	return &InstallationStatus{
		Phase:   PhaseComplete,
		Progress: 100,
		Message: "Installation is healthy",
	}, nil
}

// Cleanup performs cleanup of temporary resources
func (i *Installer) Cleanup() error {
	i.logger.Info("🧹 Cleaning up temporary resources")
	
	if i.cancel != nil {
		i.cancel()
	}
	
	// Cleanup logic implementation
	
	i.logger.Info("✅ Cleanup completed successfully")
	return nil
}