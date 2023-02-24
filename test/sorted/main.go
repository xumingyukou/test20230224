package main

import (
	"fmt"
	"github.com/zhangyunhao116/pdqsort"
	"math/rand"
	"sort"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixMicro())
}

type SortMode string

const (
	Random         SortMode = "Random"
	Sorted         SortMode = "Sorted"
	NearlySorted   SortMode = "NearlySorted"
	Reversed       SortMode = "Reversed"
	NearlyReversed SortMode = "NearlyReversed"
	AllEqual       SortMode = "AllEqual"
)

type TestSlice []int64

func (d TestSlice) Len() int           { return len(d) }
func (d TestSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d TestSlice) Less(i, j int) bool { return d[i] < d[j] }
func (d *TestSlice) Reverse() {
	for i := 0; i < d.Len()/2; i++ {
		(*d)[i], (*d)[d.Len()-i-1] = (*d)[d.Len()-i-1], (*d)[i]
	}
}

func GenerateData(mode SortMode, n int) TestSlice {
	switch mode {
	case Random:
		return generateRandomData(n)
	case Sorted:
		return generateSortedData(n)
	case NearlySorted:
		return generateNearlySortedData(n)
	case Reversed:
		return generateReversedData(n)
	case NearlyReversed:
		return generateNearlyReversedData(n)
	case AllEqual:
		return generateAllEqualData(n)
	default:
		return []int64{}
	}
}

func generateRandomData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, rand.Int63n(time.Now().UnixMicro()))
	}
	return res
}

func generateSortedData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, int64(i))
	}
	return res
}

func generateNearlySortedData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, int64(i))
	}
	for i := 0; i < 3; i++ {
		idx := rand.Intn(n)
		res[idx], res[n-idx] = res[n-idx], res[idx]
	}
	return res
}

func generateReversedData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, int64(n-i))
	}
	return res
}

func generateNearlyReversedData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, int64(n-i))
	}
	for i := 0; i < 3; i++ {
		idx := rand.Intn(n)
		res[idx], res[n-idx] = res[n-idx], res[idx]
	}
	return res
}

func generateAllEqualData(n int) TestSlice {
	var res TestSlice
	for i := 0; i < n; i++ {
		res = append(res, int64(n))
	}
	return res
}

func time_use(f func(data sort.Interface), input sort.Interface, title string) {
	start := time.Now()
	f(input)
	fmt.Println(title, time.Now().Sub(start))
}

func time_use2(input []int64, title string) {
	start := time.Now()
	pdqsort.Slice(input)
	fmt.Println(title, time.Now().Sub(start))
}

func main() {
	time_use(sort.Sort, GenerateData(Random, 1000), "sort Random")
	time_use(sort.Stable, GenerateData(Random, 1000), "Stable Random")
	time_use2(GenerateData(Random, 1000), "pdqsort Random")
	time_use(sort.Sort, GenerateData(Sorted, 1000), "sort Sorted")
	time_use(sort.Stable, GenerateData(Sorted, 1000), "Stable Sorted")
	time_use2(GenerateData(Sorted, 1000), "pdqsort Sorted")
	time_use(sort.Sort, GenerateData(NearlySorted, 1000), "sort NearlySorted")
	time_use(sort.Stable, GenerateData(NearlySorted, 1000), "Stable NearlySorted")
	time_use2(GenerateData(NearlySorted, 1000), "pdqsort NearlySorted")
	time_use(sort.Sort, GenerateData(Reversed, 1000), "sort Reversed")
	time_use(sort.Stable, GenerateData(Reversed, 1000), "Stable Reversed")
	time_use2(GenerateData(Reversed, 1000), "pdqsort Reversed")
	time_use(sort.Sort, GenerateData(NearlyReversed, 1000), "sort NearlyReversed")
	time_use(sort.Stable, GenerateData(NearlyReversed, 1000), "Stable NearlyReversed")
	time_use2(GenerateData(NearlyReversed, 1000), "pdqsort NearlyReversed")
	time_use(sort.Sort, GenerateData(AllEqual, 1000), "sort AllEqual")
	time_use(sort.Stable, GenerateData(AllEqual, 1000), "Stable AllEqual")
	time_use2(GenerateData(AllEqual, 1000), "pdqsort AllEqual")
}
