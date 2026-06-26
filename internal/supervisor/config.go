package supervisor

import (
	"fmt"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

// Config represents the supervisor configuration.
type Config struct {
	Server      *OpAMPServerConfig      `koanf:"server"`
	Agent       *AgentConfig            `koanf:"agent"`
	Artifactory *ArtifactoryConfig      `koanf:"artifactory"`
	API         *APIConfig              `koanf:"api"`
	Auth        *AuthConfig             `koanf:"auth"`
}

// OpAMPServerConfig configures the OpAMP server connection.
type OpAMPServerConfig struct {
	Endpoint string `koanf:"endpoint"`
}

// AgentConfig configures the collector agent.
type AgentConfig struct {
	Executable string `koanf:"executable"`
}

// ArtifactoryConfig configures the Artifactory connection for binary downloads.
type ArtifactoryConfig struct {
	BaseURL    string `koanf:"base_url"`
	RepoKey    string `koanf:"repo_key"`
	APIToken   string `koanf:"api_token"`
	InsecureTLS bool  `koanf:"insecure_tls"`
}

// APIConfig configures the management REST API.
type APIConfig struct {
	Listen string `koanf:"listen"`
}

// AuthConfig configures OIDC-based authentication.
type AuthConfig struct {
	Provider      string `koanf:"provider"`
	IssuerURL     string `koanf:"issuer_url"`
	ClientID      string `koanf:"client_id"`
	ClientSecret  string `koanf:"client_secret"`
	RedirectURL   string `koanf:"redirect_url"`
	AdminRole     string `koanf:"admin_role"`
	DeployerRole  string `koanf:"deployer_role"`
	ViewerRole    string `koanf:"viewer_role"`
	AdminGroup    string `koanf:"admin_group"`
	DeployerGroup string `koanf:"deployer_group"`
	ViewerGroup   string `koanf:"viewer_group"`
}

// LoadConfig loads the supervisor configuration from a file and environment variables.
func LoadConfig(configPath string) (*Config, error) {
	k := koanf.New(".")

	// Load from config file
	if configPath != "" {
		if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("cannot load config file %s: %w", configPath, err)
		}
	}

	// Override with environment variables
	if err := k.Load(env.Provider("OPAMP_", ".", func(s string) string {
		return s
	}), nil); err != nil {
		return nil, fmt.Errorf("cannot load env config: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config: %w", err)
	}

	// Set defaults
	if cfg.API.Listen == "" {
		cfg.API.Listen = ":8080"
	}

	// Override Artifactory API token from environment variable if set
	if token := os.Getenv("ARTIFACTORY_API_TOKEN"); token != "" {
		cfg.Artifactory.APIToken = token
	}

	return &cfg, nil
}