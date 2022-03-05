package jobs

import (
	"database/sql"

	configmanager "syncdata/internal/config"
	"syncdata/pkg/repository"
	"syncdata/pkg/services"

	"github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type JobSetting struct {
	JobsItem []cron.Job
}

func NewJobSetting(mq *services.RabbitWorkerManager) *JobSetting {
	config := &mysql.Config{
		User:                 configmanager.GlobalConfig.SyncDataPlatform.Mysql.Username,
		Passwd:               configmanager.GlobalConfig.SyncDataPlatform.Mysql.Password,
		Addr:                 configmanager.GlobalConfig.SyncDataPlatform.Mysql.Address,
		Net:                  "tcp",
		DBName:               configmanager.GlobalConfig.SyncDataPlatform.Mysql.Database,
		AllowNativePasswords: true,
		MultiStatements:      true,
	}
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		log.Err(err).Msgf("Failed mysql connection")
	}
	repo := repository.NewGameRecordRepository(db)
	allUserRepo := repository.NewAllUserRepository(db)
	srv := services.NewGameRecordService(repo, allUserRepo)

	roomRepo := repository.NewOnlineRoomRepository(db)
	roomSrv := services.NewOnlineRoomService(roomRepo)

	accountsRepo := repository.NewAccountsRepository(db)
	accountsSrv := services.NewAccountsService(accountsRepo)

	jobs := []cron.Job{
		&ReassignRecordJob{
			scheduleTime:  "*/30 * * * *",
			rabbitMq:      mq,
			isEnable:      true,
			recordService: srv,
		},
		&OnlineRoomJob{
			scheduleTime:  "1 0 * * *",
			rabbitMq:      mq,
			isEnable:      true,
			recordService: roomSrv,
		},
		&AccountsHistoryJob{
			scheduleTime:  "0 0 * * *",
			rabbitMq:      mq,
			isEnable:      true,
			recordService: accountsSrv,
		},
	}

	return &JobSetting{
		JobsItem: jobs,
	}
}
