package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/conf/conf"
	"github.com/conf/models"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type OwnController struct {
	beego.Controller
}

type FluentdServiceServer struct {
	Service map[string]string
	Server  []string
}

func (c *OwnController) ChangeRegexp() {
	service_string := c.GetString("service")
	regexp_string := c.GetString("regexp")
	if service_string == "" {
		fmt.Println("service = nil")
		return
	}

	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	con.Do("SET", "service:"+service_string, regexp_string)

	var fss FluentdServiceServer
	fss.Service = make(map[string]string)

	var service_name string
	service_keys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range service_keys {
		service_name = strings.Split(key, ":")[1]
		regexp_, err := redis.String(con.Do("GET", key))
		if err != nil {
			fmt.Println(err)
		}
		fss.Service[service_name] = regexp_
	}

	fss.Server = c.ServerList()

	templ, _ := template.ParseFiles(conf.Config().FileTemplate)
	var doc bytes.Buffer
	templ.Execute(&doc, fss)
	new_conf := doc.String()

	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}

func (c *OwnController) DeleteService() {
	service_string := c.GetString("service")
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	con.Do("DEL", "service:"+service_string)

	var fss FluentdServiceServer
	fss.Service = make(map[string]string)

	var service_name string
	service_keys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range service_keys {
		service_name = strings.Split(key, ":")[1]
		regexp_, _ := redis.String(con.Do("GET", service_name))
		fss.Service[service_name] = regexp_
	}

	fss.Server = c.ServerList()

	templ, _ := template.ParseFiles(conf.Config().FileTemplate)
	var doc bytes.Buffer
	templ.Execute(&doc, fss)
	new_conf := doc.String()
	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}

func (c *OwnController) Deploy() {
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	var fss FluentdServiceServer
	fss.Service = make(map[string]string)

	var service_name string
	service_keys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range service_keys {
		service_name = strings.Split(key, ":")[1]
		regexp_, _ := redis.String(con.Do("GET", service_name))
		fss.Service[service_name] = regexp_
	}

	fss.Server = c.ServerList()

	templ, _ := template.ParseFiles(conf.Config().FileTemplate)
	var doc bytes.Buffer
	templ.Execute(&doc, fss)
	new_conf := doc.String()
	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}

func (c *OwnController) ServerList() []string {
	var serverinfo models.ServerInfo
	getJson(&serverinfo)

	var new_server_conf []string
	for i := range serverinfo.Data {
		new_server_conf = append(new_server_conf, serverinfo.Data[i].Hostname)
	}

	return new_server_conf
}

func getJson(target *models.ServerInfo) error {
	r, err := http.Get(conf.Config().ServerUrl)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
