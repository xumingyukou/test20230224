package store

import (
	"archive/tar"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/pkg/errors"
	"github.com/warmplanet/proto/go/sdk"
	"github.com/warmplanet/proto/go/sdk/broker"
	"github.com/warmplanet/proto/go/sdk/skiplist"
)

/*
基于pebble实现的diskstore，配置参考conf/storage.toml

record key的基本格式：
主键(没有主键生自动成row id)：'t_' + table id[4字节] + pk/row id[8字节]  值：用户自定义字节数组
唯一索引：'i_' + table id[4字节] + index id[2字节] + 索引值 值：pk/row id
非唯一索引：'i_' + table id[4字节] + index id[2字节] + 索引值 + pk/row id 值：null
主键和二级索引都要满足字节序可比性，遍历的顺序应该和用户设置的顺序一致，比如整数索引，
应该负数在前，正数在后，字符串索引应当按照字典序。

存储在pebble中的key格式：
record key + timestamp(8字节)
每个record key都附加一个全局递增的时间戳，同一个事务中附加相同的时间戳，每个record key通过时间戳保存了多个版本
后台有定期垃圾回收，扫描所有的key并删除所有过期或者删除的key
目前仅实现获取最新数据的接口，历史版本获取暂不支持
TODO: it.Value()需要确认需不需要clone
*/
type DiskStoreConfig struct {
	DataDir            string       `toml:"data_dir" json:"data_dir"`
	BackupDir          string       `toml:"backup_dir" json:"backup_dir"`
	NoSync             bool         `toml:"no_sync" json:"no_sync"`
	BytesPerSync       sdk.ByteSize `toml:"bytes_per_sync" json:"bytes_per_sync"`
	ForceSyncWAL       bool         `toml:"force_sync_wal" json:"force_sync_wal"`
	DisableWAL         bool         `toml:"disable_wal" json:"disable_wal"`
	WALBytesPerSync    sdk.ByteSize `toml:"wal_bytes_per_sync" json:"wal_bytes_per_sync"`
	WALMinSyncInterval sdk.Duration `toml:"wal_min_sync_interval" json:"wal_min_sync_interval"`
}

// NOTE: xxx_PFX定义的常量字符串应该具有相同的长度
// 所有的表数据以此前缀开头
const TBL_PFX string = "t_"

// 所有的索引以此前缀为开头
const INDEX_PFX string = "i_"

// 所有的元数据以此前缀开头
const META_PFX string = "m_"

// 所有的自增序列以此前缀开头
const SEQ_PFX string = "s_"

// key最小和最大值，用于范围查询
var MIN_KEY1 = []byte{0}
var MAX_KEY1 = []byte{0xff}

var MIN_KEY2 = []byte{0, 0}
var MAX_KEY2 = []byte{0xff, 0xff}

var MIN_KEY4 = []byte{0, 0, 0, 0}
var MAX_KEY4 = []byte{0xff, 0xff, 0xff, 0xff}

var MIN_KEY8 = []byte{0, 0, 0, 0, 0, 0, 0, 0}
var MAX_KEY8 = []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

const SCAN_BATCH_SIZE int = 100

type DiskKvDb struct {
	// 创建对象时设置好，之后调用Open
	Cfg DiskStoreConfig
	// 设置成公开成员，方便进行调试代码开发
	Db   *pebble.DB
	Name string

	// 时间戳，确保单调递增
	ts int64
	// 用于生成table id
	ds KvSequencer

	tblCache sdk.ConcurrentMapI // *Table
	seqCache sdk.ConcurrentMapI // *Sequencer
}

// 表元信息
type DiskTable struct {
	KvDb    *DiskKvDb        `json:"-"`
	Db      *pebble.DB       `json:"-"`
	Cfg     *DiskStoreConfig `json:"-"`
	Indexer KvIndexer        `json:"-"`
	Ds      KvSequencer      `json:"-"`

	TableId uint32 `json:"tid"`   // 表唯一id
	Name    string `json:"name"`  // 表名
	Class   string `json:"class"` // 数据类型，用于辅助数据编解码
	//Indexs  []DiskIndex `json:"indexs"` // 所有的index信息
	Indexs map[uint16]DiskIndex `json:"indexs"` // 所有的index信息，key为indexId
}

// 索引元信息
type DiskIndex struct {
	Name    string `json:"name"`
	IndexId uint16 `json:"iid"`    // 索引唯一id
	Unique  bool   `json:"unique"` // 是否是唯一索引
}

type DiskSequencer struct {
	Db   *pebble.DB
	Cfg  *DiskStoreConfig
	Name string // seq名字，和map中的key一致

	upperSeq  uint64 // 上限sequence，达到后需要去db里重新批量分配
	current   uint64 // 当前序列号
	cacheSize uint64 // load sequence的间隔

	loaded uint64     // 最新加载的db里的sequence值
	load   sync.Mutex // 从db加载时加锁，加载完毕解锁
	update sync.Mutex // 更新current时加锁
}

/*
	系统元数据key结构

m_ + "t_" + 表名 -> 表元数据
m_ + "i_" + 表ID[4字节] + 索引名 -> 索引元数据
*/
const GLB_TBL_SEQ string = TBL_PFX

func (dt *DiskTable) Get(key []byte, index string) (pk []byte, data []byte, err error) {
	// 不支持根据非唯一索引Get，如果要根据非唯一索引获取数据可以使用Scan接口
	// 最大k长度为2字节前缀+4字节tableId+2字节indexId+8字节rowId+key+timestamp
	k, unique := dt.genKey(key, index)

	if !unique || len(k) == 0 {
		return nil, nil, errors.Errorf("Get item from non-unique/invalid index")
	}

	it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: k[:len(k)-8], UpperBound: k})
	defer it.Close()

	if it.Last() {
		// 主键直接返回数据
		if index == "" {
			return key, clone(it.Value()), nil
		}

		// 索引还需要再次获取，唯一索引的值是主键+时间戳
		pk := make([]byte, len(TBL_PFX)+4+len(it.Value()))
		copy(pk, TBL_PFX)
		binary.BigEndian.PutUint32(pk[len(TBL_PFX):], dt.TableId)
		copy(pk[len(TBL_PFX)+4:], it.Value())

		r := dt.getByPk(pk, true)
		return pk[len(TBL_PFX)+4 : len(pk)-8], r, nil
	}

	return nil, nil, nil
}

func (dt *DiskTable) GetValue(key []byte, index string) (pk []byte, value interface{}, err error) {
	return nil, nil, errors.Errorf("Unsupporteed GetValue call")
}

func (dt *DiskTable) Set(pk []byte, value interface{}, data []byte) error {
	batch := dt.Db.NewBatch()
	var seqTs [16]byte // 8B seq + 8B ts
	_ts := dt.KvDb.Timestamp()
	binary.BigEndian.PutUint64(seqTs[8:], uint64(_ts))

	indexOut := dt.Indexer.Index(value, data)
	pka := make([]byte, 0, len(TBL_PFX)+4+8+8)

	if len(pk) == 0 {
		indexId, pk2 := indexOut.Read(0)
		if len(pk2) != 0 && indexId == 0 {
			pk = pk2
		}
	}

	if len(pk) == 0 {
		// 没有主键，使用sequence自动生成自增id
		seq := dt.Ds.GetNextSequence(0)
		binary.BigEndian.PutUint64(seqTs[:], seq)
		pk = seqTs[:8]
	}

	pka = append(pka, TBL_PFX...)
	pka = append(pka, MIN_KEY4...)
	binary.BigEndian.PutUint32(pka[len(TBL_PFX):], dt.TableId)
	pka = append(pka, pk...)
	pka = append(pka, seqTs[8:]...)

	// 设置主键
	batch.Set(pka, data, nil)

	idxKey := make([]byte, len(INDEX_PFX)+6, len(INDEX_PFX)+6+32+8) // 32为索引+pk长度估计值，不够会自动扩展
	copy(idxKey, INDEX_PFX)
	binary.BigEndian.PutUint32(idxKey[len(INDEX_PFX):], dt.TableId)

	n := uint16(0)
	for indexId, v := indexOut.Read(n); len(v) != 0; indexId, v = indexOut.Read(n) {
		if indexId == 0 {
			n++
			continue
		}
		idxKey = idxKey[:len(INDEX_PFX)+6]
		binary.BigEndian.PutUint16(idxKey[len(INDEX_PFX)+4:], indexId)

		meta := dt.Indexs[indexId]
		if meta.IndexId == 0 {
			// 忽略主键或者无法找到的索引
			n++
			continue
		}

		if meta.Unique {
			idxKey = append(idxKey, v...)
			idxKey = append(idxKey, seqTs[8:]...)
			batch.Set(idxKey, pka[len(TBL_PFX)+4:], nil)
		} else {
			idxKey = append(idxKey, v...)
			idxKey = append(idxKey, pka[len(TBL_PFX)+4:]...)
			batch.Set(idxKey, nil, nil)
		}
		n++
	}

	err := batch.Commit(dt.Cfg.WriteOption())
	if err == nil && dt.Cfg.ForceSyncWAL {
		dt.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return err
	}

	return nil
}

func (dt *DiskTable) Delete(key []byte, index string) error {
	// 通过index找到pk，之后使用pk写入nil数据表示删除，后台垃圾回收异步删除数据
	k, unique := dt.genKey(key, index)

	if !unique || len(k) == 0 {
		return errors.Errorf("Delete on non-unique/invalid index")
	}

	var ts [8]byte
	var pk []byte

	if index == "" {
		_ts := dt.KvDb.Timestamp()
		binary.BigEndian.PutUint64(ts[:], uint64(_ts))
		pk = k
		copy(pk[len(pk)-8:], ts[:])
	} else {
		// 根据索引查询主键
		it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: k[:len(k)-8], UpperBound: k})
		defer it.Close()

		if it.Last() {
			_ts := dt.KvDb.Timestamp()
			binary.BigEndian.PutUint64(ts[:], uint64(_ts))

			// 索引还需要再次获取，唯一索引的值是主键+时间戳
			pk = make([]byte, len(TBL_PFX)+4+len(it.Value()))
			copy(pk, TBL_PFX)
			binary.BigEndian.PutUint32(pk[len(TBL_PFX):], dt.TableId)
			copy(pk[len(TBL_PFX)+4:], it.Value())
			copy(pk[len(pk)-8:], ts[:])
		} else {
			return nil
		}
	}

	// 使用最新的时间戳写入空值，表示删除
	err := dt.Db.Set(pk, nil, dt.Cfg.WriteOption())
	if err == nil && dt.Cfg.ForceSyncWAL {
		dt.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return err
	}

	return nil
}

func (dt *DiskTable) BatchSet(iter func() (pk []byte, value interface{}, data []byte)) error {
	idxKey := make([]byte, len(INDEX_PFX)+6, len(INDEX_PFX)+6+32+8) // 32为索引+pk长度估计值，不够会自动扩展
	var seqTs [16]byte
	batch := dt.Db.NewBatch()

	pka := make([]byte, 0, len(TBL_PFX)+4+8+8)
	for pk, v, rv := iter(); v != nil; pk, v, rv = iter() {
		indexOut := dt.Indexer.Index(v, rv)

		if len(pk) == 0 {
			// 第一个为主键
			indexId, pk2 := indexOut.Read(0)
			if len(pk2) != 0 && indexId == 0 {
				pk = pk2
			}
		}
		if len(pk) == 0 {
			// 没有主键，使用sequence自动生成自增id
			seq := dt.Ds.GetNextSequence(0)
			binary.BigEndian.PutUint64(seqTs[:], seq)
			pk = seqTs[:8]
		}

		_ts := dt.KvDb.Timestamp()
		binary.BigEndian.PutUint64(seqTs[8:], uint64(_ts))
		pka := pka[:0]

		pka = append(pka, TBL_PFX...)
		pka = append(pka, MIN_KEY4...)
		binary.BigEndian.PutUint32(pka[len(TBL_PFX):], dt.TableId)
		pka = append(pka, pk...)
		pka = append(pka, seqTs[8:]...)

		// 添加主键
		batch.Set(pka, rv, nil)
		// 添加索引
		copy(idxKey, INDEX_PFX)
		binary.BigEndian.PutUint32(idxKey[len(INDEX_PFX):], dt.TableId)

		n := uint16(0)
		for indexId, v := indexOut.Read(n); len(v) != 0; indexId, v = indexOut.Read(n) {
			if indexId == 0 {
				n++
				continue
			}
			idxKey = idxKey[:len(INDEX_PFX)+6]
			binary.BigEndian.PutUint16(idxKey[len(INDEX_PFX)+4:], indexId)

			meta := dt.Indexs[indexId]
			if meta.IndexId == 0 {
				// 忽略主键和无法找到的索引
				n++
				continue
			}

			// 唯一索引
			if meta.Unique {
				idxKey = append(idxKey, v...)
				idxKey = append(idxKey, seqTs[8:]...)
				batch.Set(idxKey, pka[len(TBL_PFX)+4:], nil)
			} else {
				idxKey = append(idxKey, v...)
				idxKey = append(idxKey, pka[len(TBL_PFX)+4:]...)
				batch.Set(idxKey, nil, nil)
			}
			n++
		}
	}

	err := batch.Commit(dt.Cfg.WriteOption())
	if err == nil && dt.Cfg.ForceSyncWAL {
		dt.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return err
	}

	return nil
}

// NOTE: 只能返回满足本索引和主键上唯一性约束的数据，不能返回满足全部索引上唯一性约束的数据
// 要实现需要对每一个value计算所有的唯一索引，然后只保留每个索引值上的最新值，最后将各个索引过滤后的求交集才能得到
// 此功能可以在f函数中实现，不在引擎接口支持
func (dt *DiskTable) Scan(start, end []byte, index string, f func(pk []byte, data []byte) error) error {
	k1, _ := dt.genKey(start, index)
	k2, _ := dt.genKey(end, index)

	if len(k1) == 0 || len(k2) == 0 {
		return errors.Errorf("Scan item from invalid index")
	}
	// 确保起始key的timestamp为0
	copy(k1[len(k1)-8:], MIN_KEY8)

	// 主键，只需要简单遍历即可，每个pk选择最后的版本
	if index == "" {
		var lastK []byte
		var lastV []byte
		it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: k1, UpperBound: k2})
		defer it.Close()

		for it.First(); it.Valid(); it.Next() {
			k := it.Key()
			if len(lastK) == 0 {
				lastK = append(lastK, k[:len(k)-8]...)
			}

			// 只使用key的最新值调用f
			if !bytes.Equal(k[:len(k)-8], lastK) && len(lastV) > 0 {
				if f(lastK[len(TBL_PFX)+4:], lastV) != nil {
					return nil
				}
			}
			// 考虑到大部分更改频繁的pk值都较小，每次都拷贝键值对，比it.Prev()更快速
			lastK = append(lastK[:0], k[:len(k)-8]...)
			lastV = append(lastV[:0], it.Value()...)
		}

		if len(lastV) > 0 {
			f(lastK[len(TBL_PFX)+4:], lastV)
		}
	} else {
		it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: k1, UpperBound: k2})
		// 缓存所有的primary key
		kvs := make([]scanItem, 0)
		sl := skiplist.New(skiplist.Bytes) // 用于pk去重
		var idx int
		var lastKey []byte
		var k []byte
		// 首先保存所有对应的pk
		for it.First(); it.Valid(); _, idx = it.Next(), idx+1 {
			ik := it.Key()[len(INDEX_PFX)+4+2:]
			v := it.Value()
			if len(v) == 0 {
				// 非唯一索引，需要调用unpackIndex从key中提取pk
				k = ik[dt.Indexer.UnpackIndex(index, ik, nil):]
				kvs = append(kvs, scanItem{idx: idx, pk: clone(k)})
			} else {
				// 唯一索引，值为主键+时间戳，key只保留最新值
				if bytes.HasPrefix(lastKey, ik[:len(ik)-8]) {
					kvs[idx-1].pk = nil
				}
				kvs = append(kvs, scanItem{idx: idx, pk: clone(v)})
				k = v
			}

			e := sl.Get(k[:len(k)-8])
			if e != nil {
				kvs[e.Value.(int)].pk = nil
				e.Value = idx
			} else {
				sl.Set(k[:len(k)-8], idx)
			}
			lastKey = append(lastKey[:0], ik...)
		}
		it.Close()

		// 每一批暂存kv，缓存完所有value后，按照索引序调用f
		cache := make([]scanItem, SCAN_BATCH_SIZE)

		total := 0
		pk, _ := dt.genKey([]byte{}, "")

		for len(kvs) > 0 {
			n := copy(cache, kvs)
			t := kvs[:n]
			// 对pk进行排序并合并
			sort.Slice(t, func(i int, j int) bool {
				return bytes.Compare(t[i].pk, t[j].pk) < 0
			})

			var i int
			// 忽略前面pk为空的元素
			for i = 0; i < len(t); i++ {
				if len(t[i].pk) > 0 {
					break
				}
			}
			if i == len(t) {
				kvs = kvs[n:]
				total += n
				continue
			}

			t = t[i:]
			pk = append(pk[:len(TBL_PFX)+4], t[0].pk...)
			// 每一个SCAN_BATCH_SIZE创建一个新的iter，避免使用一个iter长时间遍历增加内存和磁盘
			it = dt.Db.NewIter(&pebble.IterOptions{LowerBound: pk})

			var lastK []byte
			var lastV []byte
			for i, _ = 0, it.SeekGE(pk); i < len(t) && it.Valid(); {
				k := it.Key()
				k = k[len(TBL_PFX)+4:]
				tk := t[i].pk

				c := bytes.Compare(tk[:len(tk)-8], k[:len(k)-8])
				if c < 0 {
					if len(lastV) > 0 {
						// 反向定位到index顺序的cache，保存value
						cache[t[i].idx-total].value = clone(lastV)
					}
					i++
					// 更新kv到batch中
					lastK, lastV = lastK[:0], lastV[:0]
					continue
				} else if c > 0 {
					pk = append(pk[:len(TBL_PFX)+4], t[i].pk...)
					it.SeekGE(pk)
					lastK, lastV = lastK[:0], lastV[:0]
					continue
				}

				if bytes.Equal(tk[len(tk)-8:], k[len(k)-8:]) {
					lastK = append(lastK[:0], k[:len(k)-8]...)
					lastV = append(lastV[:0], it.Value()...)
				} else {
					lastK, lastV = lastK[:0], lastV[:0]
				}

				it.Next()
			}

			if len(lastV) > 0 {
				cache[t[i].idx-total].value = clone(lastV)
			}

			it.Close()

			for i := 0; i < n; i++ {
				if len(cache[i].value) == 0 {
					continue
				}

				if f(cache[i].pk, cache[i].value) != nil {
					return nil
				}
			}

			kvs = kvs[n:]
			total += n
		}
	}

	return nil
}

func (dt *DiskTable) BatchDelete(iter func() (key []byte, index string)) error {
	var ts [8]byte
	keys := make([][]byte, 0)
	pks := make([][]byte, 0)

	// 先收集所有genKeys生成的key，之后进行排序，再通过一个iter进行无回溯的seek
	for k, idx := iter(); len(k) != 0; k, idx = iter() {
		k, unique := dt.genKey(k, idx)
		if !unique || len(k) == 0 {
			return errors.Errorf("Delete on non-unique/invalid index")
		}

		if idx == "" {
			_ts := dt.KvDb.Timestamp()
			binary.BigEndian.PutUint64(ts[:], uint64(_ts))
			copy(k[len(k)-8:], ts[:])
			pks = append(pks, k)
		} else {
			// 添加不包含时间戳的key，方便后续迭代器SeekGE调用
			keys = append(keys, k[:len(k)-8])
		}
	}

	sort.Slice(keys, func(i int, j int) bool {
		return bytes.Compare(keys[i], keys[j]) < 0
	})

	if len(keys) > 0 {
		it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: keys[0], UpperBound: append(keys[len(keys)-1], MAX_KEY8...)})
		for i, _ := 0, it.SeekGE(keys[0]); i < len(keys) && it.Valid(); {
			k := it.Key()
			c := bytes.Compare(keys[i], k[:len(k)-8])
			if c < 0 {
				i++
				continue
			} else if c > 0 {
				it.SeekGE(keys[i])
				continue
			}

			_ts := dt.KvDb.Timestamp()
			binary.BigEndian.PutUint64(ts[:], uint64(_ts))
			// 获取唯一索引的值：主键+时间戳来构造key
			pk := make([]byte, len(TBL_PFX)+4+len(it.Value()))
			copy(pk, TBL_PFX)
			binary.BigEndian.PutUint32(pk[len(TBL_PFX):], dt.TableId)
			copy(pk[len(TBL_PFX)+4:], it.Value())
			copy(pk[len(pk)-8:], ts[:])
			pks = append(pks, pk)
			i++
		}

		it.Close()
	}

	batch := dt.Db.NewBatch()

	for _, pk := range pks {
		batch.Set(pk, nil, nil)
	}

	err := batch.Commit(dt.Cfg.WriteOption())
	if err == nil && dt.Cfg.ForceSyncWAL {
		dt.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return err
	}

	return nil
}

func (dt *DiskTable) ScanValue(start, end []byte, index string, f func(pk []byte, value interface{}) error) error {
	return errors.New("Unsupported scan operation")
}

// 垃圾回收
func (dt *DiskTable) Gc() {
	//TODO: 删除所有无用的历史数据和索引
	// 先并发删除所有pk和unique-index上的过期数据
	// 然后删除nonunique-index上的过期数据
}

func (dt *DiskTable) CreateIndex(name string, unique bool) (*DiskIndex, error) {
	iid := make([]byte, len(META_PFX)+len("i_")+4+len(name))
	copy(iid, META_PFX)
	copy(iid[len(META_PFX):], "i_")
	binary.BigEndian.PutUint32(iid[len(META_PFX)+len("i_"):], dt.TableId)
	copy(iid[len(iid)-len(name):], name)

	meta := DiskIndex{Name: name}

	// 从db加载
	v, closer, _ := dt.Db.Get(iid)

	if len(v) != 0 {
		if err := json.Unmarshal(v, &meta); err != nil {
			panic("Create index failed with error: " + err.Error())
		}

		return &meta, nil
	}

	if closer != nil {
		closer.Close()
	}

	// 创建index
	idxId := Crc16(name)

	if old, ok := dt.Indexs[idxId]; ok {
		return nil, errors.Errorf("CreateIndex failed due to name conflict: %s %s", old.Name, name)
	}

	meta.IndexId = uint16(idxId)
	meta.Unique = unique

	jo, err := json.Marshal(&meta)
	if err != nil {
		panic("Create index failed with error: " + err.Error())
	}

	// 使用merge创建，防止并发掉用后创建的覆盖之前创建的
	err = dt.Db.Merge(iid, jo, dt.Cfg.WriteOption())

	if err == nil && dt.Cfg.ForceSyncWAL {
		dt.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return nil, err
	}

	// merge完成后，尝试get，并发情况下get的可能不是我们设置的而是之前设置的
	v, closer, err = dt.Db.Get(iid)

	if closer != nil {
		defer closer.Close()
	}

	if len(v) != 0 {
		if err := json.Unmarshal(v, &meta); err != nil {
			panic("Create index failed with error: " + err.Error())
		}
		return &meta, nil
	} else {
		return nil, err
	}
}

func (dt *DiskTable) Stat() string {
	stat := make(map[string]interface{})
	//TODO: 增加表的元素数量的统计信息

	jo, _ := json.Marshal(stat)
	return string(jo)
}

// 根据index和索引生成查询的key
func (dt *DiskTable) genKey(key []byte, index string) (end []byte, unique bool) {
	// 按照非唯一索引分配最大可能空间（假设所有pk字节排序小于MAX_KEY8）
	k := make([]byte, len(INDEX_PFX)+14+len(key)+8)

	if index == "" {
		// 主键
		copy(k, TBL_PFX)
		binary.BigEndian.PutUint32(k[len(TBL_PFX):], dt.TableId)
		copy(k[len(TBL_PFX)+4:], key)

		// 查询时间戳范围0～max64的所有版本
		copy(k[len(TBL_PFX)+4+len(key):], MAX_KEY8)

		k = k[:len(TBL_PFX)+4+len(key)+8]
		return k, true
	} else {
		indexId := Crc16(index)
		meta, ok := dt.Indexs[indexId]

		if !ok || meta.Name != index {
			return nil, false
		}

		copy(k, INDEX_PFX)
		binary.BigEndian.PutUint32(k[len(INDEX_PFX):], dt.TableId)
		binary.BigEndian.PutUint16(k[len(INDEX_PFX)+4:], indexId)
		copy(k[len(INDEX_PFX)+6:], key)

		//非唯一索引，查询pk范围0～max64的所有数据
		//唯一索引，查询时间戳范围0~max64的所有版本
		copy(k[len(INDEX_PFX)+6+len(key):], MAX_KEY8)
		k = k[:len(INDEX_PFX)+6+len(key)+8]
		return k, meta.Unique
	}
}

// 根据pk获取值，如果latest为true，尝试获取pk的最新版本，如果版本和pk里的不一致，返回nil，一致返回
// 如果latest为false，尝试直接获取pk对应的历史版本，存在则返回
func (dt *DiskTable) getByPk(pk []byte, latest bool) []byte {
	if !latest {
		r, closer, _ := dt.Db.Get(pk)
		r = clone(r)

		if closer != nil {
			closer.Close()
		}

		return r
	} else {
		pk2 := clone(pk)
		copy(pk2[len(pk2)-8:], MAX_KEY8)
		it := dt.Db.NewIter(&pebble.IterOptions{LowerBound: pk, UpperBound: pk2})
		defer it.Close()

		if it.Last() {
			k := it.Key()
			// pk就是latest，返回值
			if bytes.Equal(k, pk) {
				return clone(it.Value())
			}

			// pk的值不是latest版本，返回nil
			return nil
		} else {
			return nil
		}
	}
}

// 调用前需要加ctx.load锁，加载完毕后释放
func (ds *DiskSequencer) loadSequence() {
	var v [8]byte
	var err error

	ctx := ds
	loaded := ctx.loaded
	// 上次缓存的最新从db加载的值
	if loaded > ctx.current {
		ctx.load.Unlock()

		ctx.update.Lock()
		loaded = ctx.loaded
		if loaded >= ctx.upperSeq && loaded > 0 {
			ctx.upperSeq = loaded
		} else {
			// 并发调用了load sequence，其他线程更快完成了更新
			broker.BrokerLogger.Infof("Merge conflict when load sequence %d %d", loaded, ctx.upperSeq)
		}
		ctx.update.Unlock()
		return
	}

	binary.BigEndian.PutUint64(v[:], ctx.cacheSize)
	k := []byte(SEQ_PFX + ctx.Name)

	current := ctx.current
	// 从db加载上次的current
	if current == 0 {
		seq, closer, _ := ds.Db.Get(k)
		if len(seq) != 0 {
			current = binary.BigEndian.Uint64(seq)
		}
		if closer != nil {
			closer.Close()
		}

		ctx.update.Lock()
		if ctx.current == 0 {
			ctx.current = current
		}
		ctx.update.Unlock()
	}

	err = ds.Db.Merge(k, v[:], ds.Cfg.WriteOption())

	if err == nil && ds.Cfg.ForceSyncWAL {
		ds.Db.LogData(nil, pebble.Sync)
	}

	if err == nil {
		seqInit := uint64(0)
		seq, closer, _ := ds.Db.Get(k)
		if len(seq) != 0 {
			seqInit = binary.BigEndian.Uint64(seq)
			ctx.loaded = seqInit
		}

		if closer != nil {
			closer.Close()
		}
		ctx.load.Unlock()

		ctx.update.Lock()
		if seqInit >= ctx.upperSeq && ctx.loaded > 0 {
			ctx.upperSeq = seqInit
		} else {
			// 并发调用了load sequence，其他线程更快完成了更新
			broker.BrokerLogger.Infof("Merge conflict when load sequence %d %d", ctx.upperSeq, ctx.upperSeq)
		}
		ctx.update.Unlock()
		return
	} else {
		ctx.load.Unlock()
		broker.BrokerLogger.Errorf("meta store load sequene error: %v", err)
		return
	}
}

// 返回0表示获取失败
func (ds *DiskSequencer) GetNextSequence(cacheSize uint64) uint64 {
	if cacheSize == 0 {
		cacheSize = ds.cacheSize
	}

	cur := ds
	if cur.current >= cur.upperSeq {
		// sync load
		cur.load.Lock()
		ds.loadSequence()
	}

	current := uint64(0)
	cur.update.Lock()
	defer cur.update.Unlock()

	cur.current = cur.current + 1
	current = cur.current
	if cur.cacheSize != cacheSize {
		cur.cacheSize = cacheSize
	}

	// 小于cacheSize的1/10开始预加载
	if (cur.upperSeq-cur.current) < cur.cacheSize/10 && cur.load.TryLock() {
		// async load
		go ds.loadSequence()
	}

	if current == 0 {
		panic("Get next sequence failed")
	}

	return current
}

func (ds *DiskSequencer) ResetSequence(initSeq uint64, cacheSize uint64) (err error) {
	var v [8]byte
	binary.BigEndian.PutUint64(v[:], initSeq)
	k := []byte(SEQ_PFX + ds.Name)

	if cacheSize == 0 {
		cacheSize = ds.cacheSize
	}

	// 加双重锁，注意顺序需要先load后update和其他地方一致
	ds.load.Lock()
	ds.update.Lock()

	defer ds.load.Unlock()
	defer ds.update.Unlock()

	err = ds.Db.Set(k, v[:], ds.Cfg.WriteOption())
	if err == nil && ds.Cfg.ForceSyncWAL {
		ds.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return
	}

	ds.upperSeq = initSeq
	ds.current = initSeq
	ds.loaded = 0
	ds.cacheSize = cacheSize

	return
}

func (dkd *DiskKvDb) Open() error {
	dkd.tblCache = sdk.NewCmapI()
	dkd.seqCache = sdk.NewCmapI()

	db, err := OpenDiskStore(&dkd.Cfg)
	if err != nil {
		log.Fatal(err)
	}

	dkd.Db = db
	dkd.ds, _ = dkd.CreateSequencer(GLB_TBL_SEQ, 1, nil)

	return nil
}

func (dkd *DiskKvDb) Close() error {
	return dkd.Db.Close()
}

func (dkd *DiskKvDb) CreateTable(name string, class string, indexer KvIndexer, args map[string]interface{}) (table KvTable, err error) {
	// 从缓存加载
	if t, ok := dkd.tblCache.Get(name); ok {
		return t.(KvTable), nil
	}

	ds, _ := dkd.CreateSequencer(TBL_PFX+name, 100, nil)
	tid := append([]byte(META_PFX+"t_"), []byte(name)...)
	meta := DiskTable{Db: dkd.Db, Cfg: &dkd.Cfg, Indexer: indexer, Indexs: make(map[uint16]DiskIndex), Ds: ds, KvDb: dkd}
	if meta.Indexer == nil {
		meta.Indexer = &noneIndexer{}
	}

	// 从db加载
	v, closer, _ := dkd.Db.Get(tid)
	v = clone(v)

	if closer != nil {
		closer.Close()
	}

	if len(v) != 0 {
		if err := json.Unmarshal(v, &meta); err != nil {
			panic("Create table failed with error: " + err.Error())
		}
		// 更新缓存
		dkd.tblCache.Set(name, &meta)
		return &meta, nil
	}

	// 创建table
	tblId := dkd.ds.GetNextSequence(1)

	// 目前不支持累计创建超过uint32Max数量的表
	if tblId >= math.MaxUint32 {
		panic("table id exceed uint32 limit")
	}

	meta.Class = class
	meta.Name = name
	meta.TableId = uint32(tblId)

	if indexer != nil {
		// 主键indexId固定为0
		meta.Indexs[0] = DiskIndex{Unique: true}
		for k, v := range indexer.IndexNames() {
			if k == "" {
				// pk
				continue
			}
			dt, err := meta.CreateIndex(k, v)
			if err != nil {
				panic("Create index failed: " + err.Error())
			}
			meta.Indexs[dt.IndexId] = *dt
		}
	}

	jo, err := json.Marshal(&meta)
	if err != nil {
		return nil, err
	}

	// 使用merge创建，防止并发掉用后创建的覆盖之前创建的
	err = dkd.Db.Merge(tid, jo, dkd.Cfg.WriteOption())

	if err == nil && dkd.Cfg.ForceSyncWAL {
		dkd.Db.LogData(nil, pebble.Sync)
	}

	if err != nil {
		return nil, err
	}

	// merge完成后，尝试get，并发情况下get的可能不是我们设置的而是之前设置的
	v, closer, err = dkd.Db.Get(tid)

	if closer != nil {
		defer closer.Close()
	}

	if len(v) != 0 {
		if err := json.Unmarshal(v, &meta); err != nil {
			panic("Create table failed with error: " + err.Error())
		}
		// 更新缓存
		dkd.tblCache.Set(name, &meta)
		return &meta, nil
	} else {
		panic("Write table failed with error: " + err.Error())
	}
}

func (dkd *DiskKvDb) CreateSequencer(name string, initCacheSize uint64, args map[string]interface{}) (sequencer KvSequencer, err error) {
	// 从缓存加载
	if t, ok := dkd.seqCache.Get(name); ok {
		return t.(*DiskSequencer), nil
	}

	if initCacheSize == 0 {
		initCacheSize = 1
	}

	ds := &DiskSequencer{Db: dkd.Db, Name: name, Cfg: &dkd.Cfg, cacheSize: initCacheSize}

	if dkd.seqCache.SetIfAbsent(name, ds) {
		return ds, nil
	} else {
		t, _ := dkd.seqCache.Get(name)
		return t.(*DiskSequencer), nil
	}
}

func (dkd *DiskKvDb) Backup(dstUrl string) error {
	var u *url.URL
	u, err := url.ParseRequestURI(dstUrl)

	if err != nil {
		return err
	}

	if u.Scheme == "file" {
		path := u.Path
		// file:./xxx/, file:xxx.db, file:/usr/local/xxx.db
		if path == "" {
			path = strings.TrimPrefix(dstUrl, "file:")
		}
		// 安全起见，确保path地址以配置的backupPath开头
		path, err2 := filepath.Abs(dkd.Cfg.BackupDir + "/" + path)
		if err2 != nil {
			return err2
		}

		bp, _ := filepath.Abs(dkd.Cfg.BackupDir)

		if bp == "" || !strings.HasPrefix(path, bp) {
			return errors.New("invalid backup path")
		}

		return dkd.Db.Checkpoint(path)
	} else if u.Scheme == "http" {
		// 创建临时目录
		dst := fmt.Sprintf("%s/%d", dkd.Cfg.BackupDir, time.Now().Unix())
		if err = dkd.Db.Checkpoint(dst); err != nil {
			return err
		}

		// 将所有文件打包
		if err = Tar(dst, dkd.Cfg.BackupDir); err != nil {
			return err
		}

		pr, pw := io.Pipe()
		request, err2 := http.NewRequest("POST", dstUrl, pr)
		if err2 != nil {
			return err2
		}

		finfo, err := os.Stat(dst + ".tar")
		if err != nil {
			return err
		}

		request.Header.Add("Content-Type", "application/octet-stream")
		request.Header.Add("Content-Length", strconv.FormatInt(finfo.Size(), 10))
		client := &http.Client{}
		go func() {
			f, _ := os.Open(dst + ".tar")
			_, e := io.Copy(pw, f)
			if e != nil {
				broker.BrokerLogger.Error(e)
			}
			pw.Close()
			//删除临时文件和目录
			os.Remove(dst + ".tar")
			os.RemoveAll(dst)
		}()

		response, err := client.Do(request)
		if err != nil {
			return err
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			r, _ := io.ReadAll(response.Body)
			return fmt.Errorf("response error %d: %s", response.StatusCode, string(r))
		}

		return nil
	}

	return errors.New("invalid url scheme")
}

func (dkd *DiskKvDb) Stat() string {
	stat := make(map[string]interface{})
	metrics := dkd.Db.Metrics()
	bs := sdk.NewByteSize(metrics.DiskSpaceUsage())

	stat["disk_usage"] = bs.String()
	stat["wal.size"] = sdk.NewByteSize(metrics.WAL.Size).String()
	stat["wal.physical_size"] = sdk.NewByteSize(metrics.WAL.PhysicalSize).String()
	stat["wal.files"] = metrics.WAL.Files
	stat["blockcache.size"] = sdk.NewByteSize(uint64(metrics.BlockCache.Size))
	stat["blockcache.count"] = metrics.BlockCache.Count

	jo, _ := json.Marshal(stat)
	return string(jo)
}

func (dkd *DiskKvDb) Timestamp() int64 {
	ts := time.Now().UnixNano()

retry:
	cur := atomic.LoadInt64(&dkd.ts)

	if cur >= ts {
		ts = cur + 1
	}

	if !atomic.CompareAndSwapInt64(&dkd.ts, cur, ts) {
		goto retry
	}

	return ts
}

func NewDiskKvDb(config map[string]interface{}) KvDb {
	jo, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	cfg := DiskStoreConfig{}
	if err := json.Unmarshal(jo, &cfg); err != nil {
		panic(err)
	}

	m := &DiskKvDb{tblCache: sdk.NewCmapI(), seqCache: sdk.NewCmapI(), Cfg: cfg}

	err = m.Open()
	if err != nil {
		log.Fatal(err)
	}

	return m
}

// 保留最先写入值的merger
type RetainFirstMerger struct {
	buf []byte
}

func (a *RetainFirstMerger) MergeNewer(value []byte) error {
	if a.buf == nil {
		buf := make([]byte, len(value))
		copy(buf, value)
		a.buf = buf
	}
	return nil
}

func (a *RetainFirstMerger) MergeOlder(value []byte) error {
	buf := make([]byte, len(value))
	copy(buf, value)
	a.buf = buf
	return nil
}

func (a *RetainFirstMerger) Finish(includesBase bool) ([]byte, io.Closer, error) {
	return a.buf, nil, nil
}

// 用于自增id的merger，将所有的value进行累加
type IncUint64Merger struct {
	buf [8]byte
}

func (a *IncUint64Merger) MergeNewer(value []byte) error {
	op := binary.BigEndian.Uint64(value)
	orig := binary.BigEndian.Uint64(a.buf[:])
	binary.BigEndian.PutUint64(a.buf[:], orig+op)
	return nil
}

func (a *IncUint64Merger) MergeOlder(value []byte) error {
	op := binary.BigEndian.Uint64(value)
	orig := binary.BigEndian.Uint64(a.buf[:])
	binary.BigEndian.PutUint64(a.buf[:], orig+op)
	return nil
}

func (a *IncUint64Merger) Finish(includesBase bool) ([]byte, io.Closer, error) {
	return a.buf[:], nil, nil
}

var merger = &pebble.Merger{
	Merge: func(key, value []byte) (pebble.ValueMerger, error) {
		// 对于s_类型，实现自增整数语义
		if bytes.HasPrefix(key, []byte(SEQ_PFX)) {
			res := &IncUint64Merger{}
			copy(res.buf[:], value)
			return res, nil
		} else {
			// 对于t_,i_,m_类型，通过merge实现putIfAbsent语义
			res := &RetainFirstMerger{}
			res.buf = append(res.buf, value...)
			return res, nil
		}
	},

	// 此名字会保存到db中，不能修改
	Name: "broker.merger",
}

func OpenDiskStore(cfg *DiskStoreConfig) (*pebble.DB, error) {
	option := &pebble.Options{
		DisableWAL:      cfg.DisableWAL,
		BytesPerSync:    int(cfg.BytesPerSync),
		WALBytesPerSync: int(cfg.WALBytesPerSync),
		Merger:          merger,
	}

	if cfg.WALMinSyncInterval.Duration > 0 {
		option.WALMinSyncInterval = func() time.Duration {
			return cfg.WALMinSyncInterval.Duration
		}
	}

	db, err := pebble.Open(cfg.DataDir, option)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Tar(source, target string) error {
	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.tar", filename))
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})
}

func (c *DiskStoreConfig) WriteOption() *pebble.WriteOptions {
	if !c.NoSync {
		return pebble.Sync
	} else {
		return pebble.NoSync
	}
}

func clone(s []byte) []byte {
	if s == nil {
		return nil
	}
	tmp := make([]byte, len(s))
	copy(tmp, s)
	return tmp
}

type scanItem struct {
	idx   int
	pk    []byte
	value []byte
}

type noneIndexer struct {
}

func (ni *noneIndexer) IndexNames() map[string]bool {
	return nil
}

func (ni *noneIndexer) Index(value interface{}, data []byte) MultiIndex {
	return MultiIndex{}
}

func (ni *noneIndexer) UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int {
	return 0
}
