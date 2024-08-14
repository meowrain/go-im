package args

type ContactArg struct {
	PageArg
	UserId int64 `json:"userid" form:"userid"`
	DstId  int64 `json:"dstid" form:"dstid"`
}
