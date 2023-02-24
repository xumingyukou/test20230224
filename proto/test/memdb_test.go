package test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/warmplanet/proto/go/sdk/broker/store"
)

func TestMemBuffer(t *testing.T) {
	mb := store.MemBuffer{}
	i := 1
	mb.Init(3, func(item *store.MemItem, retired bool) {
		if item.Value.(int) != i {
			t.Errorf("Retire error %d %d", item.Value, i)
		}
		i++
	})

	for i := 1; i < 100; i++ {
		mm := store.MultiIndex{}
		mm.WriteInit(2)
		mm.Write("a", []byte(strconv.Itoa(i)))
		mm.Write("b", []byte("xxx"))
		mm.Write("", []byte(strconv.Itoa(i)))
		mb.AddItem(i, nil, mm)
	}

	j := 97
	mb.IterItem(func(item *store.MemItem) error {
		if j != item.Value.(int) {
			t.Errorf("Iter error %d %d", j, item.Value)
		}
		j++
		return nil
	}, false)

	j = 99
	mb.IterItem(func(item *store.MemItem) error {
		if j != item.Value.(int) {
			t.Errorf("Iter error %d %d", j, item.Value)
		}
		j--
		return nil
	}, true)

	i = 99
	mb.DelItem(3)

	j = 97
	mb.IterItem(func(item *store.MemItem) error {
		id, pk := item.Keys.Read(0)
		if id == 0 && string(pk) != strconv.Itoa(j) {
			t.Errorf("Iter error %d %s", j, string(pk))
		}
		if j != item.Value.(int) {
			t.Errorf("Iter error %d %d", j, item.Value)
		}
		j++
		return nil
	}, false)
}

// 并发测试
func TestMemBuffer2(t *testing.T) {
	var mtx sync.Mutex
	mb := store.MemBuffer{}
	c := int32(0)
	mb.Init(20, func(item *store.MemItem, retired bool) {
		atomic.AddInt32(&c, 1)
	})

	var wg sync.WaitGroup
	wg.Add(10)

	s := time.Now().UnixMicro()
	for j := 0; j < 10000; j++ {
		mm := store.MultiIndex{}
		mm.WriteInit(2)
		mm.Write("a", []byte(strconv.Itoa(10000+j)))
		mm.Write("b", []byte("xxx"))
		mb.AddItem(10000+j, nil, mm)
	}
	fmt.Printf("cost %d\n", time.Now().UnixMicro()-s)

	for i := 0; i < 10; i++ {
		go func(i int) {
			s := time.Now().UnixMicro()
			for j := 0; j < 10000; j++ {
				mm := store.MultiIndex{}
				mm.WriteInit(2)
				mm.Write("a", []byte(strconv.Itoa(i*10000+j)))
				mm.Write("b", []byte("xxx"))
				mtx.Lock()
				mb.AddItem(i*10000+j, nil, mm)
				mtx.Unlock()
			}
			fmt.Printf("thread %d cost %d\n", i, time.Now().UnixMicro()-s)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestMultiIndex(t *testing.T) {
	var mi store.MultiIndex
	mi.WriteInit(5)

	mi.Write("aaa", []byte("a"))
	mi.Write("bbb", []byte("bbb"))
	mi.Write("ccc", []byte("ccc"))
	mi.Write("ddd", []byte("ddd"))
	mi.Write("", []byte("pppkkk"))

	mi.Write("aaa", []byte("aaa"))

	for i := uint16(0); i < 10; i++ {
		id, index := mi.Read(i)
		if len(index) == 0 {
			break
		}
		if id != 0 && id != store.Crc16(string(index)) {
			t.Error("Invalid id ")
		}
	}
}

func (*IndexerJson) IndexType(name string) int {
	if name == "name" {
		return store.MEM_TREE_INDEX
	} else {
		return store.MEM_HASH_INDEX
	}
}

func TestMemDb(t *testing.T) {
	dbConfig := map[string]interface{}{"buffer_size": 10, "index_type": "hash"}

	di := uint32(1)
	xx := func(item *store.MemItem, retired bool) {
		if di != item.Value.(*TestRow).Id {
			t.Error("Unexpected di counter")
		}
		di++
	}

	tblConfig := map[string]interface{}{"buffer_size": 10, "buffer_delete_cb": (store.MemDeleteCb)(xx)}

	db := store.NewMemKvDb(dbConfig)
	indexer := &IndexerJson{}

	tbl, err := db.CreateTable("test", "dt.dpeth", indexer, tblConfig)
	if tbl == nil || err != nil {
		t.Error("Create table test failed")
	}

	var encoder store.IndexEncoder

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

	// 前90个已经淘汰
	for i := uint32(1); i < 90; i++ {
		encoder.Reset()
		kk := encoder.PackUint32(i).Bytes(true)

		_, r, err := tbl.Get(kk, "id")
		if r != nil || err != nil {
			t.Error(string(r), err)
		}
	}

	if di != 90 {
		t.Error("Invalid delete di")
	}

	for i := uint32(90); i < 100; i++ {
		encoder.Reset()
		kk := encoder.PackUint32(i).Bytes(true)

		_, r, err := tbl.Get(kk, "id")
		if r == nil || err != nil {
			t.Error(string(r), err)
		}

		err = tbl.Delete(kk, "id")
		if err != nil {
			t.Error(err)
		}
	}

	if di != 100 {
		t.Error("Invalid delete di")
	}

	stat := make(map[string]interface{})
	json.Unmarshal([]byte(tbl.Stat()), &stat)

	if stat["mem_buffer_len"] != 0.0 || stat["non_unique_index_number"] != 0.0 || stat["unique_index_number"] != 0.0 {
		t.Errorf("Invalid stat %v", stat)
	}
}

func TestMemDbBatch(t *testing.T) {
	dbConfig := map[string]interface{}{"buffer_size": 1000, "index_type": "hash"}
	di := uint32(1)
	xx := func(item *store.MemItem, retired bool) {
		di++
	}
	tblConfig := map[string]interface{}{"buffer_size": 1000, "buffer_delete_cb": (store.MemDeleteCb)(xx)}
	db := store.NewMemKvDb(dbConfig)
	indexer := &IndexerJson{}

	tbl, err := db.CreateTable("test", "dt.dpeth", indexer, tblConfig)
	if tbl == nil || err != nil {
		t.Error("Create table test failed")
	}

	var encoder store.IndexEncoder
	i := uint32(1)
	err = tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
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
}

func TestMemDbAbnormal(t *testing.T) {
	// 测试pk冲突
	dbConfig := map[string]interface{}{"buffer_size": 1000, "index_type": "hash"}
	di := uint32(1)
	xx := func(item *store.MemItem, retired bool) {
		di++
	}
	tblConfig := map[string]interface{}{"buffer_size": 1000, "buffer_delete_cb": (store.MemDeleteCb)(xx)}
	db := store.NewMemKvDb(dbConfig)
	indexer := &IndexerJson{}

	tbl, err := db.CreateTable("test", "dt.dpeth", indexer, tblConfig)
	if tbl == nil || err != nil {
		t.Error("Create table test failed")
	}

	//var encoder store.IndexEncoder

	v := TestRow{Id: 1, Name: "abc"}
	vs, _ := json.Marshal(v)
	err = tbl.Set(nil, &v, vs)
	if err != nil {
		t.Error(err)
	}

	var encoder store.IndexEncoder
	pk, r, _ := tbl.Get(encoder.PackUint32(1).Bytes(false), "id")
	if !bytes.Equal(pk, []byte{0, 0, 0, 0, 0, 0, 0, 1}) || string(r) != `{"id":1,"name":"abc"}` {
		t.Error("Get by id error")
	}

	// 测试唯一索引覆盖
	v2 := TestRow{Id: 1, Name: "abcde"}
	vs, _ = json.Marshal(v2)
	err = tbl.Set(nil, &v2, vs)
	if err != nil {
		t.Error(err)
	}

	encoder.Reset()
	pk, r, _ = tbl.Get(encoder.PackUint32(1).Bytes(false), "id")
	if !bytes.Equal(pk, []byte{0, 0, 0, 0, 0, 0, 0, 2}) || string(r) != `{"id":1,"name":"abcde"}` {
		t.Error("Get by id error")
	}

	// 测试主键覆盖
	v2 = TestRow{Id: 2, Name: "abcdefg"}
	vs, _ = json.Marshal(v2)
	err = tbl.Set([]byte{0, 0, 0, 0, 0, 0, 0, 2}, &v2, vs)
	if err != nil {
		t.Error(err)
	}

	encoder.Reset()
	pk, r, _ = tbl.Get(encoder.PackUint32(2).Bytes(false), "id")
	if !bytes.Equal(pk, []byte{0, 0, 0, 0, 0, 0, 0, 2}) || string(r) != `{"id":2,"name":"abcdefg"}` {
		t.Error("Get by id error")
	}

	v = TestRow{Id: 3, Name: "aaa"}
	vs, _ = json.Marshal(v)
	err = tbl.Set(nil, &v, vs)
	if err != nil {
		t.Error(err)
	}

	// 测试唯一索引冲突
	v = TestRow{Id: 3, Name: "bbb"}
	vs, _ = json.Marshal(v)
	err = tbl.Set([]byte{0, 0, 0, 0, 0, 0, 0, 2}, &v, vs)
	if err == nil {
		t.Error("Should conflict")
	}

	v = TestRow{Id: 2, Name: "ccc"}
	vs, _ = json.Marshal(v)
	err = tbl.Set([]byte{0, 0, 0, 0, 0, 0, 0, 3}, &v, vs)
	if err == nil {
		t.Error("Should conflict")
	}
}

// 简单的性能测试，m1上1个主键+3个索引插入1us左右
type IndexerJson1 struct {
}

func (i *IndexerJson1) IndexNames() map[string]bool {
	return map[string]bool{"id": true, "id2": true, "nameId": true}
}

func (i *IndexerJson1) Index(value interface{}, rawValue []byte) (indexOut store.MultiIndex) {
	var encoder store.IndexEncoder
	v := value.(*TestRow)
	// 自增主键和3个唯一索引
	indexOut.WriteInit(3)
	indexOut.Write("id", encoder.PackUint32(v.Id).Bytes(false))
	//encoder.Reset()
	//indexOut.Write("id2", encoder.PackUint32(v.Id).Bytes(false))
	encoder.Reset()
	indexOut.Write("nameId", []byte(v.Name+strconv.FormatInt(int64(v.Id), 10)))

	return
}

func (i *IndexerJson1) UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int {
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

// hash大概是tree用时的1/4
func TestMemDbPerf1(t *testing.T) {
	dbConfig := map[string]interface{}{"buffer_size": 10, "index_type": "hash"}

	di := uint32(1)
	xx := func(item *store.MemItem, retired bool) {
		di++
	}

	tblConfig := map[string]interface{}{"buffer_size": 10, "buffer_delete_cb": (store.MemDeleteCb)(xx)}

	db := store.NewMemKvDb(dbConfig)
	indexer := &IndexerJson1{}

	tbl, err := db.CreateTable("test", "dt.dpeth", indexer, tblConfig)
	if tbl == nil || err != nil {
		t.Error("Create table test failed")
	}

	testRows := make([]TestRow, 0)
	testRowsRaw := make([][]byte, 0)
	for i := uint32(1); i < 10000; i++ {
		v := TestRow{Id: i, Name: "abc"}
		vs, _ := json.Marshal(v)

		testRows = append(testRows, v)
		testRowsRaw = append(testRowsRaw, vs)
	}

	start := time.Now().UnixMicro()
	for i := 0; i < len(testRows); i++ {
		err := tbl.Set(nil, &testRows[i], testRowsRaw[i])
		if err != nil {
			t.Error(err)
		}
	}
	fmt.Println(time.Now().UnixMicro() - start)

	var wg sync.WaitGroup
	wg.Add(10)

	start = time.Now().UnixMicro()
	for j := 0; j < 10; j++ {
		go func() {
			for i := 0; i < len(testRows); i++ {
				err := tbl.Set(nil, &testRows[i], testRowsRaw[i])
				if err != nil {
					t.Error(err)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println(time.Now().UnixMicro() - start)
}

func TestMapPerf1(t *testing.T) {
	testRows := make([]TestRow, 0)
	for i := uint32(1); i < 10000; i++ {
		v := TestRow{Id: i, Name: "abc"}

		testRows = append(testRows, v)
	}
	var mtx sync.Mutex

	m1 := make(map[uint32]*TestRow, 0)
	m2 := make(map[uint32]*TestRow, 0)
	m3 := make(map[string]*TestRow, 0)

	start := time.Now().UnixMicro()
	for i := 0; i < len(testRows); i++ {
		mtx.Lock()
		m1[testRows[i].Id] = &testRows[i]
		m2[testRows[i].Id] = &testRows[i]
		m3["abc"+strconv.Itoa(i)] = &testRows[i]
		mtx.Unlock()
	}
	fmt.Println(time.Now().UnixMicro() - start)
}

func TestMemDbScan(t *testing.T) {

	dbConfig := map[string]interface{}{"buffer_size": 3000, "index_type": "hash"}
	di := uint32(1)
	xx := func(item *store.MemItem, retired bool) {
		di++
	}
	tblConfig := map[string]interface{}{"buffer_size": 3000, "buffer_delete_cb": (store.MemDeleteCb)(xx)}
	db := store.NewMemKvDb(dbConfig)
	indexer := &IndexerJson{}

	tbl, err := db.CreateTable("test", "dt.dpeth", indexer, tblConfig)
	if tbl == nil || err != nil {
		t.Error("Create table test failed")
	}

	//var encoder store.IndexEncoder
	i := uint32(1)
	err = tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i > 1000 {
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

	f2 := func(pk []byte, value interface{}) error {
		v := value.(*TestRow)

		if v.Id != ii {
			fmt.Println(v.Id, ii, pk)
			//t.Errorf("Invalid value id %d %d", v.Id, ii)
		}
		ii = ii + 1
		return nil
	}

	tbl.Scan([]byte("abc"), []byte("abd"), "name", f)
	if ii != 1001 {
		t.Error("scan error")
	}

	ii = 1
	tbl.ScanValue([]byte("abc"), []byte("abd"), "name", f2)
	if ii != 1001 {
		t.Error("scan error")
	}

	i = 100
	err = tbl.BatchSet(func() (pk []byte, v interface{}, rv []byte) {
		if i >= 1000 {
			return nil, nil, nil
		}
		pk = nil
		v = &TestRow{Id: i, Name: "abd" + strconv.Itoa(int(i))}
		rv, _ = json.Marshal(v)
		i++
		return
	})

	if err != nil {
		t.Error(err)
	}

	ii = 100
	f3 := func(pk []byte, value interface{}) error {
		v := value.(*TestRow)

		if v.Name != "abd"+strconv.Itoa(int(ii)) || binary.BigEndian.Uint64(pk) != uint64(ii+901) {
			t.Errorf("Invalid value id %d %d", v.Id, ii)
		}

		ii = ii + 1
		return nil
	}

	tbl.ScanValue([]byte("abd"), []byte("abe"), "name", f3)

	stat := make(map[string]interface{})
	json.Unmarshal([]byte(tbl.Stat()), &stat)

	if stat["mem_buffer_len"] != 1000.0 || stat["non_unique_index_number"] != 1000.0 || stat["unique_index_number"] != 2000.0 {
		t.Errorf("Invalid stat %v", stat)
	}
}
