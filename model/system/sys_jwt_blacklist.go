/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:39
 */

package system

import "time"

type JwtBlacklist struct {
	ID        uint      `db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"` // 更新时间
	Jwt       string    `json:"jwt" db:"jwt"`
}
