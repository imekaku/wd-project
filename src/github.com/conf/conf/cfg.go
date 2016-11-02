package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type GolablConf struct {
	FilePath           string `json:filepath`
	ServerUrl          string `json:serverurl`
	FileTemplate       string `json:filetemplate`
	RedisConnectMethod string `json:redisconnectmethod`
	RedisAddressPort   string `json:redisaddressport`
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
