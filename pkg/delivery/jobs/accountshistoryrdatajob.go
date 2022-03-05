package jobs

import (
	"encoding/json"
	"time"

	"syncdata/pkg/services"

	"github.com/rs/zerolog/log"
)

type AccountsHistoryJob struct {
	scheduleTime  string
	isEnable      bool
	rabbitMq      *services.RabbitWorkerManager
	recordService *services.AccountsService
}

func (param *AccountsHistoryJob) Run() {
	utc, _ := time.LoadLocation("Asia/Taipei")
	// 從DB撈資料
	rows, rowsErr := param.recordService.GetAccountsRecord()

	if rowsErr != nil {
		log.Error().Err(rowsErr).Msg("get AccountsHistoryJob from sql error")
		return
	}

	// 送Queue
	if rows != nil {

		byteRecord, byteRecordErr := json.Marshal(rows)
		if byteRecordErr != nil {
			log.Error().Err(byteRecordErr).Msg("get AccountsHistoryJob json Marshal error")
			return
		}
		mqErr := param.rabbitMq.PublishMessage("accounts_history_record", byteRecord)
		if mqErr != nil {
			log.Error().Err(mqErr).Msg("publish AccountsHistoryJob to the queue error")
			return
		}
	}
	doneTime := time.Now().UTC().In(utc).Format("2006-01-02 15:04:05")
	log.Info().Msgf("AccountsHistoryJob Done, end time: %s, publish number of rows is: %d", doneTime, len(rows))

}
