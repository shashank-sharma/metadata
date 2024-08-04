package activitywatch

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/shashank-sharma/metadata/internal/logger"
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

func (as *AWService) FetchBuckets() ([]ActivityWatchBucket, error) {
	req, _ := as.Client.NewRequest("GET", "/api/0/buckets", nil, map[string]string{})
	resp, err := as.Client.Do(req)

	if err != nil {
		return nil, err
	}

	logger.Debug.Printf("User response status: %+v", resp.StatusCode)

	if resp.StatusCode/100 != 2 {
		return nil, errors.New(resp.Status)
	}

	data := []ActivityWatchBucket{}
	json.NewDecoder(resp.Body).Decode(&data)

	logger.Debug.Printf("Fetch user response: %+v", data)
	return data, nil
}
