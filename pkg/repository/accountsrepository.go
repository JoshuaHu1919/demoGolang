package repository

import (
	"database/sql"

	"syncdata/pkg/repository/entities"

	"github.com/rs/zerolog/log"
)

type AccountsRepository struct {
	db *sql.DB
}

func NewAccountsRepository(db *sql.DB) *AccountsRepository {
	return &AccountsRepository{
		db: db,
	}
}
func (s *AccountsRepository) QueryAccountsRecord() ([]entities.AccountsEntity, error) {

	querySql := "SELECT account,agent,createdate,lastlogintime FROM game_api.accounts ;"
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

	accountsRecords := make([]entities.AccountsEntity, 0)
	for res.Next() {
		onlineRoomRecord := entities.AccountsEntity{}
		err := res.Scan(&onlineRoomRecord.Account, &onlineRoomRecord.Agent, &onlineRoomRecord.CreateDate, &onlineRoomRecord.LastLoginTime)

		if err != nil {
			log.Err(err).Msgf("Failed mysql parse row")
			return nil, err
		}
		accountsRecords = append(accountsRecords, onlineRoomRecord)
	}
	return accountsRecords, nil
}
