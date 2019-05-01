package yin

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// 计算路由hash
func getRouterHash(method string, path string) string {
	hashStr := fmt.Sprintf("%s-%s-yin", method, path)
	h := md5.New()
	h.Write([]byte(hashStr))
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
