package mysql

import "testing"

func TestMysql(t *testing.T) {
	// ReadConf测试
	got, err := ReadConf("../conf/bubble.ini")
	want := "root:ch981205@tcp(127.0.0.1:3306)/bubble"
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got:%s\nwant:%s\n", got, want)
	}
	// 数据库初始化测试
	err = InitDB(got)
	defer DB.Close()
	if err != nil {
		t.Error(err)
	}
	// 用户认证测试
	users := []Users{
		{"chenyuao", "123456"},
		{"aaaaa", "bbbbb"},
	}
	errs := make([]error, len(users))
	for i := 0; i < len(users); i++ {
		errs[i] = AuthUser(users[i].UserName, users[i].Password)
	}
	if errs[0] != nil || errs[1] == nil {
		t.Error("Authority failed")
	}
}
