package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetRoutes(engine *gin.RouterGroup, resetConsumerHTTPHandler *ResetConsumerHTTPHandler) {
	engine.GET("/v1/resetConsumer", resetConsumerHTTPHandler.Reset)
	engine.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "start")
	})
}
