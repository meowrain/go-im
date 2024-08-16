package controllers

import (
	"fmt"
	"im/conf"
	"im/model"
	"im/service"
	"im/utils"
	"im/utils/meowlog"
	"math/rand"
	"net/http"
)

func LoginFunc(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	loginok := false
	userService := service.UserService{}
	var user *model.User
	var err error
	user, loginok, err = userService.Login(mobile, passwd)
	if err != nil {
		utils.RespFailed(writer, http.StatusInternalServerError, err.Error())
		return
	}
	if loginok {
		utils.RespOk(writer, user, "登录成功")
		return
	} else {
		utils.RespFailed(writer, http.StatusUnauthorized, "用户名或密码错误")
	}
}
func RegisterFunc(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	if mobile == "" || passwd == "" {
		utils.Resp(writer, nil, 400, "参数错误")
		return
	}
	nickname := fmt.Sprintf("user%06d", rand.Int31n(100000))

	userService := service.UserService{}
	user, err := userService.Register(mobile, passwd, nickname, "", model.SEX_UNKNOW)
	if err != nil {
		utils.RespFailed(writer, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Resp(writer, user, 200, "注册成功")
}
func GetUserInfo(writer http.ResponseWriter, request *http.Request) {

	err := request.ParseForm()
	if err != nil {
		return
	}
	userid := request.PostForm.Get("userid")
	user := model.User{}
	has, err := conf.DbEngine.Where("id = ?", userid).Get(&user)
	if err != nil {
		return
	}
	if !has {
		logger := meowlog.NewLogger("console", "error")
		logger.Info("用户不存在")
	}
	utils.RespOk(writer, user, "success")

}
