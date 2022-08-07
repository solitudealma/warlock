/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 11:29
 */

package system

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/solitudealma/warlock/global"
	"math/big"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/solitudealma/warlock/config"
	"github.com/solitudealma/warlock/model/system"
	"github.com/solitudealma/warlock/utils"
)

//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type UserService struct{}

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	querySql := "SELECT `username` FROM sys_users where username = ?"
	err = config.MasterDB.Get(&user, querySql, u.Username)
	if err == nil { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.NewV4()
	index, _ := rand.Int(rand.Reader, big.NewInt(4))

	avatars := []string{"https://pic.imgdb.cn/item/603750c35f4313ce25f92a35.jpg",
		"https://cdn.acwing.com/media/user/profile/photo/71847_lg_f844104f10.jpg",
		"https://pic.imgdb.cn/item/5ed326c1c2a9a83be5ca4f33.jpg",
		"https://cdn.acwing.com/media/user/profile/photo/21_lg_44ac082631.JPG"}
	querySql = "INSERT INTO sys_users (`created_at`, `updated_at`, `uuid`, `username`, `password`, " +
		"`avatar`) VALUES (?, ?, ?, ?, ?, ?)"
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	_, err = config.MasterDB.Exec(querySql, currentTime, currentTime, u.UUID, u.Username, u.Password,
		avatars[index.Uint64()])
	return u, err
}

//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (userInter *system.SysUser, err error) {

	var user system.SysUser
	querySql := "SELECT uuid, username, password, avatar FROM sys_users where username = ?"
	err = config.MasterDB.Get(&user, querySql, u.Username)
	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
	}

	return &user, err
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(username string) (user system.SysUser, err error) {
	var reqUser system.SysUser
	querySql := "SELECT `uuid`, `username`, `avatar`, `score` FROM `sys_users` where `username` = ?"
	err = config.MasterDB.Get(&reqUser, querySql, username)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return system.SysUser{}, err
	}
	return reqUser, err
}

func (userService *UserService) UpdateUserScore(username string, score int32) (err error) {
	querySql := "UPDATE sys_users set score = score +  ? where `username` = ?"
	res, err := config.MasterDB.Exec(querySql, score, username)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("%+v\n", err)
		return err
	}
	global.WlLog.Infof("update score %d row", affected)
	return nil
}
