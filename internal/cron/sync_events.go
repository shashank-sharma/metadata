package cron

import (
	"errors"
	"sort"
	"time"

	"github.com/shashank-sharma/metadata/internal/config"
	"github.com/shashank-sharma/metadata/internal/logger"
	"github.com/shashank-sharma/metadata/internal/services/activitywatch"
	"github.com/shashank-sharma/metadata/internal/services/backend"
	"github.com/shashank-sharma/metadata/internal/types"
)

func SyncAWEventJob(awService activitywatch.AWService, backendService backend.BackendService, c config.AppConfig, bucketId string) func() {
	return func() {
		var err error
		logger.LogDebug("Running Sync AWEvent Job")
		userSettings := c.Settings.UserSettings
		tempBucket := userSettings.Bucket[bucketId]
		syncTime := time.Now().UTC()

		startTimestamp := tempBucket.StartTimestamp

		// Step 1: Fetch initial timestamp
		if startTimestamp.Equal(time.Time{}) {
			logger.LogDebug("Timestamp default, fetching start timestamp")
			startTimestamp, err = findStartTimestamp(awService, bucketId)
			if err != nil {
				logger.LogError("Failed finding start timestamp")
				return
			}
			// For testing assuming I started AW 1 hour ago
			// startTimestamp = syncTime.Add(-1 * time.Hour)
			tempBucket.StartTimestamp = startTimestamp
		} else {
			logger.LogDebug("Continuing from EndTimeStamp")
			startTimestamp = tempBucket.EndTimestamp
		}

		logger.LogDebug("Starting with timestamp: ", startTimestamp)
		logger.LogDebug("LastSynced timestamp is: ", tempBucket.LastSynced)

		// Step 2: Fetch and sync events in 1-day intervals up to the current time
		for start := startTimestamp; start.Before(syncTime); start = start.AddDate(0, 0, 1) {
			end := start.AddDate(0, 0, 1)
			if end.After(syncTime) {
				end = syncTime
			}
			events, err := awService.FetchEvents(bucketId, start, end)
			sort.Sort(events)
			if err != nil {
				logger.LogWarning("Failed fetching events")
				break
			}

			// TODO: Sync only if events are greater than 0
			// TODO: Need better error handling
			if len(events) == 0 {
				logger.LogError("Failed to find any events")
				break
			}

			data, err := backendService.SyncEventData(userSettings.ProductId, bucketId, events)
			if err != nil {
				logger.LogError("Error syncing data with backend")
				break
			}
			logger.LogDebug("Synced with response: ", data)
			tempBucket.LastSynced = end
			tempBucket.EndTimestamp = events[0].Timestamp
		}

		userSettings.Bucket[bucketId] = tempBucket
		c.SettingsManager.SaveSettings(userSettings)
	}
}

func findStartTimestamp(awService activitywatch.AWService, bucketId string) (time.Time, error) {
	for id := 0; id < 10; id++ {
		logger.LogDebug("Trying for id: ", id)
		event, err := awService.FetchEventById(bucketId, id)
		if err != nil || event == (types.AWEvent{}) {
			logger.LogWarning("Failed fetching timestamp for id: ", id)
			continue
		}

		logger.LogDebug("Found timestamp: ", event.Timestamp)
		return event.Timestamp, nil
	}
	return time.Now(), errors.New("Failed finding ID")
}
