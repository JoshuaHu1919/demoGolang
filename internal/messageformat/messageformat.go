package messageformat

type PlatformInfo struct {
	Platform  string `json:"Platform"`
	AlertName string `json:"AlertName"`
}

type MonitorMessage struct {
	PlatformInfo PlatformInfo `json:"PlatformInfo"`
	TriggerRule  string       `json:"TriggerRule"`
	Timestamp    uint64       `json:"Timestamp"`
}

type SubscriptUser struct {
	PlatformInfo
	Name          string `json:"Name"`
	BroadcastType uint64 `json:"BroadcastType"`
}
