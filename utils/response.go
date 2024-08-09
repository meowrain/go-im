package utils

import (
	"encoding/json"
	"net/http"
)

type HttpResponse struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func Resp(writer http.ResponseWriter, data any, code int, msg string) {
	dt := HttpResponse{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
	res, err := json.Marshal(dt)
	if err != nil {
		errorResp := HttpResponse{
			Code: 500,
			Data: nil,
			Msg:  "服务器内部错误",
		}
		errRespBytes, _ := json.Marshal(errorResp)
		writer.Write([]byte(errRespBytes))
		return
	}
	writer.Write(res)
}

func RespFailed(writer http.ResponseWriter, code int, msg string) {
	Resp(writer, nil, code, msg)
}

func RespOk(writer http.ResponseWriter, data any, msg string) {
	Resp(writer, data, http.StatusOK, msg)
}
