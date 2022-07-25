/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:47
 */

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/warlock-backend/global"
	"github.com/warlock-backend/model/common/response"
	"github.com/warlock-backend/model/system"
)

type JwtApi struct{}

// JsonInBlacklist @Tags Jwt
// @Summary jwt加入黑名单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "jwt加入黑名单"
// @Router /jwt/jsonInBlacklist [post]
func (j *JwtApi) JsonInBlacklist(c *gin.Context) {
	token := c.Request.Header.Get("x-token")
	jwt := system.JwtBlacklist{Jwt: token}
	if err := jwtService.JsonInBlacklist(jwt); err != nil {
		global.WlLog.Errorf("jwt作废失败! err: %v", err)
		response.FailWithMessage("jwt作废失败", c)
	} else {
		response.OkWithMessage("jwt作废成功", c)
	}
}
