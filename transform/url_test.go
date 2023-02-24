package transform

import "testing"

func TestQuote(t *testing.T) {
	uri := "testmove_virtual@b4xaf7pvnoemail.com"
	uri_quoted := Quote(uri)
	if uri_quoted != "testmove_virtual%40b4xaf7pvnoemail.com" {
		t.Error("url quote failed")
	}
	uri_unquoted := Unquote(uri)
	if uri_unquoted != uri {
		t.Error("url un quote failed")
	}
	t.Log("success")
}
