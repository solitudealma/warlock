/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/8 21:26
 */

package request

import (
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
)

// CustomClaims Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.RegisteredClaims
}

type BaseClaims struct {
	UUID     uuid.UUID
	ID       uint
	Username string
}
