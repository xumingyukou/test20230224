package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func SaveToFile(filename string, dataList ...interface{}) {
	var (
		f          *os.File
		content    []byte
		contentSum []byte
		err        error
	)
	for _, data := range dataList {
		content, err = json.Marshal(data)
		if err != nil {
			fmt.Println("json dump err:", err)
			continue
		}
		contentSum = append(contentSum, content...)
		contentSum = append(contentSum, '\n')
	}

	if _, err = os.Stat(path.Dir(filename)); err != nil {
		err = os.MkdirAll(path.Dir(filename), 0777)
		if err != nil {
			fmt.Println("MkdirAll err:", err)
			return
		}
	}

	f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		fmt.Println("OpenFile err:", err)
		return
	}
	defer f.Close()
	_, err = f.Write(contentSum)
	if err != nil {
		fmt.Println("Write err:", err)
		return
	}
}
