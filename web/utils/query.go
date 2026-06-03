 package utils

 import (
 	"strconv"

 	"github.com/gin-gonic/gin"
 )

 // PathInt64 读取 URL path 参数 :id (int64)。
 func PathInt64(ctx *gin.Context, key string) int64 {
 	val := ctx.Param(key)
 	if val == "" {
 		return 0
 	}
 	out, err := strconv.ParseInt(val, 10, 64)
 	if err != nil {
 		return 0
 	}
 	return out
 }

 // PathInt 读取 URL path 参数 (int)。
 func PathInt(ctx *gin.Context, key string) int {
 	v := PathInt64(ctx, key)
 	if v > int64(^uint(0)>>1) {
 		return 0
 	}
 	return int(v)
 }

 // QueryInt 读取 query 参数 (int)。
 func QueryInt(ctx *gin.Context, key string, def int) int {
 	val := ctx.Query(key)
 	if val == "" {
 		return def
 	}
 	out, err := strconv.Atoi(val)
 	if err != nil {
 		return def
 	}
 	return out
 }

 // QueryInt64 读取 query 参数 (int64)。
 func QueryInt64(ctx *gin.Context, key string, def int64) int64 {
 	val := ctx.Query(key)
 	if val == "" {
 		return def
 	}
 	out, err := strconv.ParseInt(val, 10, 64)
 	if err != nil {
 		return def
 	}
 	return out
 }
