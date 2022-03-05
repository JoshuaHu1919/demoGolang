package repository

import (
	"database/sql"
	"fmt"
	"syncdata/pkg/repository/entities"

	"github.com/rs/zerolog/log"
)

type OnlineRoomRepository struct {
	db *sql.DB
}

func NewOnlineRoomRepository(db *sql.DB) *OnlineRoomRepository {
	return &OnlineRoomRepository{
		db: db,
	}
}
func (s *OnlineRoomRepository) QueryOnlineRoomRecord(date string) ([]entities.OnlineRoomEntity, error) {

	querySql := fmt.Sprintf(`SELECT * FROM game_statistics.online_room where createtime >='%s 00:00:00' and createtime <='%s 23:59:59';`, date, date)
	res, err := s.db.Query(querySql)
	if err != nil {
		log.Err(err).Msgf(querySql)
		return nil, err
	}
	defer func(res *sql.Rows) {
		err := res.Close()
		if err != nil {
			log.Err(err).Msgf("Failed mysql  row")
		}
	}(res)

	onlineRoomRecords := make([]entities.OnlineRoomEntity, 0)
	for res.Next() {
		onlineRoomRecord := entities.OnlineRoomEntity{}
		err := res.Scan(&onlineRoomRecord.Id, &onlineRoomRecord.RoomId, &onlineRoomRecord.Value, &onlineRoomRecord.IP, &onlineRoomRecord.CreateTime)

		if err != nil {
			log.Err(err).Msgf("Failed mysql parse row")
			return nil, err
		}
		onlineRoomRecords = append(onlineRoomRecords, onlineRoomRecord)
	}
	return onlineRoomRecords, nil
}
