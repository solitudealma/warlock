/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:54
 */

package utils

var (
	LoginVerify    = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "PasswordConfirm": {NotEmpty()}}
)
