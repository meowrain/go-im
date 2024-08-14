package model

import "time"

type Contact struct {
	Id int64 `xorm:"pk autoincr bigint(20)" form:"id" json:"id"`

	//账号拥有者id
	OwnerId int64 `xorm:"bigint(20)" form:"owner_id" json:"owner_id"`

	//添加的人或者群的id
	DstId int64 `xorm:"bigint(20)" form:"dst_id" json:"dst_id"`

	//添加的类别 user或者community
	Cate int `xorm:"int(11)" form:"cate" json:"cate"`

	//描述
	Memo string `xorm:"varchar(120)" form:"memo" json:"memo"`

	//加好友的日期
	CreateAt time.Time `xorm:"datetime" form:"create_at" json:"create_at"`
}

const (
	CONCAT_CATE_USER      = 0x01
	CONCAT_CATE_COMMUNITY = 0x02
)
