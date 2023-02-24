package base

import (
	"fmt"
	"testing"

	"github.com/warmplanet/proto/go/depth"
)

func TestMergeBitbank(t *testing.T) {
	depthCache := &OrderBook{
		Bids: []*depth.DepthLevel{
			&depth.DepthLevel{
				Price:  194027,
				Amount: 0.002,
			},
			&depth.DepthLevel{
				Price:  194014,
				Amount: 0.017,
			},
			&depth.DepthLevel{
				Price:  194013,
				Amount: 0.4634,
			},
			&depth.DepthLevel{
				Price:  194006,
				Amount: 1.644,
			},
			&depth.DepthLevel{
				Price:  194001,
				Amount: 0.0005,
			},
		},
		Asks: []*depth.DepthLevel{
			&depth.DepthLevel{
				Price:  194057,
				Amount: 0.918,
			},
			&depth.DepthLevel{
				Price:  194079,
				Amount: 0.2628,
			},
			&depth.DepthLevel{
				Price:  194080,
				Amount: 1.1211,
			},
			&depth.DepthLevel{
				Price:  194141,
				Amount: 0.017,
			},
			&depth.DepthLevel{
				Price:  194142,
				Amount: 20.4981,
			},
		},
	}

	deltaDepth := &DeltaDepthUpdate{
		Bids: []*depth.DepthLevel{
			&depth.DepthLevel{
				Price:  193509,
				Amount: 0.1,
			},
			&depth.DepthLevel{
				Price:  193287,
				Amount: 0.1,
			},
		},
		Asks: []*depth.DepthLevel{
			&depth.DepthLevel{
				Price:  194057,
				Amount: 0.459,
			},
			&depth.DepthLevel{
				Price: 194057,
			},
			&depth.DepthLevel{
				Price:  194079,
				Amount: 0.1314,
			},
			&depth.DepthLevel{
				Price: 194079,
			},
		},
	}

	d := &depth.Depth{}
	UpdateBidsAndAsks(deltaDepth, depthCache, 5, d)
	fmt.Println(d)
	fmt.Println(depthCache)
	for _, level := range depthCache.Asks {
		if level.Price == 194057 || level.Price == 194079 {
			t.Fatal("error merge")
		}
	}
}

func TestMergeDuplicate(t *testing.T) {
	deltaDepth := []*depth.DepthLevel{
		{Price: 1, Amount: 1},
		{Price: 1},
		{Price: 1, Amount: 2},
		{Price: 2, Amount: 2},
		{Price: 3, Amount: 2},
		{Price: 3},
		{Price: 5, Amount: 5},
	}
	expected := []*depth.DepthLevel{
		{Price: 1, Amount: 2},
		{Price: 2, Amount: 2},
		{Price: 3},
		{Price: 5, Amount: 5},
	}
	result := mergeDuplicate(deltaDepth)
	if len(result) != len(expected) {
		t.Fatal("incorrect length")
	}
	for i := 0; i < len(expected); i++ {
		if expected[i].Price != result[i].Price || expected[i].Amount != result[i].Amount {
			t.Fatal("incorrect price or amount")
		}
	}
}
