package contract

type OnlineRoomContract struct {
	RoomId     int     `json:"roomId"`
	Value      float64 `json:"onlineUserAvg"`
	CreateDate string  `json:"createDate"`
}
