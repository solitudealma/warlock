/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:36
 */

package service

import (
	"github.com/warlock-backend/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup system.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
