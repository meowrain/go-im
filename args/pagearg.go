package args

import (
	"fmt"
	"time"
)

type PageArg struct {
	//从哪页开始
	PageFrom int `json:"pagefrom" form:"pagefrom"`
	//每页大小
	PageSize int `json:"pagesize" form:"pagesize"`
	//关键词
	Kword string `json:"kword" form:"kword"`
	//asc：“id”  id asc
	Asc  string `json:"asc" form:"asc"`
	Desc string `json:"desc" form:"desc"`
	//
	Name string `json:"name" form:"name"`
	//
	UserId int64 `json:"userid" form:"userid"`
	//dstid
	DstID int64 `json:"dstid" form:"dstid"`
	//时间点1
	DateFrom time.Time `json:"datafrom" form:"datafrom"`
	//时间点2
	DateTo time.Time `json:"dateto" form:"dateto"`
	//
	Total int64 `json:"total" form:"total"`
}

func (p *PageArg) GetPageSize() int {
	if p.PageSize == 0 {
		return 100
	} else {
		return p.PageSize
	}

}
func (p *PageArg) GetPageFrom() int {
	if p.PageFrom < 0 {
		return 0
	} else {
		return p.PageFrom
	}
}

func (p *PageArg) GetOrderBy() string {
	if len(p.Asc) > 0 {
		return fmt.Sprintf(" %s asc", p.Asc)
	} else if len(p.Desc) > 0 {
		return fmt.Sprintf(" %s desc", p.Desc)
	} else {
		return ""
	}
}
