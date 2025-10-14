package monitoring

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

// Manager handles monitoring and logging operations
type Manager struct {
	config config.MonitoringConfig
	logger *logrus.Logger
}

// NewManager creates a new monitoring manager
func NewManager(cfg config.MonitoringConfig, logger *logrus.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

// InstallStack installs the complete monitoring stack
func (m *Manager) InstallStack(ctx context.Context) error {
	m.logger.Info("📊 Installing monitoring and logging stack")

	// Install Prometheus if enabled
	if m.config.Prometheus.Enabled {
		if err := m.installPrometheus(ctx); err != nil {
			return fmt.Errorf("prometheus installation failed: %w", err)
		}
	}

	// Install Grafana if enabled
	if m.config.Grafana.Enabled {
		if err := m.installGrafana(ctx); err != nil {
			return fmt.Errorf("grafana installation failed: %w", err)
		}
	}

	// Install ELK stack if enabled
	if m.config.ELK.Enabled {
		if err := m.installELKStack(ctx); err != nil {
			return fmt.Errorf("ELK stack installation failed: %w", err)
		}
	}

	// Setup alerting if enabled
	if m.config.Alerting.Enabled {
		if err := m.setupAlerting(ctx); err != nil {
			return fmt.Errorf("alerting setup failed: %w", err)
		}
	}

	m.logger.Info("✅ Monitoring stack installation completed successfully")
	return nil
}

// installPrometheus installs Prometheus monitoring
func (m *Manager) installPrometheus(ctx context.Context) error {
	m.logger.Info("📊 Installing Prometheus")

	// Install Prometheus using Helm
	cmd := exec.CommandContext(ctx, "helm", "install", "prometheus", "prometheus/kube-prometheus-stack",
		"--namespace", "monitoring",
		"--create-namespace",
		"--set", "prometheus.prometheusSpec.retention="+m.config.Prometheus.Retention,
		"--set", "prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage="+m.config.Prometheus.Storage)

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Prometheus installation failed")
		return fmt.Errorf("prometheus installation failed: %w", err)
	}

	// Apply custom Prometheus rules
	if err := m.applyPrometheusRules(ctx); err != nil {
		return fmt.Errorf("failed to apply Prometheus rules: %w", err)
	}

	m.logger.Info("✅ Prometheus installed successfully")
	return nil
}

// installGrafana installs Grafana dashboard
func (m *Manager) installGrafana(ctx context.Context) error {
	m.logger.Info("📈 Installing Grafana")

	// Install Grafana using Helm
	cmd := exec.CommandContext(ctx, "helm", "install", "grafana", "grafana/grafana",
		"--namespace", "monitoring",
		"--create-namespace",
		"--set", "adminPassword=admin123",
		"--set", "service.type=LoadBalancer")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Grafana installation failed")
		return fmt.Errorf("grafana installation failed: %w", err)
	}

	// Import custom dashboards
	if err := m.importGrafanaDashboards(ctx); err != nil {
		return fmt.Errorf("failed to import Grafana dashboards: %w", err)
	}

	m.logger.Info("✅ Grafana installed successfully")
	return nil
}

// installELKStack installs Elasticsearch, Logstash, and Kibana
func (m *Manager) installELKStack(ctx context.Context) error {
	m.logger.Info("📊 Installing ELK stack")

	// Install Elasticsearch
	if err := m.installElasticsearch(ctx); err != nil {
		return fmt.Errorf("elasticsearch installation failed: %w", err)
	}

	// Install Logstash
	if err := m.installLogstash(ctx); err != nil {
		return fmt.Errorf("logstash installation failed: %w", err)
	}

	// Install Kibana
	if err := m.installKibana(ctx); err != nil {
		return fmt.Errorf("kibana installation failed: %w", err)
	}

	// Install Filebeat for log collection
	if err := m.installFilebeat(ctx); err != nil {
		return fmt.Errorf("filebeat installation failed: %w", err)
	}

	m.logger.Info("✅ ELK stack installed successfully")
	return nil
}

// installElasticsearch installs Elasticsearch
func (m *Manager) installElasticsearch(ctx context.Context) error {
	m.logger.Info("🔍 Installing Elasticsearch")

	cmd := exec.CommandContext(ctx, "helm", "install", "elasticsearch", "elastic/elasticsearch",
		"--namespace", "logging",
		"--create-namespace",
		"--set", "replicas=3",
		"--set", "volumeClaimTemplate.resources.requests.storage=10Gi")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Elasticsearch installation failed")
		return fmt.Errorf("elasticsearch installation failed: %w", err)
	}

	return nil
}

// installLogstash installs Logstash
func (m *Manager) installLogstash(ctx context.Context) error {
	m.logger.Info("📊 Installing Logstash")

	cmd := exec.CommandContext(ctx, "helm", "install", "logstash", "elastic/logstash",
		"--namespace", "logging",
		"--create-namespace")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Logstash installation failed")
		return fmt.Errorf("logstash installation failed: %w", err)
	}

	return nil
}

// installKibana installs Kibana
func (m *Manager) installKibana(ctx context.Context) error {
	m.logger.Info("📊 Installing Kibana")

	cmd := exec.CommandContext(ctx, "helm", "install", "kibana", "elastic/kibana",
		"--namespace", "logging",
		"--create-namespace",
		"--set", "service.type=LoadBalancer")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Kibana installation failed")
		return fmt.Errorf("kibana installation failed: %w", err)
	}

	return nil
}

// installFilebeat installs Filebeat for log collection
func (m *Manager) installFilebeat(ctx context.Context) error {
	m.logger.Info("📊 Installing Filebeat")

	cmd := exec.CommandContext(ctx, "helm", "install", "filebeat", "elastic/filebeat",
		"--namespace", "logging",
		"--create-namespace")

	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Filebeat installation failed")
		return fmt.Errorf("filebeat installation failed: %w", err)
	}

	return nil
}

// setupAlerting configures alerting rules and notifications
func (m *Manager) setupAlerting(ctx context.Context) error {
	m.logger.Info("🚨 Setting up alerting")

	// Apply alerting rules
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/monitoring/alerts/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Alerting rules application failed")
		return fmt.Errorf("alerting rules application failed: %w", err)
	}

	// Configure notification channels
	for channel, config := range m.config.Alerting.Channels {
		if err := m.configureAlertChannel(ctx, channel, config); err != nil {
			return fmt.Errorf("failed to configure alert channel %s: %w", channel, err)
		}
	}

	m.logger.Info("✅ Alerting setup completed successfully")
	return nil
}

// applyPrometheusRules applies custom Prometheus rules
func (m *Manager) applyPrometheusRules(ctx context.Context) error {
	m.logger.Info("📊 Applying Prometheus rules")

	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/monitoring/prometheus-rules/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Prometheus rules application failed")
		return fmt.Errorf("prometheus rules application failed: %w", err)
	}

	return nil
}

// importGrafanaDashboards imports custom Grafana dashboards
func (m *Manager) importGrafanaDashboards(ctx context.Context) error {
	m.logger.Info("📈 Importing Grafana dashboards")

	// Apply dashboard ConfigMaps
	cmd := exec.CommandContext(ctx, "kubectl", "apply", "-f", "./deploy/manifests/monitoring/grafana-dashboards/")
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.WithField("output", string(output)).Error("Grafana dashboards import failed")
		return fmt.Errorf("grafana dashboards import failed: %w", err)
	}

	return nil
}

// configureAlertChannel configures a specific alert notification channel
func (m *Manager) configureAlertChannel(ctx context.Context, channel, config string) error {
	m.logger.WithFields(logrus.Fields{
		"channel": channel,
		"config":  config,
	}).Info("🚨 Configuring alert channel")

	// Implementation for specific alert channel configuration
	// This would vary based on the channel type (Slack, email, PagerDuty, etc.)
	
	return nil
}

// GetMetrics retrieves monitoring metrics
func (m *Manager) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	m.logger.Info("📊 Retrieving monitoring metrics")

	metrics := make(map[string]interface{})

	// Get Prometheus metrics
	if m.config.Prometheus.Enabled {
		prometheusMetrics, err := m.getPrometheusMetrics(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get Prometheus metrics: %w", err)
		}
		metrics["prometheus"] = prometheusMetrics
	}

	// Get Grafana status
	if m.config.Grafana.Enabled {
		grafanaStatus, err := m.getGrafanaStatus(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get Grafana status: %w", err)
		}
		metrics["grafana"] = grafanaStatus
	}

	// Get ELK stack status
	if m.config.ELK.Enabled {
		elkStatus, err := m.getELKStatus(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get ELK status: %w", err)
		}
		metrics["elk"] = elkStatus
	}

	return metrics, nil
}

// getPrometheusMetrics retrieves Prometheus metrics
func (m *Manager) getPrometheusMetrics(ctx context.Context) (map[string]interface{}, error) {
	// Implementation to query Prometheus metrics
	return map[string]interface{}{
		"status": "running",
		"targets": "healthy",
	}, nil
}

// getGrafanaStatus retrieves Grafana status
func (m *Manager) getGrafanaStatus(ctx context.Context) (map[string]interface{}, error) {
	// Implementation to check Grafana status
	return map[string]interface{}{
		"status": "running",
		"dashboards": len(m.config.Grafana.Dashboards),
	}, nil
}

// getELKStatus retrieves ELK stack status
func (m *Manager) getELKStatus(ctx context.Context) (map[string]interface{}, error) {
	// Implementation to check ELK stack status
	return map[string]interface{}{
		"elasticsearch": "running",
		"logstash": "running",
		"kibana": "running",
	}, nil
}