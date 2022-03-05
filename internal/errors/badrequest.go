package errors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InvalidInputFunc(c *gin.Context, errMsg string) {
	errorMsg := AppError{
		Code:    InvalidInput,
		Message: errMsg,
	}
	c.JSON(http.StatusBadRequest, errorMsg)
}
