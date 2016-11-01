package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type GolablConf struct {
	FilePath           string `json:file_path`
	ServerUrl          string `json:server_url`
	FileTemplate       string `json:file_template`
	RedisConnectMethod string `json:redis_connect_method`
	RedisAddressPort   string `json:redis_address_port`
}

var (
	config *GolablConf
)

func Config() *GolablConf {
	return config
}

func ParseConfig(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("os.Open err=", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("ioutil.ReadAll err=", err)
	}

	var cfg GolablConf
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
	}
	config = &cfg
}
