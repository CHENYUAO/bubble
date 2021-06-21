package handle

import (
	"testing"
)

func TestToken(t *testing.T) {
	testGroups := []string{
		"chenyuao",
		"xiaohong",
		"xiaoming",
	}
	for _, v := range testGroups {
		tokenString, err := GenerateToken(v)
		if err != nil {
			t.Error(err)
		}
		mc, err := ParseToken(tokenString)
		if err != nil {
			t.Error(err)
		}
		if mc.Username != v {
			t.Errorf("解析出的用户名不一致，username:%v,got:%v", v, mc.Username)
		}
	}
}
