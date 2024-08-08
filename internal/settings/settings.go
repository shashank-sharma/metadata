package settings

import (
	"fyne.io/fyne/v2"
)

type BaseSettings interface {
	FileName() string
}

type UserSettings struct {
	UserId    string          `json:"userid"`
	Token     string          `json:"token"`
	ProductId string          `json:"productid"`
	Bucket    map[string]bool `json:"bucket"`
}

func (us *UserSettings) FileName() string {
	return "user_settings.json"
}

type ApplicationSettings struct {
	HostName        string
	OperatingSystem string
	Arch            string
	storagePath     fyne.URI
	userTheme       fyne.Theme
}

func (as *ApplicationSettings) FileName() string {
	return "application_settings.json"
}

type Settings struct {
	ApplicationSettings *ApplicationSettings
	UserSettings        *UserSettings
}
