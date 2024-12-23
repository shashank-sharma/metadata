package controllers

import (
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/router"
	"github.com/shashank-sharma/metadata/internal/settings"
	"github.com/shashank-sharma/metadata/internal/types"
)

type FocusController struct {
	router       *router.Router
	isLogging    bool
	selectedTags []string
	metadata     string
	startTime    time.Time
	elapsedLabel *widget.Label
	tagLabels    map[string]*widget.Check
	stopRefresh  chan bool
}

func NewFocusController(router *router.Router) *FocusController {
	return &FocusController{
		router:       router,
		isLogging:    false,
		selectedTags: make([]string, 0),
	}
}

func (lc *FocusController) loadFocusState() {
	userSettings := lc.router.AppCtx.Config.Settings.UserSettings
	if userSettings.CurrentFocus != nil {
		lc.isLogging = true
		lc.selectedTags = userSettings.CurrentFocus.Tags
		lc.metadata = userSettings.CurrentFocus.Metadata
		lc.startTime = userSettings.CurrentFocus.CreatedAt
	}
}

func (lc *FocusController) saveFocusState() error {
	userSettings := lc.router.AppCtx.Config.Settings.UserSettings
	userSettings.CurrentFocus = &settings.FocusConfig{
		Tags:      lc.selectedTags,
		Metadata:  lc.metadata,
		CreatedAt: lc.startTime,
	}
	return lc.router.AppCtx.Config.SettingsManager.SaveSettings(userSettings)
}

func (lc *FocusController) clearFocusState() error {
	userSettings := lc.router.AppCtx.Config.Settings.UserSettings
	userSettings.CurrentFocus = nil
	return lc.router.AppCtx.Config.SettingsManager.SaveSettings(userSettings)
}

func (lc *FocusController) Screen(meta router.RouteMetadata) fyne.CanvasObject {
	title := widget.NewLabel(meta.Title)

	lc.tagLabels = make(map[string]*widget.Check)
	lc.stopRefresh = make(chan bool)

	lc.loadFocusState()

	existingTags := []string{"work", "entertainment", "gaming", "reading", "creative", "finance", "research", "writing", "learning"}
	tagChecks := container.NewGridWithColumns(3)

	for _, tag := range existingTags {
		currentTag := tag
		check := widget.NewCheck(currentTag, func(checked bool) {
			if checked {
				if !contains(lc.selectedTags, currentTag) {
					lc.selectedTags = append(lc.selectedTags, currentTag)
				}
			} else {
				lc.selectedTags = removeTag(lc.selectedTags, currentTag)
			}
		})
		lc.tagLabels[currentTag] = check
		tagChecks.Add(check)
	}

	selectedTagsLabel := widget.NewLabel("Selected Tags: none")
	lc.elapsedLabel = widget.NewLabel("Time Elapsed: 00:00:00")
	metadataEntry := widget.NewMultiLineEntry()
	metadataEntry.SetPlaceHolder("Enter metadata here...")

	// Initialize UI with saved state if exists
	if lc.isLogging {
		for _, tag := range lc.selectedTags {
			if check, exists := lc.tagLabels[tag]; exists {
				check.SetChecked(true)
				check.Disable()
			}
		}
		selectedTagsLabel.SetText(fmt.Sprintf("Selected Tags: %s",
			strings.Join(lc.selectedTags, ", ")))
		metadataEntry.SetText(lc.metadata)
		metadataEntry.Disable()
		go lc.updateElapsedTime()
	}

	var buttonText string
	if lc.isLogging {
		buttonText = "Stop"
	} else {
		buttonText = "Start"
	}
	toggleButton := widget.NewButton(buttonText, nil)
	toggleButton.OnTapped = func() {
		if !lc.isLogging {
			if len(lc.selectedTags) == 0 {
				lc.router.AppCtx.Notification.Show("Please select at least one tag", "error")
				return
			}

			lc.startTime = time.Now()
			lc.isLogging = true
			lc.metadata = metadataEntry.Text
			toggleButton.SetText("Stop")

			if err := lc.saveFocusState(); err != nil {
				logger.LogError("Failed to save focus state:", err)
				lc.router.AppCtx.Notification.Show("Failed to save state", "error")
				return
			}

			go lc.updateElapsedTime()

			selectedTagsLabel.SetText(fmt.Sprintf("Selected Tags: %s",
				strings.Join(lc.selectedTags, ", ")))

			for _, check := range lc.tagLabels {
				check.Disable()
			}
			metadataEntry.Disable()

		} else {
			lc.isLogging = false
			toggleButton.SetText("Start")

			lc.stopRefresh <- true

			endTime := time.Now()

			userSettings := lc.router.AppCtx.Config.Settings.UserSettings

			payload := types.TrackFocusPayload{
				User:      userSettings.UserId,
				Device:    userSettings.ProductId,
				Tags:      lc.selectedTags,
				Metadata:  lc.metadata,
				BeginDate: lc.startTime.Format(time.RFC3339),
				EndDate:   endTime.Format(time.RFC3339),
			}

			response, err := lc.router.AppCtx.BackendService.CreateTrackFocus(payload)
			if err != nil || !response {
				lc.router.AppCtx.Notification.Show("Failed to upload", "error")
				logger.LogError("Error uploading focus time: ", err.Error())
				return
			}

			if err := lc.clearFocusState(); err != nil {
				logger.LogError("Failed to clear focus state:", err)
			}

			lc.resetState(metadataEntry, selectedTagsLabel)

			for _, check := range lc.tagLabels {
				check.Enable()
			}
			metadataEntry.Enable()
		}
	}

	// Main container
	return container.NewVBox(
		title,
		widget.NewLabel(meta.Content),
		widget.NewLabel("Select Tags:"),
		tagChecks,
		selectedTagsLabel,
		lc.elapsedLabel,
		widget.NewLabel("Enter Metadata:"),
		metadataEntry,
		toggleButton,
	)
}

func (lc *FocusController) updateElapsedTime() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if lc.isLogging {
				elapsed := time.Since(lc.startTime)
				hours := int(elapsed.Hours())
				minutes := int(elapsed.Minutes()) % 60
				seconds := int(elapsed.Seconds()) % 60
				lc.elapsedLabel.SetText(fmt.Sprintf("Time Elapsed: %02d:%02d:%02d",
					hours, minutes, seconds))
			}
		case <-lc.stopRefresh:
			return
		}
	}
}

func (lc *FocusController) resetState(metadataEntry *widget.Entry, selectedTagsLabel *widget.Label) {
	lc.selectedTags = []string{}

	for _, check := range lc.tagLabels {
		check.SetChecked(false)
	}

	selectedTagsLabel.SetText("Selected Tags: none")
	lc.elapsedLabel.SetText("Time Elapsed: 00:00:00")
	metadataEntry.SetText("")
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func removeTag(slice []string, str string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if v != str {
			result = append(result, v)
		}
	}
	return result
}

func getTagFromCheck(tagLabels map[string]*widget.Check, checked bool) string {
	for tag, check := range tagLabels {
		if check.Checked == checked {
			return tag
		}
	}
	return ""
}
