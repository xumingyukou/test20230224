package test

import (
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/sdk/art"
)

func TestArt(t *testing.T) {

	tree := art.New()

	tree.Insert(art.Key("Hi, I'm Key"), "Nice to meet you, I'm Value")
	value, found := tree.Search(art.Key("Hi, I'm Key"))
	if found {
		_ = value
	}

	tree.ForEach(func(node art.Node) bool {
		return true
	})

	for it := tree.Iterator(); it.HasNext(); {
		it.Next()
	}

	tree.Delete(art.Key("Hi, I'm Key"))

	tree.Insert(art.Key("hello world 123 !"), "123")
	tree.Insert(art.Key("hello world 1"), "1")
	tree.Insert(art.Key("hello world"), "world")
	tree.Insert(art.Key("hello"), "hello")
	tree.Insert(art.Key("abc"), "abc")
	tree.Insert(art.Key("ibc"), "ibc")

	k, v, b := tree.LongestSearch(art.Key("hello world 123 "))
	if string(k) != "hello world 1" || v.(string) != "1" || !b {
		fmt.Println(string(k), "->", v, b)
		t.Error("Longest search error")
	}

	i := 0
	expectK := []string{"hello", "hello world", "hello world 1", "hello world 123 !"}
	expectV := []string{"hello", "world", "1", "123"}

	tree.SearchPrefixCb(art.Key("hello world 123 !"), func(k art.Key, v art.Value) bool {
		if string(k) != expectK[i] || v.(string) != expectV[i] {
			t.Error("Search cb error")
		}
		i++
		return true
	})
}
