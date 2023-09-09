package client

import (
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Version   string
	APIKey    string
	APIHost   string
	ApiScheme string
	Debug     bool
}

type Client struct {
	config     *Config
	httpClient *http.Client
}

func NewClient(config *Config) *Client {
	dCfg := defaultConfig()
	dCfg.merge(config)
	httpClient := &http.Client{
		Transport: &AddHeaderTransport{
			Transport: http.DefaultTransport,
			config:    dCfg,
		},
	}
	return &Client{
		config,
		httpClient,
	}
}

func defaultConfig() *Config {
	apiKey := os.Getenv("BASELIMEIO_API_KEY")
	return &Config{
		"0.0.1",
		apiKey,
		"go.baselime.io",
		"https",
		false,
	}
}

func (cfg *Config) merge(cfg2 *Config) {
	if cfg2.APIKey != "" {
		cfg.APIKey = cfg2.APIKey
	}
	if cfg2.APIHost != "" {
		cfg.APIHost = cfg2.APIHost
	}
	if cfg2.ApiScheme != "" {
		cfg.ApiScheme = cfg2.ApiScheme
	}
	cfg.Debug = cfg2.Debug
}

type AddHeaderTransport struct {
	Transport http.RoundTripper
	config    *Config
}

func (adt *AddHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = adt.config.ApiScheme
	req.URL.Host = adt.config.APIHost
	req.Header.Add("x-api-key", adt.config.APIKey)
	req.Header.Add("User-Agent", fmt.Sprintf("baselime-io-terraform-provider/%s", adt.config.Version))
	req.Header.Add("Content-Type", "application/json")
	return adt.Transport.RoundTrip(req)
}
