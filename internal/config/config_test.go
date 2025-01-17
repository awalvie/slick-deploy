package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigSuccess(t *testing.T) {
	// Set up a temporary YAML file with valid configuration data
	tempFile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(`
app:
  name: "Test App"
  image: "testapp/image"
  container_port: 8080
  env: ["ENV_VAR=VALUE"]
  port_range:
    start: 8000
    end: 9000
caddy:
  admin_api: "http://localhost:2019"
  rules: []
health_check:
  endpoint: "/health"
  timeout_seconds: 30
`)
	require.NoError(t, err)
	err = tempFile.Close()
	require.NoError(t, err)

	// Load the configuration from the temporary file
	config, err := LoadConfig(tempFile.Name())
	require.NoError(t, err)

	// Assert that the configuration values are as expected
	assert.Equal(t, "Test App", config.App.Name)
	assert.Equal(t, "testapp/image", config.App.ImageName)
	assert.Equal(t, 8080, config.App.ContainerPort)
	assert.Equal(t, "http://localhost:2019", config.Caddy.AdminAPI)
	assert.Equal(t, "/health", config.HealthCheck.Endpoint)
	assert.Equal(t, 30, config.HealthCheck.TimeoutSeconds)
}

func TestLoadConfigUnmarshallFail(t *testing.T) {
	// Set up a temporary YAML file with valid configuration data
	tempFile, err := os.CreateTemp("", "*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(`
{
	"app": "Thou shall fail",
}
`)
	require.NoError(t, err)
	err = tempFile.Close()
	require.NoError(t, err)

	// Load the configuration from the temporary file
	_, err = LoadConfig(tempFile.Name())
	assert.Error(t, err)
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yaml")
	assert.Error(t, err)
}

func TestLoadConfigDefaultValues(t *testing.T) {
	// Set up a temporary YAML file with minimal configuration data
	tempFile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(`
app:
  name: "Test App"
`)
	require.NoError(t, err)
	err = tempFile.Close()
	require.NoError(t, err)

	// Load the configuration from the temporary file
	config, err := LoadConfig(tempFile.Name())
	require.NoError(t, err)

	// Assert that the default values are set correctly
	assert.Equal(t, 8000, config.App.PortRange.Start)
	assert.Equal(t, 9000, config.App.PortRange.End)
	assert.Equal(t, "http://localhost:2019", config.Caddy.AdminAPI)
}

func TestLoadConfigRegistry(t *testing.T) {
	// Set up a temporary YAML file with valid configuration data
	tempFile, err := os.CreateTemp("", "*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())
	_, err = tempFile.WriteString(`
app:
  name: "Test App"
  image: "testapp/image"
  container_port: 8080
  registry:
    username: "testuser"
    password: TEST_REGISTRY_PASSWORD
`)
	require.NoError(t, err)
	err = tempFile.Close()
	require.NoError(t, err)

	// Set the environment variable
	os.Setenv("TEST_REGISTRY_PASSWORD", "testpassword")
	defer os.Unsetenv("TEST_REGISTRY_PASSWORD")

	// Load the configuration from the temporary file
	config, err := LoadConfig(tempFile.Name())
	require.NoError(t, err)

	// Assert that the configuration values are as expected
	assert.Equal(t, "testuser", config.App.Registry.Username)
	assert.Equal(t, "testpassword", config.App.Registry.Password)
}
