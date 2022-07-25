/**
 * @Author: SolitudeAlma
 * @Date: 2022 2022/7/12 10:50
 */

package utils

import (
	"crypto/md5"
	"encoding/hex"
)

//@function: MD5V
//@description: md5加密
//@param: str []byte
//@return: string

func MD5V(str []byte, b ...byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(b))
}
