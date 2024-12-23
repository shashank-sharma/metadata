package types

type TrackFocusPayload struct {
	User      string   `json:"user"`
	Device    string   `json:"device"`
	Tags      []string `json:"tags,omitempty"`
	Metadata  string   `json:"metadata,omitempty"`
	BeginDate string   `json:"begin_date,omitempty"`
	EndDate   string   `json:"end_date,omitempty"`
}
