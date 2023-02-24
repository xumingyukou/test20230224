package store

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/warmplanet/proto/go/sdk"
	"github.com/warmplanet/proto/go/sdk/skiplist"
)

/*
	基于循环buffer和skiplist实现的内存db，支持索引，适合作为临时数据的缓存，会自动删除老的消息
	每个表使用一个循环buffer存储数据，每个索引使用单独的hash/skiplist存储，都存储指向buffer元素的指针
	skiplist支持唯一和非唯一索引，hash只支持唯一索引
	索引名称使用Crc16映射到uint16的indexId，主键的indexId为0

	主键(没有主键生自动成row id)：pk/row id[8字节] 值：BufferItem pos
	唯一索引：索引值 值：BufferItem pos
	非唯一索引：索引值 + pk/row id 值：BufferItem pos
	主键和二级索引都要满足字节序可比性，遍历的顺序应该和用户设置的顺序一致，比如整数索引，
	应该负数在前，正数在后，字符串索引应当按照字典序。

	NOTE: 受限MultiIndex中keys的编码格式，每个pk/索引的长度不能超过255字节
*/

const MEM_HASH_INDEX int = 0 // 默认pk和唯一索引使用hash index
const MEM_TREE_INDEX int = 1

// createTabl的args参数
const MEM_ARGS_BUFFER_SIZE string = "buffer_size"        // 值应该为int类型，表示表的buffer大小
const MEM_ARGS_BUFFER_DEL_CB string = "buffer_delete_cb" // 值应该为MemDeleteCb类型，表示buffer淘汰元素调用的回调

var NULL_MI_INSTANCE MultiIndex

type MemKvDbConfig struct {
	// 此内存引擎默认buffer大小
	BufferSize int `toml:"buffer_size" json:"buffer_size"`
	// 默认索引类型, hash/tree
	IndexType string `toml:"index_type" json:"index_type"`
}

type MemIndex struct {
	Name      string
	IndexId   uint16
	Unique    bool
	TreeIndex bool // 默认基于HashIndex
}

// 循环buffer+skiplist，用于内存中缓存和索引消息
type MemTable struct {
	MemBuffer

	// value 存储的是membuffer中的postion
	ti skiplist.SkipList // treed index
	hi map[string]int32  // hash index
	// 目前读写单表都需要加锁，建议一个表由一个go routine单独操作，没有多线程冲突效率最高
	// 或者使用分表方法，提高整体并行度
	mtx sync.Mutex

	indexs map[uint16]*MemIndex

	Indexer KvIndexer
	Ds      MemSequencer

	Id    uint16 // 每个table有一个唯一id
	Name  string
	Class string
}

type MemSequencer struct {
	Name    string
	counter uint64
}

type MemDeleteCb func(item *MemItem, retired bool)

type MemItem struct {
	NextPos int32       // 下一个item在数组中的位置，没有下一个为-1
	PrevPos int32       // 前一个item在数组中的位置，没有前一个为-1
	Value   interface{} //原始消息体
	Data    []byte      // 序列化后的字节数组
	Keys    MultiIndex  // 包含pk和所有index的值
}

// 循环数组，非线程安全
type MemBuffer struct {
	buf []MemItem
	// 空闲链表，新释放的放在链表头部
	free int32

	// 当buffer满且空闲链表为空时，会将item链表的第一个淘汰，同时调用此回调函数
	OnDelete MemDeleteCb
}

func (mb *MemBuffer) Init(maxSize int32, onDelete MemDeleteCb) {
	// buf[0]为元素链表头节点
	mb.buf = make([]MemItem, 1+maxSize)
	mb.OnDelete = onDelete

	for i := int32(1); i <= maxSize; i++ {
		mb.buf[i].NextPos = (i + 1) % (maxSize + 1)
	}

	// 元素列表为空
	mb.buf[0].NextPos = 0
	mb.buf[0].PrevPos = 0
	mb.free = 1
}

func (mb *MemBuffer) AddItem(value interface{}, data []byte, keys MultiIndex) int32 {
	idx := mb.free

	// 还有空闲元素
	if idx != 0 {
		mb.free = mb.buf[idx].NextPos
		last := mb.buf[0].PrevPos
		// 添加到元素链表尾部
		mb.buf[idx] = MemItem{NextPos: 0, PrevPos: last, Value: value, Data: data, Keys: keys}
		mb.buf[last].NextPos = idx
		mb.buf[0].PrevPos = idx
		return idx
	} else {
		// 没有空闲元素，淘汰head指向的元素并插入到链尾
		first := mb.buf[0].NextPos
		last := mb.buf[0].PrevPos
		if first == 0 {
			panic("error first item")
		}

		next := mb.buf[first].NextPos
		if mb.OnDelete != nil {
			mb.OnDelete(&mb.buf[first], true)
		}
		mb.buf[0].NextPos = next
		mb.buf[0].PrevPos = first
		mb.buf[last].NextPos = first
		mb.buf[next].PrevPos = 0

		mb.buf[first] = MemItem{NextPos: 0, PrevPos: last, Value: value, Data: data, Keys: keys}
		return first
	}
}

func (mb *MemBuffer) DelItem(i int32) {
	prev := mb.buf[i].PrevPos
	next := mb.buf[i].NextPos

	mb.buf[prev].NextPos = next
	mb.buf[next].PrevPos = prev

	if mb.OnDelete != nil {
		mb.OnDelete(&mb.buf[i], false)
	}

	mb.buf[i].NextPos = mb.free
	mb.buf[i].PrevPos = 0
	mb.free = i
}

func (mb *MemBuffer) GetItem(i int32) *MemItem {
	return &mb.buf[i]
}

func (mb *MemBuffer) IterItem(cb func(item *MemItem) error, reverse bool) {
	if !reverse {
		for i := mb.buf[0].NextPos; i != 0; i = mb.buf[i].NextPos {
			if cb(&mb.buf[i]) != nil {
				break
			}
		}
	} else {
		for i := mb.buf[0].PrevPos; i != 0; i = mb.buf[i].PrevPos {
			if cb(&mb.buf[i]) != nil {
				break
			}
		}
	}
}

// NOTE：此方法一般只在统计时使用，暂无优化必要
func (mb *MemBuffer) Len() int {
	count := 0
	for i := mb.buf[0].NextPos; i != 0; i = mb.buf[i].NextPos {
		count++
	}
	return count
}

type MemKvDb struct {
	Cfg  MemKvDbConfig
	Name string
	// pk和唯一索引可以选择索引类型，非唯一索引只支持treeIndex
	IndexType int // 使用hashIndex后，table不支持scan接口

	// 用于生成table id
	ds MemSequencer

	tables sdk.ConcurrentMapI // string -> *MemTable
}

func (md *MemKvDb) Open() error {
	// 如果用户没设置db默认值
	if md.Cfg.BufferSize <= 0 {
		md.Cfg.BufferSize = 1024 // 默认buffer大小
	}

	md.tables = sdk.NewCmapI()

	return nil
}

func (md *MemKvDb) Close() error {
	return nil
}

func (md *MemKvDb) Stat() string {
	return ""
}

func (md *MemKvDb) CreateTable(name string, class string, indexer KvIndexer, args map[string]interface{}) (table KvTable, err error) {
	// 从缓存加载
	if t, ok := md.tables.Get(name); ok {
		return t.(KvTable), nil
	}

	id := md.ds.GetNextSequence(0)
	if id >= math.MaxUint16 {
		panic("table id exceed uint16 limit")
	}

	tbl := &MemTable{Name: name, Class: class, Id: uint16(id), Indexer: indexer,
		hi: make(map[string]int32), ti: *skiplist.New(skiplist.Bytes), indexs: make(map[uint16]*MemIndex)}

	var cb MemDeleteCb

	if v, ok := args[MEM_ARGS_BUFFER_DEL_CB]; ok {
		cb = func(item *MemItem, retired bool) {
			n := uint16(0)
			var pk []byte
			for idIndex := item.Keys.Read2(n); len(idIndex) != 0; idIndex = item.Keys.Read2(n) {
				indexId := binary.BigEndian.Uint16(idIndex)
				if indexId == 0 {
					pk = idIndex[2:]
				}

				tbl.delIndex(indexId, idIndex, pk)
				n++
			}
			v.(MemDeleteCb)(item, retired)
		}
	} else {
		cb = func(item *MemItem, retired bool) {
			n := uint16(0)
			var pk []byte
			for idIndex := item.Keys.Read2(n); len(idIndex) != 0; idIndex = item.Keys.Read2(n) {
				indexId := binary.BigEndian.Uint16(idIndex)
				if indexId == 0 {
					pk = idIndex[2:]
				}

				tbl.delIndex(indexId, idIndex, pk)
				n++
			}
		}
	}

	if v, ok := args[MEM_ARGS_BUFFER_SIZE]; ok {
		tbl.MemBuffer.Init(int32(v.(int)), cb)
	} else {
		tbl.MemBuffer.Init(int32(md.Cfg.BufferSize), cb)
	}

	if indexer != nil {
		memIndexer, ok := indexer.(MemKvIndexer)
		// 如果传入的是memIndexer，可以获取某个index的类型，否则使用数据库配置的默认索引类型
		f := func(name string) int {
			if ok {
				return memIndexer.IndexType(name)
			} else {
				return md.IndexType
			}
		}

		for k, v := range indexer.IndexNames() {
			indexId := Crc16(k)
			if _, ok := tbl.indexs[indexId]; ok {
				panic("Create index failed: index id conflict")
			}

			if !v && f(k) != MEM_TREE_INDEX {
				panic("Only tree index support non unique index")
			}

			// 创建索引
			tbl.indexs[indexId] = &MemIndex{Name: k, IndexId: uint16(indexId), Unique: v, TreeIndex: f(k) == MEM_TREE_INDEX}
		}

		if _, ok := tbl.indexs[0]; !ok {
			tbl.indexs[0] = &MemIndex{Unique: true, TreeIndex: f("") == MEM_TREE_INDEX}
		}
	}

	if md.tables.SetIfAbsent(name, tbl) {
		return tbl, nil
	} else {
		t, _ := md.tables.Get(name)
		return t.(*MemTable), nil
	}
}

func (md *MemKvDb) CreateSequencer(name string, initCacheSize uint64, args map[string]interface{}) (sequencer KvSequencer, err error) {
	s := &MemSequencer{Name: name}
	return s, nil
}

func (md *MemKvDb) Backup(url string) error {
	return nil
}

func (ms *MemSequencer) GetNextSequence(cacheSize uint64) uint64 {
	return atomic.AddUint64(&ms.counter, 1)
}

func (ms *MemSequencer) ResetSequence(initSeq uint64, cacheSize uint64) error {
	atomic.StoreUint64(&ms.counter, initSeq)
	return nil
}

func (mt *MemTable) get(key []byte, index string) (pk []byte, value interface{}, data []byte) {
	k := make([]byte, 2+len(key))
	indexId := Crc16(index)

	mi, ok := mt.indexs[indexId]

	if !ok {
		return nil, nil, nil
	}

	binary.BigEndian.PutUint16(k, mi.IndexId)
	copy(k[2:], key)

	pos := int32(0)
	if !mi.TreeIndex {
		if v, ok := mt.hi[string(k)]; ok {
			pos = v
		}
	} else {
		if v, ok := mt.ti.GetValue(k); ok {
			pos = v.(int32)
		}
	}

	if pos != 0 {
		item := mt.MemBuffer.GetItem(pos)
		id, pk := item.Keys.Read(0)
		if id == 0 {
			return pk, item.Value, item.Data
		} else {
			return nil, item.Value, item.Data
		}
	}

	return nil, nil, nil
}

func (mt *MemTable) Get(key []byte, index string) (pk []byte, data []byte, err error) {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	pk, _, data = mt.get(key, index)
	return
}

func (mt *MemTable) GetValue(key []byte, index string) (pk []byte, value interface{}, err error) {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	pk, value, _ = mt.get(key, index)
	return
}

// 查询索引内容
func (mt *MemTable) getIndex(indexId uint16, index []byte, pk []byte) int32 {
	mi := mt.indexs[indexId]
	if mi == nil {
		// error
		panic("Invalid index id in " + mt.Name)
	}

	if mi.TreeIndex {
		var elem *skiplist.Element
		if mi.Unique {
			elem = mt.ti.Get(index)
		} else {
			tt := make([]byte, 0, len(index)+len(pk))
			tt = append(tt, index...)
			elem = mt.ti.Get(append(tt, pk...))
		}

		if elem != nil {
			return elem.Value.(int32)
		}
		return 0
	} else {
		// hash索引不支持非唯一索引
		return mt.hi[string(index)]
	}
}

// 删除一个元素所有的索引
func (mt *MemTable) delIndex(indexId uint16, index []byte, pk []byte) {
	mi := mt.indexs[indexId]
	if mi == nil {
		// error
		//panic("Invalid index id in " + mt.Name)
		return
	}

	if mi.TreeIndex {
		if mi.Unique {
			mt.ti.Remove(index)
		} else {
			tt := make([]byte, 0, len(index)+len(pk))
			tt = append(tt, index...)
			mt.ti.Remove(append(tt, pk...))
		}
	} else {
		// hash索引不支持非唯一索引
		delete(mt.hi, string(index))
	}
}

func (mt *MemTable) addIndex(indexId uint16, index []byte, pos int32, pk []byte) {
	mi := mt.indexs[indexId]
	if mi == nil {
		// error
		panic("Invalid index id in " + mt.Name)
	}

	if mi.TreeIndex {
		if mi.Unique {
			mt.ti.Set(index, pos)
		} else {
			tt := make([]byte, 0, len(index)+len(pk))
			tt = append(tt, index...)
			mt.ti.Set(append(tt, pk...), pos)
		}
	} else {
		mt.hi[string(index)] = pos
	}
}

// 更新index，通过归并算法删除old中无用的索引，保留有用的，同时添加new中新的索引
func (mt *MemTable) update(old *MultiIndex, pos int32, new *MultiIndex) {
	var oldIndex, newIndex []byte
	var i, j uint16
	var oldPk, newPk []byte

	for oldIndex, newIndex = old.Read2(i), new.Read2(j); oldIndex != nil && newIndex != nil; oldIndex, newIndex = old.Read2(i), new.Read2(j) {
		indexIdOld, indexIdNew := binary.BigEndian.Uint16(oldIndex), binary.BigEndian.Uint16(newIndex)
		if indexIdOld == 0 {
			oldPk = oldIndex[2:]
		}

		if indexIdNew == 0 {
			newPk = newIndex[2:]
		}

		if indexIdOld < indexIdNew {
			// 删除此索引
			mt.delIndex(indexIdOld, oldIndex, oldPk)
			i++
		} else if indexIdOld == indexIdNew {
			// indexId相同但是内容不容，需要删除老的
			if !bytes.Equal(oldIndex, newIndex) {
				mt.delIndex(indexIdOld, oldIndex, oldPk)
				mt.addIndex(indexIdNew, newIndex, pos, newPk)
			}
			i, j = i+1, j+1
		} else {
			mt.addIndex(indexIdNew, newIndex, pos, newPk)
			j++
		}
	}

	for index := old.Read2(i); index != nil; index = old.Read2(i) {
		indexId := binary.BigEndian.Uint16(index)
		if indexId == 0 {
			oldPk = index[2:]
		}
		mt.delIndex(indexId, index, oldPk)
		i++
	}

	for index := new.Read2(j); index != nil; index = new.Read2(j) {
		indexId := binary.BigEndian.Uint16(index)
		if indexId == 0 {
			newPk = index[2:]
		}
		mt.addIndex(indexId, index, pos, newPk)
		j++
	}
}

func (mt *MemTable) Set(pk []byte, value interface{}, data []byte) error {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	return mt.set(pk, value, data)
}

func (mt *MemTable) set(pk []byte, value interface{}, data []byte) error {
	var seq [8]byte // 8B seq

	indexOut := mt.Indexer.Index(value, data)

	if len(pk) == 0 {
		indexId, pk := indexOut.Read(0)
		if len(pk) == 0 || indexId != 0 {
			// 没有主键，使用sequence自动生成自增id
			s := mt.Ds.GetNextSequence(0)
			binary.BigEndian.PutUint64(seq[:], s)
			pk = seq[:]
			indexOut.Write("", pk)
		}
	} else {
		indexOut.Write("", pk)
	}

	pos := int32(0)
	n := uint16(0)

	// 找到要覆盖的元素的oldPos，如果需要覆盖一个以上的元素，报错返回
	for idIndex := indexOut.Read2(n); len(idIndex) != 0; idIndex = indexOut.Read2(n) {
		indexId := binary.BigEndian.Uint16(idIndex)
		pos2 := mt.getIndex(indexId, idIndex, pk)
		if pos == 0 {
			pos = pos2
		} else if pos2 != 0 {
			if pos != pos2 {
				return errors.New("Set failed due to pk/uk conflict")
			}
		}
		n++
	}

	if pos == 0 {
		pos = mt.MemBuffer.AddItem(value, data, indexOut)
		mt.update(&NULL_MI_INSTANCE, pos, &indexOut)
	} else {
		item := mt.MemBuffer.GetItem(pos)
		mt.update(&item.Keys, pos, &indexOut)
		item.Data = data
		item.Value = value
		item.Keys = indexOut
	}

	return nil
}

func (mt *MemTable) Delete(key []byte, index string) error {
	indexId := Crc16(index)

	k := make([]byte, 2, len(key)+2)
	binary.BigEndian.PutUint16(k, indexId)
	k = append(k, key...)

	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	// 非唯一索引无法删除
	pos := mt.getIndex(indexId, k, nil)
	if pos == 0 {
		return nil
	}

	mt.MemBuffer.DelItem(pos)

	return nil
}

func (mt *MemTable) BatchSet(iter func() (pk []byte, value interface{}, data []byte)) error {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	//FIXME: 批量设置不支持原子性，中间出错前面的无法回滚，需要调用方自行处理
	for pk, value, data := iter(); value != nil; pk, value, data = iter() {
		if err := mt.set(pk, value, data); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MemTable) BatchDelete(iter func() (key []byte, index string)) error {
	k := make([]byte, 2, 16)
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	for key, index := iter(); key != nil; key, index = iter() {
		indexId := Crc16(index)

		binary.BigEndian.PutUint16(k, indexId)
		k = append(k[:2], key...)

		pos := mt.getIndex(indexId, k, nil)
		if pos != 0 {
			mt.MemBuffer.DelItem(pos)
		}
	}

	return nil
}

func (mt *MemTable) Scan(start, end []byte, index string, f func(pk []byte, data []byte) error) error {
	indexId := Crc16(index)
	k := make([]byte, 2, len(start)+2)
	binary.BigEndian.PutUint16(k, indexId)
	k = append(k, start...)

	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	elem := mt.ti.Find(k)
	if elem == nil {
		return nil
	}

	for ; elem != nil; elem = elem.Next() {
		kk := elem.Key().([]byte)
		if bytes.Compare(kk[2:], end) >= 0 {
			break
		}
		item := mt.MemBuffer.GetItem(elem.Value.(int32))
		_, pk := item.Keys.Read(0)

		if err := f(pk, item.Data); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MemTable) ScanValue(start, end []byte, index string, f func(pk []byte, value interface{}) error) error {
	indexId := Crc16(index)
	k := make([]byte, 2, len(start)+2)
	binary.BigEndian.PutUint16(k, indexId)
	k = append(k, start...)

	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	elem := mt.ti.Find(k)
	if elem == nil {
		return nil
	}

	for ; elem != nil; elem = elem.Next() {
		kk := elem.Key().([]byte)
		if bytes.Compare(kk[2:], end) >= 0 {
			break
		}
		item := mt.MemBuffer.GetItem(elem.Value.(int32))
		_, pk := item.Keys.Read(0)

		if err := f(pk, item.Value); err != nil {
			return err
		}
	}

	return nil
}

func (mt *MemTable) Stat() string {
	mt.mtx.Lock()
	defer mt.mtx.Unlock()

	stat := make(map[string]interface{})
	stat["unique_index_number"] = len(mt.hi)
	stat["non_unique_index_number"] = mt.ti.Len()
	stat["mem_buffer_len"] = mt.MemBuffer.Len()

	jo, _ := json.Marshal(stat)

	return string(jo)
}

func NewMemKvDb(config map[string]interface{}) KvDb {
	jo, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	cfg := MemKvDbConfig{}
	if err := json.Unmarshal(jo, &cfg); err != nil {
		panic(err)
	}

	indexType := MEM_HASH_INDEX
	if strings.ToLower(cfg.IndexType) == "tree" {
		indexType = MEM_TREE_INDEX
	}

	m := MemKvDb{Cfg: cfg, IndexType: indexType}
	m.Open()
	return &m
}
