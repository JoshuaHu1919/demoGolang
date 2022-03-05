package errors

import (
	"fmt"
)

// AppError 用來處理 awardCenter app 的業務邏輯錯誤
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

// AppError 返回string型態的AppError內容
func (e AppError) Error() string {
	return fmt.Sprintf("%s-%s", e.Code, e.Message)
}

// New functions create a new appError instance
func New(code, message string) AppError {
	return AppError{Code: code, Message: message}
}
