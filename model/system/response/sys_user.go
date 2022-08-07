/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:44
 */

package response

import (
	"github.com/solitudealma/warlock/model/system"
)

type SysUserResponse struct {
	User system.SysUser `json:"user"`
}

type LoginResponse struct {
	User      system.SysUser `json:"user"`
	Token     string         `json:"token"`
	ExpiresAt int64          `json:"expiresAt"`
}
