package activitywatch

import (
	"encoding/json"
	"fmt"
	"time"

	types "github.com/shashank-sharma/metadata/internal/types"
)

func (as *AWService) FetchEventById(bucket string, id int) (types.AWEvent, error) {
	req, _ := as.Client.NewRequest("GET", fmt.Sprintf("/api/0/buckets/%s/events/%d", bucket, id), nil, map[string]string{})
	resp, err := as.Client.Do(req)

	if err != nil {
		return types.AWEvent{}, err
	}

	if resp.StatusCode/100 != 2 {
		return types.AWEvent{}, err
	}

	data := types.AWEvent{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}

func (as *AWService) FetchEvents(bucket string, start time.Time, end time.Time) (types.AWEvents, error) {
	layout := "2006-01-02T15:04:05.999999"
	params := map[string]string{
		"start": start.Format(layout),
		"end":   end.Format(layout),
	}
	req, _ := as.Client.NewRequestWithParams("GET", fmt.Sprintf("/api/0/buckets/%s/events", bucket), params, map[string]string{})
	resp, err := as.Client.Do(req)

	if err != nil {
		return types.AWEvents{}, err
	}

	if resp.StatusCode/100 != 2 {
		return types.AWEvents{}, err
	}

	data := types.AWEvents{}
	json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}
