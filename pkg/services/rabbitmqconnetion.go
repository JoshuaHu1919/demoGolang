package services

import (
	"context"
	"fmt"
	"sync"
	configManager "syncdata/internal/config"

	"github.com/rs/zerolog/log"
	"github.com/streadway/amqp"
)

type RabbitWorkerManager struct {
	ctx          context.Context
	cancelFunc   context.CancelFunc
	waitGroup    sync.WaitGroup
	queueMap     map[string]*amqp.Channel
	MQConnection *amqp.Connection
}

func NewRabbitWorkerManager() *RabbitWorkerManager {

	_ctx, _cancelFunc := context.WithCancel(context.Background())
	url := fmt.Sprintf("amqp://%s:%s@%s/%s",
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.Username,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.Password,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.RabbitMQAddr,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.VirtualHosts,
	)

	_conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal().Msgf("%s", err.Error())
		defer func(_conn *amqp.Connection) {
			err := _conn.Close()
			if err != nil {
				log.Fatal().Msgf("%s", err.Error())
			}
		}(_conn)
	}

	log.Info().Msgf("init subscribe RabbitMQ success")

	return &RabbitWorkerManager{
		ctx:          _ctx,
		cancelFunc:   _cancelFunc,
		queueMap:     make(map[string]*amqp.Channel),
		MQConnection: _conn,
	}
}
func (r *RabbitWorkerManager) Restart(init chan bool) {
	for _, channel := range r.queueMap {
		err := channel.Close()
		if err != nil {
			log.Fatal().Msgf("RabbitWorkerManager: Restart: error: %s", err.Error())
		}
	}
	r.cancelFunc()
	r.waitGroup.Wait()
	defer func(MQConnection *amqp.Connection) {
		err := MQConnection.Close()
		if err != nil {
			log.Error().Msgf("%s", err)
		}
	}(r.MQConnection)

	r.ctx, r.cancelFunc = context.WithCancel(context.Background())

	url := fmt.Sprintf("amqp://%s:%s@%s/%s",
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.Username,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.Password,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.RabbitMQAddr,
		configManager.GlobalConfig.SyncDataPlatform.NotifierCenter.VirtualHosts)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal().Msgf("%s", err.Error())
		return
	}

	r.MQConnection = conn
	log.Info().Msgf("init subscribe RabbitMQ success")
	init <- true

	<-r.ctx.Done()
	log.Info().Msgf("I am dead: main worker")
}
func (r *RabbitWorkerManager) PublishMessage(queueName string, message []byte) error {
	channel, err := r.MQConnection.Channel()
	if err != nil {
		log.Err(err).Msgf("Failed to open a channel: %s", err)
		return err
	}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Err(err).Msgf("Failed to declare queue")
		return err
	}

	err = channel.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message})
	if err != nil {
		log.Err(err).Msgf("Failed to send message : %s", message)
		return err
	}

	return nil
}
