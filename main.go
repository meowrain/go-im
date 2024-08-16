package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"im/controllers"
	"log"
	"net/http"
	"strings"
)

func RegisterView() {
	//一次解析出全部模板
	tpl, err := template.ParseGlob("view/**/*")
	if nil != err {
		log.Fatal(err)
	}
	//通过for循环做好映射
	for _, v := range tpl.Templates() {
		tplname := v.Name()
		fmt.Println("HandleFunc     " + tplname)
		if strings.HasSuffix(v.Name(), ".shtml") {
			http.HandleFunc(tplname, func(w http.ResponseWriter,
				request *http.Request) {
				//
				fmt.Println("parse     " + v.Name() + "==" + tplname)
				err := tpl.ExecuteTemplate(w, tplname, nil)
				if err != nil {
					log.Fatal(err.Error())
				}
			})
		}
	}
}

func main() {
	http.HandleFunc("/user/login", controllers.LoginFunc)

	// Serve static assets
	fs := http.FileServer(http.Dir("./asset"))
	http.Handle("/asset/", http.StripPrefix("/asset", fs))
	http.HandleFunc("/user/register", controllers.RegisterFunc)
	http.HandleFunc("/user/info", controllers.GetUserInfo)
	http.HandleFunc("/contact/loadcommunity", controllers.LoadCommunities)
	http.HandleFunc("/contact/loadfriend", controllers.LoadFriend)
	http.HandleFunc("/contact/addcommunity", controllers.JoinCommunity)
	http.HandleFunc("/contact/addfriend", controllers.AddFriend)
	http.HandleFunc("/chat", controllers.Chat)
	RegisterView()
	http.ListenAndServe(":8080", nil)
}
