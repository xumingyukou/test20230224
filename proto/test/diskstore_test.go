package test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/warmplanet/proto/go/sdk"
	"github.com/warmplanet/proto/go/sdk/broker/store"
)

type TestStoreConfig struct {
	Configs []sdk.StoreConfig `toml:"storage"`
}

// sequence基本功能测试
func TestSeq1(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)

	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	seq, _ := dms.CreateSequencer("test", 100, nil)

	for i := uint64(1); i <= 10000; i++ {
		s := seq.GetNextSequence(100)
		if s != i {
			t.Errorf("Unexpected sequence %d %d", s, i)
		}
	}

	seq.ResetSequence(0, 100)

	for i := uint64(1); i <= 100; i++ {
		s := seq.GetNextSequence(1)
		if s != i {
			t.Errorf("Unexpected sequence %d %d", s, i)
		}
	}
	seq.ResetSequence(0, 1)
	dms.Close()
}

// sequence奔溃测试
func TestSeq2(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	seq, _ := dms.CreateSequencer("test", 100, nil)

	var i uint64
	for i = 1; i <= 101; i++ {
		s := seq.GetNextSequence(100)
		if s != i {
			t.Errorf("Unexpected sequence %d %d", s, i)
		}
	}

	time.Sleep(1 * time.Second)
	dms.Close()

	//下次加载应该从200开始
	dms = store.NewDiskKvDb(cc.Configs[0].StoreConfig)
	seq, _ = dms.CreateSequencer("test", 100, nil)

	s := seq.GetNextSequence(100)
	if s != 201 {
		t.Errorf("Unexpected sequence %d not 201", s)
	}

	seq.ResetSequence(0, 100)
	dms.Close()
}

// sequence并发测试
func TestSeq3(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)
	seq, _ := dms.CreateSequencer("test", 100, nil)

	rs := make([][]uint64, 5)
	for i := 0; i < 5; i++ {
		rs[i] = make([]uint64, 0)
	}

	var wg sync.WaitGroup
	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(idx int) {
			var v uint64
			for v < 1000 {
				v = seq.GetNextSequence(100)
				rs[idx] = append(rs[idx], v)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	m := make(map[uint64]bool, 0)

	for i := 0; i < 5; i++ {
		last := uint64(0)
		for _, v := range rs[i] {
			if v <= last {
				t.Errorf("not monotone increasing for g:%d %d %d", i, v, last)
			}
			if _, ok := m[v]; ok {
				t.Errorf("Dupliacate sequence %d", v)
			} else {
				m[v] = true
			}
		}
	}

	if len(m) != 1004 {
		t.Error("Invalid length 1004")
	}

	for i := uint64(1); i <= 1004; i++ {
		delete(m, i)
	}

	if len(m) != 0 {
		t.Error("Invalid length 0")
	}

	seq.ResetSequence(0, 100)

	wg.Add(100)
	// 创建100个key，并发getNextSequence
	for i := 0; i < 100; i++ {
		go func(idx int) {
			s, _ := dms.CreateSequencer("test"+strconv.Itoa(idx), 100, nil)
			for j := uint64(1); j < 1000; j++ {
				v := s.GetNextSequence(100)
				if v != j {
					t.Errorf("Mismatched sequence %d %d", v, j)
				}
			}
			s.ResetSequence(0, 100)
			wg.Done()
		}(i)
	}

	wg.Wait()
	dms.Close()
}

type IndexerNone struct {
}

func (i *IndexerNone) IndexNames() map[string]bool {
	return map[string]bool{}
}

func (i *IndexerNone) Index(value interface{}, rawValue []byte) store.MultiIndex {
	return store.MultiIndex{}
}

func (i *IndexerNone) UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int {
	return 0
}

// 测试表/索引管理
func TestDb(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	indexer := &IndexerNone{}
	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			meta, _ := dms.CreateTable("test", "dt.depth", indexer, nil)
			if meta.(*store.DiskTable).TableId != 1 {
				t.Error("Table id should be 1")
			}
			wg.Done()
		}()
	}

	wg.Wait()

	meta, _ := dms.CreateTable("test1", "dt.depth", indexer, nil)
	tbl := meta.(*store.DiskTable)

	// 创建index
	for i := 0; i < 100; i++ {
		ii := "i" + strconv.Itoa(i)
		idx, _ := tbl.CreateIndex("i"+strconv.Itoa(i), true)
		if idx.IndexId != store.Crc16(ii) {
			t.Error("Index id error")
		}
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			meta, _ := tbl.CreateIndex("i_100", false)
			if meta.IndexId != store.Crc16("i_100") || meta.Unique {
				t.Error("Index id should be 101")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	tmp, _ := tbl.CreateIndex("i19", false)
	if tmp == nil || tmp.IndexId != store.Crc16("i19") || !tmp.Unique {
		t.Error("Get i19 failed")
	}

	Db := dms.(*store.DiskKvDb).Db

	tid := append([]byte(store.META_PFX+"t_"), []byte("test")...)
	Db.Delete(tid, pebble.Sync)
	k := []byte(store.SEQ_PFX + store.GLB_TBL_SEQ)
	Db.Delete(k, pebble.Sync)

	tid = append([]byte(store.SEQ_PFX+"i_"), []byte("test")...)
	Db.Delete(tid, pebble.Sync)

	batch := Db.NewBatch()
	for i := 0; i < 101; i++ {
		k := []byte(store.META_PFX + "i_")
		k = append(k, []byte{0, 0, 0, 1}...)
		k = append(k, []byte("i"+strconv.Itoa(i))...)
		batch.Delete(k, nil)
	}

	batch.Commit(pebble.Sync)
	dms.Close()
}

type TestRow struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
}

type IndexerJson struct {
}

func (i *IndexerJson) IndexNames() map[string]bool {
	return map[string]bool{"id": true, "name": false}
}

func (i *IndexerJson) Index(value interface{}, rawValue []byte) (indexOut store.MultiIndex) {
	var encoder store.IndexEncoder

	var v *TestRow
	if value == nil {
		tmp := TestRow{}
		json.Unmarshal(rawValue, &tmp)
		v = &tmp
	} else {
		v = value.(*TestRow)
	}

	indexOut.WriteInit(2)

	indexOut.Write("id", encoder.PackUint32(v.Id).Bytes(false))
	encoder.Reset()
	indexOut.Write("name", encoder.PackString(v.Name).Bytes(true))

	return
}

func (i *IndexerJson) UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int {
	var decoder store.IndexDecoder
	decoder.FromBytes(index)

	if indexName == "name" {
		var s string
		len1 := decoder.RemainLen()
		decoder.UnpackString(&s)
		if cb != nil {
			cb(0, s)
		}
		return len1 - decoder.RemainLen()
	} else if indexName == "id" {
		var id uint32
		decoder.UnpackUint32(&id)
		if cb != nil {
			cb(0, id)
		}
		return 4
	} else {
		panic("invalid index name " + indexName)
	}
}

func TestTable(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	indexer := &IndexerJson{}
	var encoder store.IndexEncoder

	meta, _ := dms.CreateTable("test", "dt.depth", indexer, nil)
	tbl := meta.(*store.DiskTable)

	for i := uint32(1); i < 100; i++ {
		v := TestRow{Id: i, Name: "abc"}
		vs, _ := json.Marshal(v)

		err := tbl.Set(nil, &v, vs)
		if err != nil {
			t.Error(err)
		}

		encoder.Reset()
		pk, r, err := tbl.Get(encoder.PackUint32(v.Id).Bytes(false), "id")
		if err != nil || !bytes.Equal(r, vs) {
			t.Error("Get by id error")
		}
		if binary.BigEndian.Uint64(pk) != uint64(i) {
			t.Errorf("Primary key error %d %d", pk, i)
		}
		encoder.Reset()
		// 非唯一索引，get应该返回nil
		pk, r, _ = tbl.Get(encoder.PackString(v.Name).Bytes(false), "name")
		if pk != nil || r != nil {
			t.Error("Get by name should return nil")
		}
	}

	for i := uint32(1); i < 100; i++ {
		encoder.Reset()
		kk := encoder.PackUint32(i).Bytes(true)
		err := tbl.Delete(kk, "id")
		if err != nil {
			t.Error(err)
		}

		_, r, err := tbl.Get(kk, "id")
		if r != nil || err != nil {
			t.Error(r, err)
		}
	}

	s, _ := dms.CreateSequencer(store.TBL_PFX+"test", 100, nil)
	s.ResetSequence(0, 100)
	dms.Close()
}

func TestTableBatch(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	indexer := &IndexerJson{}
	var encoder store.IndexEncoder

	meta, _ := dms.CreateTable("test", "dt.depth", indexer, nil)
	tbl := meta.(*store.DiskTable)
	i := uint32(1)

	err := tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i >= 100 {
			return nil, nil, nil
		}
		pk = nil
		v = &TestRow{Id: i, Name: "abc"}
		rv, _ = json.Marshal(v)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	for i := uint32(1); i < 100; i++ {
		encoder.Reset()
		pk, r, err := tbl.Get(encoder.PackUint32(i).Bytes(false), "id")
		if err != nil || r == nil {
			t.Error("Get by id error")
		}

		tr := TestRow{}
		json.Unmarshal(r, &tr)

		if tr.Id != i {
			t.Error("Get by id error wrong id")
		}

		if binary.BigEndian.Uint64(pk) != uint64(i) {
			t.Error("Primary key error")
		}
		encoder.Reset()
		// 非唯一索引，get应该返回nil
		pk, r, _ = tbl.Get(encoder.PackString("abc").Bytes(false), "name")
		if pk != nil || r != nil {
			t.Error("Get by name should return nil")
		}
	}

	// 批量删除
	i = 1
	err = tbl.BatchDelete(func() (key []byte, index string) {
		i++
		if i >= 200 {
			return nil, ""
		}
		encoder.Reset()
		if i%2 == 0 {
			return encoder.PackUint32(i + 1000).Bytes(true), "id"
		} else {
			return encoder.PackUint32(i/2 + 1).Bytes(true), "id"
		}
	})

	if err != nil {
		t.Error(err)
	}

	for i := uint32(2); i < 100; i++ {
		encoder.Reset()
		kk := encoder.PackUint32(i).Bytes(true)

		_, r, err := tbl.Get(kk, "id")
		if r != nil || err != nil {
			t.Error(r, err)
		}
	}

	encoder.Reset()
	_, r, err := tbl.Get(encoder.PackUint32(1).Bytes(true), "id")
	if r == nil || err != nil {
		t.Error(r, err)
	}

	tbl.Delete(encoder.PackUint32(1).Bytes(true), "id")

	time.Sleep(1 * time.Second)
	s, _ := dms.CreateSequencer(store.TBL_PFX+"test", 100, nil)
	s.ResetSequence(0, 100)
	dms.Close()
}

func TestTableScan(t *testing.T) {
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	indexer := &IndexerJson{}

	meta, _ := dms.CreateTable("test", "dt.depth", indexer, nil)
	tbl := meta.(*store.DiskTable)
	i := uint32(1)

	err := tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i > 1000 {
			return nil, nil, nil
		}
		pk = nil
		v = &TestRow{Id: i, Name: "abc" + strconv.Itoa(int(i)%100)}
		rv, _ = json.Marshal(v)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	ii := uint32(1)
	f := func(pk []byte, value []byte) error {
		v := TestRow{}

		if json.Unmarshal(value, &v) != nil {
			t.Error("Unmarshal error")
		}

		if v.Id != ii {
			fmt.Println(v.Id, ii, pk)
			//t.Errorf("Invalid value id %d %d", v.Id, ii)
		}
		ii = ii + 1
		return nil
	}

	tbl.Scan([]byte{0, 0, 0, 0}, []byte{0xff, 0xff, 0xff, 0xff}, "id", f)
	i = uint32(1)
	err = tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i > 1000 {
			return nil, nil, nil
		}
		pk = nil
		v = &TestRow{Id: i, Name: "abc" + strconv.Itoa(int(i)%100)}
		v2 := &TestRow{Id: i + 1000, Name: "abc" + strconv.Itoa(int(i)%100)}
		rv, _ = json.Marshal(v2)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	ii = 1001
	tbl.Scan([]byte{0, 0, 0, 0}, []byte{0xff, 0xff, 0xff, 0xff}, "id", f)

	ii = 1
	f2 := func(pk []byte, value []byte) error {
		v := TestRow{}

		if json.Unmarshal(value, &v) != nil {
			t.Error("Unmarshal error")
		}

		//fmt.Println(v.Name, ii, pk)
		ii = ii + 1
		return nil
	}

	tbl.Scan([]byte(""), []byte("xx"), "name", f2)

	if ii != 2001 {
		t.Error("Scan non-unique error")
	}

	time.Sleep(1 * time.Second)
	s, _ := dms.CreateSequencer(store.TBL_PFX+"test", 100, nil)
	s.ResetSequence(0, 100)

	dms.Close()
}

type IndexerJson2 struct {
}

func (i *IndexerJson2) IndexNames() map[string]bool {
	return map[string]bool{"": true, "name": false}
}

func (i *IndexerJson2) Index(value interface{}, rawValue []byte) (indexOut store.MultiIndex) {
	var encoder store.IndexEncoder

	var v *TestRow
	if value == nil {
		tmp := TestRow{}
		json.Unmarshal(rawValue, &tmp)
		v = &tmp
	} else {
		v = value.(*TestRow)
	}

	indexOut.WriteInit(2)
	indexOut.Write("name", encoder.PackString(v.Name).Bytes(true))
	encoder.Reset()
	indexOut.Write("", encoder.PackUint32(v.Id).Bytes(false))
	return
}

func (i *IndexerJson2) UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int {
	var decoder store.IndexDecoder
	decoder.FromBytes(index)

	if indexName == "name" {
		var s string
		len1 := decoder.RemainLen()
		decoder.UnpackString(&s)
		if cb != nil {
			cb(0, s)
		}
		return len1 - decoder.RemainLen()
	} else if indexName == "" {
		var id uint32
		decoder.UnpackUint32(&id)
		if cb != nil {
			cb(0, id)
		}
		return 4
	} else {
		panic("invalid index name " + indexName)
	}
}

func TestTableSetPk(t *testing.T) {
	// 测试自定义indexer生成pk和Set接口传递pk
	cc := TestStoreConfig{}
	sdk.LoadConfigFile("../conf/storage.toml", &cc)
	cc.Configs[0].StoreConfig["seq_gap"] = 100
	dms := store.NewDiskKvDb(cc.Configs[0].StoreConfig)

	indexer := &IndexerJson2{}
	var encoder store.IndexEncoder

	meta, _ := dms.CreateTable("test2", "dt.depth", indexer, nil)
	tbl := meta.(*store.DiskTable)
	i := uint32(1)

	err := tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i >= 100 {
			return nil, nil, nil
		}
		pk = nil
		v = &TestRow{Id: i, Name: "abc"}
		rv, _ = json.Marshal(v)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	for i := uint32(1); i < 100; i++ {
		encoder.Reset()
		pk, r, err := tbl.Get(encoder.PackUint32(i).Bytes(false), "")
		if err != nil || r == nil {
			t.Error("Get by id error")
		}

		tr := TestRow{}
		json.Unmarshal(r, &tr)

		if tr.Id != i {
			t.Error("Get by id error wrong id")
		}

		if binary.BigEndian.Uint32(pk) != uint32(i) {
			t.Error("Primary key error")
		}
		encoder.Reset()
		// 非唯一索引，get应该返回nil
		pk, r, _ = tbl.Get(encoder.PackString("abc").Bytes(false), "name")
		if pk != nil || r != nil {
			t.Error("Get by name should return nil")
		}
	}

	i = uint32(1)
	err = tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i >= 50 {
			return nil, nil, nil
		}
		var kk [4]byte
		binary.BigEndian.PutUint32(kk[:], i)
		pk = kk[:]
		v = &TestRow{Id: i, Name: "abc"}
		v2 := &TestRow{Id: i + 1000, Name: "abc"}
		rv, _ = json.Marshal(v2)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	tr := TestRow{}
	for i := uint32(1); i < 50; i++ {
		encoder.Reset()
		pk, r, err := tbl.Get(encoder.PackUint32(i).Bytes(false), "")
		if err != nil || r == nil {
			t.Error("Get by id error")
		}

		json.Unmarshal(r, &tr)

		if tr.Id != i+1000 {
			t.Error("Get by id error wrong id")
		}

		if binary.BigEndian.Uint32(pk) != uint32(i) {
			t.Error("Primary key error")
		}
	}

	encoder.Reset()
	_, r, err := tbl.Get(encoder.PackUint32(66).Bytes(false), "")
	if err != nil || r == nil {
		t.Error("Get by id error")
	}

	json.Unmarshal(r, &tr)

	if tr.Id != 66 {
		t.Error("Get by id error wrong id")
	}

	encoder.Reset()
	tr.Id = 8888
	tt, _ := json.Marshal(&tr)

	if err = tbl.Set(encoder.PackUint32(8888).Bytes(false), &tr, tt); err != nil {
		t.Error(err)
	}

	encoder.Reset()
	_, r, err = tbl.Get(encoder.PackUint32(8888).Bytes(false), "")
	if err != nil || r == nil {
		t.Error("Get by id error")
	}

	json.Unmarshal(r, &tr)

	if tr.Id != 8888 {
		t.Error("Get by id error wrong id")
	}

	time.Sleep(1 * time.Second)
	dms.Close()
}
