package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/fluentd_conf/conf"
	"github.com/fluentd_conf/models"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type ServiceController struct {
	beego.Controller
}

type FluentdServiceServer struct {
	Services map[string]string
	Servers  []string // fluentd server
}

// GetServicesList parses to get all services and the regexps of services
func (c *ServiceController) GetServicesList() {
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		beego.Emergency(err)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		return
	}

	var fluentdSS FluentdServiceServer
	fluentdSS.Services = make(map[string]string)

	var serviceName string
	var serviceNameRegexList []map[string]string
	ServiceKeys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range ServiceKeys {
		serviceName = strings.Split(key, ":")[1]
		serviceRegexp, _ := redis.String(con.Do("GET", key))
		fluentdSS.Services[serviceName] = serviceRegexp

		// serviceNameRegexList for response content
		serviceNameRegexList = append(serviceNameRegexList,
			map[string]string{
				"service": serviceName,
				"regexp":  serviceRegexp,
			})
	}

	fluentdSS.Servers = c.GetServersList()

	templ, err := template.ParseFiles(conf.Config().FileTemplate)
	if err != nil {
		beego.Emergency(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	var doc bytes.Buffer
	templ.Execute(&doc, fluentdSS)
	newConf := doc.String()
	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	ioutil.WriteFile(conf.Config().FilePath, []byte(newConf), 0)

	c.Data["json"] = serviceNameRegexList
	c.ServeJSON()
}

// AddServiceRegexp parses to add a service regexp
func (c *ServiceController) AddServiceRegexp() {
	var requestServiceRegexp models.RequestServiceRegexp
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestServiceRegexp); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	serviceName := requestServiceRegexp.ServiceName
	serviceRegexp := requestServiceRegexp.ServiceRegexp

	if serviceName == "" {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]interface{}{"success": false, "error": "service = nil"}
		c.ServeJSON()
		return
	}
	c.UpdateServiceRegexp(serviceName, serviceRegexp)
}

func (c *ServiceController) GetServiceRegexp() {
	serviceName := c.GetString(":service")
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		beego.Emergency(err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	serviceRegexp, err := redis.String(con.Do("GET", "service:"+serviceName))
	if err != nil {
		beego.Emergency(err)
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	c.Data["json"] = map[string]interface{}{"service": serviceName, "regexp": serviceRegexp}
	c.ServeJSON()
}

func (c *ServiceController) ChangeServiceRegexp() {
	serviceName := c.GetString(":service")

	var requestServiceRegexp models.RequestServiceRegexp
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestServiceRegexp); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(400)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	serviceRegexp := requestServiceRegexp.ServiceRegexp
	c.UpdateServiceRegexp(serviceName, serviceRegexp)
}

func (c *ServiceController) DeleteService() {
	serviceName := c.GetString(":service")
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		beego.Emergency(err)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	con.Do("DEL", "service:"+serviceName)

	var fluentdSS FluentdServiceServer
	fluentdSS.Services = make(map[string]string)

	ServiceKeys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range ServiceKeys {
		serviceName = strings.Split(key, ":")[1]
		serviceRegexp, _ := redis.String(con.Do("GET", key))
		fluentdSS.Services[serviceName] = serviceRegexp
	}

	fluentdSS.Servers = c.GetServersList()

	templ, err := template.ParseFiles(conf.Config().FileTemplate)
	if err != nil {
		beego.Emergency(err)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}

	var doc bytes.Buffer
	templ.Execute(&doc, fluentdSS)
	newConf := doc.String()
	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	ioutil.WriteFile(conf.Config().FilePath, []byte(newConf), 0)
	c.Data["json"] = map[string]interface{}{"success": true}
	c.ServeJSON()
}

func (c *ServiceController) GetServersList() []string {
	var serverInfo models.ServerInfo
	getJson(&serverInfo)

	var newServerConf []string
	for i := range serverInfo.Data {
		newServerConf = append(newServerConf, serverInfo.Data[i].Hostname)
	}

	return newServerConf
}

func (c *ServiceController) UpdateServiceRegexp(serviceName, serviceRegexp string) {
	con, err := redis.Dial(conf.Config().RedisConnectMethod, conf.Config().RedisAddressPort)
	defer con.Close()
	if err != nil {
		// redis connect error
		beego.Emergency(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}

	_, err = con.Do("SET", "service:"+serviceName, serviceRegexp)
	if err != nil {
		// redis set error
		beego.Emergency(err)
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}

	// get services and regexps
	var fluentdSS FluentdServiceServer
	fluentdSS.Services = make(map[string]string)

	ServiceKeys, _ := redis.Strings(con.Do("KEYS", "service:*"))
	for _, key := range ServiceKeys {
		serviceName = strings.Split(key, ":")[1]
		serviceRegexp, err := redis.String(con.Do("GET", key))
		if err != nil {
			beego.Emergency(err)
		}
		fluentdSS.Services[serviceName] = serviceRegexp
	}

	// get servers
	fluentdSS.Servers = c.GetServersList()

	templ, err := template.ParseFiles(conf.Config().FileTemplate)
	if err != nil {
		// template parse error
		c.Ctx.ResponseWriter.WriteHeader(500)
		c.Data["json"] = map[string]interface{}{"success": false, "error": err}
		c.ServeJSON()
		return
	}
	var doc bytes.Buffer
	templ.Execute(&doc, fluentdSS)
	newConf := doc.String()

	if !Exist(conf.Config().FilePath) {
		os.Create(conf.Config().FilePath)
	}
	// get new fluentd td-agent configure file
	ioutil.WriteFile(conf.Config().FilePath, []byte(newConf), 0)
	c.Data["json"] = map[string]interface{}{"success": true}
	c.ServeJSON()
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
