package config

import (
	"encoding/json"
	"time"
)

// InstallerConfig represents the complete configuration for the K8s installer
type InstallerConfig struct {
	Installer      InstallerSettings    `json:"installer" validate:"required"`
	Artifacts      ArtifactsConfig      `json:"artifacts" validate:"required"`
	Infrastructure InfrastructureConfig `json:"infrastructure,omitempty"`
	Database       DatabaseConfig       `json:"database,omitempty"`
	Deployment     DeploymentConfig     `json:"deployment,omitempty"`
	Validation     ValidationConfig     `json:"validation,omitempty"`
	Monitoring     MonitoringConfig     `json:"monitoring,omitempty"`
	Security       SecurityConfig       `json:"security,omitempty"`
	Kubernetes     K8sConfig            `json:"kubernetes,omitempty"`
	Cloud          CloudConfig          `json:"cloud,omitempty"`
}

// InstallerSettings contains general installer configuration
type InstallerSettings struct {
	Version   string `json:"version" validate:"required,semver"`
	Workspace string `json:"workspace" validate:"required,dir"`
	Verbose   bool   `json:"verbose"`
	DryRun    bool   `json:"dryRun"`
	LogLevel  string `json:"logLevel" validate:"oneof=debug info warn error"`
	LogFormat string `json:"logFormat" validate:"oneof=json text"`
}

// ArtifactsConfig handles OCI images, Helm charts, and Terraform modules
type ArtifactsConfig struct {
	Images    ImageConfig     `json:"images"`
	Helm      HelmConfig      `json:"helm"`
	Terraform TerraformConfig `json:"terraform"`
}

// ImageConfig manages OCI image synchronization
type ImageConfig struct {
	SkipPull bool             `json:"skipPull"`
	Vendor   RegistryConfig   `json:"vendor" validate:"required"`
	Client   RegistryConfig   `json:"client"`
	Images   []ImageReference `json:"images" validate:"required,min=1,dive"`
}

// RegistryConfig contains registry authentication and settings
type RegistryConfig struct {
	Registry       string    `json:"registry" validate:"required,url"`
	URL            string    `json:"url" validate:"omitempty,url"` // Alias for Registry for backward compatibility
	Auth           AuthConfig `json:"auth"`
	EnablePipeline bool      `json:"enablePipeline"`
	Insecure       bool      `json:"insecure"`
	Timeout        string    `json:"timeout" validate:"duration"`
}

// AuthConfig supports multiple authentication methods
type AuthConfig struct {
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	KeyFile  string `json:"keyFile,omitempty" validate:"omitempty,file"`
}

// ImageReference defines an OCI image to be managed
type ImageReference struct {
	Name        string            `json:"name" validate:"required"`
	Version     string            `json:"version" validate:"required"`
	Required    bool              `json:"required"`
	PullPolicy  string            `json:"pullPolicy" validate:"oneof=Always IfNotPresent Never"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// HelmConfig manages Helm chart repositories and synchronization
type HelmConfig struct {
	Vendor GitRepoConfig `json:"vendor" validate:"required"`
	Client GitRepoConfig `json:"client"`
	Charts []HelmChart   `json:"charts,omitempty"`
}

// TerraformConfig manages Terraform module repositories
type TerraformConfig struct {
	Vendor  GitRepoConfig     `json:"vendor" validate:"required"`
	Client  GitRepoConfig     `json:"client"`
	Modules []TerraformModule `json:"modules,omitempty"`
}

// GitRepoConfig contains Git repository configuration
type GitRepoConfig struct {
	Repo       string    `json:"repo" validate:"required,url"`
	Branch     string    `json:"branch"`
	Tag        string    `json:"tag"`
	Auth       AuthConfig `json:"auth"`
	PushToRepo bool      `json:"pushToRepo"`
	LocalPath  string    `json:"localPath,omitempty"`
}

// HelmChart defines a Helm chart configuration
type HelmChart struct {
	Name      string            `json:"name" validate:"required"`
	Path      string            `json:"path" validate:"required"`
	Version   string            `json:"version"`
	Values    map[string]interface{} `json:"values,omitempty"`
	Override  string            `json:"override,omitempty" validate:"omitempty,file"`
}

// TerraformModule defines a Terraform module configuration
type TerraformModule struct {
	Name      string                 `json:"name" validate:"required"`
	Path      string                 `json:"path" validate:"required"`
	Variables map[string]interface{} `json:"variables,omitempty"`
	Outputs   []string               `json:"outputs,omitempty"`
}

// InfrastructureConfig manages infrastructure provisioning
type InfrastructureConfig struct {
	Terraform   TerraformExecution `json:"terraform"`
	HealthCheck HealthCheckConfig  `json:"healthCheck"`
}

// TerraformExecution contains Terraform execution settings
type TerraformExecution struct {
	Enabled        bool              `json:"enabled"`
	Modules        []string          `json:"modules" validate:"required_if=Enabled true,min=1"`
	Workspace      string            `json:"workspace"`
	VarFiles       []string          `json:"varFiles,omitempty" validate:"dive,file"`
	Variables      map[string]string `json:"variables,omitempty"`
	ValidateHealth bool              `json:"validateHealth"`
	AutoApprove    bool              `json:"autoApprove"`
	Parallelism    int               `json:"parallelism" validate:"min=1,max=100"`
	Timeout        string            `json:"timeout" validate:"duration"`
}

// DatabaseConfig manages database initialization and migration
type DatabaseConfig struct {
	Enabled            bool               `json:"enabled"`
	RunAsInitContainer bool               `json:"runAsInitContainer"`
	Scripts            GitRepoConfig      `json:"scripts" validate:"required_if=Enabled true"`
	Connection         DatabaseConnection `json:"connection"`
	Validation         DatabaseValidation `json:"validation"`
	Migration          MigrationConfig    `json:"migration"`
}

// DatabaseConnection contains database connection details
type DatabaseConnection struct {
	Host     string `json:"host" validate:"required_if=Enabled true"`
	Port     int    `json:"port" validate:"required_if=Enabled true,min=1,max=65535"`
	Database string `json:"database" validate:"required_if=Enabled true"`
	Username string `json:"username" validate:"required_if=Enabled true"`
	Password string `json:"password" validate:"required_if=Enabled true"`
	SSLMode  string `json:"sslMode" validate:"oneof=disable require verify-ca verify-full"`
	Timeout  string `json:"timeout" validate:"duration"`
}

// DatabaseValidation contains database validation settings
type DatabaseValidation struct {
	Enabled     bool   `json:"enabled"`
	HealthCheck string `json:"healthCheck"`
	Timeout     string `json:"timeout" validate:"duration"`
	Retries     int    `json:"retries" validate:"min=0,max=10"`
}

// MigrationConfig contains database migration settings
type MigrationConfig struct {
	Path     string `json:"path" validate:"required_if=Enabled true"`
	Tool     string `json:"tool" validate:"oneof=flyway liquibase custom"`
	Baseline bool   `json:"baseline"`
	DryRun   bool   `json:"dryRun"`
	Timeout  string `json:"timeout" validate:"duration"`
}

// DeploymentConfig manages application deployment
type DeploymentConfig struct {
	Helm       HelmDeployment   `json:"helm"`
	Kubernetes K8sConfig        `json:"kubernetes"`
	Validation DeployValidation `json:"validation"`
}

// HelmDeployment contains Helm deployment configuration
type HelmDeployment struct {
	Charts          []DeployChart `json:"charts" validate:"min=1,dive"`
	CreateNamespace bool          `json:"createNamespace"`
	Wait            bool          `json:"wait"`
	Timeout         string        `json:"timeout" validate:"duration"`
	Atomic          bool          `json:"atomic"`
	CleanupOnFail   bool          `json:"cleanupOnFail"`
}

// DeployChart defines a chart to be deployed
type DeployChart struct {
	Name        string                 `json:"name" validate:"required"`
	Path        string                 `json:"path" validate:"required"`
	Namespace   string                 `json:"namespace" validate:"required"`
	Order       int                    `json:"order" validate:"min=1"`
	Values      map[string]interface{} `json:"values,omitempty"`
	ValuesFile  string                 `json:"valuesFile,omitempty" validate:"omitempty,file"`
	HealthCheck HealthCheckConfig      `json:"healthCheck"`
	DependsOn   []string               `json:"dependsOn,omitempty"`
}

// K8sConfig contains Kubernetes-specific settings
type K8sConfig struct {
	Context     string `json:"context"`
	Namespace   string `json:"namespace"`
	ConfigPath  string `json:"configPath" validate:"omitempty,file"`
	Timeout     string `json:"timeout" validate:"duration"`
	WaitTimeout string `json:"waitTimeout" validate:"duration"`
	
	// RBAC configuration
	RBAC struct {
		Enabled     bool     `json:"enabled"`
		Roles       []string `json:"roles"`
		Bindings    []string `json:"bindings"`
	} `json:"rbac"`
	
	// Ingress configuration
	Ingress struct {
		Enabled     bool   `json:"enabled"`
		Controller  string `json:"controller"`
		Namespace   string `json:"namespace"`
		Class       string `json:"class"`
	} `json:"ingress"`
	
	// Networking configuration
	Networking struct {
		CNI     string `json:"cni"`
		CIDR    string `json:"cidr"`
		Config  map[string]interface{} `json:"config"`
	} `json:"networking"`
	
	// Storage configuration
	Storage struct {
		Class       string `json:"class"`
		Provisioner string `json:"provisioner"`
		Config      map[string]interface{} `json:"config"`
	} `json:"storage"`
}

// DeployValidation contains deployment validation settings
type DeployValidation struct {
	PodHealth     bool                `json:"podHealth"`
	ServiceHealth bool                `json:"serviceHealth"`
	HealthChecks  []HealthCheckConfig `json:"healthChecks,omitempty"`
	CustomChecks  []CustomValidation  `json:"customChecks,omitempty"`
	Timeout       string              `json:"timeout" validate:"duration"`
	RetryInterval string              `json:"retryInterval" validate:"duration"`
}

// HealthCheckConfig defines health check parameters
type HealthCheckConfig struct {
	URL             string            `json:"url" validate:"url"`
	Method          string            `json:"method" validate:"oneof=GET POST PUT HEAD"`
	Headers         map[string]string `json:"headers,omitempty"`
	ExpectedStatus  int               `json:"expectedStatus" validate:"min=100,max=599"`
	ExpectedContent string            `json:"expectedContent,omitempty"`
	Timeout         string            `json:"timeout" validate:"duration"`
	Retries         int               `json:"retries" validate:"min=0,max=10"`
	Interval        string            `json:"interval" validate:"duration"`
}

// CustomValidation defines custom validation scripts
type CustomValidation struct {
	Name         string   `json:"name" validate:"required"`
	Script       string   `json:"script" validate:"required"`
	Args         []string `json:"args,omitempty"`
	Timeout      string   `json:"timeout" validate:"duration"`
	ExpectedExit int      `json:"expectedExit"`
}

// ValidationConfig manages post-deployment validation
type ValidationConfig struct {
	Post PostValidation `json:"post"`
	E2E  E2EConfig      `json:"e2e"`
}

// PostValidation contains post-deployment validation settings
type PostValidation struct {
	Scripts      []ScriptConfig      `json:"scripts,omitempty"`
	HealthChecks []HealthCheckConfig `json:"healthChecks,omitempty"`
	CustomChecks []CustomValidation  `json:"customChecks,omitempty"`
	Parallel     bool                `json:"parallel"`
	Timeout      string              `json:"timeout" validate:"duration"`
}

// E2EConfig contains end-to-end testing configuration
type E2EConfig struct {
	Enabled   bool          `json:"enabled"`
	TestSuite string        `json:"testSuite" validate:"required_if=Enabled true"`
	Framework string        `json:"framework" validate:"oneof=pytest junit go-test custom"`
	Config    E2ETestConfig `json:"config"`
	Reporting ReportConfig  `json:"reporting"`
	Timeout   string        `json:"timeout" validate:"duration"`
}

// E2ETestConfig contains test execution configuration
type E2ETestConfig struct {
	Environment map[string]string `json:"environment,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Parallel    bool              `json:"parallel"`
	Workers     int               `json:"workers" validate:"min=1,max=20"`
	Retries     int               `json:"retries" validate:"min=0,max=5"`
}

// ReportConfig contains test reporting configuration
type ReportConfig struct {
	Format    string `json:"format" validate:"oneof=junit xml json html"`
	Output    string `json:"output" validate:"required"`
	Archive   bool   `json:"archive"`
	Upload    bool   `json:"upload"`
	UploadURL string `json:"uploadUrl,omitempty" validate:"omitempty,url"`
}

// SecurityConfig defines security configuration
type SecurityConfig struct {
	ScanImages        bool     `json:"scanImages"`
	PolicyFiles       []string `json:"policyFiles"`
	AllowedRegistries []string `json:"allowedRegistries"`
	RequiredLabels    map[string]string `json:"requiredLabels"`
	
	// Scanning configuration
	Scanning struct {
		Enabled   bool     `json:"enabled"`
		Tools     []string `json:"tools"`
		Registries []string `json:"registries"`
	} `json:"scanning"`
	
	// Authentication configuration
	Authentication struct {
		Enabled bool   `json:"enabled"`
		Method  string `json:"method"`
		Config  map[string]interface{} `json:"config"`
	} `json:"authentication"`
	
	// Encryption configuration
	Encryption struct {
		Enabled bool   `json:"enabled"`
		Method  string `json:"method"`
		KeyPath string `json:"keyPath"`
	} `json:"encryption"`
	
	// Security policies
	Policies struct {
		Enabled bool     `json:"enabled"`
		Files   []string `json:"files"`
	} `json:"policies"`
}

// MonitoringConfig defines monitoring configuration
type MonitoringConfig struct {
	Enabled    bool   `json:"enabled"`
	Namespace  string `json:"namespace"`
	
	// Prometheus configuration
	Prometheus struct {
		Enabled   bool   `json:"enabled"`
		Endpoint  string `json:"endpoint"`
		Namespace string `json:"namespace"`
		Version   string `json:"version"`
	} `json:"prometheus"`
	
	// Grafana configuration
	Grafana struct {
		Enabled   bool   `json:"enabled"`
		Endpoint  string `json:"endpoint"`
		Namespace string `json:"namespace"`
		Version   string `json:"version"`
	} `json:"grafana"`
	
	// ELK Stack configuration
	ELK struct {
		Enabled       bool   `json:"enabled"`
		Elasticsearch string `json:"elasticsearch"`
		Logstash      string `json:"logstash"`
		Kibana        string `json:"kibana"`
	} `json:"elk"`
	
	// Alerting configuration
	Alerting struct {
		Enabled bool     `json:"enabled"`
		Rules   []string `json:"rules"`
	} `json:"alerting"`
}

// CloudConfig defines cloud provider configuration
type CloudConfig struct {
	Provider string `json:"provider" validate:"required,oneof=aws azure gcp"`
	Region   string `json:"region" validate:"required"`
	
	// AWS specific
	AWS struct {
		AccessKeyID     string `json:"accessKeyId"`
		SecretAccessKey string `json:"secretAccessKey"`
		SessionToken    string `json:"sessionToken"`
		Profile         string `json:"profile"`
	} `json:"aws,omitempty"`
	
	// Azure specific
	Azure struct {
		TenantID       string `json:"tenantId"`
		ClientID       string `json:"clientId"`
		ClientSecret   string `json:"clientSecret"`
		SubscriptionID string `json:"subscriptionId"`
	} `json:"azure,omitempty"`
	
	// GCP specific
	GCP struct {
		ProjectID   string `json:"projectId"`
		ServiceAccountKey string `json:"serviceAccountKey"`
	} `json:"gcp,omitempty"`
}

// Alias for backward compatibility
type Config = InstallerConfig
type KubernetesConfig = K8sConfig
type ScriptConfig struct {
	Name    string            `json:"name" validate:"required"`
	Path    string            `json:"path" validate:"required,file"`
	Args    []string          `json:"args,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Timeout string            `json:"timeout" validate:"duration"`
	WorkDir string            `json:"workDir,omitempty"`
	Shell   string            `json:"shell" validate:"oneof=bash sh zsh fish"`
}

// InstallState tracks the state of installation steps
type InstallState struct {
	Steps     []StepState `json:"steps"`
	StartTime time.Time   `json:"startTime"`
	EndTime   *time.Time  `json:"endTime,omitempty"`
	Status    string      `json:"status" validate:"oneof=pending running completed failed paused"`
	LastError string      `json:"lastError,omitempty"`
	Resume    bool        `json:"resume"`
}

// StepState tracks individual step execution state
type StepState struct {
	Name      string     `json:"name" validate:"required"`
	Status    string     `json:"status" validate:"oneof=pending running completed failed skipped"`
	StartTime *time.Time `json:"startTime,omitempty"`
	EndTime   *time.Time `json:"endTime,omitempty"`
	Error     string     `json:"error,omitempty"`
	Output    string     `json:"output,omitempty"`
	Retries   int        `json:"retries"`
}



// ToJSON converts the config to JSON string
func (c *InstallerConfig) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON loads config from JSON string
func (c *InstallerConfig) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), c)
}

// GenerateDefaultConfig creates a sample configuration
func GenerateDefaultConfig() *InstallerConfig {
	return &InstallerConfig{
		Installer: InstallerSettings{
			Version:   "1.0.0",
			Workspace: "./workspace",
			Verbose:   false,
			DryRun:    false,
			LogLevel:  "info",
			LogFormat: "text",
		},
		Artifacts: ArtifactsConfig{
			Images: ImageConfig{
				SkipPull: false,
				Vendor: RegistryConfig{
					Registry: "ghcr.io/vendor",
					Auth:     AuthConfig{Token: "vendor_token_here"},
					Timeout:  "30s",
				},
				Client: RegistryConfig{
					Registry:       "acr.client.com",
					Auth:           AuthConfig{Username: "client_user", Password: "client_pass"},
					EnablePipeline: true,
					Timeout:        "60s",
				},
				Images: []ImageReference{
					{
						Name:       "app-backend",
						Version:    "v1.2.3",
						Required:   true,
						PullPolicy: "IfNotPresent",
					},
					{
						Name:       "app-frontend",
						Version:    "v1.2.3",
						Required:   true,
						PullPolicy: "IfNotPresent",
					},
				},
			},
			Helm: HelmConfig{
				Vendor: GitRepoConfig{
					Repo:   "https://github.com/vendor/helm-charts",
					Branch: "main",
					Auth:   AuthConfig{Token: "vendor_github_token"},
				},
				Client: GitRepoConfig{
					Repo:       "https://github.com/client/helm-charts",
					PushToRepo: true,
					Auth:       AuthConfig{Token: "client_github_token"},
				},
			},
			Terraform: TerraformConfig{
				Vendor: GitRepoConfig{
					Repo: "https://github.com/vendor/terraform-modules",
					Tag:  "v1.0.0",
					Auth: AuthConfig{Token: "vendor_github_token"},
				},
				Client: GitRepoConfig{
					Repo:       "https://github.com/client/terraform-modules",
					PushToRepo: false,
				},
			},
		},
		Infrastructure: InfrastructureConfig{
			Terraform: TerraformExecution{
				Enabled:        true,
				Modules:        []string{"cluster", "database", "networking"},
				Workspace:      "production",
				ValidateHealth: true,
				AutoApprove:    false,
				Parallelism:    10,
				Timeout:        "30m",
			},
		},
		Database: DatabaseConfig{
			Enabled:            true,
			RunAsInitContainer: true,
			Scripts: GitRepoConfig{
				Repo: "https://github.com/vendor/db-scripts",
				Tag:  "v1.0.0",
				Auth: AuthConfig{Token: "vendor_github_token"},
			},
			Validation: DatabaseValidation{
				Enabled:     true,
				HealthCheck: "SELECT 1",
				Timeout:     "30s",
				Retries:     3,
			},
		},
		Deployment: DeploymentConfig{
			Helm: HelmDeployment{
				Charts: []DeployChart{
					{
						Name:      "backend",
						Path:      "./charts/backend",
						Namespace: "app",
						Order:     1,
						HealthCheck: HealthCheckConfig{
							URL:            "http://backend:8080/health",
							Method:         "GET",
							ExpectedStatus: 200,
							Timeout:        "30s",
							Retries:        5,
							Interval:       "10s",
						},
					},
				},
				CreateNamespace: true,
				Wait:            true,
				Timeout:         "10m",
				Atomic:          true,
				CleanupOnFail:   true,
			},
			Kubernetes: K8sConfig{
				Namespace:   "default",
				Timeout:     "5m",
				WaitTimeout: "10m",
			},
		},
		Validation: ValidationConfig{
			Post: PostValidation{
				Scripts: []ScriptConfig{
					{
						Name:    "post-deploy",
						Path:    "./scripts/post-deploy.sh",
						Timeout: "5m",
						Shell:   "bash",
					},
				},
				HealthChecks: []HealthCheckConfig{
					{
						URL:            "http://app.example.com/health",
						Method:         "GET",
						ExpectedStatus: 200,
						Timeout:        "30s",
						Retries:        3,
						Interval:       "10s",
					},
				},
				Timeout: "15m",
			},
			E2E: E2EConfig{
				Enabled:   true,
				TestSuite: "./tests/e2e",
				Framework: "go-test",
				Config: E2ETestConfig{
					Parallel: false,
					Workers:  1,
					Retries:  1,
				},
				Reporting: ReportConfig{
					Format:  "junit",
					Output:  "./reports/e2e-results.xml",
					Archive: true,
				},
				Timeout: "30m",
			},
		},
	}
}