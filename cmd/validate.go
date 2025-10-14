package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/judebantony/e2e-k8s-installer/pkg/validation"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the environment and configuration",
	Long: `Perform comprehensive validation of the target environment including:
- Kubernetes cluster connectivity and permissions
- Cloud provider credentials and quotas
- Container registry access
- Network connectivity
- Resource availability
- Configuration file validation

This command helps identify issues before running the actual installation.`,
	RunE: runValidate,
}

func runValidate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Starting environment validation...")

	validator := validation.NewValidator()
	
	if err := validator.ValidateAll(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Println("✅ All validations passed!")
	return nil
}