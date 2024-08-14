package controllers

import (
	"im/args"
	"im/model"
	"im/service"
	"im/utils"
	"im/utils/meowlog"
	"net/http"
)

var contactService service.ContactService

func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	logger := meowlog.NewLogger("console", "info")
	utils.Bind(req, &arg)
	logger.Info("args:%v", arg.UserId)
	users := contactService.SearchFriend(arg.UserId)
	//fmt.Println(users)
	type Friends struct {
		Users []model.User `json:"users"`
		Len   int          `json:"len"`
	}
	utils.RespOk(w, Friends{Users: users, Len: len(users)}, "success")
}
func LoadCommunities(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	utils.Bind(req, &arg)
	communities := contactService.SearchCommunity(arg.UserId)
	type Communities struct {
		Communitys []model.Community
		Len        int
	}
	utils.RespOk(w, Communities{Communitys: communities, Len: len(communities)}, "success")
}
func JoinCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	utils.Bind(req, &arg)
	err := contactService.JoinCommunity(arg.UserId, arg.DstId)
	if err != nil {
		utils.RespFailed(w, 500, err.Error())
		return
	}
	utils.RespOk(w, nil, "success")
}

func AddFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	utils.Bind(req, &arg)
	err := contactService.AddFriend(arg.UserId, arg.DstId)
	if err != nil {
		utils.RespFailed(w, 500, err.Error())
		return
	}
	utils.RespOk(w, nil, "success")
}
