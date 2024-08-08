package activitywatch

import (
	"encoding/json"
	"errors"
	"time"
)

type ActivityWatchBucket struct {
	ID          string          `json:"id"`
	Created     time.Time       `json:"created"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Client      string          `json:"client"`
	Hostname    string          `json:"hostname"`
	Data        json.RawMessage `json:"data"`
	LastUpdated time.Time       `json:"last_updated"`
}

func (as *AWService) FetchBuckets() (map[string]ActivityWatchBucket, error) {
	req, _ := as.Client.NewRequest("GET", "/api/0/buckets", nil, map[string]string{})
	resp, err := as.Client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(resp.Status)
	}

	var bucketMap map[string]ActivityWatchBucket
	if err := json.NewDecoder(resp.Body).Decode(&bucketMap); err != nil {
		return nil, err
	}

	return bucketMap, nil
}
