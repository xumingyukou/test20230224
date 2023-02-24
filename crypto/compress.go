package crypto

import (
	"bytes"
	"clients/logger"
	"compress/gzip"
	"encoding/binary"
	"io/ioutil"
)

func GZIPDecompress(data []byte) ([]byte, error) {
	// fmt.Println(data)
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		logger.Logger.Info("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			logger.Logger.Warn("[ParseGzip]  ioutil.ReadAll error: %v", err)
			return nil, err
		}
		//fmt.Println("GZIP: ", string(undatas))
		return undatas, nil
	}
}
