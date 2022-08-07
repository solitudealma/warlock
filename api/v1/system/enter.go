/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:46
 */

package system

import (
	"github.com/solitudealma/warlock/service"
)

type ApiGroup struct {
	JwtApi
	BaseApi
}

var (
	jwtService  = service.ServiceGroupApp.SystemServiceGroup.JwtService
	userService = service.ServiceGroupApp.SystemServiceGroup.UserService
)
