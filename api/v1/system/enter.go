/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:46
 */

package system

import (
	"github.com/warlock-backend/service"
)

type ApiGroup struct {
	JwtApi
	BaseApi
	WSApi
}

var (
	jwtService  = service.ServiceGroupApp.SystemServiceGroup.JwtService
	userService = service.ServiceGroupApp.SystemServiceGroup.UserService
	wsService   = service.ServiceGroupApp.SystemServiceGroup.WSService
)
