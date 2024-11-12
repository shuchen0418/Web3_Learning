package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Print 打印
func Print(i interface{}) {
	fmt.Println("=======")
	fmt.Println(i)
	fmt.Println("=======")
}

// RetJSON 统一返回Json
func RetJSON(code, msg string, data interface{}, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}

// GetTimeUnix 获取当前时间戳
func GetTimeUnix() int64 {
	return time.Now().Unix()
}

//MD5方法

func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

//生成签名
//验证签名
