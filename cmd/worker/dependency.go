package main

import (
	"net/http"
	"syncdata/pkg/delivery/jobs"

	resetConsumerHTTP "syncdata/pkg/delivery/http"

	"syncdata/pkg/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initialize() (*gin.Engine, *services.JobManager, error) {
	manager := services.NewRabbitWorkerManager()
	schedule := services.NewJobManager()
	jobSetting := jobs.NewJobSetting(manager)
	resetConsumerHandler := resetConsumerHTTP.NewResetConsumerHTTPHandler(manager)
	return newGinEngine(resetConsumerHandler, schedule, jobSetting), schedule, nil
}

func newGinEngine(resetConsumerHandler *resetConsumerHTTP.ResetConsumerHTTPHandler, schedule *services.JobManager, jobSetting *jobs.JobSetting) *gin.Engine {
	init := make(chan bool)
	go schedule.StartWorkers(init, jobSetting.JobsItem)
	<-init
	defaultRouter := gin.Default()
	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"*", "Authorization", "Content-Type", "Origin", "Content-Length"},

		// firefox 和 safari 不支援 *, 所以需要一個一個打，但更好的是要抓 Access-Control-Request-Headers
		// https://stackoverflow.com/questions/54666673/cors-check-fails-for-firefox-but-passes-for-chrome
	}
	defaultRouter.Use(cors.New(corsConfig))
	app := defaultRouter.Group("")
	resetConsumerHTTP.SetRoutes(app, resetConsumerHandler)
	defaultRouter.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "StatusNotFound")
	})

	return defaultRouter
}
