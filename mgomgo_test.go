package mgomgo

import "testing"

func assertEqual(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Fatalf("%s != %s", s1, s2)
	}
}

func TestNewDBParamsFromURIFull(t *testing.T) {
	uri := "mongodb://user-hoge:pass01@host-aaa_bbb.co.jp:1111/db_name11/col_name00"
	params, err := NewDBParamsFromURI(uri)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, params.Host, "host-aaa_bbb.co.jp:1111")
	assertEqual(t, params.UserName, "user-hoge")
	assertEqual(t, params.Password, "pass01")
	assertEqual(t, params.Database, "db_name11")
	assertEqual(t, params.Collection, "col_name00")
}

func TestNewDBParamsFromURISimple(t *testing.T) {
	uri := "mongodb://hoge.com/db/col"
	params, err := NewDBParamsFromURI(uri)
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, params.Host, "hoge.com")
	assertEqual(t, params.Database, "db")
	assertEqual(t, params.Collection, "col")
}
