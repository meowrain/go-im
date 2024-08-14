package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"im/conf"
	"im/model"
	"im/utils"
	"math/rand"
	"time"
)

type UserService struct{}

// Register 用户注册
func (userService *UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (*model.User, error) {
	salt := fmt.Sprintf("%06d", rand.Int31n(10000))
	passwd := utils.MakePasswd(plainpwd, salt)
	user := &model.User{
		Mobile:   mobile,
		Passwd:   passwd,
		Nickname: nickname,
		Avatar:   avatar,
		Sex:      sex,
		Salt:     salt,
		CreateAt: time.Now(),
		Token:    fmt.Sprintf("%08d", rand.Int31n(100000000)),
	}
	tmp := &model.User{}
	_, err := conf.DbEngine.Where("mobile =?", mobile).Get(tmp)
	if err != nil {
		return tmp, fmt.Errorf("数据库错误")
	}
	if tmp.Id > 0 {
		return tmp, fmt.Errorf("手机号已注册")
	} else {
		_, err = conf.DbEngine.Insert(user)
		if err != nil {
			return user, err
		}
	}
	return user, nil
}

// Login 用户登录函数
func (userService *UserService) Login(mobile, password string) (*model.User, bool, error) {
	user := &model.User{}
	_, err := conf.DbEngine.Where("mobile = ?", mobile).Get(user)
	if err != nil {
		return nil, false, err
	}

	if user.Id <= 0 {
		return nil, false, fmt.Errorf("用户不存在")
	}

	if isValid := utils.ValidatePasswd(password, user.Salt, user.Passwd); isValid {
		str := fmt.Sprintf("%d", time.Now().Unix())
		token := utils.MD5Encode(str)
		user.Token = token
		_, err := conf.DbEngine.Cols("token").Update(user)
		if err != nil {
			return nil, false, err
		}
		return user, true, nil
	} else {
		return nil, false, fmt.Errorf("密码错误")
	}
}

// Find 寻找用户
func (s *UserService) Find(userId int64) model.User {
	tmp := model.User{}
	conf.DbEngine.ID(userId).Get(&tmp)
	return tmp
}
