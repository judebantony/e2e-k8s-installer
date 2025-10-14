package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	verbose    bool
	dryRun     bool
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "e2e-k8s-installer",
	Short: "Enterprise-grade E2E Kubernetes installer for airgapped environments",
	Long: `E2E Kubernetes Installer is a comprehensive solution designed to facilitate 
the deployment of applications in Kubernetes clusters, particularly in airgapped 
environments where internet access is restricted.

Features:
- Multi-cloud support (AWS, Azure, GCP, On-premises)
- Airgapped environment deployment
- Infrastructure provisioning with Terraform
- Application deployment with Helm charts
- Comprehensive monitoring and logging
- End-to-end testing and validation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.e2e-k8s-installer.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "perform a dry run without making changes")
	rootCmd.PersistentFlags().StringVar(&configPath, "config-path", "", "path to configuration directory")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
	viper.BindPFlag("config-path", rootCmd.PersistentFlags().Lookup("config-path"))

	// Add subcommands that we know work
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(packagePullCmd)
	rootCmd.AddCommand(provisionInfraCmd)
	rootCmd.AddCommand(dbMigrateCmd)
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(postValidateCmd)
	// Temporary placeholder for e2e-test command  
	tempE2ECmd := &cobra.Command{
		Use:   "e2e-test",
		Short: "Execute end-to-end testing suite (placeholder)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("e2e-test command not yet fully implemented")
		},
	}
	rootCmd.AddCommand(tempE2ECmd)
	rootCmd.AddCommand(installCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".e2e-k8s-installer")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}