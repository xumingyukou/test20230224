package store

import (
	"encoding/binary"
)

// 将多个IndexId和index值封装到一个字节数组，方便序列化和存储，只能一次写入，多次读取
// indexId使用crc16根据索引名字生成，使用者需要保证添加的索引名称在KvInderer.IndexNames()中
type MultiIndex struct {
	MaxNum uint16 // 最大索引个数
	n      uint16 // 当前写入个数
	d      []byte // ([offset](2字节)){n个，按照indexId从小到大排列}([indexLen](1字节)[indexId](2字节)index_data){n个}
}

// NOTE: 对于有符号整型或者变长索引，请使用IndexEncoder进行编解码，否则数据查询会错误
type KvIndexer interface {
	// 返回: 索引名->是否是唯一索引的键值对，如果索引名字为空，代表是主键
	IndexNames() map[string]bool
	// 输入：value,data: 对象和编码后的字节数组
	// 输出：MultiIndex，如果indexId=0表示主键
	Index(value interface{}, data []byte) MultiIndex
	// 解码index，对于联合索引，解码出一个元素调用一次cb，最后返回解码的整个长度
	UnpackIndex(indexName string, index []byte, cb func(idx int, v interface{})) int
}

type MemKvIndexer interface {
	KvIndexer
	// 返回某个index应该使用MEM_HASH_INDEX还是MEM_TREE_INDEX
	IndexType(name string) int
}

type KvSequencer interface {
	GetNextSequence(cacheSize uint64) uint64
	ResetSequence(initSeq uint64, cacheSize uint64) error
}

// key-value引擎存储的数据库操作接口
type KvDb interface {
	// 打开此db后才能使用
	Open() error

	// 关闭此db，关闭后数据刷盘且此对象无法使用
	Close() error

	// 创建表，name：表名，class：表数据类型, indexer：索引生成器，args：引擎自定义参数
	CreateTable(name string, class string, indexer KvIndexer, args map[string]interface{}) (table KvTable, err error)

	// 创建发号器，用于生成单调递增序列号
	// name: 发号器名称 initCacheSize：初始cache大小
	CreateSequencer(name string, initCacheSize uint64, args map[string]interface{}) (sequencer KvSequencer, err error)

	// 备份数据库到特定url，最少应当支持file://和http://
	Backup(url string) error

	// 返回库的统计信息
	Stat() string
}

// key-value引擎存储的表操作接口
// value表示序列化前的对象，data表示序列化后的对象
type KvTable interface {
	// 根据指定的index查询原始数据，如果index为空，表示根据主键进行查询
	Get(key []byte, index string) (pk []byte, data []byte, err error)

	// 根据指定的index查询解码后的数据，引擎可根据自身情况分别实现Get/GetValue
	GetValue(key []byte, index string) (pk []byte, value interface{}, err error)

	// 插入数据，如果设置了indexer需要输入value，data是否设置取决于引擎保存形式
	// 如果pk不为nil，优先使用pk作为主键，否则使用indexer生成主键
	Set(pk []byte, value interface{}, data []byte) error

	// 删除index上匹配key的所有数据
	Delete(key []byte, index string) error

	// 批量修改，需要符合原子性
	BatchSet(iter func() (pk []byte, value interface{}, data []byte)) error

	// 批量删除，需要符合原子性
	BatchDelete(iter func() (key []byte, index string)) error

	// 检索范围内所有的key-value，并调用回调函数，区间[start,end)
	Scan(start, end []byte, index string, f func(pk []byte, data []byte) error) error

	// 检索并返回内部数据表示
	ScanValue(start, end []byte, index string, f func(pk []byte, value interface{}) error) error

	// 返回表的统计信息
	Stat() string
}

func NewKvDb(mst string, config map[string]interface{}) KvDb {
	if mst == "memory" {
		return NewMemKvDb(config)
	} else if mst == "disk" {
		return NewDiskKvDb(config)
	} else {
		panic("Invalid kv db type")
	}
}

// 写之前需要调用
func (mi *MultiIndex) WriteInit(indexNumber uint16) {
	// [tag数量](2字节) + ([offset](2字节)+[长度](2字节)){n个} + ([indexId](2字节)+index){n个}
	mi.MaxNum = indexNumber + 1 // 多分配一个tag，用于后续添加pk
	mi.d = make([]byte, 2+4*(indexNumber+1), 128)
}

// 从读取前需要调用，初始化
func (mi *MultiIndex) ReadInit(d []byte) {
	mi.d = d
	mi.MaxNum = binary.BigEndian.Uint16(d)
	mi.n = mi.MaxNum
}

func (mi *MultiIndex) Write(name string, index []byte) {
	if mi.n >= mi.MaxNum || len(mi.d) > 65535 {
		panic("Too many index adding/too long index")
	}
	indexId := Crc16(name)
	var i int
	limit := 2 + int(mi.n)*4

	// indexId需要按照从小到大，相同的indexId会覆盖老的
	for i = 2; i < limit; i += 4 {
		o := binary.BigEndian.Uint16(mi.d[i:])
		id := binary.BigEndian.Uint16(mi.d[o:])
		if indexId < id {
			copy(mi.d[i+4:2+mi.MaxNum*4], mi.d[i:])
			break
		} else if indexId == id {
			mi.n--
			break
		}
	}

	binary.BigEndian.PutUint16(mi.d[i:], uint16(len(mi.d)))
	binary.BigEndian.PutUint16(mi.d[i+2:], uint16(2+len(index)))

	mi.d = append(mi.d, byte(indexId>>8), byte(indexId&0xff))
	mi.d = append(mi.d, index...)
	mi.n++
	binary.BigEndian.PutUint16(mi.d, mi.n)
}

//  单独返回indexId和index
func (mi *MultiIndex) Read(n uint16) (indexId uint16, index []byte) {
	if n >= mi.n {
		return 0, nil
	}
	hb := mi.d[2+4*n:]
	o := binary.BigEndian.Uint16(hb)
	l := binary.BigEndian.Uint16(hb[2:])

	indexId = binary.BigEndian.Uint16(mi.d[o:])
	index = mi.d[o+2 : o+l]

	return
}

// 返回indexId+index的连续字节序列
func (mi *MultiIndex) Read2(n uint16) (idIndex []byte) {
	if n >= mi.n {
		return nil
	}
	hb := mi.d[2+4*n:]
	o := binary.BigEndian.Uint16(hb)
	l := binary.BigEndian.Uint16(hb[2:])

	idIndex = mi.d[o : o+l]

	return
}

func (mi *MultiIndex) Bytes() []byte {
	return mi.d
}
