package repository

import (
	"database/sql"
	"fmt"
	"strconv"
	"syncdata/pkg/repository/entities"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/rs/zerolog/log"
)

type AllUserRepository struct {
	db *sql.DB
}

func NewAllUserRepository(db *sql.DB) *AllUserRepository {
	return &AllUserRepository{
		db: db,
	}
}

type allUserEntity struct {
	Account   string `db:"Account"`
	ChannelId int    `db:"ChannelId"`
	CellScore int    `db:"CellScore"`
	WinNum    int    `db:"WinNum"`
	LostNum   int    `db:"LostNum"`
	WinGold   int    `db:"WinGold"`
	LostGold  int    `db:"LostGold"`
}

func (s *AllUserRepository) QueryAllUserRecord() ([]entities.GroupByAllUserEntity, error) {
	utc, _ := time.LoadLocation("Asia/Taipei")
	nd := time.Now().UTC().In(utc)
	wd := nd.Add(-time.Hour * time.Duration(720))

	startTime := fmt.Sprintf("%d-%02d-%02d 00:00:00",
		wd.Year(), wd.Month(), wd.Day())

	endTime := fmt.Sprintf("%d-%02d-%02d 23:59:59",
		nd.Year(), nd.Month(), nd.Day())

	var querySql = fmt.Sprintf("SELECT Account,ChannelID, CellScore,WinNum,LostNum,WinGold,LostGold FROM KYStatisUsers.statis_all_users  where UpdateTime between '%s' and '%s'", startTime, endTime)
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

	allUserRecords := make([]allUserEntity, 0)
	for res.Next() {
		allUserRecord := allUserEntity{}
		err := res.Scan(&allUserRecord.Account, &allUserRecord.ChannelId, &allUserRecord.CellScore, &allUserRecord.WinNum, &allUserRecord.LostNum, &allUserRecord.WinGold, &allUserRecord.LostGold)

		if err != nil {
			log.Err(err).Msgf("Failed mysql parse row")
			return nil, err
		}
		allUserRecords = append(allUserRecords, allUserRecord)
	}
	//聚合
	var result []entities.GroupByAllUserEntity
	linq.From(allUserRecords).GroupByT(
		func(p allUserEntity) string {
			return p.Account + strconv.Itoa(p.ChannelId)
		},
		func(p allUserEntity) allUserEntity { return p },
	).OrderByT(
		func(g linq.Group) string { return g.Key.(string) },
	).SelectIndexedT(func(index int, wordGroup linq.Group) entities.GroupByAllUserEntity {
		dataItem := entities.GroupByAllUserEntity{}
		for _, objItem := range wordGroup.Group {
			obj := objItem.(allUserEntity)
			dataItem.Account = obj.Account
			dataItem.CellScore = obj.CellScore
			dataItem.WinNum = obj.WinNum
			dataItem.LostNum = obj.LostNum
			dataItem.WinGold = obj.WinGold
			dataItem.LostGold = obj.LostGold
		}

		return dataItem
	}).ToSlice(&result)

	return result, nil
}
