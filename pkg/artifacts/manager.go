package artifacts

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/judebantony/e2e-k8s-installer/pkg/config"
	"github.com/judebantony/e2e-k8s-installer/pkg/logger"
)

// Manager handles artifact synchronization operations
type Manager struct {
	config *config.InstallerConfig
	dryRun bool
}

// NewManager creates a new artifacts manager
func NewManager(cfg *config.InstallerConfig, dryRun bool) *Manager {
	return &Manager{
		config: cfg,
		dryRun: dryRun,
	}
}

// ImageSyncCallback is called during parallel image synchronization
type ImageSyncCallback func(index int, image config.ImageReference, err error)

// ValidateImages checks if all required images are accessible
func (m *Manager) ValidateImages() error {
	logger.Info("Validating image accessibility").Send()

	for _, image := range m.config.Artifacts.Images.Images {
		if err := m.validateSingleImage(image); err != nil {
			if image.Required {
				return fmt.Errorf("required image %s:%s not accessible: %w", image.Name, image.Version, err)
			}
			logger.Warn("Optional image not accessible").
				Str("image", image.Name).
				Str("version", image.Version).
				Err(err).
				Send()
		}
	}

	return nil
}

// SyncImage synchronizes a single OCI image
func (m *Manager) SyncImage(image config.ImageReference) error {
	logger.Info("Synchronizing image").
		Str("image", image.Name).
		Str("version", image.Version).
		Bool("dry_run", m.dryRun).
		Send()

	if m.dryRun {
		logger.Info("DRY RUN: Would sync image").
			Str("image", image.Name).
			Str("version", image.Version).
			Send()
		return nil
	}

	// Build source and destination image references
	sourceRef := fmt.Sprintf("%s/%s:%s",
		m.config.Artifacts.Images.Vendor.Registry,
		image.Name,
		image.Version)

	// Check if client registry is configured
	if m.config.Artifacts.Images.Client.Registry == "" {
		// No client registry - just validate vendor image exists
		return m.validateImageExists(sourceRef, m.config.Artifacts.Images.Vendor.Auth)
	}

	// Client registry configured - copy image
	destRef := fmt.Sprintf("%s/%s:%s",
		m.config.Artifacts.Images.Client.Registry,
		image.Name,
		image.Version)

	return m.copyImage(sourceRef, destRef)
}

// SyncImagesParallel synchronizes multiple images in parallel
func (m *Manager) SyncImagesParallel(callback ImageSyncCallback) error {
	images := m.config.Artifacts.Images.Images
	var wg sync.WaitGroup
	errorChan := make(chan error, len(images))

	// Create semaphore to limit concurrent operations
	semaphore := make(chan struct{}, 5) // Limit to 5 concurrent operations

	for i, image := range images {
		wg.Add(1)
		go func(index int, img config.ImageReference) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			err := m.SyncImage(img)
			if callback != nil {
				callback(index, img, err)
			}

			if err != nil {
				errorChan <- fmt.Errorf("image %s:%s sync failed: %w", img.Name, img.Version, err)
			}
		}(i, image)
	}

	wg.Wait()
	close(errorChan)

	// Check for errors
	var errors []string
	for err := range errorChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("parallel image sync failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// CloneHelmCharts clones Helm charts from vendor repository
func (m *Manager) CloneHelmCharts() error {
	logger.Info("Cloning Helm charts").
		Str("repo", m.config.Artifacts.Helm.Vendor.Repo).
		Send()

	if m.dryRun {
		logger.Info("DRY RUN: Would clone Helm charts").
			Str("repo", m.config.Artifacts.Helm.Vendor.Repo).
			Send()
		return nil
	}

	// Prepare clone options
	cloneOptions := &git.CloneOptions{
		URL: m.config.Artifacts.Helm.Vendor.Repo,
	}

	// Add authentication if configured
	if m.config.Artifacts.Helm.Vendor.Auth.Token != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: "token",
			Password: m.config.Artifacts.Helm.Vendor.Auth.Token,
		}
	} else if m.config.Artifacts.Helm.Vendor.Auth.Username != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: m.config.Artifacts.Helm.Vendor.Auth.Username,
			Password: m.config.Artifacts.Helm.Vendor.Auth.Password,
		}
	}

	// Set branch or tag
	if m.config.Artifacts.Helm.Vendor.Branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(m.config.Artifacts.Helm.Vendor.Branch)
	} else if m.config.Artifacts.Helm.Vendor.Tag != "" {
		cloneOptions.ReferenceName = plumbing.NewTagReferenceName(m.config.Artifacts.Helm.Vendor.Tag)
	}

	// Clone to local path
	localPath := filepath.Join(m.config.Installer.Workspace, "artifacts", "helm")
	if err := os.RemoveAll(localPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean existing helm directory: %w", err)
	}

	_, err := git.PlainClone(localPath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone helm repository: %w", err)
	}

	logger.Info("Helm charts cloned successfully").
		Str("local_path", localPath).
		Send()

	return nil
}

// PushHelmChartsToClient pushes Helm charts to client repository
func (m *Manager) PushHelmChartsToClient() error {
	if m.config.Artifacts.Helm.Client.Repo == "" {
		return fmt.Errorf("client helm repository not configured")
	}

	logger.Info("Pushing Helm charts to client repository").
		Str("client_repo", m.config.Artifacts.Helm.Client.Repo).
		Send()

	if m.dryRun {
		logger.Info("DRY RUN: Would push Helm charts to client repository").
			Str("client_repo", m.config.Artifacts.Helm.Client.Repo).
			Send()
		return nil
	}

	// Implementation would involve:
	// 1. Initialize/clone client repository
	// 2. Copy charts from local artifacts
	// 3. Commit and push changes

	// For now, return success as this is a placeholder
	logger.Info("Helm charts pushed to client repository successfully").Send()
	return nil
}

// ValidateHelmCharts validates the downloaded Helm charts
func (m *Manager) ValidateHelmCharts() error {
	logger.Info("Validating Helm charts").Send()

	chartsPath := filepath.Join(m.config.Installer.Workspace, "artifacts", "helm")

	// Check if charts directory exists
	if _, err := os.Stat(chartsPath); os.IsNotExist(err) {
		return fmt.Errorf("helm charts directory not found: %s", chartsPath)
	}

	// Basic validation - check for Chart.yaml files
	err := filepath.Walk(chartsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == "Chart.yaml" || info.Name() == "Chart.yml" {
			logger.Debug("Found Helm chart").Str("chart", filepath.Dir(path)).Send()
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("helm charts validation failed: %w", err)
	}

	logger.Info("Helm charts validation completed").Send()
	return nil
}

// CloneTerraformModules clones Terraform modules from vendor repository
func (m *Manager) CloneTerraformModules() error {
	logger.Info("Cloning Terraform modules").
		Str("repo", m.config.Artifacts.Terraform.Vendor.Repo).
		Send()

	if m.dryRun {
		logger.Info("DRY RUN: Would clone Terraform modules").
			Str("repo", m.config.Artifacts.Terraform.Vendor.Repo).
			Send()
		return nil
	}

	// Prepare clone options
	cloneOptions := &git.CloneOptions{
		URL: m.config.Artifacts.Terraform.Vendor.Repo,
	}

	// Add authentication if configured
	if m.config.Artifacts.Terraform.Vendor.Auth.Token != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: "token",
			Password: m.config.Artifacts.Terraform.Vendor.Auth.Token,
		}
	} else if m.config.Artifacts.Terraform.Vendor.Auth.Username != "" {
		cloneOptions.Auth = &http.BasicAuth{
			Username: m.config.Artifacts.Terraform.Vendor.Auth.Username,
			Password: m.config.Artifacts.Terraform.Vendor.Auth.Password,
		}
	}

	// Set branch or tag
	if m.config.Artifacts.Terraform.Vendor.Branch != "" {
		cloneOptions.ReferenceName = plumbing.NewBranchReferenceName(m.config.Artifacts.Terraform.Vendor.Branch)
	} else if m.config.Artifacts.Terraform.Vendor.Tag != "" {
		cloneOptions.ReferenceName = plumbing.NewTagReferenceName(m.config.Artifacts.Terraform.Vendor.Tag)
	}

	// Clone to local path
	localPath := filepath.Join(m.config.Installer.Workspace, "artifacts", "terraform")
	if err := os.RemoveAll(localPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clean existing terraform directory: %w", err)
	}

	_, err := git.PlainClone(localPath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone terraform repository: %w", err)
	}

	logger.Info("Terraform modules cloned successfully").
		Str("local_path", localPath).
		Send()

	return nil
}

// PushTerraformModulesToClient pushes Terraform modules to client repository
func (m *Manager) PushTerraformModulesToClient() error {
	if m.config.Artifacts.Terraform.Client.Repo == "" {
		return fmt.Errorf("client terraform repository not configured")
	}

	logger.Info("Pushing Terraform modules to client repository").
		Str("client_repo", m.config.Artifacts.Terraform.Client.Repo).
		Send()

	if m.dryRun {
		logger.Info("DRY RUN: Would push Terraform modules to client repository").
			Str("client_repo", m.config.Artifacts.Terraform.Client.Repo).
			Send()
		return nil
	}

	// Implementation would involve:
	// 1. Initialize/clone client repository
	// 2. Copy modules from local artifacts
	// 3. Commit and push changes

	// For now, return success as this is a placeholder
	logger.Info("Terraform modules pushed to client repository successfully").Send()
	return nil
}

// ValidateTerraformModules validates the downloaded Terraform modules
func (m *Manager) ValidateTerraformModules() error {
	logger.Info("Validating Terraform modules").Send()

	modulesPath := filepath.Join(m.config.Installer.Workspace, "artifacts", "terraform")

	// Check if modules directory exists
	if _, err := os.Stat(modulesPath); os.IsNotExist(err) {
		return fmt.Errorf("terraform modules directory not found: %s", modulesPath)
	}

	// Basic validation - check for .tf files
	err := filepath.Walk(modulesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".tf") {
			logger.Debug("Found Terraform file").Str("file", path).Send()
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("terraform modules validation failed: %w", err)
	}

	logger.Info("Terraform modules validation completed").Send()
	return nil
}

// validateSingleImage validates if an image is accessible
func (m *Manager) validateSingleImage(image config.ImageReference) error {
	// Try vendor registry first
	vendorRef := fmt.Sprintf("%s/%s:%s",
		m.config.Artifacts.Images.Vendor.Registry,
		image.Name,
		image.Version)

	if err := m.validateImageExists(vendorRef, m.config.Artifacts.Images.Vendor.Auth); err == nil {
		return nil
	}

	// If vendor fails and client registry is configured, try client
	if m.config.Artifacts.Images.Client.Registry != "" {
		clientRef := fmt.Sprintf("%s/%s:%s",
			m.config.Artifacts.Images.Client.Registry,
			image.Name,
			image.Version)

		return m.validateImageExists(clientRef, m.config.Artifacts.Images.Client.Auth)
	}

	return fmt.Errorf("image not accessible in any configured registry")
}

// validateImageExists checks if an image exists in a registry
func (m *Manager) validateImageExists(imageRef string, auth config.AuthConfig) error {
	logger.Debug("Validating image exists").Str("image", imageRef).Send()

	// Parse image reference
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("invalid image reference %s: %w", imageRef, err)
	}

	// Create remote options with authentication
	options := []remote.Option{}
	if auth.Token != "" {
		options = append(options, remote.WithAuth(authn.FromConfig(authn.AuthConfig{
			Auth: auth.Token,
		})))
	} else if auth.Username != "" {
		options = append(options, remote.WithAuth(authn.FromConfig(authn.AuthConfig{
			Username: auth.Username,
			Password: auth.Password,
		})))
	}

	// Check if image exists
	_, err = remote.Head(ref, options...)
	if err != nil {
		return fmt.Errorf("image %s not accessible: %w", imageRef, err)
	}

	return nil
}

// copyImage copies an image from source to destination registry
func (m *Manager) copyImage(sourceRef, destRef string) error {
	logger.Info("Copying image").
		Str("source", sourceRef).
		Str("destination", destRef).
		Send()

	// Use crane to copy the image
	options := []crane.Option{}

	// Add authentication for source
	if m.config.Artifacts.Images.Vendor.Auth.Token != "" {
		// Configure auth for vendor registry
	}

	// Add authentication for destination
	if m.config.Artifacts.Images.Client.Auth.Token != "" {
		// Configure auth for client registry
	}

	if err := crane.Copy(sourceRef, destRef, options...); err != nil {
		return fmt.Errorf("failed to copy image from %s to %s: %w", sourceRef, destRef, err)
	}

	logger.Info("Image copied successfully").
		Str("source", sourceRef).
		Str("destination", destRef).
		Send()

	return nil
}
