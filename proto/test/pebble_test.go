package test

import (
	"encoding/binary"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/warmplanet/proto/go/sdk"
)

func TestSimple(t *testing.T) {
	option := &pebble.Options{
		DisableWAL:      true,
		WALBytesPerSync: 64 * 1024,
	}

	db, err := pebble.Open("db", option)
	if err != nil {
		log.Fatal(err)
	}
	key := make([]byte, 8)
	var i uint64
	v := []byte("godbless")

	for i = 0; i < 100000; i++ {
		binary.BigEndian.PutUint64(key, i)
		if err := db.Set(key, v, pebble.NoSync); err != nil {
			log.Fatal(err)
		}
	}

	binary.BigEndian.PutUint64(key, 7777)
	value, closer, err := db.Get(key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(key, value)
	if err := closer.Close(); err != nil {
		log.Fatal(err)
	}

	if err := db.Close(); err != nil {
		log.Fatal(err)
	}
}

func upsert(exist bool, valueInMap int64, newValue int64) int64 {
	if !exist {
		return newValue
	}

	return valueInMap + 1
}

func TestCmap(t *testing.T) {
	m := sdk.NewCmap()

	for i := 0; i < 10; i++ {
		go m.Upsert("abc", 1, upsert)
	}

	time.Sleep(time.Second * 1)

	if v, _ := m.Get("abc"); v != 10 {
		t.Error("cmap get failed")
	}
}
