package model

import "time"

type Community struct {
	Id int64 `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`
	//名称
	Name string `xorm:"varchar(30)" form:"name" json:"name"`
	//群主ID
	OwnerId int64 `xorm:"bigint(20)" form:"owner_id" json:"owner_id"`

	//群Logo
	Icon string `xorm:"varchar(250)" form:"icon" json:"icon"`

	//类型
	Cate int `xorm:"int(11)" form:"cate" json:"cate"`

	//描述
	Memo string `xorm:"varchar(120)" form:"mmmo" json:"mmmo"`

	//创建时间
	CreateAt time.Time `xorm:"datetime" form:"createat" json:"createat"`
}

const (
	COMMUNITY_CATE_COM = 0x01
)
