package services

import (
	"context"
	"reflect"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type JobManager struct {
	cron *cron.Cron
}

func NewJobManager() *JobManager {
	utc8, _ := time.LoadLocation("Asia/Taipei")
	schedule := cron.New(cron.WithLocation(utc8))
	return &JobManager{
		cron: schedule,
	}
}

func (r *JobManager) StartWorkers(init chan bool, jobs []cron.Job) {
	for _, job := range jobs {
		rj := reflect.ValueOf(job)
		f := reflect.Indirect(rj).FieldByName("scheduleTime")
		isEnableObj := reflect.Indirect(rj).FieldByName("isEnable")
		scheduleTime := f.String()
		isEnable := isEnableObj.Bool()
		if isEnable {
			_, err := r.cron.AddJob(scheduleTime, job)
			if err != nil {
				panic(err)
			}
		}
	}

	log.Info().Msgf("init job schedule success")
	init <- true
	r.cron.Start()
}

func (r *JobManager) Stop() context.Context {
	return r.cron.Stop()
}
