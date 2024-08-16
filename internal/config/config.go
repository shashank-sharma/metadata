package config

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"

	"fyne.io/fyne/v2"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/settings"
)

//go:embed config.json
var embeddedFiles embed.FS

// AppConfig holds the application configuration.
type AppConfig struct {
	BackendEndpoint string
	AWEndpoint      string
	Settings        *settings.Settings
	SettingsManager *settings.SettingsManager
	Debug           bool
	Prod            bool
}

func LoadConfig(uri fyne.URI) (AppConfig, error) {
	config := AppConfig{
		BackendEndpoint: "http://localhost:8090/",
		AWEndpoint:      "http://localhost:5600/",
		Debug:           true,
		Prod:            false,
	}
	settingsManager := &settings.SettingsManager{StorageRoot: uri}
	appSettings, err := initSettings(settingsManager)
	if err != nil {
		return AppConfig{}, err
	}

	fileData, err := embeddedFiles.ReadFile("config.json")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			logger.LogWarning("config.json not found, using default configuration")
		} else {
			logger.LogError("error reading config.json:", err)
			return AppConfig{}, err
		}
	} else {
		if err := json.Unmarshal(fileData, &config); err != nil {
			logger.LogError("error parsing config.json:", err)
			return AppConfig{}, err
		}
	}

	config.Settings = &appSettings
	config.SettingsManager = settingsManager
	return config, nil
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
