/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:45
 */

package v1

import (
	"github.com/solitudealma/warlock/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
