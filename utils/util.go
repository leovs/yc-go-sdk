// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package utils

import (
	"crypto/md5"
	"fmt"
	"strconv"
)

// MD5 MD5加密
func MD5(str string) string {
	data := []byte(str) //切片
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str
}

// String2Int 字符串转Int
func String2Int(intStr string) (intNum int) {
	intNum, _ = strconv.Atoi(intStr)
	return
}

// Ternary 弥补三元运算
func Ternary(c bool, s interface{}, p interface{}) interface{} {
	if c {
		return s
	} else {
		return p
	}
}
