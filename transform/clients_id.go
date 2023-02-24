package transform

import (
	"strconv"
	"strings"
)

// IdToClientId 序列化clientId：id + 'P' + producer
func IdToClientId(producer []byte, id int64) string {
	v := make([]byte, 0, len(producer)*2+8)
	v = strconv.AppendInt(v, id, 10)
	v = append(v, 'P')
	v = EncodeToBuf(v, producer)
	return string(v)
}

// ClientIdToId 反序列化ClientId为producer、id
func ClientIdToId(clientId string) (producer []byte, id int64) {
	idx := strings.IndexByte(clientId, 'P')
	if idx < 0 {
		return nil, 0
	}

	var err error
	id, err = strconv.ParseInt(clientId[:idx], 10, 64)
	if err != nil {
		return nil, 0
	}

	producer, err = DecodeString(clientId[idx+1:])
	if err != nil {
		return nil, 0
	}
	return
}
