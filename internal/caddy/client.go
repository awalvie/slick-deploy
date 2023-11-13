package caddy

import "net/http"

type CaddyConfig struct {
	Apps AppConfig `json:"apps"`
}

// AppConfig represents the top-level app configuration.
// AppsConfig holds the apps configurations.
type AppConfig struct {
	HTTP HTTPConfig `json:"http"`
}

// HTTPConfig holds HTTP-specific configurations.
type HTTPConfig struct {
	Servers map[string]ServerConfig `json:"servers"`
}

// ServerConfig defines the configuration for a server.
type ServerConfig struct {
	Listen []string      `json:"listen"`
	Routes []RouteConfig `json:"routes"`
}

// RouteConfig represents the configuration for a route.
type RouteConfig struct {
	Handle   []HandlerConfig `json:"handle"`
	Match    []MatchConfig   `json:"match"`
	Terminal bool            `json:"terminal"`
}

// HandlerConfig defines a handler in the route.
type HandlerConfig struct {
	Handler   string           `json:"handler"`
	Routes    []RouteConfig    `json:"routes,omitempty"`
	Upstreams []UpstreamConfig `json:"upstreams,omitempty"`
}

// MatchConfig represents the match conditions for a route.
type MatchConfig struct {
	Host []string `json:"host,omitempty"`
	Path []string `json:"path,omitempty"`
}

// UpstreamConfig represents the configuration for an upstream server.
type UpstreamConfig struct {
	Dial string `json:"dial"`
}

type CaddyClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewCaddyClient(baseURL string) *CaddyClient {
	return &CaddyClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// func DefaultConfig() CaddyConfig {
// 	return CaddyConfig{
// 		Apps: AppConfig{
// 			HTTP: HTTPConfig{
// 				Servers: {},
// 			},
// 		},
// 	}
// }
