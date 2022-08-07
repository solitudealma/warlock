/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:46
 */

package system

import (
	"github.com/satori/go.uuid"
)

type SysUser struct {
	UUID     uuid.UUID `json:"uuid" db:"uuid"`         // 用户UUID
	Username string    `json:"username" db:"username"` // 用户登录名
	Password string    `json:"-" db:"password"`        // 用户登录密码
	Avatar   string    `json:"avatar" db:"avatar"`     // 用户头像
	Score    uint32    `json:"score" db:"score"`       // 用户分数
}
