package settings

import (
	"encoding/json"
	"io"
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"github.com/shashank-sharma/metadata/internal/logger"
)

// SettingsManager is responsible for managing reading and writing settings.
type SettingsManager struct {
	StorageRoot fyne.URI
}

func (manager *SettingsManager) LoadSettings(s BaseSettings) error {
	fileURI, _ := storage.Child(manager.StorageRoot, s.FileName())
	logger.LogDebug("Reading file: %s", fileURI)
	r, err := storage.Reader(fileURI)
	if err != nil {
		return err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, s)
}

func (manager *SettingsManager) SaveSettings(s BaseSettings) error {
	fileURI, _ := storage.Child(manager.StorageRoot, s.FileName())
	w, err := storage.Writer(fileURI)
	if err != nil {
		return err
	}
	defer w.Close()

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (manager *SettingsManager) ResetSettings(s BaseSettings) error {
	fileURI, _ := storage.Child(manager.StorageRoot, s.FileName())
	return storage.Delete(fileURI)
}

func (manager *SettingsManager) SettingsExists(s BaseSettings) (bool, error) {
	fileURI, _ := storage.Child(manager.StorageRoot, s.FileName())
	f, err := storage.CanRead(fileURI)
	if f == false || err != nil {
		return false, err
	}

	r, err := storage.Reader(fileURI)
	if err != nil {
		logger.LogWarning("Can't open existing file")
		return false, nil
	}
	defer r.Close()
	return f, nil
}

func (manager *SettingsManager) InitializeSettings(s BaseSettings) error {
	exists, err := manager.SettingsExists(s)
	if err != nil {
		return err
	}

	if !exists {
		logger.LogDebug("Not exist: ", s.FileName())
		switch v := s.(type) {
		case *ApplicationSettings:
			hName, err := os.Hostname()
			if err != nil {
				logger.Error.Printf("error retrieving hostname: %v", err)
				hName = "unknown"
			}

			v.HostName = hName
			v.OperatingSystem = runtime.GOOS
			v.Arch = runtime.GOARCH
		}

		err = manager.SaveSettings(s)
		if err != nil {
			return err
		}
	} else {
		err = manager.LoadSettings(s)
		if err != nil {
			return err
		}
	}

	return nil
}
