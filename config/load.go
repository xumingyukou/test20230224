package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// 从配置文件加载静态配置
func LoadConfigFromFile(file string, v interface{}) error {
	r, err := os.Open(file)
	fmt.Println(file)
	if err != nil {
		fmt.Println("read error")
		return err
	}

	d := toml.NewDecoder(r)
	//d.SetStrict(true)
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
