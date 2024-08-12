package types

import "time"

type AWEvent struct {
	Id        int64       `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Duration  float64     `json:"duration"`
	Data      AWEventData `json:"data"`
}

type AWEventData struct {
	Title string `json:"title"`
	App   string `json:"app"`
}

type AWEvents []AWEvent

func (e AWEvents) Len() int {
	return len(e)
}
func (e AWEvents) Less(i, j int) bool {
	return e[i].Timestamp.After(e[j].Timestamp)
}
func (e AWEvents) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
