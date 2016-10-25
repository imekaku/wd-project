package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/wd-project/fluentd-conf-change-api/conf"
	"io/ioutil"
	"os"
	"regexp"
)

type OwnController struct {
	beego.Controller
}

func (c *OwnController) ChangeRegexp() {
	service_string := c.GetString("service")
	regexp_string := c.GetString("regexp")

	if service_string == "" {
		fmt.Println("service = nil")
		return
	}
	f, err := os.Open(conf.Config().FilePath)
	defer f.Close()
	if err != nil {
		fmt.Println("controllers ChangeRegexp os.Open err =", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("controllers ChangeRegexp ioutil.ReadAll err =", err)
	}

	r_service := regexp.MustCompile(`(?s)<match reformed\.docker\.` + service_string + `\.\*\*>(.*?)</match>`)
	service_match := r_service.FindAllString(string(data), 1)

	r_regexp := regexp.MustCompile(`regexp1 log .*`)
	new_service_match := r_regexp.ReplaceAllString(service_match[0], "regexp1 log "+regexp_string)

	new_conf := r_service.ReplaceAllString(string(data), new_service_match)
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}

func (c *OwnController) DeleteService() {
	service_string := c.GetString("service")
	if service_string == "" {
		fmt.Println("service = nil")
		return
	}
	f, err := os.Open(conf.Config().FilePath)
	defer f.Close()
	if err != nil {
		fmt.Println("controllers DeleteService os.Open err =", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("controllers DeleteService ioutil.ReadAll err =", err)
	}

	r_service := regexp.MustCompile(`(?s)<match reformed\.docker\.` + service_string + `\.\*\*>(.*?)</match>\n`)

	new_conf := r_service.ReplaceAllString(string(data), "")
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}

func (c *OwnController) AddService() {
	service_string := c.GetString("service")
	regexp_string := c.GetString("regexp")

	if service_string == "" {
		fmt.Println("service = nil")
		return
	}

	f, err := os.Open(conf.Config().FilePath)
	defer f.Close()
	if err != nil {
		fmt.Println("controllers ChangeRegexp os.Open err =", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("controllers ChangeRegexp ioutil.ReadAll err =", err)
	}

	r_service := regexp.MustCompile(`(?s)<match docker\.\*\*>(.*?)</match>`)
	service_match := r_service.FindAllString(string(data), 1)
	new_service_match := service_match[0] + "\n\n<match reformed.docker." + service_string + ".**>\n  type grep\n  regexp1 log " + regexp_string + "\n  add_tag_prefix regexp\n</match>\n"

	new_conf := r_service.ReplaceAllString(string(data), new_service_match)
	ioutil.WriteFile(conf.Config().FilePath, []byte(new_conf), 0)
}
