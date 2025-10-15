package progress

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

// ProgressManager manages multiple progress indicators with enterprise-scale features
type ProgressManager struct {
	spinners       map[string]*pterm.SpinnerPrinter
	progressBars   map[string]*pterm.ProgressbarPrinter
	areas          map[string]*pterm.AreaPrinter
	operations     map[string]*OperationProgress
	mutex          sync.RWMutex
	startTime      time.Time
	enterpriseMode bool
}

// OperationProgress tracks detailed progress for enterprise operations
type OperationProgress struct {
	ID          string
	Name        string
	Description string
	StartTime   time.Time
	EndTime     *time.Time
	Status      OperationStatus
	Progress    int
	Total       int
	SubSteps    []SubStep
	Metadata    map[string]interface{}
	Duration    time.Duration
	ErrorMsg    string
}

// SubStep represents a sub-operation within a main operation
type SubStep struct {
	Name        string
	Status      OperationStatus
	StartTime   time.Time
	EndTime     *time.Time
	Duration    time.Duration
	Progress    int
	Total       int
	Description string
}

// OperationStatus represents the status of an operation
type OperationStatus string

const (
	StatusPending    OperationStatus = "pending"
	StatusRunning    OperationStatus = "running"
	StatusCompleted  OperationStatus = "completed"
	StatusFailed     OperationStatus = "failed"
	StatusSkipped    OperationStatus = "skipped"
	StatusCancelled  OperationStatus = "cancelled"
	StatusWarning    OperationStatus = "warning"
)

// ProgressMetrics holds overall progress metrics
type ProgressMetrics struct {
	TotalOperations     int
	CompletedOperations int
	FailedOperations    int
	SkippedOperations   int
	OverallProgress     float64
	EstimatedTimeLeft   time.Duration
	ElapsedTime         time.Duration
	Throughput          float64
}

// NewProgressManager creates a new progress manager with enterprise features
func NewProgressManager() *ProgressManager {
	return &ProgressManager{
		spinners:       make(map[string]*pterm.SpinnerPrinter),
		progressBars:   make(map[string]*pterm.ProgressbarPrinter),
		areas:          make(map[string]*pterm.AreaPrinter),
		operations:     make(map[string]*OperationProgress),
		startTime:      time.Now(),
		enterpriseMode: true,
	}
}

// EnableEnterpriseMode enables advanced enterprise features
func (pm *ProgressManager) EnableEnterpriseMode() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	pm.enterpriseMode = true
}

// StartOperation starts tracking a new operation with enterprise features
func (pm *ProgressManager) StartOperation(id, name, description string, total int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	operation := &OperationProgress{
		ID:          id,
		Name:        name,
		Description: description,
		StartTime:   time.Now(),
		Status:      StatusRunning,
		Progress:    0,
		Total:       total,
		SubSteps:    []SubStep{},
		Metadata:    make(map[string]interface{}),
	}

	pm.operations[id] = operation

	if pm.enterpriseMode {
		pm.displayEnterpriseProgressUnsafe()
	}
}

// UpdateOperationProgress updates the progress of an operation
func (pm *ProgressManager) UpdateOperationProgress(id string, progress int, status OperationStatus, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if operation, exists := pm.operations[id]; exists {
		operation.Progress = progress
		operation.Status = status
		operation.Description = message
		operation.Duration = time.Since(operation.StartTime)

		if status == StatusCompleted || status == StatusFailed {
			now := time.Now()
			operation.EndTime = &now
		}

		pm.operations[id] = operation

		if pm.enterpriseMode {
			pm.displayEnterpriseProgressUnsafe()
		}
	}
}

// AddSubStep adds a sub-step to an operation
func (pm *ProgressManager) AddSubStep(operationID, stepName, description string, total int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if operation, exists := pm.operations[operationID]; exists {
		subStep := SubStep{
			Name:        stepName,
			Status:      StatusRunning,
			StartTime:   time.Now(),
			Progress:    0,
			Total:       total,
			Description: description,
		}

		operation.SubSteps = append(operation.SubSteps, subStep)
		pm.operations[operationID] = operation

		if pm.enterpriseMode {
			pm.displayEnterpriseProgressUnsafe()
		}
	}
}

// UpdateSubStep updates a sub-step within an operation
func (pm *ProgressManager) UpdateSubStep(operationID, stepName string, progress int, status OperationStatus) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if operation, exists := pm.operations[operationID]; exists {
		for i, subStep := range operation.SubSteps {
			if subStep.Name == stepName {
				operation.SubSteps[i].Progress = progress
				operation.SubSteps[i].Status = status
				operation.SubSteps[i].Duration = time.Since(subStep.StartTime)

				if status == StatusCompleted || status == StatusFailed {
					now := time.Now()
					operation.SubSteps[i].EndTime = &now
				}
				break
			}
		}

		pm.operations[operationID] = operation

		if pm.enterpriseMode {
			pm.displayEnterpriseProgressUnsafe()
		}
	}
}

// CompleteOperation marks an operation as complete
func (pm *ProgressManager) CompleteOperation(id string, status OperationStatus, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if operation, exists := pm.operations[id]; exists {
		operation.Status = status
		operation.Description = message
		now := time.Now()
		operation.EndTime = &now
		operation.Duration = time.Since(operation.StartTime)

		if status == StatusCompleted {
			operation.Progress = operation.Total
		}

		pm.operations[id] = operation

		if pm.enterpriseMode {
			pm.displayEnterpriseProgressUnsafe()
		}
	}
}

// GetProgressMetrics returns overall progress metrics
func (pm *ProgressManager) GetProgressMetrics() ProgressMetrics {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.getProgressMetricsUnsafe()
}

// getProgressMetricsUnsafe returns metrics without acquiring mutex (for internal use)
func (pm *ProgressManager) getProgressMetricsUnsafe() ProgressMetrics {
	metrics := ProgressMetrics{
		TotalOperations: len(pm.operations),
		ElapsedTime:     time.Since(pm.startTime),
	}

	totalProgress := 0
	totalPossible := 0

	for _, operation := range pm.operations {
		switch operation.Status {
		case StatusCompleted:
			metrics.CompletedOperations++
			totalProgress += operation.Total
		case StatusFailed:
			metrics.FailedOperations++
		case StatusSkipped:
			metrics.SkippedOperations++
		default:
			totalProgress += operation.Progress
		}
		totalPossible += operation.Total
	}

	if totalPossible > 0 {
		metrics.OverallProgress = float64(totalProgress) / float64(totalPossible) * 100
	}

	// Calculate throughput (operations per second)
	if metrics.ElapsedTime.Seconds() > 0 {
		metrics.Throughput = float64(metrics.CompletedOperations) / metrics.ElapsedTime.Seconds()
	}

	// Estimate time left based on current throughput
	remainingOps := metrics.TotalOperations - metrics.CompletedOperations
	if metrics.Throughput > 0 && remainingOps > 0 {
		metrics.EstimatedTimeLeft = time.Duration(float64(remainingOps)/metrics.Throughput) * time.Second
	}

	return metrics
}

// displayEnterpriseProgress displays a comprehensive enterprise progress view
func (pm *ProgressManager) displayEnterpriseProgress() {
	if !pm.enterpriseMode {
		return
	}

	metrics := pm.GetProgressMetrics()
	
	// Create enterprise progress display
	content := pm.buildEnterpriseProgressContent(metrics)
	
	// Update or create the enterprise progress area
	if area, exists := pm.areas["enterprise"]; exists {
		area.Update(content)
	} else {
		area, _ := pterm.DefaultArea.Start()
		pm.areas["enterprise"] = area
		area.Update(content)
	}
}

// displayEnterpriseProgressUnsafe displays progress without acquiring mutex (for internal use)
func (pm *ProgressManager) displayEnterpriseProgressUnsafe() {
	if !pm.enterpriseMode {
		return
	}

	metrics := pm.getProgressMetricsUnsafe()
	
	// Create enterprise progress display
	content := pm.buildEnterpriseProgressContentUnsafe(metrics)
	
	// Update or create the enterprise progress area
	if area, exists := pm.areas["enterprise"]; exists {
		area.Update(content)
	} else {
		area, _ := pterm.DefaultArea.Start()
		pm.areas["enterprise"] = area
		area.Update(content)
	}
}

// buildEnterpriseProgressContent builds the enterprise progress display content
func (pm *ProgressManager) buildEnterpriseProgressContent(metrics ProgressMetrics) string {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	return pm.buildEnterpriseProgressContentUnsafe(metrics)
}

// buildEnterpriseProgressContentUnsafe builds content without acquiring mutex (for internal use)
func (pm *ProgressManager) buildEnterpriseProgressContentUnsafe(metrics ProgressMetrics) string {
	var content strings.Builder

	// Header with branding
	content.WriteString(pterm.DefaultHeader.Sprint("🏢 Enterprise Kubernetes Installer"))
	content.WriteString("\n\n")

	// Overall progress bar
	progressBar := pm.createProgressBar(int(metrics.OverallProgress), 100)
	content.WriteString(fmt.Sprintf("📊 Overall Progress: %s %.1f%%\n", progressBar, metrics.OverallProgress))
	content.WriteString("\n")

	// Metrics dashboard
	content.WriteString("📈 Execution Metrics:\n")
	content.WriteString(fmt.Sprintf("   ⏱️  Elapsed Time: %s\n", formatDuration(metrics.ElapsedTime)))
	
	if metrics.EstimatedTimeLeft > 0 {
		content.WriteString(fmt.Sprintf("   ⏳ Estimated Time Left: %s\n", formatDuration(metrics.EstimatedTimeLeft)))
	}
	
	content.WriteString(fmt.Sprintf("   🎯 Operations: %d total, %d completed, %d failed\n", 
		metrics.TotalOperations, metrics.CompletedOperations, metrics.FailedOperations))
	
	if metrics.Throughput > 0 {
		content.WriteString(fmt.Sprintf("   🚀 Throughput: %.2f ops/sec\n", metrics.Throughput))
	}
	content.WriteString("\n")

	// Operation details
	content.WriteString("🔄 Operation Status:\n")
	for _, operation := range pm.operations {
		content.WriteString(pm.formatOperationLine(operation))
	}

	return content.String()
}

// formatOperationLine formats a single operation line with progress and status
func (pm *ProgressManager) formatOperationLine(operation *OperationProgress) string {
	var line strings.Builder

	// Status icon
	statusIcon := pm.getStatusIcon(operation.Status)
	
	// Progress calculation
	progressPercent := 0.0
	if operation.Total > 0 {
		progressPercent = float64(operation.Progress) / float64(operation.Total) * 100
	}
	
	// Duration formatting
	duration := operation.Duration
	if operation.EndTime != nil {
		duration = operation.EndTime.Sub(operation.StartTime)
	}

	// Main operation line
	line.WriteString(fmt.Sprintf("   %s %s", statusIcon, operation.Name))
	
	if operation.Status == StatusRunning {
		progressBar := pm.createProgressBar(operation.Progress, operation.Total)
		line.WriteString(fmt.Sprintf(" %s %.1f%%", progressBar, progressPercent))
	}
	
	line.WriteString(fmt.Sprintf(" (%s)", formatDuration(duration)))
	
	if operation.ErrorMsg != "" {
		line.WriteString(fmt.Sprintf(" - %s", pterm.Red(operation.ErrorMsg)))
	}
	
	line.WriteString("\n")

	// Sub-steps (if any)
	for _, subStep := range operation.SubSteps {
		subProgressPercent := 0.0
		if subStep.Total > 0 {
			subProgressPercent = float64(subStep.Progress) / float64(subStep.Total) * 100
		}
		
		subStatusIcon := pm.getStatusIcon(subStep.Status)
		subDuration := subStep.Duration
		if subStep.EndTime != nil {
			subDuration = subStep.EndTime.Sub(subStep.StartTime)
		}

		line.WriteString(fmt.Sprintf("     └─ %s %s", subStatusIcon, subStep.Name))
		
		if subStep.Status == StatusRunning && subStep.Total > 0 {
			subProgressBar := pm.createProgressBar(subStep.Progress, subStep.Total)
			line.WriteString(fmt.Sprintf(" %s %.1f%%", subProgressBar, subProgressPercent))
		}
		
		line.WriteString(fmt.Sprintf(" (%s)\n", formatDuration(subDuration)))
	}

	return line.String()
}

// createProgressBar creates a visual progress bar
func (pm *ProgressManager) createProgressBar(current, total int) string {
	if total <= 0 {
		return "[████████████████████] 100%"
	}
	
	percent := float64(current) / float64(total)
	if percent > 1.0 {
		percent = 1.0
	}
	
	width := 20
	filled := int(percent * float64(width))
	
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += "]"
	
	return pterm.NewStyle(pterm.FgCyan).Sprint(bar)
}

// getStatusIcon returns the appropriate icon for operation status
func (pm *ProgressManager) getStatusIcon(status OperationStatus) string {
	switch status {
	case StatusCompleted:
		return pterm.Green("✅")
	case StatusFailed:
		return pterm.Red("❌")
	case StatusRunning:
		return pterm.Yellow("🔄")
	case StatusPending:
		return pterm.LightWhite("⏳")
	case StatusSkipped:
		return pterm.Yellow("⏭️")
	case StatusCancelled:
		return pterm.Red("🚫")
	case StatusWarning:
		return pterm.Yellow("⚠️")
	default:
		return pterm.LightWhite("❓")
	}
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	} else if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	} else {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
}

// FormatDuration is a public function to format duration in a human-readable way
func FormatDuration(d time.Duration) string {
	return formatDuration(d)
}

// StartSpinner starts a spinner with the given ID and message
func (pm *ProgressManager) StartSpinner(id, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	spinner, _ := pterm.DefaultSpinner.WithText(message).Start()
	pm.spinners[id] = spinner
}

// UpdateSpinner updates an existing spinner message
func (pm *ProgressManager) UpdateSpinner(id, message string) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if spinner, exists := pm.spinners[id]; exists {
		spinner.UpdateText(message)
	}
}

// SuccessSpinner marks a spinner as successful and stops it
func (pm *ProgressManager) SuccessSpinner(id, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Success(message)
		delete(pm.spinners, id)
	}
}

// FailSpinner marks a spinner as failed and stops it
func (pm *ProgressManager) FailSpinner(id, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Fail(message)
		delete(pm.spinners, id)
	}
}

// WarningSpinner marks a spinner as warning and stops it
func (pm *ProgressManager) WarningSpinner(id, message string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	if spinner, exists := pm.spinners[id]; exists {
		spinner.Warning(message)
		delete(pm.spinners, id)
	}
}

// StartProgressBar starts a progress bar with the given ID, title, and total
func (pm *ProgressManager) StartProgressBar(id, title string, total int) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	progressBar, _ := pterm.DefaultProgressbar.WithTitle(title).WithTotal(total).Start()
	pm.progressBars[id] = progressBar
}

// UpdateProgressBar updates the progress bar with current value
func (pm *ProgressManager) UpdateProgressBar(id string, current int) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if progressBar, exists := pm.progressBars[id]; exists {
		progressBar.Current = current
	}
}

// IncrementProgressBar increments the progress bar by 1
func (pm *ProgressManager) IncrementProgressBar(id string) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if progressBar, exists := pm.progressBars[id]; exists {
		progressBar.Increment()
	}
}

// CompleteProgressBar completes and stops the progress bar
func (pm *ProgressManager) CompleteProgressBar(id string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	if progressBar, exists := pm.progressBars[id]; exists {
		_, _ = progressBar.Stop()
		delete(pm.progressBars, id)
	}
}

// StartArea starts a dynamic text area
func (pm *ProgressManager) StartArea(id string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	area, _ := pterm.DefaultArea.Start()
	pm.areas[id] = area
}

// UpdateArea updates the content of a text area
func (pm *ProgressManager) UpdateArea(id, content string) {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()
	
	if area, exists := pm.areas[id]; exists {
		area.Update(content)
	}
}

// StopArea stops and clears a text area
func (pm *ProgressManager) StopArea(id string) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
	if area, exists := pm.areas[id]; exists {
		area.Stop()
		delete(pm.areas, id)
	}
}

// StopAll stops all active progress indicators
func (pm *ProgressManager) StopAll() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	
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
	
	// Clear all operations
	pm.operations = make(map[string]*OperationProgress)
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

// Enhanced convenience functions for enterprise progress patterns

// ShowStepProgress shows a step-based progress indicator with percentage
func ShowStepProgress(steps []string, currentStep int) {
	pm := GetProgressManager()

	// Create a progress display with percentage
	content := "\n"
	progressPercent := 0.0
	if len(steps) > 0 {
		progressPercent = float64(currentStep) / float64(len(steps)) * 100
	}

	// Add progress header with percentage
	content += pterm.DefaultHeader.Sprintf("📋 Installation Progress: %.1f%% (%d/%d)", 
		progressPercent, currentStep, len(steps)) + "\n\n"

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
	
	// Add visual progress bar
	progressBar := createProgressBarString(currentStep, len(steps))
	content += fmt.Sprintf("\n%s\n", progressBar)
	content += "\n"

	pm.UpdateArea("steps", content)
}

// ShowImagePullProgress shows progress for pulling multiple images with enhanced metrics
func ShowImagePullProgress(images []string, completed []bool) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("📦 Container Image Management") + "\n\n"

	completedCount := 0
	for i, image := range images {
		var symbol string
		var color pterm.Color
		if len(completed) > i && completed[i] {
			symbol = "✅"
			color = pterm.FgGreen
			completedCount++
		} else {
			symbol = "🔄"
			color = pterm.FgYellow
		}

		// Enhanced image display with size info
		content += fmt.Sprintf("  %s %s\n",
			symbol,
			pterm.NewStyle(color).Sprintf("%s", image))
	}

	// Enhanced summary with progress bar
	progress := 0.0
	if len(images) > 0 {
		progress = float64(completedCount) / float64(len(images)) * 100
	}
	
	progressBar := createProgressBarString(completedCount, len(images))
	content += fmt.Sprintf("\n%s\n", progressBar)
	content += fmt.Sprintf("📊 Pull Progress: %d/%d (%.1f%%) completed\n",
		completedCount, len(images), progress)

	pm.UpdateArea("images", content)
}

// ShowHealthCheckProgress shows health check progress with enhanced monitoring
func ShowHealthCheckProgress(checks map[string]string) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🏥 System Health Monitoring") + "\n\n"

	healthyCount := 0
	totalChecks := len(checks)

	for service, status := range checks {
		var symbol string
		var color pterm.Color
		var statusText string
		
		switch status {
		case "healthy":
			symbol = "✅"
			color = pterm.FgGreen
			statusText = "HEALTHY"
			healthyCount++
		case "unhealthy":
			symbol = "❌"
			color = pterm.FgRed
			statusText = "UNHEALTHY"
		case "checking":
			symbol = "🔄"
			color = pterm.FgYellow
			statusText = "CHECKING"
		case "degraded":
			symbol = "⚠️"
			color = pterm.FgYellow
			statusText = "DEGRADED"
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
			statusText = "PENDING"
		}

		content += fmt.Sprintf("  %s %-20s %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s:", service),
			pterm.NewStyle(color).Sprintf("%s", statusText))
	}

	// Health summary with percentage
	healthPercent := 0.0
	if totalChecks > 0 {
		healthPercent = float64(healthyCount) / float64(totalChecks) * 100
	}
	
	progressBar := createProgressBarString(healthyCount, totalChecks)
	content += fmt.Sprintf("\n%s\n", progressBar)
	content += fmt.Sprintf("📊 System Health: %.1f%% (%d/%d services healthy)\n",
		healthPercent, healthyCount, totalChecks)

	pm.UpdateArea("health", content)
}

// ShowTerraformProgress shows Terraform execution progress with enhanced details
func ShowTerraformProgress(modules []string, status map[string]string) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🏗️ Infrastructure Provisioning") + "\n\n"

	completedCount := 0
	failedCount := 0

	for _, module := range modules {
		moduleStatus := status[module]
		var symbol string
		var color pterm.Color
		var statusText string
		
		switch moduleStatus {
		case "completed":
			symbol = "✅"
			color = pterm.FgGreen
			statusText = "DEPLOYED"
			completedCount++
		case "running":
			symbol = "🔄"
			color = pterm.FgYellow
			statusText = "DEPLOYING"
		case "failed":
			symbol = "❌"
			color = pterm.FgRed
			statusText = "FAILED"
			failedCount++
		case "planned":
			symbol = "📋"
			color = pterm.FgCyan
			statusText = "PLANNED"
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
			statusText = "PENDING"
		}

		content += fmt.Sprintf("  %s %-25s %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s:", module),
			pterm.NewStyle(color).Sprintf("%s", statusText))
	}

	// Infrastructure summary
	totalModules := len(modules)
	progressPercent := 0.0
	if totalModules > 0 {
		progressPercent = float64(completedCount) / float64(totalModules) * 100
	}
	
	progressBar := createProgressBarString(completedCount, totalModules)
	content += fmt.Sprintf("\n%s\n", progressBar)
	content += fmt.Sprintf("📊 Infrastructure: %.1f%% complete (%d/%d modules)\n",
		progressPercent, completedCount, totalModules)
	
	if failedCount > 0 {
		content += fmt.Sprintf("⚠️  Failed Modules: %d\n", failedCount)
	}

	pm.UpdateArea("terraform", content)
}

// ShowTestProgress shows test execution progress with detailed results
func ShowTestProgress(testSuites []string, results map[string]TestResult) {
	pm := GetProgressManager()

	content := pterm.DefaultHeader.Sprint("🧪 Test Suite Execution") + "\n\n"

	totalPassed := 0
	totalFailed := 0
	totalTests := 0
	completedSuites := 0

	for _, suite := range testSuites {
		result := results[suite]
		var symbol string
		var color pterm.Color
		var statusText string
		
		switch result.Status {
		case "passed":
			symbol = "✅"
			color = pterm.FgGreen
			statusText = fmt.Sprintf("PASSED (%d/%d)", result.Passed, result.Total)
			completedSuites++
		case "failed":
			symbol = "❌"
			color = pterm.FgRed
			statusText = fmt.Sprintf("FAILED (%d/%d passed)", result.Passed, result.Total)
			completedSuites++
		case "running":
			symbol = "🔄"
			color = pterm.FgYellow
			statusText = fmt.Sprintf("RUNNING (%d/%d)", result.Passed, result.Total)
		default:
			symbol = "⏳"
			color = pterm.FgLightWhite
			statusText = "PENDING"
		}

		content += fmt.Sprintf("  %s %-20s %s\n",
			symbol,
			pterm.NewStyle(pterm.FgLightWhite).Sprintf("%s:", suite),
			pterm.NewStyle(color).Sprintf("%s", statusText))

		totalPassed += result.Passed
		totalFailed += result.Failed
		totalTests += result.Total
	}

	// Test summary with metrics
	suiteProgress := 0.0
	if len(testSuites) > 0 {
		suiteProgress = float64(completedSuites) / float64(len(testSuites)) * 100
	}
	
	testProgress := 0.0
	if totalTests > 0 {
		testProgress = float64(totalPassed) / float64(totalTests) * 100
	}
	
	progressBar := createProgressBarString(completedSuites, len(testSuites))
	content += fmt.Sprintf("\n%s\n", progressBar)
	content += fmt.Sprintf("📊 Test Suites: %.1f%% complete (%d/%d suites)\n",
		suiteProgress, completedSuites, len(testSuites))
	content += fmt.Sprintf("📊 Test Cases: %.1f%% passed (%d/%d tests)\n",
		testProgress, totalPassed, totalTests)
	
	if totalFailed > 0 {
		content += fmt.Sprintf("❌ Failed Tests: %d\n", totalFailed)
	}

	pm.UpdateArea("tests", content)
}

// createProgressBarString creates a visual progress bar string
func createProgressBarString(current, total int) string {
	if total <= 0 {
		return pterm.NewStyle(pterm.FgCyan).Sprint("[████████████████████] 100%")
	}
	
	percent := float64(current) / float64(total)
	if percent > 1.0 {
		percent = 1.0
	}
	
	width := 30
	filled := int(percent * float64(width))
	
	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += fmt.Sprintf("] %.1f%%", percent*100)
	
	return pterm.NewStyle(pterm.FgCyan).Sprint(bar)
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

// ShowBanner displays an enhanced enterprise banner with the installer information
func ShowBanner(version string) {
	// Create enterprise banner using simple text styling
	banner := pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("╔══════════════════════════════════════╗\n") +
		pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("║    ") + pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold).Sprint("KUBERNETES INSTALLER") + pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("        ║\n") +
		pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("║        ") + pterm.NewStyle(pterm.FgYellow).Sprint("Enterprise Edition") + pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("         ║\n") +
		pterm.NewStyle(pterm.FgCyan, pterm.Bold).Sprint("╚══════════════════════════════════════╝")

	pterm.DefaultCenter.Println(banner)
	
	// Enterprise subtitle
	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(
		pterm.NewStyle(pterm.FgLightMagenta, pterm.Bold).Sprint("Enterprise Kubernetes Installation Platform") + "\n" +
		pterm.NewStyle(pterm.FgGray).Sprintf("Version: %s | Build: Enterprise", version) + "\n" +
		pterm.NewStyle(pterm.FgGray).Sprintf("Runtime: %s", time.Now().Format("2006-01-02 15:04:05 MST")))
	
	// Add separator
	pterm.Println()
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgCyan).Sprint("═══════════════════════════════════════"))
	pterm.Println()
}

// ShowEnterpriseWelcome displays a comprehensive enterprise welcome screen
func ShowEnterpriseWelcome(version string, environment string) {
	ShowBanner(version)
	
	// Environment information
	pterm.DefaultSection.Println("Environment Information")
	
	info := [][]string{
		{"Environment", environment},
		{"Installer Version", version},
		{"Build Type", "Enterprise"},
		{"Session ID", fmt.Sprintf("k8s-%d", time.Now().Unix())},
		{"Start Time", time.Now().Format("2006-01-02 15:04:05 MST")},
	}
	
	pterm.DefaultTable.WithHasHeader().WithData(
		append([][]string{{"Property", "Value"}}, info...),
	).Render()
	
	pterm.Println()
}

// ShowSummary displays an enhanced installation summary with enterprise metrics
func ShowSummary(steps []string, results map[string]string, duration time.Duration) {
	pterm.DefaultSection.Println("🏢 Enterprise Installation Summary")

	successCount := 0
	failedCount := 0
	skippedCount := 0
	warningCount := 0

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
			skippedCount++
		case "warning":
			symbol = "⚠️"
			color = pterm.FgYellow
			warningCount++
		default:
			symbol = "❓"
			color = pterm.FgLightWhite
		}

		pterm.Printf("  %s %-30s %s\n",
			symbol,
			step,
			pterm.NewStyle(color, pterm.Bold).Sprintf("%-10s", strings.ToUpper(result)))
	}

	pterm.Println()

	// Enhanced summary statistics with enterprise metrics
	totalSteps := len(steps)
	successRate := 0.0
	if totalSteps > 0 {
		successRate = float64(successCount) / float64(totalSteps) * 100
	}
	
	// Create metrics table
	metrics := [][]string{
		{"Total Steps", fmt.Sprintf("%d", totalSteps)},
		{"Successful", fmt.Sprintf("%d (%.1f%%)", successCount, successRate)},
		{"Failed", fmt.Sprintf("%d", failedCount)},
		{"Skipped", fmt.Sprintf("%d", skippedCount)},
		{"Warnings", fmt.Sprintf("%d", warningCount)},
		{"Duration", formatDuration(duration)},
		{"Success Rate", fmt.Sprintf("%.1f%%", successRate)},
	}
	
	if duration.Seconds() > 0 {
		throughput := float64(successCount) / duration.Seconds()
		metrics = append(metrics, []string{"Throughput", fmt.Sprintf("%.2f steps/sec", throughput)})
	}
	
	pterm.DefaultTable.WithHasHeader().WithData(
		append([][]string{{"Metric", "Value"}}, metrics...),
	).Render()

	pterm.Println()

	// Final status with enterprise styling
	if failedCount == 0 {
		pterm.DefaultBox.WithTitle("🎉 Installation Status").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgGreen)).
			Println(pterm.Green("✅ INSTALLATION COMPLETED SUCCESSFULLY\n\n") +
				pterm.LightGreen(fmt.Sprintf("All %d steps completed in %s", successCount, formatDuration(duration))))
	} else {
		pterm.DefaultBox.WithTitle("❌ Installation Status").
			WithTitleTopCenter().
			WithBoxStyle(pterm.NewStyle(pterm.FgRed)).
			Println(pterm.Red("❌ INSTALLATION COMPLETED WITH ERRORS\n\n") +
				pterm.LightRed(fmt.Sprintf("%d steps failed out of %d total", failedCount, totalSteps)))
	}
	
	// Add enterprise footer
	pterm.Println()
	pterm.DefaultCenter.Println(pterm.NewStyle(pterm.FgGray).Sprint("Enterprise Kubernetes Installer - Powered by Go"))
}
