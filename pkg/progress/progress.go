package progress

import (
	"fmt"
	"time"

	"github.com/pterm/pterm"
)

// ProgressManager manages multiple progress indicators
type ProgressManager struct {
	spinners     map[string]*pterm.SpinnerPrinter
	progressBars map[string]*pterm.ProgressbarPrinter
	areas        map[string]*pterm.AreaPrinter
}

// NewProgressManager creates a new progress manager
func NewProgressManager() *ProgressManager {
	return &ProgressManager{
		spinners:     make(map[string]*pterm.SpinnerPrinter),
		progressBars: make(map[string]*pterm.ProgressbarPrinter),
		areas:        make(map[string]*pterm.AreaPrinter),
	}
}

// StartSpinner starts a spinner with the given ID and message
func (pm *ProgressManager) StartSpinner(id, message string) {
	spinner, _ := pterm.DefaultSpinner.WithText(message).Start()
	pm.spinners[id] = spinner
}

// UpdateSpinner updates an existing spinner message
func (pm *ProgressManager) UpdateSpinner(id, message string) {
	if spinner, exists := pm.spinners[id]; exists {
		spinner.UpdateText(message)
	}
}

// SuccessSpinner marks a spinner as successful and stops it
func (pm *ProgressManager) SuccessSpinner(id, message string) {
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Success(message)
		delete(pm.spinners, id)
	}
}

// FailSpinner marks a spinner as failed and stops it
func (pm *ProgressManager) FailSpinner(id, message string) {
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Fail(message)
		delete(pm.spinners, id)
	}
}

// WarningSpinner marks a spinner as warning and stops it
func (pm *ProgressManager) WarningSpinner(id, message string) {
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Warning(message)
		delete(pm.spinners, id)
	}
}

// StartProgressBar starts a progress bar with the given ID, title, and total
func (pm *ProgressManager) StartProgressBar(id, title string, total int) {
	progressBar, _ := pterm.DefaultProgressbar.WithTitle(title).WithTotal(total).Start()
	pm.progressBars[id] = progressBar
}

// UpdateProgressBar updates the progress bar with current value
func (pm *ProgressManager) UpdateProgressBar(id string, current int) {
	if progressBar, exists := pm.progressBars[id]; exists {
		progressBar.Current = current
	}
}

// IncrementProgressBar increments the progress bar by 1
func (pm *ProgressManager) IncrementProgressBar(id string) {
	if progressBar, exists := pm.progressBars[id]; exists {
		progressBar.Increment()
	}
}

// CompleteProgressBar completes and stops the progress bar
func (pm *ProgressManager) CompleteProgressBar(id string) {
	if progressBar, exists := pm.progressBars[id]; exists {
		_, _ = progressBar.Stop()
		delete(pm.progressBars, id)
	}
}

// StartArea starts a dynamic text area
func (pm *ProgressManager) StartArea(id string) {
	area, _ := pterm.DefaultArea.Start()
	pm.areas[id] = area
}

// UpdateArea updates the content of a text area
func (pm *ProgressManager) UpdateArea(id, content string) {
	if area, exists := pm.areas[id]; exists {
		area.Update(content)
	}
}

// StopArea stops and clears a text area
func (pm *ProgressManager) StopArea(id string) {
	if area, exists := pm.areas[id]; exists {
		area.Stop()
		delete(pm.areas, id)
	}
}

// StopAll stops all active progress indicators
func (pm *ProgressManager) StopAll() {
	// Stop all spinners
	for id, spinner := range pm.spinners {
		spinner.Stop()
		delete(pm.spinners, id)
	}

	// Stop all progress bars
	for id, progressBar := range pm.progressBars {
		_, _ = progressBar.Stop()
		delete(pm.progressBars, id)
	}

	// Stop all areas
	for id, area := range pm.areas {
		area.Stop()
		delete(pm.areas, id)
	}
}

// Global progress manager instance
var globalProgressManager *ProgressManager

// InitGlobalProgressManager initializes the global progress manager
func InitGlobalProgressManager() {
	globalProgressManager = NewProgressManager()
}

// GetProgressManager returns the global progress manager
func GetProgressManager() *ProgressManager {
	if globalProgressManager == nil {
		globalProgressManager = NewProgressManager()
	}
	return globalProgressManager
}

// Convenience functions for common progress patterns

// ShowStepProgress shows a step-based progress indicator
func ShowStepProgress(steps []string, currentStep int) {
	pm := GetProgressManager()

	// Create a progress display
	content := "\n"
	for i, step := range steps {
		var symbol string
		var color pterm.Color
		if i < currentStep {
			symbol = "✅"
			color = pterm.FgGreen
		} else if i == currentStep {
			symbol = "🔄"
			color = pterm.FgYellow
		} else {
			symbol = "⏳"
			color = pterm.FgLightWhite
		}

		content += fmt.Sprintf("  %s %s %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%d.", i+1),
			pterm.NewStyle(color).Sprintf("%s", step))
	}
	content += "\n"

	pm.UpdateArea("steps", content)
}

// ShowImagePullProgress shows progress for pulling multiple images
func ShowImagePullProgress(images []string, completed []bool) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("📦 Image Pull Progress") + "\n\n"

	completedCount := 0
	for i, image := range images {
		var symbol string
		var color pterm.Color
		if len(completed) > i && completed[i] {
			symbol = "✅"
			color = pterm.FgGreen
			completedCount++
		} else {
			symbol = "⏳"
			color = pterm.FgYellow
		}

		content += fmt.Sprintf("  %s %s\n",
			symbol,
			pterm.NewStyle(color).Sprintf("%s", image))
	}

	// Add summary
	progress := float64(completedCount) / float64(len(images)) * 100
	content += fmt.Sprintf("\n📊 Progress: %d/%d (%.1f%%)\n",
		completedCount, len(images), progress)

	pm.UpdateArea("images", content)
}

// ShowHealthCheckProgress shows health check progress
func ShowHealthCheckProgress(checks map[string]string) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🏥 Health Check Status") + "\n\n"

	for service, status := range checks {
		var symbol string
		var color pterm.Color
		switch status {
		case "healthy":
			symbol = "✅"
			color = pterm.FgGreen
		case "unhealthy":
			symbol = "❌"
			color = pterm.FgRed
		case "checking":
			symbol = "🔄"
			color = pterm.FgYellow
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
		}

		content += fmt.Sprintf("  %s %s: %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s", service),
			pterm.NewStyle(color).Sprintf("%s", status))
	}

	pm.UpdateArea("health", content)
}

// ShowTerraformProgress shows Terraform execution progress
func ShowTerraformProgress(modules []string, status map[string]string) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🏗️ Infrastructure Progress") + "\n\n"

	for _, module := range modules {
		moduleStatus := status[module]
		var symbol string
		var color pterm.Color
		switch moduleStatus {
		case "completed":
			symbol = "✅"
			color = pterm.FgGreen
		case "running":
			symbol = "🔄"
			color = pterm.FgYellow
		case "failed":
			symbol = "❌"
			color = pterm.FgRed
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
		}

		content += fmt.Sprintf("  %s %s: %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s", module),
			pterm.NewStyle(color).Sprintf("%s", moduleStatus))
	}

	pm.UpdateArea("terraform", content)
}

// ShowTestProgress shows test execution progress
func ShowTestProgress(testSuites []string, results map[string]TestResult) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🧪 Test Execution Progress") + "\n\n"

	for _, suite := range testSuites {
		result := results[suite]
		var symbol string
		var color pterm.Color
		switch result.Status {
		case "passed":
			symbol = "✅"
			color = pterm.FgGreen
		case "failed":
			symbol = "❌"
			color = pterm.FgRed
		case "running":
			symbol = "🔄"
			color = pterm.FgYellow
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
		}

		statusText := result.Status
		if result.Total > 0 {
			statusText = fmt.Sprintf("%s (%d/%d)", result.Status, result.Passed, result.Total)
		}

		content += fmt.Sprintf("  %s %s: %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s", suite),
			pterm.NewStyle(color).Sprintf("%s", statusText))
	}

	pm.UpdateArea("tests", content)
}

// TestResult represents test execution results
type TestResult struct {
	Status string
	Passed int
	Failed int
	Total  int
}

// Success messages and formatting
func ShowSuccess(message string) {
	pterm.Success.Println(message)
}

func ShowError(message string) {
	pterm.Error.Println(message)
}

func ShowWarning(message string) {
	pterm.Warning.Println(message)
}

func ShowInfo(message string) {
	pterm.Info.Println(message)
}

// ShowBanner displays a banner with the installer information
func ShowBanner(version string) {
	banner, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("K8s", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("Installer", pterm.NewStyle(pterm.FgLightMagenta))).
		Srender()

	pterm.DefaultCenter.Println(banner)
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		pterm.LightMagenta("Enterprise Kubernetes Installation Tool\n") +
			pterm.Gray("Version: "+version))
	pterm.Println()
}

// ShowSummary displays installation summary
func ShowSummary(steps []string, results map[string]string, duration time.Duration) {
	pterm.DefaultSection.Println("Installation Summary")

	successCount := 0
	failedCount := 0

	for _, step := range steps {
		result := results[step]
		var symbol string
		var color pterm.Color
		switch result {
		case "success":
			symbol = "✅"
			color = pterm.FgGreen
			successCount++
		case "failed":
			symbol = "❌"
			color = pterm.FgRed
			failedCount++
		case "skipped":
			symbol = "⏭️"
			color = pterm.FgYellow
		default:
			symbol = "❓"
			color = pterm.FgLightWhite
		}

		pterm.Printf("  %s %s: %s\n",
			symbol,
			step,
			pterm.NewStyle(color).Sprintf("%s", result))
	}

	pterm.Println()

	// Summary statistics
	totalSteps := len(steps)
	pterm.Printf("📊 Results: %s successful, %s failed, %s total\n",
		pterm.Green(fmt.Sprintf("%d", successCount)),
		pterm.Red(fmt.Sprintf("%d", failedCount)),
		pterm.LightWhite(fmt.Sprintf("%d", totalSteps)))

	pterm.Printf("⏱️ Total Duration: %s\n", pterm.Cyan(duration.String()))

	if failedCount == 0 {
		pterm.Success.Println("🎉 Installation completed successfully!")
	} else {
		pterm.Error.Println("❌ Installation completed with errors")
	}
}
