package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	systemThemeName = "system default"
)

type SettingStorage struct {
	UserId    string `json:"userid"`
	Token     string `json:"token"`
	ProductId string `json:"productid"`
}

type Settings struct {
	HostName        string
	OperatingSystem string
	Arch            string
	fyneSettings    app.SettingsSchema
	storagePath     fyne.URI
	userTheme       fyne.Theme
	settingStorage  SettingStorage
}

func NewSettings(app fyne.App) *Settings {
	s := &Settings{}
	s.storagePath, _ = storage.Child(app.Storage().RootURI(), "settings.json")
	hName, err := os.Hostname()
	if err != nil {
		log.Fatalf("error retrieving hostname: %v", err)
		hName = "unknown"
	}

	s.HostName = hName
	s.OperatingSystem = runtime.GOARCH
	s.Arch = runtime.GOARCH
	s.load()
	if s.fyneSettings.Scale == 0 {
		s.fyneSettings.Scale = 1
	}
	return s
}

func (s *Settings) GetToken() string {
	return s.settingStorage.Token
}

func (s *Settings) SetToken(token string) {
	s.settingStorage.Token = token
	s.save()
}

func (s *Settings) GetUserId() string {
	return s.settingStorage.UserId
}

func (s *Settings) SetUserId(userId string) {
	s.settingStorage.UserId = userId
	s.save()
}

func (s *Settings) GetProductId() string {
	return s.settingStorage.ProductId
}

func (s *Settings) SetProductId(productId string) {
	s.settingStorage.ProductId = productId
	s.save()
}

func (s *Settings) ShouldSync() bool {
	return s.settingStorage.ProductId != "" &&
		s.settingStorage.Token != "" &&
		s.settingStorage.UserId != ""
}

func (s *Settings) chooseTheme(name string) {
	if name == systemThemeName {
		name = ""
	}
	s.fyneSettings.ThemeName = name
}

func (s *Settings) load() {
	fmt.Println("Storange path", s.storagePath)
	err := s.loadFromFile(s.storagePath)
	if err != nil {
		fyne.LogError("Settings load error:", err)
	}
	fmt.Println("Loaded: ", s.settingStorage)
}

func (s *Settings) Reset() {
	path := s.storagePath
	v, err := storage.CanRead(path)
	if v == false || err != nil {
		fmt.Println("Error: ", err)
	}

	err = storage.Delete(path)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	return
}

func (s *Settings) loadFromFile(path fyne.URI) error {
	v, err := storage.CanRead(path)
	if v == false || err != nil {
		fmt.Println("Error E21: ", err)
		return err
	}
	reader, err := storage.Reader(path)
	if err != nil {
		fmt.Println("Error: E22: ", err)
		return err
	}
	defer reader.Close()

	decode := json.NewDecoder(reader)

	return decode.Decode(&s.settingStorage)
}

func (s *Settings) save() error {
	fmt.Println("Saving to file: ", s.fyneSettings.StoragePath())
	return s.saveToFile(s.storagePath)
}

func (s *Settings) saveToFile(path fyne.URI) error {
	data, err := json.Marshal(&s.settingStorage)
	if err != nil {
		return err
	}
	v, err := storage.CanWrite(path)
	if v == false || err != nil {
		fmt.Println("Error E1: ", err)
		return err
	}
	storer, err := storage.Writer(path)

	if err != nil {
		fmt.Println("Error: E2: ", err)
		return err
	}
	defer storer.Close()

	b := []byte(data)
	_, err = storer.Write(b)
	if err != nil {
		fmt.Println("Error: E3")
		return err
	}
	return nil
}

func (s *Settings) LoadAppearanceScreen(w fyne.Window) {
	s.userTheme = fyne.CurrentApp().Settings().Theme()
	if s.userTheme == nil {
		s.userTheme = theme.DefaultTheme()
	}

	def := s.fyneSettings.ThemeName
	themeNames := []string{"dark", "light"}
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		themeNames = append(themeNames, systemThemeName)
		if s.fyneSettings.ThemeName == "" {
			def = systemThemeName
		}
	}
	themes := widget.NewSelect(themeNames, s.chooseTheme)
	themes.SetSelected(def)
}
