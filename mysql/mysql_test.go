package mysql

import "testing"

func TestMysql(t *testing.T) {
	got, err := ReadConf("../conf/bubble.ini")
	want := "root:ch981205@tcp(127.0.0.1:3306)/bubble"
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("got:%s\nwant:%s\n", got, want)
	}
	err = InitDB(got)
	defer DB.Close()
	if err != nil {
		t.Error(err)
	}
}
