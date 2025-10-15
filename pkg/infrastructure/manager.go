package infrastructure

import (
	"fmt"

	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
	"github.com/judebantony/e2e-k8s-installer/pkg/makefile"
	"github.com/judebantony/e2e-k8s-installer/pkg/terraform"
)

// Manager handles infrastructure provisioning operations
type Manager struct {
	config        *config.InfrastructureConfig
	terraformMgr  *terraform.Manager
	makefileMgr   *makefile.Manager
	provisionMode string
}

// ProvisionMode constants
const (
	ProvisionModeTerraform = "terraform"
	ProvisionModeMakefile  = "makefile"
	ProvisionModeHybrid    = "hybrid"
)

// NewManager creates a new infrastructure manager
func NewManager(infraConfig *config.InfrastructureConfig) (*Manager, error) {
	if infraConfig == nil {
		return nil, fmt.Errorf("infrastructure configuration is required")
	}

	mgr := &Manager{
		config:        infraConfig,
		provisionMode: infraConfig.ProvisionMode,
	}

	// Default to terraform mode if not specified
	if mgr.provisionMode == "" {
		mgr.provisionMode = ProvisionModeTerraform
	}

	// Initialize managers based on provision mode
	switch mgr.provisionMode {
	case ProvisionModeTerraform:
		if !infraConfig.Terraform.Enabled {
			return nil, fmt.Errorf("terraform mode selected but terraform is not enabled in configuration")
		}
		tfMgr, err := terraform.NewManager(infraConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create terraform manager: %w", err)
		}
		mgr.terraformMgr = tfMgr

	case ProvisionModeMakefile:
		if !infraConfig.Makefile.Enabled {
			return nil, fmt.Errorf("makefile mode selected but makefile is not enabled in configuration")
		}
		makeMgr, err := makefile.NewManager(&infraConfig.Makefile)
		if err != nil {
			return nil, fmt.Errorf("failed to create makefile manager: %w", err)
		}
		mgr.makefileMgr = makeMgr

	case ProvisionModeHybrid:
		// Initialize both managers for hybrid mode
		if infraConfig.Terraform.Enabled {
			tfMgr, err := terraform.NewManager(infraConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to create terraform manager: %w", err)
			}
			mgr.terraformMgr = tfMgr
		}

		if infraConfig.Makefile.Enabled {
			makeMgr, err := makefile.NewManager(&infraConfig.Makefile)
			if err != nil {
				return nil, fmt.Errorf("failed to create makefile manager: %w", err)
			}
			mgr.makefileMgr = makeMgr
		}

		if mgr.terraformMgr == nil && mgr.makefileMgr == nil {
			return nil, fmt.Errorf("hybrid mode requires at least one of terraform or makefile to be enabled")
		}

	default:
		return nil, fmt.Errorf("unsupported provision mode: %s", mgr.provisionMode)
	}

	return mgr, nil
}

// Init initializes the infrastructure provisioning environment
func (m *Manager) Init(dryRun bool) error {
	logger.Info("Initializing infrastructure provisioning").
		Str("mode", m.provisionMode).
		Bool("dryRun", dryRun).
		Send()

	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.initTerraform(dryRun)
	case ProvisionModeMakefile:
		return m.initMakefile(dryRun)
	case ProvisionModeHybrid:
		return m.initHybrid(dryRun)
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// Plan creates an execution plan for infrastructure changes
func (m *Manager) Plan(dryRun bool) error {
	logger.Info("Creating infrastructure plan").
		Str("mode", m.provisionMode).
		Bool("dryRun", dryRun).
		Send()

	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.planTerraform(dryRun)
	case ProvisionModeMakefile:
		return m.planMakefile(dryRun)
	case ProvisionModeHybrid:
		return m.planHybrid(dryRun)
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// Apply applies the infrastructure changes
func (m *Manager) Apply(dryRun bool) error {
	logger.Info("Applying infrastructure changes").
		Str("mode", m.provisionMode).
		Bool("dryRun", dryRun).
		Send()

	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.applyTerraform(dryRun)
	case ProvisionModeMakefile:
		return m.applyMakefile(dryRun)
	case ProvisionModeHybrid:
		return m.applyHybrid(dryRun)
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// Destroy destroys the infrastructure
func (m *Manager) Destroy(dryRun bool) error {
	logger.Info("Destroying infrastructure").
		Str("mode", m.provisionMode).
		Bool("dryRun", dryRun).
		Send()

	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.destroyTerraform(dryRun)
	case ProvisionModeMakefile:
		return m.destroyMakefile(dryRun)
	case ProvisionModeHybrid:
		return m.destroyHybrid(dryRun)
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// Validate validates the infrastructure configuration
func (m *Manager) Validate(dryRun bool) error {
	logger.Info("Validating infrastructure configuration").
		Str("mode", m.provisionMode).
		Bool("dryRun", dryRun).
		Send()

	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.validateTerraform(dryRun)
	case ProvisionModeMakefile:
		return m.validateMakefile(dryRun)
	case ProvisionModeHybrid:
		return m.validateHybrid(dryRun)
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// Terraform-specific methods
func (m *Manager) initTerraform(dryRun bool) error {
	if m.terraformMgr == nil {
		return fmt.Errorf("terraform manager not initialized")
	}
	if dryRun {
		logger.Info("DRY RUN: Terraform initialization skipped").Send()
		return nil
	}
	return m.terraformMgr.Init()
}

func (m *Manager) planTerraform(dryRun bool) error {
	if m.terraformMgr == nil {
		return fmt.Errorf("terraform manager not initialized")
	}
	if dryRun {
		logger.Info("DRY RUN: Terraform plan skipped").Send()
		return nil
	}
	_, err := m.terraformMgr.Plan(false)
	return err
}

func (m *Manager) applyTerraform(dryRun bool) error {
	if m.terraformMgr == nil {
		return fmt.Errorf("terraform manager not initialized")
	}
	if dryRun {
		logger.Info("DRY RUN: Terraform apply skipped").Send()
		return nil
	}
	return m.terraformMgr.Apply(false)
}

func (m *Manager) destroyTerraform(dryRun bool) error {
	if m.terraformMgr == nil {
		return fmt.Errorf("terraform manager not initialized")
	}
	if dryRun {
		logger.Info("DRY RUN: Terraform destroy skipped").Send()
		return nil
	}
	return m.terraformMgr.Apply(true) // destroy=true
}

func (m *Manager) validateTerraform(dryRun bool) error {
	if m.terraformMgr == nil {
		return fmt.Errorf("terraform manager not initialized")
	}
	if dryRun {
		logger.Info("DRY RUN: Terraform validation skipped").Send()
		return nil
	}
	// Use Plan with validate-only semantics
	_, err := m.terraformMgr.Plan(false)
	return err
}

// Makefile-specific methods
func (m *Manager) initMakefile(dryRun bool) error {
	if m.makefileMgr == nil {
		return fmt.Errorf("makefile manager not initialized")
	}
	return m.makefileMgr.Init(dryRun)
}

func (m *Manager) planMakefile(dryRun bool) error {
	if m.makefileMgr == nil {
		return fmt.Errorf("makefile manager not initialized")
	}
	return m.makefileMgr.Plan(dryRun)
}

func (m *Manager) applyMakefile(dryRun bool) error {
	if m.makefileMgr == nil {
		return fmt.Errorf("makefile manager not initialized")
	}
	return m.makefileMgr.Apply(dryRun)
}

func (m *Manager) destroyMakefile(dryRun bool) error {
	if m.makefileMgr == nil {
		return fmt.Errorf("makefile manager not initialized")
	}
	return m.makefileMgr.Destroy(dryRun)
}

func (m *Manager) validateMakefile(dryRun bool) error {
	if m.makefileMgr == nil {
		return fmt.Errorf("makefile manager not initialized")
	}
	return m.makefileMgr.Validate(dryRun)
}

// Hybrid mode methods (execute both terraform and makefile)
func (m *Manager) initHybrid(dryRun bool) error {
	// Execute makefile init first, then terraform
	if m.makefileMgr != nil {
		if err := m.initMakefile(dryRun); err != nil {
			return fmt.Errorf("makefile init failed: %w", err)
		}
	}

	if m.terraformMgr != nil {
		if err := m.initTerraform(dryRun); err != nil {
			return fmt.Errorf("terraform init failed: %w", err)
		}
	}

	return nil
}

func (m *Manager) planHybrid(dryRun bool) error {
	// Execute makefile plan first, then terraform
	if m.makefileMgr != nil {
		if err := m.planMakefile(dryRun); err != nil {
			return fmt.Errorf("makefile plan failed: %w", err)
		}
	}

	if m.terraformMgr != nil {
		if err := m.planTerraform(dryRun); err != nil {
			return fmt.Errorf("terraform plan failed: %w", err)
		}
	}

	return nil
}

func (m *Manager) applyHybrid(dryRun bool) error {
	// Execute makefile apply first, then terraform
	if m.makefileMgr != nil {
		if err := m.applyMakefile(dryRun); err != nil {
			return fmt.Errorf("makefile apply failed: %w", err)
		}
	}

	if m.terraformMgr != nil {
		if err := m.applyTerraform(dryRun); err != nil {
			return fmt.Errorf("terraform apply failed: %w", err)
		}
	}

	return nil
}

func (m *Manager) destroyHybrid(dryRun bool) error {
	// Execute terraform destroy first, then makefile
	var errs []error

	if m.terraformMgr != nil {
		if err := m.destroyTerraform(dryRun); err != nil {
			errs = append(errs, fmt.Errorf("terraform destroy failed: %w", err))
		}
	}

	if m.makefileMgr != nil {
		if err := m.destroyMakefile(dryRun); err != nil {
			errs = append(errs, fmt.Errorf("makefile destroy failed: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("hybrid destroy failed with %d errors: %v", len(errs), errs)
	}

	return nil
}

func (m *Manager) validateHybrid(dryRun bool) error {
	// Validate both makefile and terraform
	if m.makefileMgr != nil {
		if err := m.validateMakefile(dryRun); err != nil {
			return fmt.Errorf("makefile validation failed: %w", err)
		}
	}

	if m.terraformMgr != nil {
		if err := m.validateTerraform(dryRun); err != nil {
			return fmt.Errorf("terraform validation failed: %w", err)
		}
	}

	return nil
}

// GetProvisionMode returns the current provision mode
func (m *Manager) GetProvisionMode() string {
	return m.provisionMode
}

// GetMakefileManager returns the makefile manager (if available)
func (m *Manager) GetMakefileManager() *makefile.Manager {
	return m.makefileMgr
}

// GetTerraformManager returns the terraform manager (if available)
func (m *Manager) GetTerraformManager() *terraform.Manager {
	return m.terraformMgr
}

// GetInfo returns information about the infrastructure manager
func (m *Manager) GetInfo() *ManagerInfo {
	info := &ManagerInfo{
		ProvisionMode:     m.provisionMode,
		TerraformEnabled:  m.terraformMgr != nil,
		MakefileEnabled:   m.makefileMgr != nil,
		HealthCheckConfig: m.config.HealthCheck,
	}

	if m.makefileMgr != nil {
		info.MakefileInfo = m.makefileMgr.GetMakefileInfo()
	}

	return info
}

// RunHealthChecks runs health checks on the infrastructure
func (m *Manager) RunHealthChecks() error {
	switch m.provisionMode {
	case ProvisionModeTerraform:
		return m.terraformMgr.RunHealthChecks()
	case ProvisionModeMakefile:
		// For Makefile mode, we can run a health check target if defined
		if m.makefileMgr != nil && m.config.Makefile.Targets.HealthCheck != "" {
			return m.makefileMgr.ExecuteTarget(m.config.Makefile.Targets.HealthCheck, false)
		}
		return nil
	case ProvisionModeHybrid:
		// For hybrid mode, run both terraform and makefile health checks
		if m.terraformMgr != nil {
			if err := m.terraformMgr.RunHealthChecks(); err != nil {
				return fmt.Errorf("terraform health check failed: %w", err)
			}
		}
		if m.makefileMgr != nil && m.config.Makefile.Targets.HealthCheck != "" {
			if err := m.makefileMgr.ExecuteTarget(m.config.Makefile.Targets.HealthCheck, false); err != nil {
				return fmt.Errorf("makefile health check failed: %w", err)
			}
		}
		return nil
	default:
		return fmt.Errorf("unsupported provision mode: %s", m.provisionMode)
	}
}

// ManagerInfo contains information about the infrastructure manager
type ManagerInfo struct {
	ProvisionMode     string                   `json:"provisionMode"`
	TerraformEnabled  bool                     `json:"terraformEnabled"`
	MakefileEnabled   bool                     `json:"makefileEnabled"`
	HealthCheckConfig config.HealthCheckConfig `json:"healthCheckConfig"`
	MakefileInfo      *makefile.MakefileInfo   `json:"makefileInfo,omitempty"`
}
