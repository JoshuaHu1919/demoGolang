package repository

import (
	"database/sql"
	"fmt"

	"syncdata/pkg/repository/entities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

type GameRecordRepository struct {
	db *sql.DB
}

func NewGameRecordRepository(db *sql.DB) *GameRecordRepository {
	return &GameRecordRepository{
		db: db,
	}
}

func (s *GameRecordRepository) QueryGameRecord(querySql string) ([]entities.GameRecordEntity, error) {
	res, err := s.db.Query(querySql)
	if err != nil {
		log.Err(err).Msgf("Failed mysql connection")
		return nil, err
	}
	defer func(res *sql.Rows) {
		err := res.Close()
		if err != nil {
			log.Err(err).Msgf("Failed mysql  row")
		}
	}(res)

	gameRecords := make([]entities.GameRecordEntity, 0)
	for res.Next() {
		gameRecord := entities.GameRecordEntity{}
		err := res.Scan(&gameRecord.ChannelId, &gameRecord.Accounts, &gameRecord.KindId, &gameRecord.ServerId, &gameRecord.CellScore, &gameRecord.Profit,
			&gameRecord.CreateTime, &gameRecord.IsNew, &gameRecord.LastLoginTime, &gameRecord.RegisterTime,
			&gameRecord.HistoryCellScore, &gameRecord.HistoryProfit, &gameRecord.HistoryGameNum)

		if err != nil {
			log.Err(err).Msgf("Failed mysql parse row")
			return nil, err
		}
		gameRecords = append(gameRecords, gameRecord)
	}
	return gameRecords, nil
}

func (s *GameRecordRepository) Exists(datetime string) ([]string, error) {
	sqlStr := fmt.Sprintf(`SELECT TABLE_NAME FROM (SELECT t.* FROM KYDB_NEW.GameInfo as game INNER JOIN information_schema.TABLES as t ON t.TABLE_SCHEMA REGEXP CONCAT('^',game.GameParameter,'_record$')) as gameTemp 
				WHERE DATEDIFF(STR_TO_DATE('%s','%%Y-%%m-%%d'),STR_TO_DATE(gameTemp.TABLE_NAME,'gameRecord%%Y%%m%%d')) = 0`, datetime)
	res, err := s.db.Query(sqlStr)
	if err != nil {
		log.Err(err).Msgf("Failed mysql connection")
		return nil, err
	}
	defer func(res *sql.Rows) {
		err := res.Close()
		if err != nil {
			log.Err(err).Msgf("Failed mysql row")
		}
	}(res)

	var tableNames []string
	for res.Next() {
		var tableName string
		err := res.Scan(&tableName)

		if err != nil {
			log.Err(err).Msgf("Failed mysql parse row")
			return nil, err
		}
		tableNames = append(tableNames, tableName)
	}
	return tableNames, nil
}
