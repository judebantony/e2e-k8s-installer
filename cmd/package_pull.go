package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
	"github.com/judebantony/e2e-k8s-installer/pkg/progress"
	"github.com/judebantony/e2e-k8s-installer/pkg/artifacts"
)

// packagePullCmd represents the package-pull command
var packagePullCmd = &cobra.Command{
	Use:   "package-pull",
	Short: "Synchronize OCI images, Helm charts, and Terraform modules",
	Long: `The package-pull command synchronizes all required artifacts for installation:

1. OCI Images:
   - Check if images exist in client registry (if skipPull=true)
   - Pull from vendor registry with authentication
   - Push to client registry or use vendor directly
   
2. Helm Charts:
   - Clone vendor helm repository
   - Push to client repository (if configured)
   - Keep local copy for deployment
   
3. Terraform Modules:
   - Clone vendor terraform repository  
   - Push to client repository (if configured)
   - Validate terraform modules

All operations include progress tracking and detailed logging.

Example:
  k8s-installer package-pull --config installer-config.json
  k8s-installer package-pull --images-only
  k8s-installer package-pull --dry-run`,
	RunE: runPackagePull,
}

var (
	packagePullConfig     string
	packagePullImagesOnly bool
	packagePullHelmOnly   bool
	packagePullTfOnly     bool
	packagePullDryRun     bool
	packagePullParallel   bool
)

func init() {
	rootCmd.AddCommand(packagePullCmd)
	
	packagePullCmd.Flags().StringVarP(&packagePullConfig, "config", "c", "installer-config.json", "Configuration file path")
	packagePullCmd.Flags().BoolVar(&packagePullImagesOnly, "images-only", false, "Only pull OCI images")
	packagePullCmd.Flags().BoolVar(&packagePullHelmOnly, "helm-only", false, "Only pull Helm charts")
	packagePullCmd.Flags().BoolVar(&packagePullTfOnly, "terraform-only", false, "Only pull Terraform modules")
	packagePullCmd.Flags().BoolVarP(&packagePullDryRun, "dry-run", "n", false, "Show what would be done without actually doing it")
	packagePullCmd.Flags().BoolVarP(&packagePullParallel, "parallel", "p", true, "Enable parallel processing")
}

func runPackagePull(cmd *cobra.Command, args []string) error {
	// Initialize progress manager
	progress.InitGlobalProgressManager()
	pm := progress.GetProgressManager()
	
	// Load configuration
	cfg, err := config.LoadConfig(packagePullConfig)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Initialize logger based on config
	logConfig := logger.Config{
		Level:  logger.LogLevel(cfg.Installer.LogLevel),
		Format: logger.LogFormat(cfg.Installer.LogFormat),
	}
	logger.InitGlobalLogger(logConfig)
	
	progress.ShowBanner("1.0.0")
	
	// Start overall progress tracking
	pm.StartArea("package-pull")
	
	// Determine steps based on flags
	steps := []string{}
	if !packagePullHelmOnly && !packagePullTfOnly {
		steps = append(steps, "Synchronize OCI Images")
	}
	if !packagePullImagesOnly && !packagePullTfOnly {
		steps = append(steps, "Synchronize Helm Charts")
	}
	if !packagePullImagesOnly && !packagePullHelmOnly {
		steps = append(steps, "Synchronize Terraform Modules")
	}
	steps = append(steps, "Package pull complete")
	
	currentStep := 0
	progress.ShowStepProgress(steps, currentStep)
	
	logger.Info("Starting package pull").
		Str("config", packagePullConfig).
		Bool("dry_run", packagePullDryRun).
		Bool("parallel", packagePullParallel).
		Send()
	
	// Create artifacts manager
	artifactsManager := artifacts.NewManager(cfg, packagePullDryRun)
	
	// Step 1: Synchronize OCI Images
	if !packagePullHelmOnly && !packagePullTfOnly {
		logger.StepStart("sync-images")
		
		pm.StartSpinner("images", "Synchronizing OCI images...")
		
		if err := syncImages(artifactsManager, cfg, pm); err != nil {
			pm.FailSpinner("images", "Image synchronization failed")
			logger.StepFailed("sync-images", err)
			return fmt.Errorf("image synchronization failed: %w", err)
		}
		
		pm.SuccessSpinner("images", "OCI images synchronized successfully")
		logger.StepComplete("sync-images", 0)
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
	}
	
	// Step 2: Synchronize Helm Charts
	if !packagePullImagesOnly && !packagePullTfOnly {
		logger.StepStart("sync-helm")
		
		pm.StartSpinner("helm", "Synchronizing Helm charts...")
		
		if err := syncHelmCharts(artifactsManager, cfg, pm); err != nil {
			pm.FailSpinner("helm", "Helm chart synchronization failed")
			logger.StepFailed("sync-helm", err)
			return fmt.Errorf("helm chart synchronization failed: %w", err)
		}
		
		pm.SuccessSpinner("helm", "Helm charts synchronized successfully")
		logger.StepComplete("sync-helm", 0)
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
	}
	
	// Step 3: Synchronize Terraform Modules
	if !packagePullImagesOnly && !packagePullHelmOnly {
		logger.StepStart("sync-terraform")
		
		pm.StartSpinner("terraform", "Synchronizing Terraform modules...")
		
		if err := syncTerraformModules(artifactsManager, cfg, pm); err != nil {
			pm.FailSpinner("terraform", "Terraform module synchronization failed")
			logger.StepFailed("sync-terraform", err)
			return fmt.Errorf("terraform module synchronization failed: %w", err)
		}
		
		pm.SuccessSpinner("terraform", "Terraform modules synchronized successfully")
		logger.StepComplete("sync-terraform", 0)
		currentStep++
		progress.ShowStepProgress(steps, currentStep)
	}
	
	// Complete
	currentStep++
	progress.ShowStepProgress(steps, currentStep)
	
	// Stop progress area
	pm.StopArea("package-pull")
	
	// Show success message
	progress.ShowSuccess("🎉 Package pull completed successfully!")
	
	return nil
}

func syncImages(manager *artifacts.Manager, cfg *config.InstallerConfig, pm *progress.ProgressManager) error {
	if cfg.Artifacts.Images.SkipPull {
		logger.Info("Skipping image pull as configured").Send()
		return manager.ValidateImages()
	}
	
	images := cfg.Artifacts.Images.Images
	completed := make([]bool, len(images))
	
	// Start image progress area
	pm.StartArea("images")
	progress.ShowImagePullProgress(extractImageNames(images), completed)
	
	// Start progress bar
	pm.StartProgressBar("image-progress", "Pulling Images", len(images))
	
	// Process images
	if packagePullParallel {
		return manager.SyncImagesParallel(func(index int, image config.ImageReference, err error) {
			if err == nil {
				completed[index] = true
				logger.Info("Image synchronized").
					Str("image", image.Name).
					Str("version", image.Version).
					Send()
			} else {
				logger.Error("Image synchronization failed").
					Str("image", image.Name).
					Str("version", image.Version).
					Err(err).
					Send()
			}
			
			pm.IncrementProgressBar("image-progress")
			progress.ShowImagePullProgress(extractImageNames(images), completed)
		})
	} else {
		for i, image := range images {
			if err := manager.SyncImage(image); err != nil {
				return fmt.Errorf("failed to sync image %s:%s: %w", image.Name, image.Version, err)
			}
			
			completed[i] = true
			pm.IncrementProgressBar("image-progress")
			progress.ShowImagePullProgress(extractImageNames(images), completed)
			
			logger.Info("Image synchronized").
				Str("image", image.Name).
				Str("version", image.Version).
				Send()
		}
	}
	
	pm.CompleteProgressBar("image-progress")
	pm.StopArea("images")
	
	return nil
}

func syncHelmCharts(manager *artifacts.Manager, cfg *config.InstallerConfig, pm *progress.ProgressManager) error {
	logger.Info("Synchronizing Helm charts").
		Str("vendor_repo", cfg.Artifacts.Helm.Vendor.Repo).
		Bool("push_to_client", cfg.Artifacts.Helm.Client.PushToRepo).
		Send()
	
	// Clone vendor repository
	if err := manager.CloneHelmCharts(); err != nil {
		return fmt.Errorf("failed to clone Helm charts: %w", err)
	}
	
	// Push to client repository if configured
	if cfg.Artifacts.Helm.Client.PushToRepo {
		if err := manager.PushHelmChartsToClient(); err != nil {
			return fmt.Errorf("failed to push Helm charts to client repository: %w", err)
		}
	}
	
	// Validate charts
	if err := manager.ValidateHelmCharts(); err != nil {
		return fmt.Errorf("helm chart validation failed: %w", err)
	}
	
	return nil
}

func syncTerraformModules(manager *artifacts.Manager, cfg *config.InstallerConfig, pm *progress.ProgressManager) error {
	logger.Info("Synchronizing Terraform modules").
		Str("vendor_repo", cfg.Artifacts.Terraform.Vendor.Repo).
		Bool("push_to_client", cfg.Artifacts.Terraform.Client.PushToRepo).
		Send()
	
	// Clone vendor repository
	if err := manager.CloneTerraformModules(); err != nil {
		return fmt.Errorf("failed to clone Terraform modules: %w", err)
	}
	
	// Push to client repository if configured
	if cfg.Artifacts.Terraform.Client.PushToRepo {
		if err := manager.PushTerraformModulesToClient(); err != nil {
			return fmt.Errorf("failed to push Terraform modules to client repository: %w", err)
		}
	}
	
	// Validate modules
	if err := manager.ValidateTerraformModules(); err != nil {
		return fmt.Errorf("terraform module validation failed: %w", err)
	}
	
	return nil
}

func extractImageNames(images []config.ImageReference) []string {
	names := make([]string, len(images))
	for i, img := range images {
		names[i] = fmt.Sprintf("%s:%s", img.Name, img.Version)
	}
	return names
}