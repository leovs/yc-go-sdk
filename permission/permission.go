// Copyright 2023 ztlcloud.com
// leovs @2023.12.12

package permission

import (
	"bytes"
	"fmt"
	"github.com/leovs/yc-go-sdk/utils"
	"math"
	"strings"
)

func Permission2Map(preStr string) *map[string]int {
	var m map[string]int
	m = make(map[string]int)
	pres := strings.Split(preStr, ",")
	for _, per := range pres {
		s := strings.Split(per, ":")
		if len(s) == 2 {
			m[s[0]] = utils.String2Int(s[1])
		}
	}
	return &m
}

func Map2Permission(preMap *map[string]int) string {
	var result bytes.Buffer
	for key, per := range *preMap {
		result.WriteString(fmt.Sprintf("%s:%d,", key, per))
	}
	return result.String()
}

func mergePermission(per1 *map[string]int, per2 *map[string]int) {
	for key, per := range *per2 {
		orgVal := (*per1)[key]
		if (orgVal & per) == 0 {
			(*per1)[key] = orgVal ^ per
		} else {
			(*per1)[key] = int(math.Max(float64(orgVal), float64(per)))
		}
	}
}

func MergePermission(rootPermissions string, permissions ...string) string {
	if len(permissions) > 0 {
		root := Permission2Map(rootPermissions)
		for _, permission := range permissions {
			mergePermission(root, Permission2Map(permission))
		}
		return Map2Permission(root)
	}
	return ""
}
