package controllers

import (
	"encoding/json"
	"im/model"
	"im/service"
	"im/utils/meowlog"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
)

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

const (
	//点对点单聊
	CMD_SINGLE_MSG = 10

	//群聊消息
	CMD_GROUP_MSG = 11

	//心跳消息
	CMD_HEART = 0
)
const (
	//文本样式
	MEDIA_TYPE_TEXT = iota + 1

	//图文消息
	MEDIA_TYPE_NEWS

	//语音
	MEDIA_TYPE_VOICE

	//纯图片
	MEDIA_TYPE_IMAGE

	// 红包
	MEDIA_TYPE_REDPACKAGE

	//emoji表情样式
	MEDIA_TYPE_EMOJ

	//超链接样式
	MEDIA_TYPE_LINK

	// 视频样式
	MEDIA_TYPE_VIDEO

	//名片样式
	MEDIA_TYPE_CONTACT

	//自定义，前端做解析就可以了
	MEDIA_TYPE_UDEF = 100
)

var clientMap map[int64]*Node = make(map[int64]*Node, 0)
var logger = meowlog.NewLogger("console", "error")
var rwlocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isValid := checkToken(userId, token)
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValid
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		logger.Error("%v", err.Error())
		return
	}
	// 打印连接成功的消息
	logger.Info("WebSocket connection established for userId: %d", userId)
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()
	go sendproc(node)
	go recvproc(node)

}

func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error("%v", err.Error())
				return
			}
		}
	}
}
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			logger.Error("%v", err.Error())
			return
		}
		logger.Debug("recv<==%s", data)
		dispatch(data)
	}
}

func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}
func checkToken(userId int64, token string) bool {
	userService := service.UserService{}
	user := userService.Find(userId)
	return user.Token == token
}
func dispatch(data []byte) {
	msg := model.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		logger.Error("%v", err.Error())
		return
	}
	checkMsgMedia := func(mediaType int) string {
		switch mediaType {
		case MEDIA_TYPE_CONTACT: //名片样式
			sendMsg(msg.DstId, data)
			return "contact"
		case MEDIA_TYPE_IMAGE: //图片
			return "image"
		case MEDIA_TYPE_EMOJ: //EMOJI
			return "emoji"
		case MEDIA_TYPE_LINK: //LINK
			return "link"
		case MEDIA_TYPE_VIDEO: //VIDEO
			return "video"
		case MEDIA_TYPE_NEWS: //图文消息
			sendMsg(msg.DstId, data)
			return "news"
		case MEDIA_TYPE_REDPACKAGE: //红包
			sendMsg(msg.DstId, data)
			return "redpackage"
		case MEDIA_TYPE_VOICE: //语音
			sendMsg(msg.DstId, data)
			return "voice"
		case MEDIA_TYPE_TEXT: //文字
			sendMsg(msg.DstId, data)
			return "text"
		case MEDIA_TYPE_UDEF: //未定义
			return "undefined"
		default:
			return "don't know what's this"
		}
	}
	checkMsgType := func(msg *model.Message) {
		var mediaType string
		switch msg.Cmd {
		case CMD_SINGLE_MSG:
			mediaType = checkMsgMedia(msg.Media)
			logger.Debug("single:%v: %v to %v %v", mediaType, msg.UserId, msg.DstId, msg.Content)
		case CMD_GROUP_MSG:
			mediaType = checkMsgMedia(msg.Media)
			logger.Debug("group:%v: %v to %v %v", mediaType, msg.UserId, msg.DstId, msg.Content)
		case CMD_HEART:
			mediaType = checkMsgMedia(msg.Media)
			logger.Debug("heart from: %v", msg.UserId)
		default:

		}
	}

	checkMsgType(&msg)
}
