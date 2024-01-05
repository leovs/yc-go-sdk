// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package cache

import (
	"fmt"
	_const "github.com/leovs/yc-go-sdk/const"
	"github.com/leovs/yc-go-sdk/utils"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// GenInstanceId 生成实例ID
func GenInstanceId() string {
	charList := []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().Unix())
	length := 8
	str := make([]byte, 0)
	for i := 0; i < length; i++ {
		str = append(str, charList[rand.Intn(len(charList))])
	}
	return string(str)
}

// GenSearchCacheKey 生成搜索缓存key
func GenSearchCacheKey(tableName string, sql string, vars ...interface{}) string {
	buf := strings.Builder{}
	buf.WriteString(sql)
	for _, v := range vars {
		pv := reflect.ValueOf(v)
		if pv.Kind() == reflect.Ptr {
			buf.WriteString(fmt.Sprintf(" {value} %v", pv.Elem()))
		} else {
			buf.WriteString(fmt.Sprintf(" {value} %v", v))
		}
	}
	return fmt.Sprintf("%s:%s:%s", _const.GormCachePrefix, tableName, utils.MD5(buf.String()))
}
