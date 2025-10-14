package cloud

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/judebantony/e2e-k8s-installer/pkg/config"
)

// Manager interface defines cloud provider operations
type Manager interface {
	ValidateAccess(ctx context.Context) error
	ProvisionInfrastructure(ctx context.Context) error
	DestroyInfrastructure(ctx context.Context) error
	GetClusterInfo(ctx context.Context) (*ClusterInfo, error)
}

// ClusterInfo contains information about the provisioned cluster
type ClusterInfo struct {
	ClusterName     string            `json:"clusterName"`
	Endpoint        string            `json:"endpoint"`
	Version         string            `json:"version"`
	NodeGroups      []NodeGroup       `json:"nodeGroups"`
	NetworkConfig   NetworkConfig     `json:"networkConfig"`
	SecurityGroups  []string          `json:"securityGroups"`
	Tags            map[string]string `json:"tags"`
}

// NodeGroup represents a group of worker nodes
type NodeGroup struct {
	Name         string `json:"name"`
	InstanceType string `json:"instanceType"`
	MinSize      int    `json:"minSize"`
	MaxSize      int    `json:"maxSize"`
	DesiredSize  int    `json:"desiredSize"`
}

// NetworkConfig represents network configuration
type NetworkConfig struct {
	VPC        string   `json:"vpc"`
	Subnets    []string `json:"subnets"`
	CIDR       string   `json:"cidr"`
	PublicAccess bool   `json:"publicAccess"`
}

// NewManager creates appropriate cloud manager based on provider
func NewManager(cfg config.CloudConfig, logger *logrus.Logger) (Manager, error) {
	switch cfg.Provider {
	case "aws":
		return NewAWSManager(cfg, logger)
	case "azure":
		return NewAzureManager(cfg, logger)
	case "gcp":
		return NewGCPManager(cfg, logger)
	case "onprem":
		return NewOnPremManager(cfg, logger)
	default:
		return nil, fmt.Errorf("unsupported cloud provider: %s", cfg.Provider)
	}
}