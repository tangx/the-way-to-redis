package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Server() {
	r := gin.Default()

	r.GET("/sms/get/:phone", handlerGetSms)
	r.POST("/sms/verify/:phone", handlerVerifySmsCode)
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func handlerGetSms(c *gin.Context) {
	phone := c.Param("phone")

	if len(phone) == 0 {
		return
	}

	// 3.1. 从 redis 中获取验证码 `userSmsCodeTimes:13912341234` 发送次数
	times := getUserSmsCodeTimes(phone)

	// 3.2. 判断验证码发送次数是超过 3次， 超过则退出。
	if times > 2 {
		c.JSON(401, map[string]interface{}{
			"msg": "发送验证码次数超过 3 次， 请明天再试。",
		})
		return
	}

	// 1.1. 生成六位数验证码
	code := genSmsCode()

	// 1.2. 验证码保存到 redis 中，`userSmsCodeValue:13912341234`
	// 过期时间 2 分钟。
	if err := setUserSmsCode(phone, code); err != nil {
		c.JSON(500, map[string]interface{}{
			"msg": "生成或验证码失败， 请重试",
			"err": err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"msg":  "success",
		"code": code,
	})
}

func handlerVerifySmsCode(c *gin.Context) {
	phone := c.Param("phone")
	inputCode := c.Query("code")

	if inputCode == "" {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"msg": "request forbidden",
			"err": "code is required",
		})

		return
	}

	// 2.1. 从 redis 中获取验证码。
	smscode := getUserSmsCode(phone)

	// 2.2 校验 sms code， 错误返回 403
	if inputCode != smscode {
		c.JSON(http.StatusForbidden, map[string]interface{}{
			"msg": "request forbidden",
			"err": "input code error",
		})

		return
	}

	// 2.2 成功返回 200， 并删除 redis 中的验证码
	_ = delUserSmsCode(phone)
	c.JSON(http.StatusOK,
		map[string]interface{}{
			"msg": "success",
		},
	)

	return

}

// genSmsCode 生成 6 位数随机验证码
func genSmsCode() string {
	code := ""
	rand.Seed(time.Now().UnixMicro())

	for i := 0; i < 6; i++ {
		code += fmt.Sprint(rand.Intn(10))
	}

	return code
}
