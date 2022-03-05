package jobs

import (
	"encoding/json"
	"fmt"

	configmanager "syncdata/internal/config"
	"syncdata/pkg/repository/entities"
	"syncdata/pkg/services"
	"time"

	"github.com/rs/zerolog/log"
)

type ReassignRecordJob struct {
	scheduleTime  string
	isEnable      bool
	rabbitMq      *services.RabbitWorkerManager
	recordService *services.GameRecordService
}

type reassignRecord struct {
	Platform       string                      `json:"platform"`
	Date           string                      `json:"date"`
	Time           string                      `json:"time"`
	GameRecordList []entities.GameRecordEntity `json:"gameRecordList"`
}

func (param *ReassignRecordJob) Run() {
	utc, _ := time.LoadLocation("Asia/Taipei")
	nd := time.Now().UTC().In(utc).Add(-time.Minute * time.Duration(1))
	wd := nd.Add(-time.Minute * time.Duration(29))

	startTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:00",
		wd.Year(), wd.Month(), wd.Day(),
		wd.Hour(), wd.Minute())

	endTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d:59",
		nd.Year(), nd.Month(), nd.Day(),
		nd.Hour(), nd.Minute())

	//測試資料時間區間
	//startTime := "2021-12-28 10:00:00"
	//endTime := "2021-12-28 10:29:59"

	runTime := time.Now().UTC().In(utc).Format("2006-01-02 15:04:05")
	log.Info().Msgf("ReassignRecordJob Run, start time: %s, time range of game record: [%s] ~ [%s]", runTime, startTime, endTime)
	// 從DB撈資料
	rows, err := param.recordService.Get(startTime, endTime)

	if err != nil {
		log.Error().Err(err).Msg("get game_record from sql error")
		return
	}
	// 送Queue
	if rows != nil {
		gameRecord := reassignRecord{
			Platform:       configmanager.GlobalConfig.SyncDataPlatform.Platform,
			Date:           fmt.Sprintf("%d%02d%02d", wd.Year(), wd.Month(), wd.Day()),
			Time:           fmt.Sprintf("%02d:%02d", wd.Hour(), wd.Minute()),
			GameRecordList: rows,
		}
		byteRecord, err := json.Marshal(gameRecord)
		if err != nil {
			log.Error().Err(err).Msg("get game_record json Marshal error")
			return
		}
		err = param.rabbitMq.PublishMessage("game_record", byteRecord)
		if err != nil {
			log.Error().Err(err).Msg("publish game_record to the queue error")
			return
		}
	}
	doneTime := time.Now().UTC().In(utc).Format("2006-01-02 15:04:05")
	log.Info().Msgf("ReassignRecordJob Done, end time: %s, publish number of rows is: %d", doneTime, len(rows))
}
