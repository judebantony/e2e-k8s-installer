package cmd

import (
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

func loadDeployConfig(configPath string) (*config.DeploymentConfig, error) {
	config := &config.DeploymentConfig{
		Helm: config.HelmDeployment{
			Charts: []config.DeployChart{
				{Name: "backend", Path: "./charts/backend", Namespace: "app", Order: 1},
				{Name: "frontend", Path: "./charts/frontend", Namespace: "app", Order: 2},
			},
			CreateNamespace: true,
			Wait:            true,
			Timeout:         "10m",
			Atomic:          true,
			CleanupOnFail:   true,
		},
		Kubernetes: config.K8sConfig{
			Namespace:   "default",
			Timeout:     "5m",
			WaitTimeout: "10m",
		},
		Validation: config.DeployValidation{
			PodHealth:     true,
			ServiceHealth: true,
			Timeout:       "5m",
			RetryInterval: "30s",
		},
	}
	return config, nil
}
