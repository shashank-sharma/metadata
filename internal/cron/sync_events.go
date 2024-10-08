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

func SyncAWEventJob(awService activitywatch.AWService, backendService backend.BackendService, c config.AppConfig, bucketId string) func(cronInfo *CronInfo) {
	return func(cronInfo *CronInfo) {
		var err error
		logger.LogDebug("Running Sync AWEvent Job new")
		userSettings := c.Settings.UserSettings
		tempBucket := userSettings.Bucket[bucketId]
		syncTime := time.Now().UTC()

		startTimestamp := tempBucket.StartTimestamp

		// Step 1: Fetch initial timestamp
		if startTimestamp.Equal(time.Time{}) {
			logger.LogDebug("Timestamp default, fetching start timestamp")
			startTimestamp, err = findStartTimestamp(awService, bucketId)
			if err != nil {
				logger.LogError("Failed finding start timestamp: ", err)
				currentFailedCount, _ := cronInfo.FailedCount.Get()
				cronInfo.FailedCount.Set(currentFailedCount + 1)
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
		for start := startTimestamp; start.Before(syncTime); start = start.Add(4 * time.Hour) {
			end := start.Add(4 * time.Hour)
			if end.After(syncTime) {
				end = syncTime
			}
			events, err := awService.FetchEvents(bucketId, start, end)
			sort.Sort(events)
			if err != nil {
				logger.LogWarning("Failed fetching events: ", err)
				currentFailedCount, _ := cronInfo.FailedCount.Get()
				cronInfo.FailedCount.Set(currentFailedCount + 1)
				return
			}

			// If ActivityWatch closes all of sudden then for
			// given timestamp, it is possible events count is 0
			// hence it is not a failure and recovery is required
			// TODO: Sync only if events are greater than 0
			// TODO: Need better error handling
			if len(events) == 0 {
				logger.LogInfo("Failed to find any events")
				tempBucket.EndTimestamp = end
			} else {
				logger.LogWarning("Found events: ", len(events))
				data, err := backendService.SyncEventData(userSettings.ProductId, bucketId, events)
				if err != nil {
					logger.LogError("Error syncing data with backend: ", err)
					currentFailedCount, _ := cronInfo.FailedCount.Get()
					cronInfo.FailedCount.Set(currentFailedCount + 1)
					return
				}
				logger.LogDebug("Synced with response: ", data)
				tempBucket.EndTimestamp = events[0].Timestamp
			}

			tempBucket.LastSynced = end
			userSettings.Bucket[bucketId] = tempBucket
			c.SettingsManager.SaveSettings(userSettings)
			currentSuccessCount, _ := cronInfo.SuccessCount.Get()
			cronInfo.SuccessCount.Set(currentSuccessCount + 1)
		}
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
