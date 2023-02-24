package transform

import (
	"fmt"
	"testing"
)

func TestStringToX(t *testing.T) {
	a := StringToX[float64]("123.102123")
	b, ok := a.(float64)
	if !ok {
		t.Fatal("convert float64 error")
	}
	fmt.Println(b)
}
