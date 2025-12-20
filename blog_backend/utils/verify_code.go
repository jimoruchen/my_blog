package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenVerifyCode 生成6位数字验证码
func GenVerifyCode() string {
	return fmt.Sprintf("%06v", rand.Int31n(1000000))
}
