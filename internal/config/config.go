// Package config provides configuration to the manager.
// It provides a globa configuration object (Settings) that should not change after application startup.
package config

import (
	"github.com/kelseyhightower/envconfig"
)

// ConfigEnvironmentPrefix is the prefix used for environment variable names
const ConfigEnvironmentPrefix = ""

var settings *Settings

// Settings provides all available configuration settings for the manager
// Each field can be configured via an environment prefixed with ConfigEnvironmentPrefix
// e.g. if the prefix is "GFG", GFG_FOO_SETTING maps to FooSetting in the struct
// See https://github.com/kelseyhightower/envconfig for implementation details
type Settings struct {
	// LogMode sets the log mode. Valid values are "auto", "dev", or "prod" (Default: "auto").
	LogMode         string `envconfig:"LOG_MODE" default:"auto"`
	Host            string `envconfig:"HOST" default:""`
	HTTPPort        string `envconfig:"PORT" default:"8923"`
	GRPCPort        string `envconfig:"GRPC_PORT" default:"8924"`
	HealthCheckPort string `envconfig:"HEALTHCHECK_PORT" default:"8925"`

	DatabaseURI      string `envconfig:"MONGO_URI" default:"mongodb://localhost:27017"`
	DatabaseName     string `envconfig:"MONGO_DATABASE_NAME" default:"couchconnections"`
	DatabaseUsername string `envconfig:"MONGO_USERNAME"`
	DatabasePassword string `envconfig:"MONGO_PASSWORD"`

	Auth0ClientID     string `envconfig:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `envconfig:"AUTH0_CLIENT_SECRET"`
	Auth0CallbackURL  string `envconfig:"AUTH0_CALLBACK_URL"`
	Auth0Domain       string `envconfig:"AUTH0_DOMAIN"`

	AuthJwksURL          string `envconfig:"AUTH_JWKS_CONFIG" default:"https://livingroompresentation.eu.auth0.com/.well-known/jwks.json"`
	AuthUserInfoEndpoint string `envconfig:"AUTH_USER_INFO_ENDPOINT" default:"https://livingroompresentation.eu.auth0.com/userinfo"`
}

// Init parses configuration from the environment. This should be called only once per application startup (typically in Main)
func init() {
	var loadSettings Settings
	err := envconfig.Process(ConfigEnvironmentPrefix, &loadSettings)
	if err == nil {
		settings = &loadSettings
	} else {
		panic("Error initializing configuration settings: " + err.Error())
	}
}

// Get returns the current configuration settings. Will panic if Init() is not called first.
func Get() Settings {
	if settings == nil {
		panic("Configuration settings not initialized.")
	}
	return *settings
}
