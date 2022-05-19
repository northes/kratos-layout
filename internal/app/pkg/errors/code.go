package errors

import (
	"net/http"
)

type ErrorCode int

const (
	SuccessCode          ErrorCode = 10000 // 成功
	ServerErrorCode      ErrorCode = 10001 // 服务器出错
	ClientErrorCode      ErrorCode = 10002 // 客户端请求错误
	ValidateErrorCode    ErrorCode = 10003 // 参数校验错误
	UnauthorizedCode     ErrorCode = 10004 // 未经授权
	PermissionDeniedCode ErrorCode = 10005 // 没有权限
	ResourceNotFoundCode ErrorCode = 10006 // 资源不存在
	TooManyRequestCode   ErrorCode = 10007 // 请求过于频繁
)

var httpStatusCode = map[ErrorCode]int{
	SuccessCode:          http.StatusOK,
	ServerErrorCode:      http.StatusInternalServerError,
	ClientErrorCode:      http.StatusBadRequest,
	ValidateErrorCode:    http.StatusBadRequest,
	UnauthorizedCode:     http.StatusUnauthorized,
	PermissionDeniedCode: http.StatusForbidden,
	ResourceNotFoundCode: http.StatusNotFound,
	TooManyRequestCode:   http.StatusTooManyRequests,
}

func (c ErrorCode) HTTPStatusCode() int {
	return httpStatusCode[c]
}

var statusCodeText = map[ErrorCode]string{
	SuccessCode:          "OK",
	ServerErrorCode:      "服务器出错",
	ClientErrorCode:      "客户端请求错误",
	ValidateErrorCode:    "参数校验错误",
	UnauthorizedCode:     "未经授权",
	PermissionDeniedCode: "暂无权限",
	ResourceNotFoundCode: "资源不存在",
	TooManyRequestCode:   "请求过于频繁",
}

func (c ErrorCode) String() string {
	return statusCodeText[c]
}
