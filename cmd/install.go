package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/judebantony/e2e-k8s-installer/pkg/installer"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

var (
	installInteractive bool
	installConfigFile  string
	installNamespace   string
	installProvider    string
	installRegistry    string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the application stack to Kubernetes cluster",
	Long: `Install the complete application stack including infrastructure provisioning,
application deployment, monitoring, and security components.

Examples:
  # Interactive installation
  e2e-k8s-installer install --interactive

  # Install using configuration file
  e2e-k8s-installer install --config-file ./config.yaml

  # Install to specific namespace and cloud provider
  e2e-k8s-installer install --namespace myapp --provider aws

  # Dry run installation
  e2e-k8s-installer install --dry-run --config-file ./config.yaml`,
	RunE: runInstall,
}

func init() {
	installCmd.Flags().BoolVarP(&installInteractive, "interactive", "i", false, "run installation in interactive mode")
	installCmd.Flags().StringVarP(&installConfigFile, "config-file", "c", "", "path to configuration file")
	installCmd.Flags().StringVarP(&installNamespace, "namespace", "n", "default", "kubernetes namespace for installation")
	installCmd.Flags().StringVarP(&installProvider, "provider", "p", "", "cloud provider (aws, azure, gcp, onprem)")
	installCmd.Flags().StringVarP(&installRegistry, "registry", "r", "", "container registry URL")
}

func runInstall(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Starting E2E Kubernetes Installation...")

	// Load configuration
	cfg, err := config.LoadConfig(installConfigFile, installInteractive)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Override with command line flags
	if installNamespace != "" {
		cfg.Kubernetes.Namespace = installNamespace
	}
	if installProvider != "" {
		cfg.Cloud.Provider = installProvider
	}
	if installRegistry != "" {
		cfg.Registry.URL = installRegistry
	}

	// Create installer instance
	inst, err := installer.NewInstaller(cfg)
	if err != nil {
		return fmt.Errorf("failed to create installer: %w", err)
	}

	// Run installation
	if err := inst.Install(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	fmt.Println("✅ Installation completed successfully!")
	return nil
}