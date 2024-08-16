package model

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           //消息ID
	UserId  int64  `json:"userid,omitempty" form:"userid"`   //发送消息的用户
	DstId   int64  `json:"dstid,omitempty" form:"dstid"`     //发送消息到的用户id
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         //发送信息的类型，群聊或者私聊
	Media   int    `json:"media,omitempty" form:"media"`     // 按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` //发送的消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         //图片预览连接
	Url     string `json:"url,omitempty" form:"url"`         //服务的	URL
	Memo    string `json:"memo,omitempty" form:"memo"`       // 简单的描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   //其它的附加数据，比如语音的长度，红包的金额
}
