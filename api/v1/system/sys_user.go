/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:43
 */

package system

import (
	"github.com/solitudealma/warlock/global"
	"github.com/solitudealma/warlock/model/common/response"
	"github.com/solitudealma/warlock/model/system"
	systemReq "github.com/solitudealma/warlock/model/system/request"
	systemRes "github.com/solitudealma/warlock/model/system/response"
	"github.com/solitudealma/warlock/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type BaseApi struct{}

// Login @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body systemReq.Login true "用户名, 密码, 验证码"
// @Success 200 {object} response.Response{data=systemRes.LoginResponse,msg=string} "返回包括用户信息,token,过期时间"
// @Router /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	_ = c.ShouldBindQuery(&l)
	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	u := &system.SysUser{Username: l.Username, Password: l.Password}
	if user, err := userService.Login(u); err != nil {
		global.WlLog.Errorf("登陆失败! 用户名不存在或者密码错误!, err: %v", err)
		response.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		b.TokenNext(c, *user)
	}
}

// TokenNext 登录以后签发jwt
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	j := &utils.JWT{SigningKey: []byte(global.WlConfig.JWT.SigningKey)} // 唯一签名
	claims := j.CreateClaims(systemReq.BaseClaims{
		UUID:     user.UUID,
		Username: user.Username,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.WlLog.Errorf("获取token失败! err: %v", err)
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.WlConfig.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		global.WlLog.Infof("error: %v", err)
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.WlLog.Errorf("设置登录状态失败! err: %v", err)
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		global.WlLog.Info("GetRedisJWT")
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.WlLog.Errorf("设置登录状态失败! err: %v", err)
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}

		global.WlLog.Info("none")
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "登录成功", c)
	}
}

// Register @Tags SysUser
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body systemReq.Register true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {object} response.Response{data=systemRes.SysUserResponse,msg=string} "用户注册账号,返回包括用户信息"
// @Router /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	_ = c.ShouldBindQuery(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if strings.Compare(r.Password, r.PasswordConfirm) != 0 {
		response.FailWithMessage("两次输入的密码不一致", c)
	}

	user := &system.SysUser{Username: r.Username, Password: r.Password}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.WlLog.Errorf("注册失败! err: %v", err)
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败", c)
	} else {
		response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
	}
}

func (b *BaseApi) GetUserInfo(ctx *gin.Context) {
	platform := ctx.Query("platform")
	if platform == "ACAPP" {
		//getinfo_acapp(ctx)
	} else if platform == "WEB" {
		username := utils.GetUsername(ctx)
		if ReqUser, err := userService.GetUserInfo(username); err != nil {
			global.WlLog.Errorf("获取失败! err: %v", err)
			response.FailWithMessage("获取失败", ctx)
		} else {
			response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "success", ctx)
		}
	}
}
