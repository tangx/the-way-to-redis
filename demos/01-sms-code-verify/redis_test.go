package main

import (
	"fmt"
	"testing"
)

func Test_getUserSmsCode(t *testing.T) {
	phone := `13912341234`
	times := getUserSmsCodeTimes(phone)

	fmt.Println(times)
}

func Test_ttlSeconds(t *testing.T) {
	sec := ttlSeconds()
	fmt.Println(sec)
}
