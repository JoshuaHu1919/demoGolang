package jobs

import (
	"encoding/json"
	"fmt"
	"time"

	"syncdata/pkg/services"

	"github.com/rs/zerolog/log"
)

type OnlineRoomJob struct {
	scheduleTime  string
	isEnable      bool
	rabbitMq      *services.RabbitWorkerManager
	recordService *services.OnlineRoomService
}

func (param *OnlineRoomJob) Run() {
	utc, _ := time.LoadLocation("Asia/Taipei")
	dd, _ := time.ParseDuration("-24h")
	wd := time.Now().UTC().In(utc).Add(dd)

	startDate := fmt.Sprintf("%d-%02d-%02d",
		wd.Year(), wd.Month(), wd.Day(),
	)

	log.Info().Msgf("OnlineRoomJob Run, start time: %s", utc.String())
	// 從DB撈資料
	rows, rowsErr := param.recordService.GetOnlineRoomRecord(startDate)

	if rowsErr != nil {
		log.Error().Err(rowsErr).Msg("get OnlineRoomJob from sql error")
		return
	}

	// 送Queue
	if rows != nil {

		byteRecord, byteRecordErr := json.Marshal(rows)
		if byteRecordErr != nil {
			log.Error().Err(byteRecordErr).Msg("get OnlineRoomJob json Marshal error")
			return
		}
		mqErr := param.rabbitMq.PublishMessage("onlineRoom_record", byteRecord)
		if mqErr != nil {
			log.Error().Err(mqErr).Msg("publish OnlineRoomJob to the queue error")
			return
		}
	}
	doneTime := time.Now().UTC().In(utc).Format("2006-01-02 15:04:05")
	log.Info().Msgf("OnlineRoomJob Done, end time: %s, publish number of rows is: %d", doneTime, len(rows))

}
