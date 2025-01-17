package config

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type PortRange struct {
	Start int `yaml:"start"`
	End   int `yaml:"end"`
}

type RegistryConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type App struct {
	Name          string         `yaml:"name"`
	ImageName     string         `yaml:"image"`
	Registry      RegistryConfig `yaml:"registry"`
	ContainerPort int            `yaml:"container_port"`
	Network       string         `yaml:"network"`
	ENV           []string       `yaml:"env"`
	PortRange     PortRange      `yaml:"port_range"`
}

type ReverseProxy struct {
	Path string `yaml:"path"`
	To   string `yaml:"to"`
}

type Rule struct {
	Match        string         `yaml:"match"`
	Tls          string         `yaml:"tls"`
	ReverseProxy []ReverseProxy `yaml:"reverse_proxy"`
}

type CaddyConfig struct {
	AdminAPI string `yaml:"admin_api"`
	Rules    []Rule `yaml:"rules"`
}

type HealthCheck struct {
	Endpoint        string `yaml:"endpoint"`
	TimeoutSeconds  int    `yaml:"timeout_seconds"`
	IntervalSeconds int    `yaml:"interval_seconds"`
	MaxRetries      int    `yaml:"max_retries"`
}

type DeploymentConfig struct {
	App         App         `yaml:"app"`
	Caddy       CaddyConfig `yaml:"caddy"`
	HealthCheck HealthCheck `yaml:"health_check"`
}

func replaceEnvVariables(input string) (string, error) {
	re := regexp.MustCompile(`\{env\.([a-zA-Z_][a-zA-Z0-9_]*)\}`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the variable name
		varName := strings.TrimPrefix(match, "{env.")
		varName = strings.TrimSuffix(varName, "}")
		envValue, exists := os.LookupEnv(varName)
		if exists {
			return envValue
		}

		return match
	}), nil
}

func LoadConfig(path string) (DeploymentConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return DeploymentConfig{}, fmt.Errorf("error opening config file: %v", err)
	}

	// close the file when we're done
	defer f.Close()

	// Read the file content
	data, _ := io.ReadAll(f)

	// Create a default deployment config
	c := DeploymentConfig{
		App: App{
			PortRange: PortRange{
				Start: 8000,
				End:   9000,
			},
		},
		Caddy: CaddyConfig{
			AdminAPI: "http://localhost:2019",
		},
		HealthCheck: HealthCheck{
			TimeoutSeconds:  5,
			IntervalSeconds: 5,
			MaxRetries:      3,
		},
	}

	// Override the default config with the config file
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return c, err
	}

	for i, rule := range c.Caddy.Rules {
		newTlsValue, err := replaceEnvVariables(rule.Tls)
		if err != nil {
			return c, err
		}

		c.Caddy.Rules[i].Tls = newTlsValue
	}

	if c.App.Registry.Username != "" && c.App.Registry.Password != "" {
		envValue, exists := os.LookupEnv(c.App.Registry.Password)
		if exists {
			c.App.Registry.Password = envValue
		}
	}

	return c, nil
}
