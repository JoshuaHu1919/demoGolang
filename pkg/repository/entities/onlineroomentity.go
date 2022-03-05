package entities

type OnlineRoomEntity struct {
	Id         int    `db:"id"`
	RoomId     string `db:"roomId"`
	Value      int    `db:"value"`
	IP         string `db:"ip"`
	CreateTime string `db:"createtime"`
}
