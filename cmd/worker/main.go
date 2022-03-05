package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	configmanager "syncdata/internal/config"
	"syncdata/internal/logformat"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = logformat.InitLogger()
	defer func() {
		if r := recover(); r != nil {
			// unknown error
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("unknown error: %v", err)
			}
			log.Fatal().Msgf("%v", err)
			time.Sleep(3 * time.Second)
		}
	}()
	ginEngine, jobManager, err := initialize()
	if err != nil {
		log.Panic().Msgf("main: initialize failed: %v", err)
		return
	}

	// start http server
	httpServer := &http.Server{
		Addr:    configmanager.GlobalConfig.SyncDataPlatform.HTTPBind,
		Handler: ginEngine,
	}
	go func() {
		// service connection
		log.Info().Msgf("main: Listening and serving HTTP on %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Panic().Msgf("main: http server listen failed: %v", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	log.Info().Msgf("main: shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Panic().Msgf("main: http shutdown error: %v", err)
	} else {
		log.Info().Msgf("main: http gracefully stopped")
	}

	jctx := jobManager.Stop()
	log.Info().Msgf("main: shutting down cron...")
	select {
	case <-time.After(10 * time.Second):
		log.Panic().Msgf("main: cron forced to shutdown...")
	case <-jctx.Done():
		log.Info().Msgf("main: cron exiting...")
	}
	log.Info().Msgf("main: server shut down successfully")
}
