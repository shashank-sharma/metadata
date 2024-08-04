package config

import (
	"os"

	"fyne.io/fyne/v2"
	"github.com/shashank-sharma/metadata/internal/settings"
)

// AppConfig holds the application configuration.
type AppConfig struct {
	BackendEndpoint string
	AWEndpoint      string
	Settings        settings.Settings
	SettingsManager *settings.SettingsManager
	Debug           bool
	Prod            bool
}

func LoadConfig(uri fyne.URI) (AppConfig, error) {
	settingsManager := &settings.SettingsManager{StorageRoot: uri}
	appSettings, err := initSettings(settingsManager)
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{
		BackendEndpoint: getEnv("BACKEND_ENDPOINT", ""),
		AWEndpoint:      getEnv("AW_ENDPOINT", ""),
		Settings:        appSettings,
		SettingsManager: settingsManager,
	}, nil
}

func initSettings(settingsManager *settings.SettingsManager) (settings.Settings, error) {
	appSettings := &settings.ApplicationSettings{}
	err := settingsManager.InitializeSettings(appSettings)
	if err != nil {
		return settings.Settings{}, err
	}

	userSettings := &settings.UserSettings{}
	err = settingsManager.InitializeSettings(userSettings)
	if err != nil {
		return settings.Settings{}, err
	}

	return settings.Settings{
		ApplicationSettings: appSettings,
		UserSettings:        userSettings,
	}, nil
}

/*
func Get() *AppConfig {
	if globalConfig == nil {
		log.Fatal("Config has not been loaded")
	}
	return globalConfig
}
*/

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
