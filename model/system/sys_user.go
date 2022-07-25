/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:46
 */

package system

import (
	"github.com/satori/go.uuid"
	"time"
)

type SysUser struct {
	ID        uint      `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`   // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`   // 更新时间
	DeleteAt  time.Time `json:"-" db:"deleted_at default:''"` // 删除时间
	UUID      uuid.UUID `json:"uuid" db:"uuid"`               // 用户UUID
	Username  string    `json:"username" db:"username"`       // 用户登录名
	Password  string    `json:"-" db:"password"`              // 用户登录密码
	Avatar    string    `json:"avatar" db:"avatar"`           // 用户头像
}
