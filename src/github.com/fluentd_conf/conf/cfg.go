package conf

import (
	"encoding/json"
	"github.com/astaxie/beego"
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
		beego.Emergency("os.Open err=", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Emergency("ioutil.ReadAll err=", err)
	}

	var cfg GolablConf
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		beego.Emergency("json.Unmarshal err=", err)
	}
	config = &cfg
}
