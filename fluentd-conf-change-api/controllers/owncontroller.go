package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"regexp"
)

type OwnController struct {
	beego.Controller
}

func (c *OwnController) ChangeRegexp() {
	server_string := c.GetString("server")
	regexp_string := c.GetString("regexp")

	if server_string == "" {
		fmt.Println("server = nil")
		return
	}
	f, err := os.Open("./views/conf.conf")
	defer f.Close()
	if err != nil {
		fmt.Println("controllers ChangeRegexp os.Open err =", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("controllers OpenFile ioutil.ReadAll err =", err)
	}

	r_server := regexp.MustCompile(`(?s)<match reformed\.docker\.` + server_string + `\.\*\*>(.*?)</match>`)
	service_match := r_server.FindAllString(string(data), 1)

	r_regexp := regexp.MustCompile(`regexp1 log .*`)
	new_service_match := r_regexp.ReplaceAllString(service_match[0], "regexp1 log "+regexp_string)

	new_conf := r_server.ReplaceAllString(string(data), new_service_match)
	ioutil.WriteFile("./views/conf.conf", []byte(new_conf), 0)
}
