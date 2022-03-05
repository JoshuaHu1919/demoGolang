package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"syncdata/pkg/services"
)

type ResetConsumerHTTPHandler struct {
	RabbitWorkerManager *services.RabbitWorkerManager
}

func NewResetConsumerHTTPHandler(rabbitWorkerManager *services.RabbitWorkerManager) *ResetConsumerHTTPHandler {

	return &ResetConsumerHTTPHandler{
		rabbitWorkerManager,
	}
}

func (parm *ResetConsumerHTTPHandler) Reset(ctx *gin.Context) {
	init := make(chan bool)
	go parm.RabbitWorkerManager.Restart(init)
	<-init
	ctx.JSON(http.StatusOK, "already reset")
}
