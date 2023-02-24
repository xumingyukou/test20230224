package sdk

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/pelletier/go-toml/v2"
)

// 用于将配置中的字节大小字符串转换成uint
type ByteSize uint64

// Byte size size suffixes.
const (
	B  ByteSize = 1
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
)

// Used for returning long unit form of string representation.
var longUnitMap = map[ByteSize]string{
	B:  "byte",
	KB: "kilobyte",
	MB: "megabyte",
	GB: "gigabyte",
	TB: "terabyte",
	PB: "petabyte",
	EB: "exabyte",
}

// Used for returning string representation.
var shortUnitMap = map[ByteSize]string{
	B:  "B",
	KB: "KB",
	MB: "MB",
	GB: "GB",
	TB: "TB",
	PB: "PB",
	EB: "EB",
}

// Used to convert user input to ByteSize
var unitMap = map[string]ByteSize{
	"B":     B,
	"BYTE":  B,
	"BYTES": B,

	"KB":        KB,
	"KILOBYTE":  KB,
	"KILOBYTES": KB,

	"MB":        MB,
	"MEGABYTE":  MB,
	"MEGABYTES": MB,

	"GB":        GB,
	"GIGABYTE":  GB,
	"GIGABYTES": GB,

	"TB":        TB,
	"TERABYTE":  TB,
	"TERABYTES": TB,

	"PB":        PB,
	"PETABYTE":  PB,
	"PETABYTES": PB,

	"EB":       EB,
	"EXABYTE":  EB,
	"EXABYTES": EB,
}

var (
	// Use long units, such as "megabytes" instead of "MB".
	LongUnits bool = false

	// String format of bytesize output. The unit of measure will be appended
	// to the end. Uses the same formatting options as the fmt package.
	Format string = "%.2f"
)

// Parse parses a byte size string. A byte size string is a number followed by
// a unit suffix, such as "1024B" or "1 MB". Valid byte units are "B", "KB",
// "MB", "GB", "TB", "PB" and "EB". You can also use the long
// format of units, such as "kilobyte" or "kilobytes".
func Parse(s string) (ByteSize, error) {
	// Remove leading and trailing whitespace
	s = strings.TrimSpace(s)

	split := make([]string, 0)
	for i, r := range s {
		if !unicode.IsDigit(r) && r != '.' {
			// Split the string by digit and size designator, remove whitespace
			split = append(split, strings.TrimSpace(string(s[:i])))
			split = append(split, strings.TrimSpace(string(s[i:])))
			break
		}
	}

	// Check to see if we split successfully
	if len(split) != 2 {
		return 0, errors.New("unrecognized size suffix")
	}

	// Check for MB, MEGABYTE, and MEGABYTES
	unit, ok := unitMap[strings.ToUpper(split[1])]
	if !ok {
		return 0, errors.New("Unrecognized size suffix " + split[1])

	}

	value, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return 0, err
	}

	bytesize := ByteSize(value * float64(unit))
	return bytesize, nil

}

// Satisfy the flag package  Value interface.
func (b *ByteSize) Set(s string) error {
	bs, err := Parse(s)
	if err != nil {
		return err
	}
	*b = bs
	return nil
}

// Satisfy the pflag package Value interface.
func (b *ByteSize) Type() string { return "byte_size" }

// Satisfy the encoding.TextUnmarshaler interface.
func (b *ByteSize) UnmarshalText(text []byte) error {
	return b.Set(string(text))
}

// Satisfy the flag package Getter interface.
func (b *ByteSize) Get() interface{} { return ByteSize(*b) }

// NewByteSize returns a new ByteSize type set to s.
func NewByteSize(s uint64) ByteSize {
	return ByteSize(s)
}

// Returns a string representation of b with the specified formatting and units.
func (b ByteSize) Format(format string, unit string, longUnits bool) string {
	return b.format(format, unit, longUnits)
}

// String returns the string form of b using the package global Format and
// LongUnits options.
func (b ByteSize) String() string {
	return b.format(Format, "", LongUnits)
}

func (b ByteSize) format(format string, unit string, longUnits bool) string {
	var unitSize ByteSize
	if unit != "" {
		var ok bool
		unitSize, ok = unitMap[strings.ToUpper(unit)]
		if !ok {
			return "Unrecognized unit: " + unit
		}
	} else {
		switch {
		case b >= EB:
			unitSize = EB
		case b >= PB:
			unitSize = PB
		case b >= TB:
			unitSize = TB
		case b >= GB:
			unitSize = GB
		case b >= MB:
			unitSize = MB
		case b >= KB:
			unitSize = KB
		default:
			unitSize = B
		}
	}

	if longUnits {
		var s string
		value := fmt.Sprintf(format, float64(b)/float64(unitSize))
		if printS, _ := strconv.ParseFloat(strings.TrimSpace(value), 64); printS > 0 && printS != 1 {
			s = "s"
		}
		return fmt.Sprintf(format+longUnitMap[unitSize]+s, float64(b)/float64(unitSize))
	}
	return fmt.Sprintf(format+shortUnitMap[unitSize], float64(b)/float64(unitSize))
}

// 用于将配置中的时长字符串转换成time.Duration结构
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type TtNatsConfig struct {
	Servers             []string `toml:"servers"`               //["nats://127.0.0.1:4222", "nats://127.0.0.1:4223"]
	ConnectionName      string   `toml:"connection_name"`       // 连接名称，用于nats系统监控和debug
	NkeysSeed           string   `toml:"nkeys_seed"`            // 用于nats认证的nkey token文件地址
	ConnectTimeout      Duration `toml:"connect_timeout"`       // 连接最大超时时间
	ReconnectWait       Duration `toml:"reconnect_wait"`        // 重连时间间隔
	ReconnectBufSize    int      `toml:"reconnect_buf_size"`    // 发送端重连buffer，用于连接断开时缓存消息，重连成功后会一次性发送之前缓存的所有消息
	MaxReconnects       int      `toml:"max_reconnects"`        // 最大重连次数，超过后彻底关闭连接
	PingInterval        Duration `toml:"ping_interval"`         // ping间隔时间
	MaxPingsOutstanding int      `toml:"max_pings_outstanding"` // 最大未回应ping数量
}

// msg/state持久化存储配置
type StoreConfig struct {
	Name        string                 `toml:"name"`
	StoreType   string                 `toml:"store_type"` // memory | disk
	StoreConfig map[string]interface{} `toml:"config"`     // 针对store的特殊配置
}

type PubConfig struct {
	HeartBeatInterval Duration `toml:"heartbeat_interval"`
	Persist           bool     `toml:"persist"`     // 开启后会持久化消息，并在每次重启后继续之前发送的序列号
	UniqueName        string   `toml:"unique_name"` //生产者唯一标识，如果不设置，创建Pub时会自动生成一个
}

type SubConfig struct {
	WorkerPoolSize int `toml:"worker_pool_size"` // 消费者回调的最大并发度，0表示不使用worker pool，并发度=订阅的subject数
	QueueSize      int `toml:"queue_size"`       // 使用workerPool时最大队列长度，队列满提交任务会阻塞
}

// 从配置文件加载静态配置
func LoadConfigFile(file string, v interface{}) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}

	d := toml.NewDecoder(r)
	d.DisallowUnknownFields()
	err = d.Decode(v)

	if err != nil {
		var details *toml.StrictMissingError
		if errors.As(err, &details) {
			fmt.Printf("%v", details.String())
		}
		return err
	}

	return nil
}

func LoadConfigFileWithoutStrict(file string, v interface{}) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}

	d := toml.NewDecoder(r)
	err = d.Decode(v)

	if err != nil {
		var details *toml.StrictMissingError
		if errors.As(err, &details) {
			fmt.Printf("%v", details.String())
		}
		return err
	}

	return nil
}

// 从字符串加载
func LoadConfigString(s string, v interface{}) error {
	r := strings.NewReader(s)
	d := toml.NewDecoder(r)
	d.DisallowUnknownFields()
	err := d.Decode(v)

	if err != nil {
		var details *toml.StrictMissingError
		if errors.As(err, &details) {
			fmt.Printf("%v", details.String())
		}
		return err
	}

	return nil
}
