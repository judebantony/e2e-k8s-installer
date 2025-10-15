package makefile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
)

// Manager handles Makefile-based infrastructure operations
type Manager struct {
	config     *config.MakefileExecution
	workingDir string
	makePath   string
	env        []string
}

// NewManager creates a new Makefile manager
func NewManager(makefileConfig *config.MakefileExecution) (*Manager, error) {
	if makefileConfig == nil {
		return nil, fmt.Errorf("makefile configuration is required")
	}

	if !makefileConfig.Enabled {
		return nil, fmt.Errorf("makefile execution is not enabled")
	}

	// Determine working directory
	workingDir := makefileConfig.WorkingDirectory
	if workingDir == "" {
		workingDir = "."
	}

	// Resolve Makefile path
	makefilePath := makefileConfig.MakefilePath
	if makefilePath == "" {
		makefilePath = filepath.Join(workingDir, "Makefile")
	}

	// Check if Makefile exists
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("makefile not found at path: %s", makefilePath)
	}

	// Check if make command is available
	makePath, err := exec.LookPath("make")
	if err != nil {
		return nil, fmt.Errorf("make command not found in PATH: %w", err)
	}

	// Prepare environment variables
	env := os.Environ()
	for key, value := range makefileConfig.Environment {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	// Add makefile variables as environment variables
	for key, value := range makefileConfig.Variables {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	return &Manager{
		config:     makefileConfig,
		workingDir: workingDir,
		makePath:   makePath,
		env:        env,
	}, nil
}

// ExecuteTarget executes a specific Makefile target
func (m *Manager) ExecuteTarget(target string, dryRun bool) error {
	if target == "" {
		return fmt.Errorf("target cannot be empty")
	}

	logger.Info("Executing Makefile target").
		Str("target", target).
		Str("workingDir", m.workingDir).
		Bool("dryRun", dryRun).
		Send()

	// Prepare command arguments
	args := []string{}

	// Add makefile path if not default
	if m.config.MakefilePath != "" && !strings.HasSuffix(m.config.MakefilePath, "Makefile") {
		args = append(args, "-f", m.config.MakefilePath)
	}

	// Add parallel flag if enabled
	if m.config.Parallel {
		args = append(args, "-j")
	}

	// Add keep going flag if enabled
	if m.config.KeepGoing {
		args = append(args, "-k")
	}

	// Add dry run flag if enabled or requested
	if dryRun || m.config.DryRun {
		args = append(args, "-n")
	}

	// Add target
	args = append(args, target)

	// Create command context with timeout
	ctx := context.Background()
	if m.config.Timeout != "" {
		timeout, err := time.ParseDuration(m.config.Timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout duration: %w", err)
		}
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	// Execute command
	cmd := exec.CommandContext(ctx, m.makePath, args...)
	cmd.Dir = m.workingDir
	cmd.Env = m.env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logger.Info("Running make command").
		Str("command", fmt.Sprintf("%s %s", m.makePath, strings.Join(args, " "))).
		Str("workingDir", m.workingDir).
		Send()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make target '%s' failed: %w", target, err)
	}

	logger.Info("Makefile target completed successfully").
		Str("target", target).
		Send()

	return nil
}

// Init executes the init target
func (m *Manager) Init(dryRun bool) error {
	target := m.config.Targets.Init
	if target == "" {
		target = "init"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Plan executes the plan target
func (m *Manager) Plan(dryRun bool) error {
	target := m.config.Targets.Plan
	if target == "" {
		target = "plan"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Apply executes the apply target
func (m *Manager) Apply(dryRun bool) error {
	target := m.config.Targets.Apply
	if target == "" {
		target = "apply"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Destroy executes the destroy target
func (m *Manager) Destroy(dryRun bool) error {
	target := m.config.Targets.Destroy
	if target == "" {
		target = "destroy"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Validate executes the validate target
func (m *Manager) Validate(dryRun bool) error {
	target := m.config.Targets.Validate
	if target == "" {
		target = "validate"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Clean executes the clean target
func (m *Manager) Clean(dryRun bool) error {
	target := m.config.Targets.Clean
	if target == "" {
		target = "clean"
	}
	return m.ExecuteTarget(target, dryRun)
}

// Format executes the format target
func (m *Manager) Format(dryRun bool) error {
	target := m.config.Targets.Format
	if target == "" {
		target = "fmt"
	}
	return m.ExecuteTarget(target, dryRun)
}

// ExecuteCustomTarget executes a custom target
func (m *Manager) ExecuteCustomTarget(targetName string, dryRun bool) error {
	target, exists := m.config.Targets.Custom[targetName]
	if !exists {
		target = targetName // Use the target name directly if not in custom mapping
	}
	return m.ExecuteTarget(target, dryRun)
}

// ListTargets lists all available targets in the Makefile
func (m *Manager) ListTargets() ([]string, error) {
	cmd := exec.Command(m.makePath, "-f", m.config.MakefilePath, "-p")
	cmd.Dir = m.workingDir
	cmd.Env = m.env

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list makefile targets: %w", err)
	}

	var targets []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, ":") && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, ".") {
			target := strings.Split(line, ":")[0]
			target = strings.TrimSpace(target)
			if target != "" && !strings.Contains(target, " ") {
				targets = append(targets, target)
			}
		}
	}

	return targets, nil
}

// GetMakefileInfo returns information about the Makefile configuration
func (m *Manager) GetMakefileInfo() *MakefileInfo {
	return &MakefileInfo{
		MakefilePath:     m.config.MakefilePath,
		WorkingDirectory: m.workingDir,
		Targets:          m.config.Targets,
		Environment:      m.config.Environment,
		Variables:        m.config.Variables,
		Parallel:         m.config.Parallel,
		KeepGoing:        m.config.KeepGoing,
		DryRun:           m.config.DryRun,
		Timeout:          m.config.Timeout,
	}
}

// MakefileInfo contains information about the Makefile configuration
type MakefileInfo struct {
	MakefilePath     string                 `json:"makefilePath"`
	WorkingDirectory string                 `json:"workingDirectory"`
	Targets          config.MakefileTargets `json:"targets"`
	Environment      map[string]string      `json:"environment"`
	Variables        map[string]string      `json:"variables"`
	Parallel         bool                   `json:"parallel"`
	KeepGoing        bool                   `json:"keepGoing"`
	DryRun           bool                   `json:"dryRun"`
	Timeout          string                 `json:"timeout"`
}
