/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 22:42
 */

package router

import (
	"github.com/warlock-backend/router/system"
)

type RouterGroup struct {
	System system.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
