/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:41
 */

package request

// Register User register structure
type Register struct {
	Username        string `form:"username" json:"userName" db:"username"`
	Password        string `form:"password" json:"passWord" db:"password"`
	PasswordConfirm string `form:"password_confirm" json:"passwordConfirm"`
}

// Login User login structure
type Login struct {
	Username string `form:"username" json:"username" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"` // 密码
}
