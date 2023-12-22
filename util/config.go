package util

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"github.com/shashank-sharma/metadata/settings"
)

var rootConfig *Config
var rootApp fyne.App
var rootSetting *settings.Settings
var rootWindow fyne.Window

type Storage struct {
	SuccessfulPing binding.Int
	FailedPing     binding.Int
}

type Config struct {
	ServiceHost string
	ServicePort string
	Debug       bool
	Prod        bool
	Storage     *Storage
}

func Init() {
	rootApp = app.NewWithID("dashboard-metadata")
	rootSetting = settings.NewSettings(rootApp)
	rootWindow = rootApp.NewWindow("Metadata Manager")
	rootConfig = &Config{}
	rootConfig.Storage = &Storage{}
	rootConfig.Prod = true
	if os.Getenv("PROD") != "" && strings.ToLower(os.Getenv("PROD")) == "false" {
		rootConfig.Prod = false
	}
	rootConfig.Storage.SuccessfulPing = binding.NewInt()
	rootConfig.Storage.SuccessfulPing.Set(0)
	rootConfig.Storage.FailedPing = binding.NewInt()
	rootConfig.Storage.FailedPing.Set(0)
	if !rootConfig.Prod {
		fmt.Println("Production false")
		rootConfig.ServiceHost = os.Getenv("SERVICE_HOST")
		rootConfig.ServicePort = os.Getenv("SERVICE_PORT")
	} else {
		fmt.Println("Production true")
		rootConfig.ServiceHost = "139.59.46.84"
		rootConfig.ServicePort = "80"
	}

	rootConfig.Debug = strings.ToLower(os.Getenv("DEBUG")) == "true"
	fmt.Println("config =", rootConfig)
}

func GetApp() fyne.App {
	return rootApp
}

func GetWindow() fyne.Window {
	return rootWindow
}

func GetSetting() *settings.Settings {
	return rootSetting
}

func GetStorage() *Storage {
	return rootConfig.Storage
}

func GetConfig() *Config {
	return rootConfig
}
