package config

// SettingsConfig struct represents the config for the settings.
type SettingsConfig struct {
	Logging bool   `yaml:"logging"`
	Token   string `yaml:"token"`
}

// config represents the main config for the application.
type config struct {
	Settings SettingsConfig `yaml:"settings"`
}

var Config = ConfigParser{}.getDefaultConfig()
