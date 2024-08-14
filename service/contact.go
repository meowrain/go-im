package service

import (
	"errors"
	"im/conf"
	"im/model"
	"time"
)

type ContactService struct{}

// 自动添加好友
func (service *ContactService) AddFriend(userid, dstid int64) error {
	if userid == dstid {
		return errors.New("不能添加自己为好友")
	}
	tmp := model.Contact{}
	conf.DbEngine.Where("owner_id = ?", userid).And("dst_id = ?", dstid).And("cate = ?", model.CONCAT_CATE_USER).Get(&tmp)
	if tmp.Id > 0 {
		return errors.New("该用户已经被添加过了")
	}
	var userExists bool
	conf.DbEngine.Table("user").Where("id = ?", dstid).Exist(&userExists)
	if !userExists {
		return errors.New("添加的用户不存在！！")
	}
	session := conf.DbEngine.NewSession()
	session.Begin()
	_, e2 := session.InsertOne(model.Contact{
		OwnerId:  userid,
		DstId:    dstid,
		Cate:     model.CONCAT_CATE_USER,
		CreateAt: time.Now(),
	})
	_, e3 := session.InsertOne(model.Contact{
		OwnerId:  dstid,
		DstId:    userid,
		Cate:     model.CONCAT_CATE_USER,
		CreateAt: time.Now(),
	})
	if e2 != nil {
		session.Rollback()
		session.Close()
		return e2
	}
	if e3 != nil {
		session.Rollback()
		session.Close()
		return e3
	}
	session.Commit()

	return nil

}

// SearchCommunity 方法用于根据用户ID搜索群
func (service *ContactService) SearchCommunity(userId int64) []model.Community {
	// 初始化一个空的Contact切片，用于存储查询到的联系人信息
	concats := make([]model.Contact, 0)
	// 初始化一个空的int64切片，用于存储社区ID
	comIds := make([]int64, 0)
	// 使用ORM查询方式，根据用户ID和社区类别查询联系人，并将结果存储到concats切片中
	conf.DbEngine.Where("owner_id = ? and cate = ?", userId, model.CONCAT_CATE_COMMUNITY).Find(&concats)

	// 遍历concats切片，将每个联系人的目标群ID添加到comIds切片中
	for _, v := range concats {
		comIds = append(comIds, v.DstId)
	}
	// 初始化一个空的Community切片，用于存储查询到的社区信息
	coms := make([]model.Community, 0)
	// 如果comIds切片为空，说明没有找到对应的社区，直接返回空的社区切片
	if len(comIds) == 0 {
		return coms
	}
	// 使用ORM查询方式，根据社区ID列表查询社区信息，并将结果存储到coms切片中
	conf.DbEngine.In("id", comIds).Find(&coms)
	// 返回查询到的社区信息切片
	return coms
}

// JoinCommunity 加群
func (service *ContactService) JoinCommunity(userId int64, communityId int64) error {
	cot := model.Contact{
		OwnerId: userId,
		DstId:   communityId,
		Cate:    model.CONCAT_CATE_COMMUNITY,
	}
	conf.DbEngine.Get(&cot)
	if cot.Id == 0 {
		cot.CreateAt = time.Now()
		_, err := conf.DbEngine.InsertOne(cot)
		return err
	} else {
		return nil
	}
}

func (service *ContactService) CreateCommunity(comm model.Community) (model.Community, error) {
	//检查要创建的群名是否为空
	if len(comm.Name) == 0 {
		return model.Community{}, errors.New("缺少群名称")
	}
	//检查用户是否登录
	if comm.OwnerId == 0 {
		return model.Community{}, errors.New("请先登录")
	}
	// 计算用户已经有的群组数量
	com := model.Community{OwnerId: comm.OwnerId}
	num, err := conf.DbEngine.Count(&com)
	if err != nil {
		return model.Community{}, err
	}
	// 设置可以创建的最大群数量
	if num >= 5 {
		return model.Community{}, errors.New("一个用户最多只能创建5个群")
	}
	//设置创建时间
	comm.CreateAt = time.Now()
	session := conf.DbEngine.NewSession()
	defer session.Close()

	err = session.Begin()
	if err != nil {
		return model.Community{}, err
	}
	// 插入新的社区群组
	_, err = session.InsertOne(&comm)
	if err != nil {
		session.Rollback()
		return model.Community{}, err
	}

	contact := model.Contact{
		OwnerId:  comm.OwnerId,
		DstId:    comm.Id,
		Cate:     model.CONCAT_CATE_COMMUNITY,
		CreateAt: time.Now(),
	}
	// 插入联系人记录
	_, err = session.InsertOne(&contact)
	if err != nil {
		session.Rollback() // 插入失败，回滚事务
		return model.Community{}, err
	}

	err = session.Commit()
	if err != nil {
		return model.Community{}, err
	}

	return comm, nil
}

// 查找朋友
func (service *ContactService) SearchFriend(userId int64) []model.User {
	conconts := make([]model.Contact, 0)
	objIds := make([]int64, 0)
	conf.DbEngine.Where("owner_id = ? and cate = ?", userId, model.CONCAT_CATE_USER).Find(&conconts)
	for _, v := range conconts {
		objIds = append(objIds, v.DstId)
	}
	coms := make([]model.User, 0)
	if len(objIds) == 0 {
		return coms
	}
	conf.DbEngine.In("id", objIds).Find(&coms)
	return coms
}
